package mcserver

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/internal/deployment"
	"github.com/GustaMantovani/Admine/server_handler/internal/mc_server/models"
	"github.com/GustaMantovani/Admine/server_handler/pkg"
	"github.com/gorcon/rcon"
)

type DockerMinecraftServer struct {
	DockerCompose *pkg.DockerCompose
	FullConfig    *config.Config
}

func NewDockerMinecraftServer(compose *pkg.DockerCompose, cfg *config.Config) *DockerMinecraftServer {
	return &DockerMinecraftServer{
		DockerCompose: compose,
		FullConfig:    cfg,
	}
}

func (d *DockerMinecraftServer) mcCfg() config.MinecraftServerConfig {
	return d.FullConfig.MinecraftServer
}

func (d *DockerMinecraftServer) Start(ctx context.Context) error {
	if err := deployment.GenerateDockerCompose(d.FullConfig); err != nil {
		return fmt.Errorf("failed to generate docker-compose file: %w", err)
	}
	return d.DockerCompose.Up(true)
}

func (d *DockerMinecraftServer) Stop(ctx context.Context) error {
	done := make(chan error, 1)

	if _, err := d.ExecuteCommand(ctx, "/stop"); err != nil {
		return err
	}

	go func() {
		err := pkg.StreamContainerLogs(ctx, d.mcCfg().Docker.ContainerName, func(line string) {
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
	time.Sleep(time.Duration(time.Duration.Seconds(1)));
	return d.Start(ctx)
}

// getServerUptime gets the server uptime using Docker exec
func (d *DockerMinecraftServer) getServerUptime(ctx context.Context) string {
	if d.mcCfg().Docker.ServiceName == "" {
		return "N/A - No Service Name"
	}

	results, err := d.DockerCompose.ExecStructured([]string{"sh", "-c", "stat -c %Y /proc/1"}, d.mcCfg().Docker.ServiceName)
	if err != nil || len(results) == 0 {
		return "N/A - Cannot Query Container"
	}

	startTimeStr := strings.TrimSpace(results[d.mcCfg().Docker.ServiceName])
	startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
	if err != nil {
		slog.Error("invalid timestamp", "err", err)
		return "N/A - Invalid Timestamp"
	}

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
			tpsRegex := regexp.MustCompile(`Mean TPS:\s*([0-9]+\.?[0-9]*)`)
			if matches := tpsRegex.FindStringSubmatch(tpsResponse); len(matches) > 1 {
				if parsedTPS, parseErr := strconv.ParseFloat(matches[1], 64); parseErr == nil {
					tps = parsedTPS
				}
			}
		}
	} else {
		if msptResult, err := d.ExecuteCommand(ctx, "mspt"); err == nil {
			msptResponse := msptResult.Output
			msptRegex := regexp.MustCompile(`([0-9]+\.?[0-9]*)\s*ms`)
			if matches := msptRegex.FindStringSubmatch(msptResponse); len(matches) > 1 {
				if mspt, parseErr := strconv.ParseFloat(matches[1], 64); parseErr == nil && mspt > 0 {
					tps = 1000.0 / mspt
					if tps > 20.0 {
						tps = 20.0
					}
				}
			}
		}
	}

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
	if versionResult, err := d.ExecuteCommand(ctx, "version"); err == nil {
		versionResponse := versionResult.Output

		slog.Debug("Version command response", "response", versionResponse)

		patterns := []string{
			`MC:\s*([0-9]+\.[0-9]+\.?[0-9]*)`,
			`version\s+([0-9]+\.[0-9]+\.?[0-9]*)`,
			`Minecraft\s+([0-9]+\.[0-9]+\.?[0-9]*)`,
			`([0-9]+\.[0-9]+\.?[0-9]*)\s+server`,
			`running\s+.*?([0-9]+\.[0-9]+\.?[0-9]*)`,
		}

		for _, pattern := range patterns {
			versionRegex := regexp.MustCompile(pattern)
			if matches := versionRegex.FindStringSubmatch(versionResponse); len(matches) > 1 {
				version := matches[1]
				if regexp.MustCompile(`^[0-9]+\.[0-9]+\.?[0-9]*$`).MatchString(version) {
					return version
				}
			}
		}
	}

	// Try to get from server.properties (itzg stores at /data/server.properties)
	if d.mcCfg().Docker.ServiceName != "" {
		if results, err := d.DockerCompose.ExecStructured([]string{"sh", "-c", "grep -E '^minecraft-version=' /data/server.properties 2>/dev/null || echo 'N/A'"}, d.mcCfg().Docker.ServiceName); err == nil {
			if output := strings.TrimSpace(results[d.mcCfg().Docker.ServiceName]); output != "N/A" && output != "" {
				parts := strings.Split(output, "=")
				if len(parts) > 1 {
					version := strings.TrimSpace(parts[1])
					if regexp.MustCompile(`^[0-9]+\.[0-9]+\.?[0-9]*$`).MatchString(version) {
						return version
					}
				}
			}
		}
	}

	// Try to get from VERSION env var set by itzg image
	if d.mcCfg().Docker.ServiceName != "" {
		if results, err := d.DockerCompose.ExecStructured([]string{"sh", "-c", "printenv VERSION 2>/dev/null || echo ''"}, d.mcCfg().Docker.ServiceName); err == nil {
			if version := strings.TrimSpace(results[d.mcCfg().Docker.ServiceName]); version != "" {
				if regexp.MustCompile(`^[0-9]+\.[0-9]+\.?[0-9]*$`).MatchString(version) {
					return version
				}
			}
		}
	}

	return "N/A - Version Unknown"
}

// getJavaVersion gets the Java version from the container
func (d *DockerMinecraftServer) getJavaVersion(ctx context.Context) string {
	if d.mcCfg().Docker.ServiceName == "" {
		return "N/A - No Service Name"
	}

	results, err := d.DockerCompose.ExecStructured([]string{"java", "-version"}, d.mcCfg().Docker.ServiceName)
	if err != nil || len(results) == 0 {
		return "N/A - Cannot Query Java"
	}

	output := results[d.mcCfg().Docker.ServiceName]
	versionRegex := regexp.MustCompile(`version\s*"([^"]+)"`)
	if matches := versionRegex.FindStringSubmatch(output); len(matches) > 1 {
		return matches[1]
	}

	return "N/A - Java Version Unknown"
}

// getModEngine attempts to detect the mod engine/server type
func (d *DockerMinecraftServer) getModEngine(ctx context.Context) string {
	if d.mcCfg().Docker.ServiceName != "" {
		envChecks := []struct {
			envVar string
			engine string
		}{
			{"FABRIC_LOADER_VERSION", "Fabric"},
			{"FORGE_VERSION", "Forge"},
			{"NEOFORGE_VERSION", "NeoForge"},
			{"QUILT_VERSION", "Quilt"},
		}

		for _, check := range envChecks {
			cmd := fmt.Sprintf("printenv %s 2>/dev/null || echo ''", check.envVar)
			if results, err := d.DockerCompose.ExecStructured([]string{"sh", "-c", cmd}, d.mcCfg().Docker.ServiceName); err == nil {
				if value := strings.TrimSpace(results[d.mcCfg().Docker.ServiceName]); value != "" {
					return fmt.Sprintf("%s %s", check.engine, value)
				}
			}
		}
	}

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

	if _, err := d.ExecuteCommand(ctx, "mods"); err == nil {
		return "Modded - Type Unknown"
	}

	if _, err := d.ExecuteCommand(ctx, "forge"); err == nil {
		return "Forge"
	}

	return "N/A - Engine Unknown"
}

// getMaxPlayers gets the maximum player count from server configuration
func (d *DockerMinecraftServer) getMaxPlayers(ctx context.Context, listResponse string) int {
	maxPlayersRegex := regexp.MustCompile(`max\s+of\s+([0-9]+)`)
	if matches := maxPlayersRegex.FindStringSubmatch(listResponse); len(matches) > 1 {
		if maxPlayers, err := strconv.Atoi(matches[1]); err == nil {
			return maxPlayers
		}
	}

	// itzg stores server.properties at /data/server.properties
	if d.mcCfg().Docker.ServiceName != "" {
		if results, err := d.DockerCompose.ExecStructured([]string{"sh", "-c", "grep -E '^max-players=' /data/server.properties 2>/dev/null || echo 'max-players=N/A'"}, d.mcCfg().Docker.ServiceName); err == nil {
			if output := strings.TrimSpace(results[d.mcCfg().Docker.ServiceName]); output != "max-players=N/A" && output != "" {
				parts := strings.Split(output, "=")
				if len(parts) > 1 {
					if maxPlayers, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						return maxPlayers
					}
				}
			}
		}
	}

	return -1
}

