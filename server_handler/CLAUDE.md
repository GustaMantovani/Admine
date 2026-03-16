# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build          # Compile to ./bin/server_handler
make run            # Build and run with server_handler_config.yaml
make dev            # Run without compiling (go run)
make test           # Run all tests with -v
make test-coverage  # Generate coverage.html
make check          # vet + golangci-lint + staticcheck
make fmt            # go fmt ./...
make setup          # Install golangci-lint and staticcheck
```

Run a single test or package:
```bash
go test -v ./internal/api/handlers/ -run TestGetStatus_Success
go test -v ./internal/pubsub/
```

Accept an optional config path as first argument:
```bash
./bin/server_handler /path/to/custom_config.yaml
```

## Architecture

`cmd/server_handler/main.go` is pure wiring — it constructs every object via explicit constructors and injects dependencies. There is no global state.

The two main entry points for behaviour are:
- **`internal/pubsub/handler.go` (`EventHandler`)** — routes Redis pub/sub messages (`server_on`, `server_off`, `server_down`, `restart`, `command`) to `MinecraftServer` methods.
- **`internal/api/handlers/`** — stateless Gin handlers; each receives its dependencies via `New*Handler(srv, ...)`.

**`MinecraftServer` interface** (`internal/server/server.go`) is the central seam. The only implementation is `dockerMinecraftServer` (`internal/server/docker.go`), which wraps RCON for commands and `docker.DockerCompose` for lifecycle.

**Compose generation** (`internal/deployment/`): on every `Start`, a Go template (`docker-compose.yaml.tmpl`, embedded via `//go:embed`) is rendered with the config and written to `docker.compose_output_path`. Edit the template to change the Compose structure; no Go code changes required. The template supports both ZeroTier and Tailscale sidecars via `{{ .ZeroTier.Enabled }}` / `{{ .Tailscale.Enabled }}` conditionals — only one should be enabled at a time.

**VPN sidecar**: `StartUpInfo()` in `internal/server/docker.go` detects which provider is enabled and uses the appropriate function to retrieve the node identifier after startup. For ZeroTier: `docker.GetZeroTierNodeID()` (parses `zerotier-cli info`). For Tailscale: `docker.GetTailscaleNodeKey()` (parses `tailscale status --json`, returns `Self.PublicKey` in `nodekey:XXXX` format — accepted directly by the Tailscale API).

**`AdmineMessage`** (`internal/pubsub/message.go`) is the envelope for all pub/sub traffic: `{origin, tags, message}`. `EventHandler` ignores its own messages via the `origin` field match.

## Testing

Tests use `testutils.MockMinecraftServer` (`internal/testutils/mocks.go`) built with `testify/mock`. HTTP handler tests create a `gin.TestMode` context directly — no running server needed. Pub/sub handler tests likewise inject the mock directly.
