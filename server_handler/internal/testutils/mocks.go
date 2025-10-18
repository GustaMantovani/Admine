package testutils

import (
	"context"

	mcserver "github.com/GustaMantovani/Admine/server_handler/internal/mc_server"
	mcmodels "github.com/GustaMantovani/Admine/server_handler/internal/mc_server/models"
	pubsubmodels "github.com/GustaMantovani/Admine/server_handler/internal/pubsub/models"
	"github.com/stretchr/testify/mock"
)

// MockMinecraftServer is a shared mock implementation of the MinecraftServer interface for testing
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

// MockPubSubService is a shared mock implementation of PubSubService for testing
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

// AsInterface returns the mock as a MinecraftServer interface
func (m *MockMinecraftServer) AsInterface() mcserver.MinecraftServer {
	return m
}
