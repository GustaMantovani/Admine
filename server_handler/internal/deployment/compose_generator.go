package deployment

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
)

// composeMinecraftData holds Minecraft image configuration for the template.
type composeMinecraftData struct {
	Type                string
	Version             string
	RconPassword        string
	FabricLoaderVersion string
	ForgeVersion        string
	ModpackURL          string
	JavaVersion         string
	ExtraEnv            map[string]string
}

// composeZeroTierData holds ZeroTier sidecar configuration for the template.
type composeZeroTierData struct {
	Enabled       bool
	NetworkID     string
	ContainerName string
	ApiSecret     string
}

// composeTailscaleData holds Tailscale sidecar configuration for the template.
type composeTailscaleData struct {
	Enabled       bool
	AuthKey       string
	Hostname      string
	ContainerName string
}

// composeDockerData holds Docker runtime configuration for the template.
type composeDockerData struct {
	ContainerName string
	DataPath      string
}

// composeTemplateData is the root data structure passed to the docker-compose template.
type composeTemplateData struct {
	Minecraft composeMinecraftData
	ZeroTier  composeZeroTierData
	Tailscale composeTailscaleData
	Docker    composeDockerData
}

// GenerateDockerCompose renders the docker-compose template using cfg and writes
// the result to cfg.MinecraftServer.Docker.ComposeOutputPath, creating any missing
// parent directories automatically.
func GenerateDockerCompose(cfg *config.Config) error {
	mc := cfg.MinecraftServer

	data := composeTemplateData{
		Minecraft: composeMinecraftData{
			Type:                mc.Image.Type,
			Version:             mc.Image.Version,
			RconPassword:        mc.RconPassword,
			FabricLoaderVersion: mc.Image.FabricLoaderVersion,
			ForgeVersion:        mc.Image.ForgeVersion,
			ModpackURL:          mc.Image.ModpackURL,
			JavaVersion:         mc.Image.JavaVersion,
			ExtraEnv:            mc.Image.ExtraEnv,
		},
		ZeroTier: composeZeroTierData{
			Enabled:       mc.ZeroTier.Enabled,
			NetworkID:     mc.ZeroTier.NetworkID,
			ContainerName: mc.ZeroTier.ContainerName,
			ApiSecret:     mc.ZeroTier.ApiSecret,
		},
		Tailscale: composeTailscaleData{
			Enabled:       mc.Tailscale.Enabled,
			AuthKey:       mc.Tailscale.AuthKey,
			Hostname:      mc.Tailscale.Hostname,
			ContainerName: mc.Tailscale.ContainerName,
		},
		Docker: composeDockerData{
			ContainerName: mc.Docker.ContainerName,
			DataPath:      mc.Docker.DataPath,
		},
	}

	tmpl, err := template.New("docker-compose").Parse(composeTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse compose template: %w", err)
	}

	outputPath := mc.Docker.ComposeOutputPath
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("failed to create output directory for compose file: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create compose output file %q: %w", outputPath, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to render compose template: %w", err)
	}

	return nil
}
