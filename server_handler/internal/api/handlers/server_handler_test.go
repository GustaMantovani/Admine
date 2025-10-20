package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GustaMantovani/Admine/server_handler/internal/api/models"
	mc_models "github.com/GustaMantovani/Admine/server_handler/internal/mc_server/models"
	"github.com/GustaMantovani/Admine/server_handler/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestHandler creates a new test handler with necessary setup
func setupTestHandler() *ServerHandler {
	testutils.SetupGinTestMode()
	return NewApiHandler()
}

// TestGetInfo_Success tests successful retrieval of server info
func TestGetInfo_Success(t *testing.T) {
	// Setup
	handler := setupTestHandler()
	mockServer := new(testutils.MockMinecraftServer)
	ctx := context.Background()

	expectedInfo := mc_models.NewServerInfo(
		"1.20.1",
		"Java 17",
		"Forge",
		"12345678",
		20,
	)

	mockServer.On("Info", ctx).Return(expectedInfo, nil)

	// Create test context with mock
	testutils.SetupTestContextForAPI(mockServer)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/info", nil)

	// Execute
	handler.GetInfo(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response mc_models.ServerInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedInfo.MinecraftVersion, response.MinecraftVersion)
	assert.Equal(t, expectedInfo.JavaVersion, response.JavaVersion)
	assert.Equal(t, expectedInfo.ModEngine, response.ModEngine)
	assert.Equal(t, expectedInfo.MaxPlayers, response.MaxPlayers)
	assert.Equal(t, expectedInfo.Seed, response.Seed)

	mockServer.AssertExpectations(t)
}

// TestGetInfo_InfoError tests GetInfo when Info() returns an error
func TestGetInfo_InfoError(t *testing.T) {
	// Setup
	handler := setupTestHandler()
	mockServer := new(testutils.MockMinecraftServer)
	ctx := context.Background()

	expectedError := errors.New("failed to retrieve server info")
	mockServer.On("Info", ctx).Return(nil, expectedError)

	// Create test context with mock
	testutils.SetupTestContextForAPI(mockServer)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/info", nil)

	// Execute
	handler.GetInfo(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "Failed to get server info")
	assert.Contains(t, response.Message, expectedError.Error())

	mockServer.AssertExpectations(t)
}

// TestGetStatus_Success tests successful retrieval of server status
func TestGetStatus_Success(t *testing.T) {
	// Setup
	handler := setupTestHandler()
	mockServer := new(testutils.MockMinecraftServer)
	ctx := context.Background()

	expectedStatus := mc_models.NewServerStatus(
		mc_models.HealthHealthy,
		mc_models.StatusOnline,
		"Server is running smoothly",
		"2h 30m",
		19.8,
	)

	mockServer.On("Status", ctx).Return(expectedStatus, nil)

	// Create test context with mock
	testutils.SetupTestContextForAPI(mockServer)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/status", nil)

	// Execute
	handler.GetStatus(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response mc_models.ServerStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedStatus.Health, response.Health)
	assert.Equal(t, expectedStatus.Status, response.Status)
	assert.Equal(t, expectedStatus.Description, response.Description)
	assert.Equal(t, expectedStatus.Uptime, response.Uptime)
	assert.Equal(t, expectedStatus.TPS, response.TPS)

	mockServer.AssertExpectations(t)
}

