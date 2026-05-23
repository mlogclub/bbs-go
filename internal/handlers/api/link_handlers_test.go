package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestLinkTopLinksOrdersBySortNo(t *testing.T) {
	db := setupAPILinkTestDB(t)
	mustCreateAPILink(t, db, &models.Link{
		Model:      models.Model{Id: 1},
		Title:      "Second",
		Url:        "https://second.example.com",
		SortNo:     2,
		Status:     constants.StatusOk,
		CreateTime: time.Now().UnixMilli(),
	})
	mustCreateAPILink(t, db, &models.Link{
		Model:      models.Model{Id: 2},
		Title:      "First",
		Url:        "https://first.example.com",
		SortNo:     1,
		Status:     constants.StatusOk,
		CreateTime: time.Now().UnixMilli(),
	})
	mustCreateAPILink(t, db, &models.Link{
		Model:      models.Model{Id: 3},
		Title:      "Deleted",
		Url:        "https://deleted.example.com",
		SortNo:     0,
		Status:     constants.StatusDeleted,
		CreateTime: time.Now().UnixMilli(),
	})

	links := getTopLinks(t)
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %#v", links)
	}
	if links[0].Title != "First" || links[1].Title != "Second" {
		t.Fatalf("expected links ordered by sortNo, got %#v", links)
	}
}

func setupAPILinkTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:api_link_test_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
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

func mustCreateAPILink(t *testing.T, db *gorm.DB, link *models.Link) {
	t.Helper()

	if err := db.Create(link).Error; err != nil {
		t.Fatalf("create link: %v", err)
	}
}

func getTopLinks(t *testing.T) []struct {
	Title string `json:"title"`
} {
	t.Helper()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/api/link/top_links", nil)

	LinkTopLinks(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var result struct {
		Success bool `json:"success"`
		Data    []struct {
			Title string `json:"title"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	if !result.Success {
		t.Fatalf("expected success response, got %s", recorder.Body.String())
	}
	return result.Data
}
