package admin

import (
	"bytes"
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

func TestBadgeUpdateSortUpdatesSortNoBySubmittedOrder(t *testing.T) {
	db := setupAdminBadgeTestDB(t)
	mustCreateBadge(t, db, &models.Badge{
		Model:      models.Model{Id: 1},
		Name:       "first",
		Title:      "First",
		SortNo:     0,
		Status:     constants.StatusOk,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	})
	mustCreateBadge(t, db, &models.Badge{
		Model:      models.Model{Id: 2},
		Name:       "second",
		Title:      "Second",
		SortNo:     1,
		Status:     constants.StatusOk,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	})
	mustCreateBadge(t, db, &models.Badge{
		Model:      models.Model{Id: 3},
		Name:       "third",
		Title:      "Third",
		SortNo:     2,
		Status:     constants.StatusOk,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	})

	postBadgeUpdateSort(t, []int64{3, 1, 2})

	got := map[int64]int{}
	var badges []models.Badge
	if err := db.Find(&badges).Error; err != nil {
		t.Fatalf("load badges: %v", err)
	}
	for _, badge := range badges {
		got[badge.Id] = badge.SortNo
	}

	expected := map[int64]int{3: 0, 1: 1, 2: 2}
	for id, sortNo := range expected {
		if got[id] != sortNo {
			t.Fatalf("expected badge %d sortNo %d, got %d", id, sortNo, got[id])
		}
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

func mustCreateBadge(t *testing.T, db *gorm.DB, badge *models.Badge) {
	t.Helper()

	if err := db.Create(badge).Error; err != nil {
		t.Fatalf("create badge: %v", err)
	}
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

func postBadgeUpdateSort(t *testing.T, ids []int64) {
	t.Helper()

	payload, err := json.Marshal(ids)
	if err != nil {
		t.Fatalf("marshal ids: %v", err)
	}

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/badge/update_sort", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	BadgeUpdateSort(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

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