func (d *DockerMinecraftServer) Info(ctx context.Context) (*models.ServerInfo, error) {
	if _, err := d.ExecuteCommand(ctx, "list"); err != nil {
		return &models.ServerInfo{
			MinecraftVersion: "N/A - Server Offline",
			JavaVersion:      "N/A - Server Offline",
			ModEngine:        "N/A - Server Offline",
			Seed:             "N/A - Server Offline",
			MaxPlayers:       -1,
		}, nil
	}

	minecraftVersion := d.FullConfig.MinecraftServer.Image.Version
	javaVersion := d.getJavaVersion(ctx)
	modEngine := d.getModEngine(ctx)

	seed := "N/A - Seed Hidden"
	if seedResult, err := d.ExecuteCommand(ctx, "seed"); err == nil {
		seedResponse := seedResult.Output
		seedRegex := regexp.MustCompile(`Seed:\s*\[([^\]]+)\]`)
		if matches := seedRegex.FindStringSubmatch(seedResponse); len(matches) > 1 {
			seed = matches[1]
		} else if strings.TrimSpace(seedResponse) != "" {
			seed = strings.TrimSpace(seedResponse)
		}
	}

	maxPlayers := -1
	if listResult, err := d.ExecuteCommand(ctx, "list"); err == nil {
		maxPlayers = d.getMaxPlayers(ctx, listResult.Output)
	}

	return models.NewServerInfo(
		minecraftVersion,
		javaVersion,
		modEngine,
		seed,
		maxPlayers,
	), nil
}

