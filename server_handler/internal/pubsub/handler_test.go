package pubsub_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/internal/pubsub"
	"github.com/GustaMantovani/Admine/server_handler/internal/server"
	"github.com/GustaMantovani/Admine/server_handler/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testCfg() config.MinecraftServerConfig {
	return config.MinecraftServerConfig{
		ServerOnTimeout:          5 * time.Second,
		ServerOffTimeout:         5 * time.Second,
		ServerCommandExecTimeout: 5 * time.Second,
	}
}

func newTestHandler(srv *testutils.MockMinecraftServer, ps *testutils.MockPubSubService) *pubsub.EventHandler {
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel
	return pubsub.NewEventHandler(srv, ps, "test_server", "test_server_channel", testCfg(), ctx)
}

func TestManageCommand_ServerNotInitialized(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	handler := pubsub.NewEventHandler(nil, mockPubSub, "test_server", "test_server_channel", testCfg(), context.Background())

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Server not initialized"
	})).Return(nil)

	msg := pubsub.NewAdmineMessage("origin", []string{"server_on"}, "")
	err := handler.ManageCommand(msg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "minecraft server is not initialized")
	mockPubSub.AssertExpectations(t)
}

func TestManageCommand_InvalidTag(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Invalid tag."
	})).Return(nil)

	msg := pubsub.NewAdmineMessage("origin", []string{"invalid_tag"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
}

func TestServerUp_Success(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Starting server"
	})).Return(nil).Once()

	mockServer.On("Start", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	mockServer.On("StartUpInfo", mock.AnythingOfType("*context.timerCtx")).Return("zerotier:1234567890").Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("server_on") && msg.Message == "zerotier:1234567890"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"server_on"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestServerUp_StartFailure(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Starting server"
	})).Return(nil).Once()

	mockServer.On("Start", mock.AnythingOfType("*context.timerCtx")).Return(errors.New("start failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Failed to start server: start failed"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"server_on"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestServerOff_Success(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Stopping server"
	})).Return(nil).Once()

	mockServer.On("Stop", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("server_off") && msg.Message == "Server stopped successfully"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"server_off"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestServerOff_StopFailure(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Stopping server"
	})).Return(nil).Once()

	mockServer.On("Stop", mock.AnythingOfType("*context.timerCtx")).Return(errors.New("stop failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Error stopping server: stop failed"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"server_off"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestServerDown_Success(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Removing server"
	})).Return(nil).Once()

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "/stop").
		Return(&server.CommandResult{Output: "Stopping server..."}, nil).Once()
	mockServer.On("Down", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("server_off") && msg.Message == "Server removed successfully"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"server_down"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestServerDown_StopCommandFailure(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Removing server"
	})).Return(nil).Once()

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "/stop").
		Return((*server.CommandResult)(nil), errors.New("command failed")).Once()
	mockServer.On("Down", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("server_off") && msg.Message == "Server removed successfully"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"server_down"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestServerDown_DownFailure(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Removing server"
	})).Return(nil).Once()

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "/stop").
		Return(&server.CommandResult{Output: "Stopping..."}, nil).Once()
	mockServer.On("Down", mock.AnythingOfType("*context.timerCtx")).Return(errors.New("down failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Error removing server: down failed"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"server_down"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestRestart_Success(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Restarting server"
	})).Return(nil).Once()

	mockServer.On("Stop", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	mockServer.On("Start", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	mockServer.On("StartUpInfo", mock.AnythingOfType("*context.timerCtx")).Return("zerotier:1234567890").Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("server_on") && msg.Message == "zerotier:1234567890"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"restart"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestRestart_StopFailure(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Restarting server"
	})).Return(nil).Once()

	mockServer.On("Stop", mock.AnythingOfType("*context.timerCtx")).Return(errors.New("stop failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Failed to restart server: stop failed"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"restart"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestRestart_StartFailure(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Restarting server"
	})).Return(nil).Once()

	mockServer.On("Stop", mock.AnythingOfType("*context.timerCtx")).Return(nil).Once()
	mockServer.On("Start", mock.AnythingOfType("*context.timerCtx")).Return(errors.New("start failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Failed to start server after stop: start failed"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"restart"}, "")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestCommand_SuccessWithOutput(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "list").
		Return(&server.CommandResult{Output: "There are 0 of a max of 20 players online"}, nil).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("command_result") && msg.Message == "There are 0 of a max of 20 players online"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"command"}, "list")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestCommand_SuccessWithoutOutput(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "clear").
		Return(&server.CommandResult{Output: ""}, nil).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("command_result") && msg.Message == "Command executed successfully"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"command"}, "clear")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}

func TestCommand_Failure(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestHandler(mockServer, mockPubSub)

	mockServer.On("ExecuteCommand", mock.AnythingOfType("*context.timerCtx"), "invalid").
		Return((*server.CommandResult)(nil), errors.New("command failed")).Once()

	mockPubSub.On("Publish", "test_server_channel", mock.MatchedBy(func(msg *pubsub.AdmineMessage) bool {
		return msg.HasTag("notification") && msg.Message == "Failed to execute command: command failed"
	})).Return(nil).Once()

	msg := pubsub.NewAdmineMessage("origin", []string{"command"}, "invalid")
	err := handler.ManageCommand(msg)

	assert.NoError(t, err)
	mockPubSub.AssertExpectations(t)
	mockServer.AssertExpectations(t)
}
