package testutils

import (
	"context"
	"io"

	"github.com/GustaMantovani/Admine/server_handler/internal/pubsub"
	"github.com/GustaMantovani/Admine/server_handler/internal/server"
	"github.com/stretchr/testify/mock"
)

// MockMinecraftServer is a mock implementation of server.MinecraftServer for testing
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

func (m *MockMinecraftServer) Status(ctx context.Context) (*server.ServerStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*server.ServerStatus), args.Error(1)
}

func (m *MockMinecraftServer) Info(ctx context.Context) (*server.ServerInfo, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*server.ServerInfo), args.Error(1)
}

func (m *MockMinecraftServer) Logs(ctx context.Context, n int) ([]string, error) {
	args := m.Called(ctx, n)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockMinecraftServer) StartUpInfo(ctx context.Context) string {
	args := m.Called(ctx)
	return args.String(0)
}

func (m *MockMinecraftServer) ExecuteCommand(ctx context.Context, command string) (*server.CommandResult, error) {
	args := m.Called(ctx, command)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*server.CommandResult), args.Error(1)
}

func (m *MockMinecraftServer) InstallMod(ctx context.Context, fileName string, modData io.Reader) (*server.ModInstallResult, error) {
	args := m.Called(ctx, fileName, modData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*server.ModInstallResult), args.Error(1)
}

func (m *MockMinecraftServer) ListMods(ctx context.Context) (*server.ModListResult, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*server.ModListResult), args.Error(1)
}

func (m *MockMinecraftServer) RemoveMod(ctx context.Context, fileName string) (*server.ModInstallResult, error) {
	args := m.Called(ctx, fileName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*server.ModInstallResult), args.Error(1)
}

// MockPubSubService is a mock implementation of pubsub.PubSubService for testing
type MockPubSubService struct {
	mock.Mock
}

func (m *MockPubSubService) Publish(topic string, msg *pubsub.AdmineMessage) error {
	args := m.Called(topic, msg)
	return args.Error(0)
}

func (m *MockPubSubService) Subscribe(topics ...string) (<-chan *pubsub.AdmineMessage, error) {
	args := m.Called(topics)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(chan *pubsub.AdmineMessage), args.Error(1)
}

func (m *MockPubSubService) Close() error {
	args := m.Called()
	return args.Error(0)
}
