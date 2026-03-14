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
	time.Sleep(time.Duration(time.Duration.Seconds(1)))
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

func (d *DockerMinecraftServer) Info(_ context.Context) (*models.ServerInfo, error) {
	img := d.mcCfg().Image

	// JavaVersion: parse tag "java21" → "21"; empty means the image default is used.
	javaVersion := "N/A - Not tracked in config"
	if img.JavaVersion != "" {
		javaVersion = strings.TrimPrefix(img.JavaVersion, "java")
	}

	// ModEngine: prefer pinned loader versions, then extra_env, then TYPE.
	modEngine := img.Type
	switch {
	case img.FabricLoaderVersion != "":
		modEngine = "Fabric " + img.FabricLoaderVersion
	case img.ForgeVersion != "":
		modEngine = "Forge " + img.ForgeVersion
	default:
		if v, ok := img.ExtraEnv["NEOFORGE_VERSION"]; ok && v != "" {
			modEngine = "NeoForge " + v
		} else if v, ok := img.ExtraEnv["QUILT_LOADER_VERSION"]; ok && v != "" {
			modEngine = "Quilt " + v
		}
	}

	// MaxPlayers: read from extra_env MAX_PLAYERS if set.
	maxPlayers := -1
	if v, ok := img.ExtraEnv["MAX_PLAYERS"]; ok {
		if n, err := strconv.Atoi(v); err == nil {
			maxPlayers = n
		}
	}

	return models.NewServerInfo(
		img.Version,
		javaVersion,
		modEngine,
		"N/A - Seed Hidden",
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
