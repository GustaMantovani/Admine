package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GustaMantovani/Admine/server_handler/internal/api/models"
	mcmodels "github.com/GustaMantovani/Admine/server_handler/internal/mc_server/models"
	"github.com/GustaMantovani/Admine/server_handler/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	testutils.SetupGinTestMode()
	m.Run()
}

// TestPostInstallMod_FileUploadSuccess tests successful .jar upload
func TestPostInstallMod_FileUploadSuccess(t *testing.T) {
	// Setup with full context (includes Config for goroutine)
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	_, cancel := testutils.SetupTestContext(t, mockServer)
	defer cancel()

	// Allow any async PubSub calls
	mockPubSub.On("Publish", mock.Anything, mock.Anything).Return(nil)
	mockServer.On("InstallMod", mock.Anything, "test-mod.jar", mock.Anything).
		Return(mcmodels.NewModInstallResult("test-mod.jar", true, "Installed"), nil).Maybe()

	handler := NewModHandler(mockPubSub)

	// Create multipart form with a .jar file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test-mod.jar")
	assert.NoError(t, err)
	part.Write([]byte("fake jar content"))
	writer.Close()

	// Create request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	// Assert 202 Accepted
	assert.Equal(t, http.StatusAccepted, w.Code)

	var response models.ModInstallResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "accepted", response.Status)
	assert.Contains(t, response.Message, "test-mod.jar")

	// Give goroutine time to complete
	time.Sleep(100 * time.Millisecond)
}

// TestPostInstallMod_FileUploadInvalidExtension tests rejection of non-.jar files
func TestPostInstallMod_FileUploadInvalidExtension(t *testing.T) {
	// Setup
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	_, cancel := testutils.SetupTestContext(t, mockServer)
	defer cancel()

	handler := NewModHandler(mockPubSub)

	// Create multipart form with a .txt file
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

	// Assert 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ModInstallResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, ".jar")
}

// TestPostInstallMod_URLSuccess tests successful URL-based install request
func TestPostInstallMod_URLSuccess(t *testing.T) {
	// Setup
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	_, cancel := testutils.SetupTestContext(t, mockServer)
	defer cancel()

	// Allow any async PubSub/Install calls
	mockPubSub.On("Publish", mock.Anything, mock.Anything).Return(nil)
	mockServer.On("InstallMod", mock.Anything, mock.Anything, mock.Anything).
		Return(mcmodels.NewModInstallResult("cool-mod.jar", true, "Installed"), nil).Maybe()

	handler := NewModHandler(mockPubSub)

	// Create JSON body with URL
	reqBody := models.ModInstallRequest{URL: "https://example.com/mods/cool-mod.jar"}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	// Assert 202 Accepted
	assert.Equal(t, http.StatusAccepted, w.Code)

	var response models.ModInstallResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "accepted", response.Status)
	assert.Contains(t, response.Message, "cool-mod.jar")
}

// TestPostInstallMod_URLInvalidExtension tests rejection of non-.jar URLs
func TestPostInstallMod_URLInvalidExtension(t *testing.T) {
	// Setup
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	_, cancel := testutils.SetupTestContext(t, mockServer)
	defer cancel()

	handler := NewModHandler(mockPubSub)

	reqBody := models.ModInstallRequest{URL: "https://example.com/mods/readme.txt"}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	// Assert 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.ModInstallResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
}

// TestPostInstallMod_ServerNotInitialized tests when MinecraftServer is nil
func TestPostInstallMod_ServerNotInitialized(t *testing.T) {
	// Setup
	mockPubSub := new(testutils.MockPubSubService)
	_, cancel := testutils.SetupTestContext(t, nil)
	defer cancel()

	handler := NewModHandler(mockPubSub)

	reqBody := models.ModInstallRequest{URL: "https://example.com/mods/mod.jar"}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	// Assert 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ModInstallResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "error", response.Status)
	assert.Contains(t, response.Message, "not initialized")
}

// TestPostInstallMod_URLMissingField tests JSON request without URL
func TestPostInstallMod_URLMissingField(t *testing.T) {
	// Setup
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	_, cancel := testutils.SetupTestContext(t, mockServer)
	defer cancel()

	handler := NewModHandler(mockPubSub)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/mods", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router := gin.New()
	router.POST("/api/v1/mods", handler.PostInstallMod)
	router.ServeHTTP(w, req)

	// Assert 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetListMods_Success tests successful listing of mods
func TestGetListMods_Success(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	_, cancel := testutils.SetupTestContext(t, mockServer)
	defer cancel()

	modList := mcmodels.NewModListResult([]string{"mod-a.jar", "mod-b.jar"})
	mockServer.On("ListMods", mock.Anything).Return(modList, nil)

	handler := NewModHandler(mockPubSub)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/mods", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/v1/mods", handler.GetListMods)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response mcmodels.ModListResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 2, response.Total)
	assert.Contains(t, response.Mods, "mod-a.jar")
	assert.Contains(t, response.Mods, "mod-b.jar")
}

// TestGetListMods_ServerNotInitialized tests listing when server is nil
func TestGetListMods_ServerNotInitialized(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	_, cancel := testutils.SetupTestContext(t, nil)
	defer cancel()

	handler := NewModHandler(mockPubSub)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/mods", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.GET("/api/v1/mods", handler.GetListMods)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestDeleteRemoveMod_Success tests successful mod removal
func TestDeleteRemoveMod_Success(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	_, cancel := testutils.SetupTestContext(t, mockServer)
	defer cancel()

	mockServer.On("RemoveMod", mock.Anything, "test-mod.jar").
		Return(mcmodels.NewModInstallResult("test-mod.jar", true, "Mod removed successfully"), nil)

	handler := NewModHandler(mockPubSub)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/mods/test-mod.jar", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.DELETE("/api/v1/mods/:filename", handler.DeleteRemoveMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response mcmodels.ModInstallResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

// TestDeleteRemoveMod_InvalidExtension tests rejection of non-.jar removal
func TestDeleteRemoveMod_InvalidExtension(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	mockServer := new(testutils.MockMinecraftServer)
	_, cancel := testutils.SetupTestContext(t, mockServer)
	defer cancel()

	handler := NewModHandler(mockPubSub)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/mods/readme.txt", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.DELETE("/api/v1/mods/:filename", handler.DeleteRemoveMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestDeleteRemoveMod_ServerNotInitialized tests removal when server is nil
func TestDeleteRemoveMod_ServerNotInitialized(t *testing.T) {
	mockPubSub := new(testutils.MockPubSubService)
	_, cancel := testutils.SetupTestContext(t, nil)
	defer cancel()

	handler := NewModHandler(mockPubSub)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/mods/test-mod.jar", nil)
	w := httptest.NewRecorder()

	router := gin.New()
	router.DELETE("/api/v1/mods/:filename", handler.DeleteRemoveMod)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
