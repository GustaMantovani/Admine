package mcserver

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/internal/mc_server/models"
	"github.com/GustaMantovani/Admine/server_handler/pkg"
	"github.com/gorcon/rcon"
)

type DockerMinecraftServer struct {
	DockerCompose *pkg.DockerCompose
	DockerConfig  config.DockerConfig
}

func NewDockerMinecraftServer(compose *pkg.DockerCompose, dockerConfig config.DockerConfig) *DockerMinecraftServer {
	return &DockerMinecraftServer{
		DockerCompose: compose,
		DockerConfig:  dockerConfig,
	}
}

func (d *DockerMinecraftServer) Start(ctx context.Context) error {
	return d.DockerCompose.Up(true)
}

func (d *DockerMinecraftServer) Stop(ctx context.Context) error {
	done := make(chan error, 1)

	if _, err := d.ExecuteCommand(ctx, "/stop"); err != nil {
		return err
	}

	go func() {
		err := pkg.StreamContainerLogs(ctx, d.DockerConfig.ContainerName, func(line string) {
			slog.Debug("Container line:", "line", line)
			if strings.Contains(line, "All dimensions are saved") {
				done <- nil
			}
		})
		if err != nil {
			done <- err
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		if err != nil {
			return err
		}
	}

	return d.DockerCompose.Stop()
}

func (d *DockerMinecraftServer) Down(ctx context.Context) error {
	return d.DockerCompose.Down()
}

func (d *DockerMinecraftServer) Restart(ctx context.Context) error {
	if err := d.Stop(ctx); err != nil {
		return err
	}
	return d.Start(ctx)
}

// getServerUptime gets the server uptime using Docker exec
func (d *DockerMinecraftServer) getServerUptime(ctx context.Context) string {
	if d.DockerConfig.ServiceName == "" {
		return "N/A - No Service Name"
	}

	// Try to get container start time using docker inspect
	results, err := d.DockerCompose.ExecStructured([]string{"sh", "-c", "stat -c %Y /proc/1"}, d.DockerConfig.ServiceName)
	if err != nil || len(results) == 0 {
		return "N/A - Cannot Query Container"
	}

	println(results[d.DockerConfig.ServiceName])
	startTimeStr := strings.TrimSpace(results[d.DockerConfig.ServiceName])
	startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
	if err != nil {
		slog.Error("invalid timestamp", "err", err)
		return "N/A - Invalid Timestamp"
	}

	// Calculate uptime
	uptime := time.Since(time.Unix(startTime, 0))

	days := int(uptime.Hours()) / 24
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

func (d *DockerMinecraftServer) Status(ctx context.Context) (*models.ServerStatus, error) {
	// Test connectivity by executing list command
	listResult, err := d.ExecuteCommand(ctx, "list")
	if err != nil {
		return models.NewServerStatus(
			models.HealthUnknown,
			models.StatusOffline,
			"Server is offline - cannot connect via RCON",
			"N/A - Server Offline",
			0.0,
		), nil
	}

	listResponse := listResult.Output

	tps := 20.0 // Default
	if tpsResult, err := d.ExecuteCommand(ctx, "forge tps"); err == nil {
		tpsResponse := tpsResult.Output
		if strings.Contains(tpsResponse, "TPS") {
			// Parse TPS from response like "Mean tick time: 1.23 ms. Mean TPS: 19.84"
			tpsRegex := regexp.MustCompile(`Mean TPS:\s*([0-9]+\.?[0-9]*)`)
			if matches := tpsRegex.FindStringSubmatch(tpsResponse); len(matches) > 1 {
				if parsedTPS, parseErr := strconv.ParseFloat(matches[1], 64); parseErr == nil {
					tps = parsedTPS
				}
			}
		}
	} else {
		// Try vanilla TPS command
		if msptResult, err := d.ExecuteCommand(ctx, "mspt"); err == nil {
			msptResponse := msptResult.Output
			// Parse MSPT (milliseconds per tick) and convert to TPS
			// TPS = 1000 / MSPT (since 1 second = 1000ms and ideal is 20 TPS = 50ms per tick)
			msptRegex := regexp.MustCompile(`([0-9]+\.?[0-9]*)\s*ms`)
			if matches := msptRegex.FindStringSubmatch(msptResponse); len(matches) > 1 {
				if mspt, parseErr := strconv.ParseFloat(matches[1], 64); parseErr == nil && mspt > 0 {
					tps = 1000.0 / mspt
					if tps > 20.0 {
						tps = 20.0 // Cap at theoretical maximum
					}
				}
			}
		}
	}

	// Get server uptime using Docker exec
	uptime := d.getServerUptime(ctx)

	return models.NewServerStatus(
		models.HealthHealthy,
		models.StatusOnline,
		"Server is online - "+listResponse,
		uptime,
		tps,
	), nil
}

// getMinecraftVersion attempts to get the Minecraft version from various sources
func (d *DockerMinecraftServer) getMinecraftVersion(ctx context.Context) string {
	// Try version command
	if versionResult, err := d.ExecuteCommand(ctx, "version"); err == nil {
		versionResponse := versionResult.Output

		// Log the raw response for debugging
		slog.Debug("Version command response", "response", versionResponse)

		// Try multiple patterns for different server types
		patterns := []string{
			`MC:\s*([0-9]+\.[0-9]+\.?[0-9]*)`,        // CraftBukkit/Spigot: (MC: 1.20.1)
			`version\s+([0-9]+\.[0-9]+\.?[0-9]*)`,    // Some servers: version 1.20.1
			`Minecraft\s+([0-9]+\.[0-9]+\.?[0-9]*)`,  // Vanilla: Minecraft 1.20.1
			`([0-9]+\.[0-9]+\.?[0-9]*)\s+server`,     // Pattern: 1.20.1 server
			`running\s+.*?([0-9]+\.[0-9]+\.?[0-9]*)`, // Pattern: running ... 1.20.1
		}

		for _, pattern := range patterns {
			versionRegex := regexp.MustCompile(pattern)
			if matches := versionRegex.FindStringSubmatch(versionResponse); len(matches) > 1 {
				version := matches[1]
				// Validate that it's actually a version number, not a timestamp or other data
				if regexp.MustCompile(`^[0-9]+\.[0-9]+\.?[0-9]*$`).MatchString(version) {
					return version
				}
			}
		}
	}

	// Try to get from server.properties file
	if d.DockerConfig.ServiceName != "" {
		if results, err := d.DockerCompose.ExecStructured([]string{"sh", "-c", "grep -E '^minecraft-version=' /data/server.properties 2>/dev/null || echo 'N/A'"}, d.DockerConfig.ServiceName); err == nil {
			if output := strings.TrimSpace(results[d.DockerConfig.ServiceName]); output != "N/A" && output != "" {
				parts := strings.Split(output, "=")
				if len(parts) > 1 {
					version := strings.TrimSpace(parts[1])
					// Validate version format
					if regexp.MustCompile(`^[0-9]+\.[0-9]+\.?[0-9]*$`).MatchString(version) {
						return version
					}
				}
			}
		}
	}

	return "N/A - Version Unknown"
}

// getJavaVersion gets the Java version from the container
func (d *DockerMinecraftServer) getJavaVersion(ctx context.Context) string {
	if d.DockerConfig.ServiceName == "" {
		return "N/A - No Service Name"
	}

	results, err := d.DockerCompose.ExecStructured([]string{"java", "-version"}, d.DockerConfig.ServiceName)
	if err != nil || len(results) == 0 {
		return "N/A - Cannot Query Java"
	}

	output := results[d.DockerConfig.ServiceName]
	// Parse Java version from output like 'openjdk version "17.0.2" 2022-01-18'
	versionRegex := regexp.MustCompile(`version\s*"([^"]+)"`)
	if matches := versionRegex.FindStringSubmatch(output); len(matches) > 1 {
		return matches[1]
	}

	return "N/A - Java Version Unknown"
}

// getModEngine attempts to detect the mod engine/server type
func (d *DockerMinecraftServer) getModEngine(ctx context.Context) string {
	// Try version command to detect server type
	if versionResult, err := d.ExecuteCommand(ctx, "version"); err == nil {
		versionResponse := versionResult.Output
		versionLower := strings.ToLower(versionResponse)

		if strings.Contains(versionLower, "forge") {
			return "Forge"
		} else if strings.Contains(versionLower, "fabric") {
			return "Fabric"
		} else if strings.Contains(versionLower, "craftbukkit") || strings.Contains(versionLower, "bukkit") {
			return "CraftBukkit"
		} else if strings.Contains(versionLower, "spigot") {
			return "Spigot"
		} else if strings.Contains(versionLower, "paper") {
			return "Paper"
		} else if strings.Contains(versionLower, "purpur") {
			return "Purpur"
		} else if strings.Contains(versionLower, "mohist") {
			return "Mohist"
		} else if strings.Contains(versionLower, "vanilla") {
			return "Vanilla"
		}
	}

	// Try mods command to see if it's a modded server
	if _, err := d.ExecuteCommand(ctx, "mods"); err == nil {
		return "Modded - Type Unknown"
	}

	// Try forge command
	if _, err := d.ExecuteCommand(ctx, "forge"); err == nil {
		return "Forge"
	}

	return "N/A - Engine Unknown"
}

// getMaxPlayers gets the maximum player count from server configuration
func (d *DockerMinecraftServer) getMaxPlayers(ctx context.Context, listResponse string) int {
	// Parse from list command response like "There are 2 of a max of 20 players online"
	maxPlayersRegex := regexp.MustCompile(`max\s+of\s+([0-9]+)`)
	if matches := maxPlayersRegex.FindStringSubmatch(listResponse); len(matches) > 1 {
		if maxPlayers, err := strconv.Atoi(matches[1]); err == nil {
			return maxPlayers
		}
	}

	// Try to get from server.properties file
	if d.DockerConfig.ServiceName != "" {
		if results, err := d.DockerCompose.ExecStructured([]string{"sh", "-c", "grep -E '^max-players=' /data/server.properties 2>/dev/null || echo 'max-players=N/A'"}, d.DockerConfig.ServiceName); err == nil {
			if output := strings.TrimSpace(results[d.DockerConfig.ServiceName]); output != "max-players=N/A" && output != "" {
				parts := strings.Split(output, "=")
				if len(parts) > 1 {
					if maxPlayers, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						return maxPlayers
					}
				}
			}
		}
	}

	return -1 // Indicates unknown
}

func (d *DockerMinecraftServer) Info(ctx context.Context) (*models.ServerInfo, error) {
	// Test connectivity by executing a simple command
	if _, err := d.ExecuteCommand(ctx, "list"); err != nil {
		return &models.ServerInfo{
			MinecraftVersion: "N/A - Server Offline",
			JavaVersion:      "N/A - Server Offline",
			ModEngine:        "N/A - Server Offline",
			Seed:             "N/A - Server Offline",
			MaxPlayers:       -1,
		}, nil
	}

	// Get Minecraft version
	minecraftVersion := d.getMinecraftVersion(ctx)

	// Get Java version
	javaVersion := d.getJavaVersion(ctx)

	// Get mod engine/server type
	modEngine := d.getModEngine(ctx)

	// Get seed
	seed := "N/A - Seed Hidden"
	if seedResult, err := d.ExecuteCommand(ctx, "seed"); err == nil {
		seedResponse := seedResult.Output
		// Parse seed from response like "Seed: [1234567890]"
		seedRegex := regexp.MustCompile(`Seed:\s*\[([^\]]+)\]`)
		if matches := seedRegex.FindStringSubmatch(seedResponse); len(matches) > 1 {
			seed = matches[1]
		} else if strings.TrimSpace(seedResponse) != "" {
			// If regex doesn't match, use the full response (minus whitespace)
			seed = strings.TrimSpace(seedResponse)
		}
	}

	// Get max players from list command
	maxPlayers := -1
	if listResult, err := d.ExecuteCommand(ctx, "list"); err == nil {
		listResponse := listResult.Output
		maxPlayers = d.getMaxPlayers(ctx, listResponse)
	}

	return models.NewServerInfo(
		minecraftVersion,
		javaVersion,
		modEngine,
		seed,
		maxPlayers,
	), nil
}

func (d *DockerMinecraftServer) StartUpInfo(ctx context.Context) string {
	id, err := pkg.GetZeroTierNodeID(d.DockerConfig.ContainerName)
	if err != nil {
		return ""
	}

	return id
}

func (d *DockerMinecraftServer) ExecuteCommand(ctx context.Context, command string) (*models.CommandResult, error) {
	conn, err := rcon.Dial(d.DockerConfig.RconAddress, d.DockerConfig.RconPassword)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	response, err := conn.Execute(command)
	if err != nil {
		return nil, err
	}

	slog.Debug(response)

	return models.NewCommandResultWithOutput(response), nil
}
