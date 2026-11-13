package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bbs-go/internal/pkg/config"
	"github.com/gin-gonic/gin"
)

func TestResourceHTTPErrorStatuses(t *testing.T) {
	previousConfig := config.Instance
	config.Instance = &config.Config{Language: config.DefaultLanguage}
	t.Cleanup(func() { config.Instance = previousConfig })
	setupTopicHandlerCategoryTestDB(t)
	app := gin.New()
	app.GET("/category", Category)
	app.GET("/article/:id", ArticleDetail)
	app.GET("/tag/:id", TagDetail)
	for _, tc := range []struct {
		path   string
		status int
	}{
		{"/category?categoryId=999999", http.StatusNotFound},
		{"/article/invalid", http.StatusBadRequest},
		{"/tag/invalid", http.StatusBadRequest},
	} {
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if rec.Code != tc.status {
			t.Fatalf("%s status=%d want %d", tc.path, rec.Code, tc.status)
		}
		var result struct {
			ErrorCode int `json:"errorCode"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if result.ErrorCode != tc.status {
			t.Fatalf("%s business code=%d want %d", tc.path, result.ErrorCode, tc.status)
		}
	}
}