func (d *DockerMinecraftServer) Logs(ctx context.Context, n int) ([]string, error) {
	if d.mcCfg().Docker.ServiceName != "" {
		logs, err := d.DockerCompose.ReadLastServiceLogs(uint(n), d.mcCfg().Docker.ServiceName)
		if err == nil {
			return logs, nil
		}

		slog.Warn("Failed to read logs for configured service, retrying without service filter", "service", d.mcCfg().Docker.ServiceName, "error", err)
		fallbackLogs, fallbackErr := d.DockerCompose.ReadLastServiceLogs(uint(n))
		if fallbackErr == nil {
			return fallbackLogs, nil
		}

		return nil, fmt.Errorf("failed to read logs for service %q: %w; fallback failed: %v", d.mcCfg().Docker.ServiceName, err, fallbackErr)
	}

	return d.DockerCompose.ReadLastServiceLogs(uint(n))
}

func (d *DockerMinecraftServer) StartUpInfo(ctx context.Context) string {
	ztContainerName := d.mcCfg().ZeroTier.ContainerName
	if !d.mcCfg().ZeroTier.Enabled || ztContainerName == "" {
		return ""
	}

	id, err := pkg.GetZeroTierNodeID(ztContainerName)
	if err != nil {
		return ""
	}

	return id
}

func (d *DockerMinecraftServer) ExecuteCommand(ctx context.Context, command string) (*models.CommandResult, error) {
	conn, err := rcon.Dial(d.mcCfg().RconAddress, d.mcCfg().RconPassword)
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

func (d *DockerMinecraftServer) InstallMod(ctx context.Context, fileName string, modData io.Reader) (*models.ModInstallResult, error) {
	serviceName := d.mcCfg().Docker.ServiceName

	tmpFile, err := os.CreateTemp("", "admine-mod-*.jar")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, modData); err != nil {
		return nil, fmt.Errorf("failed to write mod data to temp file: %w", err)
	}
	tmpFile.Close()

	// itzg/minecraft-server stores mods at /data/mods/
	destPath := fmt.Sprintf("%s:/data/mods/%s", serviceName, fileName)
	baseArgs := []string{"compose", "-f", d.DockerCompose.File, "cp", tmpFile.Name(), destPath}
	cmd := exec.CommandContext(ctx, "docker", baseArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Failed to copy mod to container", "error", err, "output", string(output))
		return models.NewModInstallResult(fileName, false, fmt.Sprintf("Failed to copy mod: %s", string(output))), err
	}

	slog.Info("Mod installed successfully", "file", fileName, "service", serviceName)
	return models.NewModInstallResult(fileName, true, "Mod installed successfully"), nil
}

func (d *DockerMinecraftServer) ListMods(ctx context.Context) (*models.ModListResult, error) {
	serviceName := d.mcCfg().Docker.ServiceName

	// itzg/minecraft-server stores mods at /data/mods/
	results, err := d.DockerCompose.ExecStructured([]string{"sh", "-c", "ls /data/mods/"}, serviceName)
	if err != nil {
		slog.Error("Failed to list mods", "error", err)
		return nil, fmt.Errorf("failed to list mods: %w", err)
	}

	output := strings.TrimSpace(results[serviceName])
	lines := strings.Split(output, "\n")
	var mods []string
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name != "" && strings.HasSuffix(strings.ToLower(name), ".jar") {
			mods = append(mods, name)
		}
	}

	slog.Info("Listed mods", "count", len(mods), "service", serviceName)
	return models.NewModListResult(mods), nil
}

func (d *DockerMinecraftServer) RemoveMod(ctx context.Context, fileName string) (*models.ModInstallResult, error) {
	serviceName := d.mcCfg().Docker.ServiceName

	// itzg/minecraft-server stores mods at /data/mods/
	modPath := fmt.Sprintf("/data/mods/%s", fileName)
	results, err := d.DockerCompose.ExecStructured([]string{"sh", "-c", fmt.Sprintf("rm %s", modPath)}, serviceName)
	if err != nil {
		output := results[serviceName]
		slog.Error("Failed to remove mod", "error", err, "output", output)
		return models.NewModInstallResult(fileName, false, fmt.Sprintf("Failed to remove mod: %s", output)), err
	}

	slog.Info("Mod removed successfully", "file", fileName, "service", serviceName)
	return models.NewModInstallResult(fileName, true, "Mod removed successfully"), nil
}
