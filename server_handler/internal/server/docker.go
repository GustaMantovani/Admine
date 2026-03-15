package server

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
	"github.com/GustaMantovani/Admine/server_handler/internal/docker"
	"github.com/gorcon/rcon"
)

type dockerMinecraftServer struct {
	compose *docker.DockerCompose
	cfg     config.MinecraftServerConfig
}

func newDockerMinecraftServer(compose *docker.DockerCompose, cfg config.MinecraftServerConfig) *dockerMinecraftServer {
	return &dockerMinecraftServer{
		compose: compose,
		cfg:     cfg,
	}
}

func (d *dockerMinecraftServer) Start(ctx context.Context) error {
	fullCfg := &config.Config{MinecraftServer: d.cfg}
	if err := deployment.GenerateDockerCompose(fullCfg); err != nil {
		return fmt.Errorf("failed to generate docker-compose file: %w", err)
	}
	return d.compose.Up(true)
}

func (d *dockerMinecraftServer) Stop(ctx context.Context) error {
	done := make(chan error, 1)

	if _, err := d.ExecuteCommand(ctx, "/stop"); err != nil {
		slog.Warn("RCON stop failed, falling back to docker compose down", "error", err)
		return d.compose.Down()
	}

	go func() {
		err := docker.StreamContainerLogs(ctx, d.cfg.Docker.ContainerName, func(line string) {
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

	return d.compose.Stop()
}

func (d *dockerMinecraftServer) Down(ctx context.Context) error {
	return d.compose.Down()
}

func (d *dockerMinecraftServer) Restart(ctx context.Context) error {
	if err := d.Stop(ctx); err != nil {
		return err
	}
	return d.Start(ctx)
}

func (d *dockerMinecraftServer) getServerUptime(ctx context.Context) string {
	if d.cfg.Docker.ServiceName == "" {
		return "N/A - No Service Name"
	}

	results, err := d.compose.ExecStructured([]string{"sh", "-c", "stat -c %Y /proc/1"}, d.cfg.Docker.ServiceName)
	if err != nil || len(results) == 0 {
		return "N/A - Cannot Query Container"
	}

	startTimeStr := strings.TrimSpace(results[d.cfg.Docker.ServiceName])
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

func (d *dockerMinecraftServer) Status(ctx context.Context) (*ServerStatus, error) {
	listResult, err := d.ExecuteCommand(ctx, "list")
	if err != nil {
		return NewServerStatus(
			HealthUnknown,
			StatusOffline,
			"Server is offline - cannot connect via RCON",
			"N/A - Server Offline",
			0.0,
		), nil
	}

	listResponse := listResult.Output

	tps := 20.0
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

	return NewServerStatus(
		HealthHealthy,
		StatusOnline,
		"Server is online - "+listResponse,
		uptime,
		tps,
	), nil
}

func (d *dockerMinecraftServer) Info(ctx context.Context) (*ServerInfo, error) {
	img := d.cfg.Image

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

	return NewServerInfo(
		img.Version,
		d.getJavaVersion(),
		modEngine,
		d.getSeed(ctx),
		d.getMaxPlayers(ctx),
	), nil
}

func (d *dockerMinecraftServer) getSeed(ctx context.Context) string {
	seedResult, err := d.ExecuteCommand(ctx, "seed")
	if err != nil {
		return "N/A - Server Offline"
	}
	seedRegex := regexp.MustCompile(`Seed:\s*\[([^\]]+)\]`)
	if matches := seedRegex.FindStringSubmatch(seedResult.Output); len(matches) > 1 {
		return matches[1]
	}
	if trimmed := strings.TrimSpace(seedResult.Output); trimmed != "" {
		return trimmed
	}
	return "N/A - Seed Hidden"
}

func (d *dockerMinecraftServer) getJavaVersion() string {
	if d.cfg.Docker.ServiceName == "" {
		return "N/A - No Service Name"
	}
	results, err := d.compose.ExecStructured([]string{"java", "-version"}, d.cfg.Docker.ServiceName)
	if err != nil || len(results) == 0 {
		return "N/A - Cannot Query Java"
	}
	output := results[d.cfg.Docker.ServiceName]
	versionRegex := regexp.MustCompile(`version\s*"([^"]+)"`)
	if matches := versionRegex.FindStringSubmatch(output); len(matches) > 1 {
		return matches[1]
	}
	return "N/A - Java Version Unknown"
}

func (d *dockerMinecraftServer) getMaxPlayers(ctx context.Context) int {
	listResult, err := d.ExecuteCommand(ctx, "list")
	if err != nil {
		return -1
	}
	maxRegex := regexp.MustCompile(`max of (\d+)`)
	if matches := maxRegex.FindStringSubmatch(listResult.Output); len(matches) > 1 {
		if n, err := strconv.Atoi(matches[1]); err == nil {
			return n
		}
	}
	return -1
}

func (d *dockerMinecraftServer) Logs(ctx context.Context, n int) ([]string, error) {
	if d.cfg.Docker.ServiceName != "" {
		logs, err := d.compose.ReadLastServiceLogs(uint(n), d.cfg.Docker.ServiceName)
		if err == nil {
			return logs, nil
		}

		slog.Warn("Failed to read logs for configured service, retrying without service filter", "service", d.cfg.Docker.ServiceName, "error", err)
		fallbackLogs, fallbackErr := d.compose.ReadLastServiceLogs(uint(n))
		if fallbackErr == nil {
			return fallbackLogs, nil
		}

		return nil, fmt.Errorf("failed to read logs for service %q: %w; fallback failed: %v", d.cfg.Docker.ServiceName, err, fallbackErr)
	}

	return d.compose.ReadLastServiceLogs(uint(n))
}

func (d *dockerMinecraftServer) StartUpInfo(ctx context.Context) string {
	ztContainerName := d.cfg.ZeroTier.ContainerName
	if !d.cfg.ZeroTier.Enabled || ztContainerName == "" {
		return ""
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		id, err := docker.GetZeroTierNodeID(ztContainerName)
		if err == nil {
			return id
		}
		slog.Warn("ZeroTier not ready yet, retrying", "container", ztContainerName, "error", err)

		select {
		case <-ctx.Done():
			slog.Error("Timed out waiting for ZeroTier node ID", "container", ztContainerName)
			return ""
		case <-ticker.C:
		}
	}
}

func (d *dockerMinecraftServer) ExecuteCommand(ctx context.Context, command string) (*CommandResult, error) {
	conn, err := rcon.Dial(d.cfg.RconAddress, d.cfg.RconPassword)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	response, err := conn.Execute(command)
	if err != nil {
		return nil, err
	}

	slog.Debug(response)

	return NewCommandResultWithOutput(response), nil
}

func (d *dockerMinecraftServer) InstallMod(ctx context.Context, fileName string, modData io.Reader) (*ModInstallResult, error) {
	serviceName := d.cfg.Docker.ServiceName

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

	destPath := fmt.Sprintf("%s:/data/mods/%s", serviceName, fileName)
	baseArgs := []string{"compose", "-f", d.compose.File, "cp", tmpFile.Name(), destPath}
	cmd := exec.CommandContext(ctx, "docker", baseArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Failed to copy mod to container", "error", err, "output", string(output))
		return NewModInstallResult(fileName, false, fmt.Sprintf("Failed to copy mod: %s", string(output))), err
	}

	slog.Info("Mod installed successfully", "file", fileName, "service", serviceName)
	return NewModInstallResult(fileName, true, "Mod installed successfully"), nil
}

func (d *dockerMinecraftServer) ListMods(ctx context.Context) (*ModListResult, error) {
	serviceName := d.cfg.Docker.ServiceName

	results, err := d.compose.ExecStructured([]string{"sh", "-c", "ls /data/mods/"}, serviceName)
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
	return NewModListResult(mods), nil
}

func (d *dockerMinecraftServer) RemoveMod(ctx context.Context, fileName string) (*ModInstallResult, error) {
	serviceName := d.cfg.Docker.ServiceName

	modPath := fmt.Sprintf("/data/mods/%s", fileName)
	results, err := d.compose.ExecStructured([]string{"sh", "-c", fmt.Sprintf("rm %s", modPath)}, serviceName)
	if err != nil {
		output := results[serviceName]
		slog.Error("Failed to remove mod", "error", err, "output", output)
		return NewModInstallResult(fileName, false, fmt.Sprintf("Failed to remove mod: %s", output)), err
	}

	slog.Info("Mod removed successfully", "file", fileName, "service", serviceName)
	return NewModInstallResult(fileName, true, "Mod removed successfully"), nil
}
