package pubsub

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal"
	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	mcserver "github.com/GustaMantovani/Admine/server_handler/internal/mc_server"
	mcmodels "github.com/GustaMantovani/Admine/server_handler/internal/mc_server/models"
	pubsubmodels "github.com/GustaMantovani/Admine/server_handler/internal/pubsub/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPubSubService is a mock implementation of PubSubService
type MockPubSubService struct {
	mock.Mock
}

func (m *MockPubSubService) Publish(topic string, msg *pubsubmodels.AdmineMessage) error {
	args := m.Called(topic, msg)
	return args.Error(0)
}

func (m *MockPubSubService) Subscribe(topics ...string) (<-chan *pubsubmodels.AdmineMessage, error) {
	args := m.Called(topics)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(chan *pubsubmodels.AdmineMessage), args.Error(1)
}

func (m *MockPubSubService) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockMinecraftServer is a mock implementation of MinecraftServer interface
type MockMinecraftServer struct {
	mock.Mock
}

func (m *MockMinecraftServer) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMinecraftServer) Stop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMinecraftServer) Down(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMinecraftServer) Restart(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMinecraftServer) Status(ctx context.Context) (*mcmodels.ServerStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mcmodels.ServerStatus), args.Error(1)
}

func (m *MockMinecraftServer) Info(ctx context.Context) (*mcmodels.ServerInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mcmodels.ServerInfo), args.Error(1)
}

func (m *MockMinecraftServer) StartUpInfo(ctx context.Context) string {
	args := m.Called(ctx)
	return args.String(0)
}

func (m *MockMinecraftServer) ExecuteCommand(ctx context.Context, command string) (*mcmodels.CommandResult, error) {
	args := m.Called(ctx, command)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mcmodels.CommandResult), args.Error(1)
}

// setupTestContext creates a test AppContext with mocks
func setupTestContext(t *testing.T, mockServer *MockMinecraftServer) (*internal.AppContext, context.CancelFunc) {
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

// TestManageCommand_ServerNotInitialized tests the case when MinecraftServer is nil
func TestManageCommand_ServerNotInitialized(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	_, cancel := setupTestContext(t, nil) // nil server
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Server not initialized"
	})).Return(nil)

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"server_on"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "minecraft server is not initialized")
	mockPubSub.AssertExpectations(t)
}

// TestManageCommand_InvalidTag tests handling of invalid tags
func TestManageCommand_InvalidTag(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Invalid tag."
	})).Return(nil)

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"invalid_tag"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
}

// TestServerUp_Success tests successful server startup
func TestServerUp_Success(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Starting server"
	})).Return(nil).Once()

	mockServer.On("Start", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	mockServer.On("StartUpInfo", mock.AnythingOfType("*context.timerCtx")).Return("zerotier:1234567890").Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("server_on") && msg.Message == "zerotier:1234567890"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"server_on"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestServerUp_StartFailure tests server startup failure
func TestServerUp_StartFailure(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Starting server"
	})).Return(nil).Once()

	mockServer.On("Start", mock.AnythingOfType("*context.timerCtx")).Return(errors.New("start failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Failed to start server: start failed"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"server_on"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestServerOff_Success tests successful server shutdown
func TestServerOff_Success(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Stopping server"
	})).Return(nil).Once()

	mockServer.On("Stop", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("server_off") && msg.Message == "Server stopped successfully"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"server_off"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestServerOff_StopFailure tests server shutdown failure
func TestServerOff_StopFailure(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Stopping server"
	})).Return(nil).Once()

	mockServer.On("Stop", mock.AnythingOfType("*context.timerCtx")).Return(errors.New("stop failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Error stopping server: stop failed"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"server_off"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestServerDown_Success tests successful server removal
func TestServerDown_Success(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Removing server"
	})).Return(nil).Once()

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "/stop").
		Return(&mcmodels.CommandResult{Output: "Stopping server..."}, nil).Once()

	mockServer.On("Down", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("server_off") && msg.Message == "Server removed successfully"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"server_down"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestServerDown_StopCommandFailure tests server removal when stop command fails
func TestServerDown_StopCommandFailure(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Removing server"
	})).Return(nil).Once()

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "/stop").
		Return((*mcmodels.CommandResult)(nil), errors.New("command failed")).Once()

	mockServer.On("Down", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("server_off") && msg.Message == "Server removed successfully"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"server_down"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestServerDown_DownFailure tests server removal failure
func TestServerDown_DownFailure(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Removing server"
	})).Return(nil).Once()

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "/stop").
		Return(&mcmodels.CommandResult{Output: "Stopping..."}, nil).Once()

	mockServer.On("Down", mock.AnythingOfType("*context.timerCtx")).Return(errors.New("down failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Error removing server: down failed"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"server_down"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestRestart_Success tests successful server restart
func TestRestart_Success(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Restarting server"
	})).Return(nil).Once()

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "/stop").
		Return(&mcmodels.CommandResult{Output: "Stopping..."}, nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"restart"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestRestart_Failure tests server restart failure
func TestRestart_Failure(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Restarting server"
	})).Return(nil).Once()

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "/stop").
		Return((*mcmodels.CommandResult)(nil), errors.New("restart failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Failed to restart server: restart failed"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"restart"}, "")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestCommand_SuccessWithOutput tests successful command execution with output
func TestCommand_SuccessWithOutput(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "list").
		Return(&mcmodels.CommandResult{Output: "There are 0 of a max of 20 players online"}, nil).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("command_result") && msg.Message == "There are 0 of a max of 20 players online"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"command"}, "list")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestCommand_SuccessWithoutOutput tests successful command execution without output
func TestCommand_SuccessWithoutOutput(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "clear").
		Return(&mcmodels.CommandResult{Output: ""}, nil).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("command_result") && msg.Message == "Command executed successfully"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"command"}, "clear")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

// TestCommand_Failure tests command execution failure
func TestCommand_Failure(t *testing.T) {
	// Setup
	mockPubSub := new(MockPubSubService)
	mockServer := new(MockMinecraftServer)
	_, cancel := setupTestContext(t, mockServer)
	defer cancel()

	handler := NewEventHandler(mockPubSub)

	// Set expectations
	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "invalid").
		Return((*mcmodels.CommandResult)(nil), errors.New("command failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsubmodels.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Failed to execute command: command failed"
	})).Return(nil).Once()

	// Execute
	msg := pubsubmodels.NewAdmineMessage([]string{"command"}, "invalid")
	err := handler.ManageCommand(msg)

	// Assert
	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}
