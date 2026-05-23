package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestLinkCreateDefaultsStatusToOk(t *testing.T) {
	db := setupAdminLinkTestDB(t)

	link := postLinkCreate(t, "title=Example&url=https%3A%2F%2Fexample.com&summary=site&status=1")

	if link.Status != constants.StatusOk {
		t.Fatalf("expected response status %d, got %d", constants.StatusOk, link.Status)
	}

	var saved models.Link
	if err := db.First(&saved, "id = ?", link.Id).Error; err != nil {
		t.Fatalf("load saved link: %v", err)
	}
	if saved.Status != constants.StatusOk {
		t.Fatalf("expected saved status %d, got %d", constants.StatusOk, saved.Status)
	}
}

func TestLinkRemoveDeletesLink(t *testing.T) {
	db := setupAdminLinkTestDB(t)
	mustCreateLink(t, db, &models.Link{
		Model:      models.Model{Id: 1},
		Title:      "Example",
		Url:        "https://example.com",
		Status:     constants.StatusOk,
		CreateTime: time.Now().UnixMilli(),
	})

	postLinkRemove(t, "ids=1")

	var count int64
	if err := db.Model(&models.Link{}).Where("id = ?", 1).Count(&count).Error; err != nil {
		t.Fatalf("count links: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected link to be deleted, count=%d", count)
	}
}

func setupAdminLinkTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:admin_link_test_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	sqls.SetDB(db)
	if err := db.AutoMigrate(&models.Link{}); err != nil {
		t.Fatalf("auto migrate links: %v", err)
	}
	return db
}

func mustCreateLink(t *testing.T, db *gorm.DB, link *models.Link) {
	t.Helper()

	if err := db.Create(link).Error; err != nil {
		t.Fatalf("create link: %v", err)
	}
}

func postLinkCreate(t *testing.T, body string) models.Link {
	t.Helper()

	recorder := runAdminLinkHandler(t, http.MethodPost, "/api/admin/link/create", body, LinkCreate)

	var result struct {
		Success bool        `json:"success"`
		Data    models.Link `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	if !result.Success {
		t.Fatalf("expected success response, got %s", recorder.Body.String())
	}
	return result.Data
}

func postLinkRemove(t *testing.T, body string) {
	t.Helper()

	recorder := runAdminLinkHandler(t, http.MethodPost, "/api/admin/link/delete", body, LinkRemove)

	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	if !result.Success {
		t.Fatalf("expected success response, got %s", recorder.Body.String())
	}
}

func runAdminLinkHandler(
	t *testing.T,
	method string,
	path string,
	body string,
	handler gin.HandlerFunc,
) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx.Request = req

	handler(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	return recorder
}
