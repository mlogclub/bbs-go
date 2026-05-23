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

func TestBadgeCreateDefaultsStatusToOk(t *testing.T) {
	db := setupAdminBadgeTestDB(t)

	badge := postBadgeCreate(t, "name=starter&title=Starter&status=1")

	if badge.Status != constants.StatusOk {
		t.Fatalf("expected response status %d, got %d", constants.StatusOk, badge.Status)
	}

	var saved models.Badge
	if err := db.First(&saved, "id = ?", badge.Id).Error; err != nil {
		t.Fatalf("load saved badge: %v", err)
	}
	if saved.Status != constants.StatusOk {
		t.Fatalf("expected saved status %d, got %d", constants.StatusOk, saved.Status)
	}
}

func setupAdminBadgeTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:admin_badge_test_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
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
	if err := db.AutoMigrate(&models.Badge{}); err != nil {
		t.Fatalf("auto migrate badges: %v", err)
	}
	return db
}

func postBadgeCreate(t *testing.T, body string) models.Badge {
	t.Helper()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/badge/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx.Request = req

	BadgeCreate(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var result struct {
		Success bool         `json:"success"`
		Data    models.Badge `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	if !result.Success {
		t.Fatalf("expected success response, got %s", recorder.Body.String())
	}
	return result.Data
}
