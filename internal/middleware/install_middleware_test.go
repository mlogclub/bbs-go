package middleware

import (
	"bbs-go/internal/pkg/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInstallMiddlewareUsesDedicatedStatusForUninstalledSite(t *testing.T) {
	previousConfig := config.Instance
	config.Instance = &config.Config{Installed: false}
	t.Cleanup(func() { config.Instance = previousConfig })

	app := gin.New()
	app.Use(InstallMiddleware)
	app.GET("/api/topic/topics", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/topic/topics", nil))

	if recorder.Code != http.StatusPreconditionRequired {
		t.Fatalf("status=%d want %d", recorder.Code, http.StatusPreconditionRequired)
	}

	var response struct {
		ErrorCode int `json:"errorCode"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.ErrorCode != -1 {
		t.Fatalf("errorCode=%d want -1", response.ErrorCode)
	}
}

func TestInstallMiddlewareAllowsInstallStatusEndpoint(t *testing.T) {
	previousConfig := config.Instance
	config.Instance = &config.Config{Installed: false}
	t.Cleanup(func() { config.Instance = previousConfig })

	app := gin.New()
	app.Use(InstallMiddleware)
	app.GET("/api/install/status", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/install/status", nil))

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status=%d want %d", recorder.Code, http.StatusNoContent)
	}
}
