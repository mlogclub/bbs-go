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

func TestTaskConfigCreateDefaultsStatusToOk(t *testing.T) {
	db := setupAdminTaskConfigTestDB(t)

	task := postTaskConfigCreate(t, "title=Daily&description=Daily&status=1")

	if task.Status != constants.StatusOk {
		t.Fatalf("expected response status %d, got %d", constants.StatusOk, task.Status)
	}

	var saved models.TaskConfig
	if err := db.First(&saved, "id = ?", task.Id).Error; err != nil {
		t.Fatalf("load saved task config: %v", err)
	}
	if saved.Status != constants.StatusOk {
		t.Fatalf("expected saved status %d, got %d", constants.StatusOk, saved.Status)
	}
}

func TestTaskConfigUpdateSortUpdatesSortNoBySubmittedOrder(t *testing.T) {
	db := setupAdminTaskConfigTestDB(t)
	mustCreateTaskConfig(t, db, &models.TaskConfig{
		Model:      models.Model{Id: 1},
		Title:      "First",
		SortNo:     0,
		Status:     constants.StatusOk,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	})
	mustCreateTaskConfig(t, db, &models.TaskConfig{
		Model:      models.Model{Id: 2},
		Title:      "Second",
		SortNo:     1,
		Status:     constants.StatusOk,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	})
	mustCreateTaskConfig(t, db, &models.TaskConfig{
		Model:      models.Model{Id: 3},
		Title:      "Third",
		SortNo:     2,
		Status:     constants.StatusOk,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	})

	postTaskConfigUpdateSort(t, []int64{3, 1, 2})

	got := map[int64]int{}
	var tasks []models.TaskConfig
	if err := db.Find(&tasks).Error; err != nil {
		t.Fatalf("load task configs: %v", err)
	}
	for _, task := range tasks {
		got[task.Id] = task.SortNo
	}

	expected := map[int64]int{3: 0, 1: 1, 2: 2}
	for id, sortNo := range expected {
		if got[id] != sortNo {
			t.Fatalf("expected task %d sortNo %d, got %d", id, sortNo, got[id])
		}
	}
}

func setupAdminTaskConfigTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:admin_task_config_test_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
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
	if err := db.AutoMigrate(&models.TaskConfig{}); err != nil {
		t.Fatalf("auto migrate task configs: %v", err)
	}
	return db
}

func mustCreateTaskConfig(t *testing.T, db *gorm.DB, task *models.TaskConfig) {
	t.Helper()

	if err := db.Create(task).Error; err != nil {
		t.Fatalf("create task config: %v", err)
	}
}

func postTaskConfigCreate(t *testing.T, body string) models.TaskConfig {
	t.Helper()

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/task-config/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx.Request = req

	TaskConfigCreate(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var result struct {
		Success bool              `json:"success"`
		Data    models.TaskConfig `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	if !result.Success {
		t.Fatalf("expected success response, got %s", recorder.Body.String())
	}
	return result.Data
}

func postTaskConfigUpdateSort(t *testing.T, ids []int64) {
	t.Helper()

	payload, err := json.Marshal(ids)
	if err != nil {
		t.Fatalf("marshal ids: %v", err)
	}

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/task-config/update_sort", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	TaskConfigUpdateSort(ctx)

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