// TestGetStatus_ServerNotInitialized tests GetStatus when server is nil
func TestGetStatus_ServerNotInitialized(t *testing.T) {
	// Setup
	handler := setupTestHandler()

	// Create test context with nil server
	testutils.SetupTestContextForAPIWithNilServer()

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/status", nil)

	// Execute
	handler.GetStatus(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Minecraft server not initialized", response.Message)
}

// TestGetStatus_StatusError tests GetStatus when Status() returns an error
func TestGetStatus_StatusError(t *testing.T) {
	// Setup
	handler := setupTestHandler()
	mockServer := new(testutils.MockMinecraftServer)
	ctx := context.Background()

	expectedError := errors.New("failed to retrieve server status")
	mockServer.On("Status", ctx).Return(nil, expectedError)

	// Create test context with mock
	testutils.SetupTestContextForAPI(mockServer)

	// Create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/status", nil)

	// Execute
	handler.GetStatus(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "Failed to get server status")
	assert.Contains(t, response.Message, expectedError.Error())

	mockServer.AssertExpectations(t)
}

// TestPostCommand_Success tests successful command execution
func TestPostCommand_Success(t *testing.T) {
	// Setup
	handler := setupTestHandler()
	mockServer := new(testutils.MockMinecraftServer)
	ctx := context.Background()

	testCommand := "list"
	expectedResult := mc_models.NewCommandResultWithOutput("There are 3 players online: Player1, Player2, Player3")

	mockServer.On("ExecuteCommand", ctx, testCommand).Return(expectedResult, nil)

	// Create test context with mock
	testutils.SetupTestContextForAPI(mockServer)

	// Create request
	commandPayload := models.Command{Command: testCommand}
	jsonPayload, _ := json.Marshal(commandPayload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBuffer(jsonPayload))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.PostCommand(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response mc_models.CommandResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult.Output, response.Output)
	assert.Nil(t, response.ExitCode)

	mockServer.AssertExpectations(t)
}

// TestPostCommand_InvalidJSON tests PostCommand with invalid JSON
func TestPostCommand_InvalidJSON(t *testing.T) {
	// Setup
	handler := setupTestHandler()

	// Create test context
	testutils.SetupTestContextForAPIWithNilServer()

	// Create request with invalid JSON
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBufferString("{invalid json}"))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.PostCommand(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "Invalid request")
}

// TestPostCommand_MissingCommandField tests PostCommand with missing command field
func TestPostCommand_MissingCommandField(t *testing.T) {
	// Setup
	handler := setupTestHandler()

	// Create test context
	testutils.SetupTestContextForAPIWithNilServer()

	// Create request with missing command field
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBufferString("{}"))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.PostCommand(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "Invalid request")
}

// TestPostCommand_ServerNotInitialized tests PostCommand when server is nil
func TestPostCommand_ServerNotInitialized(t *testing.T) {
	// Setup
	handler := setupTestHandler()

	// Create test context with nil server
	testutils.SetupTestContextForAPIWithNilServer()

	// Create request
	commandPayload := models.Command{Command: "list"}
	jsonPayload, _ := json.Marshal(commandPayload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBuffer(jsonPayload))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.PostCommand(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Minecraft server not initialized", response.Message)
}

// TestPostCommand_ExecuteCommandError tests PostCommand when ExecuteCommand returns an error
func TestPostCommand_ExecuteCommandError(t *testing.T) {
	// Setup
	handler := setupTestHandler()
	mockServer := new(testutils.MockMinecraftServer)
	ctx := context.Background()

	testCommand := "invalid-command"
	expectedError := errors.New("command execution failed")

	mockServer.On("ExecuteCommand", ctx, testCommand).Return(nil, expectedError)

	// Create test context with mock
	testutils.SetupTestContextForAPI(mockServer)

	// Create request
	commandPayload := models.Command{Command: testCommand}
	jsonPayload, _ := json.Marshal(commandPayload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBuffer(jsonPayload))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.PostCommand(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "Failed to execute command")
	assert.Contains(t, response.Message, expectedError.Error())

	mockServer.AssertExpectations(t)
}

// TestPostCommand_WithExitCode tests successful command execution with exit code
func TestPostCommand_WithExitCode(t *testing.T) {
	// Setup
	handler := setupTestHandler()
	mockServer := new(testutils.MockMinecraftServer)
	ctx := context.Background()

	testCommand := "say Hello World"
	exitCode := 0
	expectedResult := mc_models.NewCommandResultWithExitCode("Broadcast message", exitCode)

	mockServer.On("ExecuteCommand", ctx, testCommand).Return(expectedResult, nil)

	// Create test context with mock
	testutils.SetupTestContextForAPI(mockServer)

	// Create request
	commandPayload := models.Command{Command: testCommand}
	jsonPayload, _ := json.Marshal(commandPayload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBuffer(jsonPayload))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	handler.PostCommand(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response mc_models.CommandResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult.Output, response.Output)
	assert.NotNil(t, response.ExitCode)
	assert.Equal(t, exitCode, *response.ExitCode)

	mockServer.AssertExpectations(t)
}
