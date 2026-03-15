package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal/api/models"
	"github.com/GustaMantovani/Admine/server_handler/internal/config"
	"github.com/GustaMantovani/Admine/server_handler/internal/server"
	"github.com/GustaMantovani/Admine/server_handler/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func testModCfg() config.MinecraftServerConfig {
	return config.MinecraftServerConfig{ModInstallTimeout: 5 * time.Second}
}

func newTestModHandler(srv server.MinecraftServer, ps *testutils.MockPubSubService) *ModHandler {
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel // context lives for the duration of the test
	return NewModHandler(srv, ps, "test_server", "test_server_channel", testModCfg(), ctx)
}

func TestPostInstallMod_FileUploadSuccess(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	mockPubSub.On("Publish", mock.Anything, mock.Anything).Return(nil)
	mockServer.On("InstallMod", mock.Anything, "test-mod.jar", mock.Anything).
		Return(server.NewModInstallResult("test-mod.jar", true, "Installed"), nil).Maybe()

	handler := newTestModHandler(mockServer, mockPubSub)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test-mod.jar")
	assert.NoError(t, err)
	part.Write([]byte("fake jar content"))
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	var response models.ModInstallResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "accepted", response.Status)
	assert.Contains(t, response.Message, "test-mod.jar")

	time.Sleep(100 * time.Millisecond)
}

func TestPostInstallMod_FileUploadInvalidExtension(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestModHandler(mockServer, mockPubSub)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "not-a-mod.txt")
	assert.NoError(t, err)
	part.Write([]byte("not a jar"))
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ModInstallResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, ".jar")
}

func TestPostInstallMod_URLSuccess(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	mockPubSub.On("Publish", mock.Anything, mock.Anything).Return(nil)
	mockServer.On("InstallMod", mock.Anything, mock.Anything, mock.Anything).
		Return(server.NewModInstallResult("cool-mod.jar", true, "Installed"), nil).Maybe()

	handler := newTestModHandler(mockServer, mockPubSub)

	reqBody := models.ModInstallRequest{URL: "https://example.com/mods/cool-mod.jar"}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)

	var response models.ModInstallResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "accepted", response.Status)
	assert.Contains(t, response.Message, "cool-mod.jar")
}

func TestPostInstallMod_URLInvalidExtension(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestModHandler(mockServer, mockPubSub)

	reqBody := models.ModInstallRequest{URL: "https://example.com/mods/readme.txt"}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ModInstallResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "error", response.Status)
}

func TestPostInstallMod_ServerNotInitialized(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	handler := newTestModHandler(nil, mockPubSub)

	reqBody := models.ModInstallRequest{URL: "https://example.com/mods/mod.jar"}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ModInstallResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, "not initialized")
}

func TestPostInstallMod_URLMissingField(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestModHandler(mockServer, mockPubSub)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetListMods_Success(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)

	modList := server.NewModListResult([]string{"mod-a.jar", "mod-b.jar"})
	mockServer.On("ListMods", mock.Anything).Return(modList, nil)

	handler := newTestModHandler(mockServer, mockPubSub)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/mods", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/v1/mods", handler.GetListMods)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response server.ModListResult
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, 2, response.Total)
	assert.Contains(t, response.Mods, "mod-a.jar")
	assert.Contains(t, response.Mods, "mod-b.jar")
}

func TestGetListMods_ServerNotInitialized(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	handler := newTestModHandler(nil, mockPubSub)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/mods", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/v1/mods", handler.GetListMods)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteRemoveMod_Success(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)

	mockServer.On("RemoveMod", mock.Anything, "test-mod.jar").
		Return(server.NewModInstallResult("test-mod.jar", true, "Mod removed successfully"), nil)

	handler := newTestModHandler(mockServer, mockPubSub)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/mods/test-mod.jar", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.DELETE("/api/v1/mods/:filename", handler.DeleteRemoveMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response server.ModInstallResult
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.True(t, response.Success)
}

func TestDeleteRemoveMod_InvalidExtension(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	handler := newTestModHandler(mockServer, mockPubSub)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/mods/readme.txt", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.DELETE("/api/v1/mods/:filename", handler.DeleteRemoveMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteRemoveMod_ServerNotInitialized(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	handler := newTestModHandler(nil, mockPubSub)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/mods/test-mod.jar", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.DELETE("/api/v1/mods/:filename", handler.DeleteRemoveMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
