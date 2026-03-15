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
	"github.com/GustaMantovani/Admine/server_handler/internal/server"
	"github.com/GustaMantovani/Admine/server_handler/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetInfo_Success(t *testing.T) {
	mockServer := new(testutils.MockMinecraftServer)
	handler := NewServerHandler(mockServer)

	expected := server.NewServerInfo("1.20.1", "Java 17", "Forge", "12345678", 20)
	mockServer.On("Info", context.Background()).Return(expected, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/info", nil)

	handler.GetInfo(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response server.ServerInfo
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, expected.MinecraftVersion, response.MinecraftVersion)
	assert.Equal(t, expected.JavaVersion, response.JavaVersion)
	assert.Equal(t, expected.ModEngine, response.ModEngine)
	assert.Equal(t, expected.MaxPlayers, response.MaxPlayers)
	assert.Equal(t, expected.Seed, response.Seed)

	mockServer.AssertExpectations(t)
}

func TestGetInfo_InfoError(t *testing.T) {
	mockServer := new(testutils.MockMinecraftServer)
	handler := NewServerHandler(mockServer)

	expectedError := errors.New("failed to retrieve server info")
	mockServer.On("Info", context.Background()).Return(nil, expectedError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/info", nil)

	handler.GetInfo(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response.Message, "Failed to get server info")
	assert.Contains(t, response.Message, expectedError.Error())

	mockServer.AssertExpectations(t)
}

func TestGetInfo_ServerNotInitialized(t *testing.T) {
	handler := NewServerHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/info", nil)

	handler.GetInfo(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Minecraft server not initialized", response.Message)
}

func TestGetStatus_Success(t *testing.T) {
	mockServer := new(testutils.MockMinecraftServer)
	handler := NewServerHandler(mockServer)

	expected := server.NewServerStatus(server.HealthHealthy, server.StatusOnline, "Server is running", "2h 30m", 19.8)
	mockServer.On("Status", context.Background()).Return(expected, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/status", nil)

	handler.GetStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response server.ServerStatus
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, expected.Health, response.Health)
	assert.Equal(t, expected.Status, response.Status)
	assert.Equal(t, expected.Description, response.Description)
	assert.Equal(t, expected.Uptime, response.Uptime)
	assert.Equal(t, expected.TPS, response.TPS)

	mockServer.AssertExpectations(t)
}

func TestGetStatus_ServerNotInitialized(t *testing.T) {
	handler := NewServerHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/status", nil)

	handler.GetStatus(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Minecraft server not initialized", response.Message)
}

func TestGetStatus_StatusError(t *testing.T) {
	mockServer := new(testutils.MockMinecraftServer)
	handler := NewServerHandler(mockServer)

	expectedError := errors.New("failed to retrieve server status")
	mockServer.On("Status", context.Background()).Return(nil, expectedError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/status", nil)

	handler.GetStatus(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response.Message, "Failed to get server status")
	assert.Contains(t, response.Message, expectedError.Error())

	mockServer.AssertExpectations(t)
}

func TestGetLogs_Success(t *testing.T) {
	mockServer := new(testutils.MockMinecraftServer)
	handler := NewServerHandler(mockServer)

	expectedLogs := []string{
		"[12:00:01] [Server thread/INFO]: Starting minecraft server",
		"[12:00:05] [Server thread/INFO]: Done",
	}
	mockServer.On("Logs", context.Background(), 50).Return(expectedLogs, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/logs?n=50", nil)

	handler.GetLogs(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.LogsResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, expectedLogs, response.Lines)
	assert.Equal(t, len(expectedLogs), response.Total)

	mockServer.AssertExpectations(t)
}

func TestGetLogs_DefaultN(t *testing.T) {
	mockServer := new(testutils.MockMinecraftServer)
	handler := NewServerHandler(mockServer)

	expectedLogs := []string{"line1", "line2"}
	mockServer.On("Logs", context.Background(), 100).Return(expectedLogs, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/logs", nil)

	handler.GetLogs(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockServer.AssertExpectations(t)
}

func TestGetLogs_InvalidN(t *testing.T) {
	tests := []struct {
		name string
		n    string
	}{
		{name: "non-integer", n: "abc"},
		{name: "zero", n: "0"},
		{name: "too-large", n: "101"},
		{name: "negative", n: "-5"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := new(testutils.MockMinecraftServer)
			handler := NewServerHandler(mockServer)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", "/logs?n="+tc.n, nil)

			handler.GetLogs(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response models.ErrorResponse
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
			assert.Contains(t, response.Message, "Invalid query param 'n'")

			mockServer.AssertNotCalled(t, "Logs", mock.Anything, mock.Anything)
		})
	}
}

func TestGetLogs_ServerNotInitialized(t *testing.T) {
	handler := NewServerHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/logs?n=10", nil)

	handler.GetLogs(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Minecraft server not initialized", response.Message)
}

func TestGetLogs_LogsError(t *testing.T) {
	mockServer := new(testutils.MockMinecraftServer)
	handler := NewServerHandler(mockServer)

	expectedError := errors.New("failed to retrieve logs")
	mockServer.On("Logs", context.Background(), 25).Return(nil, expectedError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/logs?n=25", nil)

	handler.GetLogs(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response.Message, "Failed to get server logs")
	assert.Contains(t, response.Message, expectedError.Error())

	mockServer.AssertExpectations(t)
}

func TestPostCommand_Success(t *testing.T) {
	mockServer := new(testutils.MockMinecraftServer)
	handler := NewServerHandler(mockServer)

	testCommand := "list"
	expectedResult := server.NewCommandResultWithOutput("There are 3 players online")
	mockServer.On("ExecuteCommand", context.Background(), testCommand).Return(expectedResult, nil)

	payload, _ := json.Marshal(models.Command{Command: testCommand})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBuffer(payload))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PostCommand(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response server.CommandResult
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, expectedResult.Output, response.Output)
	assert.Nil(t, response.ExitCode)

	mockServer.AssertExpectations(t)
}

func TestPostCommand_InvalidJSON(t *testing.T) {
	handler := NewServerHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBufferString("{invalid json}"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PostCommand(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ErrorResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response.Message, "Invalid request")
}

func TestPostCommand_MissingCommandField(t *testing.T) {
	handler := NewServerHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBufferString("{}"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PostCommand(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostCommand_ServerNotInitialized(t *testing.T) {
	handler := NewServerHandler(nil)

	payload, _ := json.Marshal(models.Command{Command: "list"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBuffer(payload))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PostCommand(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Minecraft server not initialized", response.Message)
}

func TestPostCommand_ExecuteCommandError(t *testing.T) {
	mockServer := new(testutils.MockMinecraftServer)
	handler := NewServerHandler(mockServer)

	testCommand := "invalid-command"
	expectedError := errors.New("command execution failed")
	mockServer.On("ExecuteCommand", context.Background(), testCommand).Return(nil, expectedError)

	payload, _ := json.Marshal(models.Command{Command: testCommand})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBuffer(payload))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PostCommand(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Contains(t, response.Message, "Failed to execute command")
	assert.Contains(t, response.Message, expectedError.Error())

	mockServer.AssertExpectations(t)
}

func TestPostCommand_WithExitCode(t *testing.T) {
	mockServer := new(testutils.MockMinecraftServer)
	handler := NewServerHandler(mockServer)

	testCommand := "say Hello World"
	exitCode := 0
	expectedResult := server.NewCommandResultWithExitCode("Broadcast message", exitCode)
	mockServer.On("ExecuteCommand", context.Background(), testCommand).Return(expectedResult, nil)

	payload, _ := json.Marshal(models.Command{Command: testCommand})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/command", bytes.NewBuffer(payload))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.PostCommand(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response server.CommandResult
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, expectedResult.Output, response.Output)
	assert.NotNil(t, response.ExitCode)
	assert.Equal(t, exitCode, *response.ExitCode)

	mockServer.AssertExpectations(t)
}
