package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal"
	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	mcserver "github.com/GustaMantovani/Admine/server_handler/internal/mc_server"
	"github.com/gin-gonic/gin"
)

// SetupTestContext creates a test AppContext with mocks for event handler tests
func SetupTestContext(t *testing.T, mockServer *MockMinecraftServer) (*internal.AppContext, context.CancelFunc) {
	mainCtx, cancel := context.WithCancel(context.Background())

	mockContext := &internal.AppContext{
		MainCtx: &mainCtx,
		Config: &config.Config{
			App: config.AppConfig{
				SelfOriginName: "test_server",
			},
			PubSub: config.PubSubConfig{
				AdmineChannelsMap: config.AdmineChannelsMap{
					ServerChannel:  "test_server_channel",
					CommandChannel: "test_command_channel",
					VpnChannel:     "test_vpn_channel",
				},
			},
			MinecraftServer: config.MinecraftServerConfig{
				ServerOnTimeout:          5 * time.Second,
				ServerOffTimeout:         5 * time.Second,
				ServerCommandExecTimeout: 5 * time.Second,
			},
		},
	}

	if mockServer != nil {
		var server mcserver.MinecraftServer = mockServer
		mockContext.MinecraftServer = &server
	}

	// Set the instance for internal.Get()
	internal.SetInstanceForTest(mockContext)

	// Cleanup function
	t.Cleanup(func() {
		cancel()
	})

	return mockContext, cancel
}

// SetupTestContextForAPI creates a test AppContext for API handler tests
func SetupTestContextForAPI(mockServer *MockMinecraftServer) *internal.AppContext {
	ctx := context.Background()

	var mcServer mcserver.MinecraftServer
	if mockServer != nil {
		mcServer = mockServer
	}

	appCtx := &internal.AppContext{
		MinecraftServer: &mcServer,
		MainCtx:         &ctx,
	}
	internal.SetInstanceForTest(appCtx)

	return appCtx
}

// SetupTestContextForAPIWithNilServer creates a test AppContext with nil server for API handler tests
func SetupTestContextForAPIWithNilServer() *internal.AppContext {
	ctx := context.Background()

	appCtx := &internal.AppContext{
		MinecraftServer: nil,
		MainCtx:         &ctx,
	}
	internal.SetInstanceForTest(appCtx)

	return appCtx
}

// SetupGinTestMode sets gin to test mode (call once per test file)
func SetupGinTestMode() {
	gin.SetMode(gin.TestMode)
}
