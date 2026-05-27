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
	"bbs-go/internal/pkg/idcodec"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestUserListFiltersForbiddenUsers(t *testing.T) {
	db := setupAdminUserTestDB(t)
	now := time.Now().UnixMilli()
	mustCreateUser(t, db, &models.User{
		Model:            models.Model{Id: 1},
		Nickname:         "normal",
		ForbiddenEndTime: 0,
	})
	mustCreateUser(t, db, &models.User{
		Model:            models.Model{Id: 2},
		Nickname:         "forbidden",
		ForbiddenEndTime: now + 3600_000,
	})
	mustCreateUser(t, db, &models.User{
		Model:            models.Model{Id: 3},
		Nickname:         "forever",
		ForbiddenEndTime: -1,
	})
	mustCreateUser(t, db, &models.User{
		Model:            models.Model{Id: 4},
		Nickname:         "expired",
		ForbiddenEndTime: now - 3600_000,
	})

	users := postUserList(t, "forbidden=true")

	gotIDs := make([]int64, 0, len(users))
	for _, user := range users {
		gotIDs = append(gotIDs, int64(user["id"].(float64)))
		if forbidden, ok := user["forbidden"].(bool); !ok || !forbidden {
			t.Fatalf("expected only forbidden users, got %#v", user)
		}
	}
	wantIDs := []int64{3, 2}
	if len(gotIDs) != len(wantIDs) {
		t.Fatalf("expected ids %v, got %v", wantIDs, gotIDs)
	}
	for i := range wantIDs {
		if gotIDs[i] != wantIDs[i] {
			t.Fatalf("expected ids %v, got %v", wantIDs, gotIDs)
		}
	}
}

func TestUserListFiltersNonForbiddenUsers(t *testing.T) {
	db := setupAdminUserTestDB(t)
	now := time.Now().UnixMilli()
	mustCreateUser(t, db, &models.User{
		Model:            models.Model{Id: 1},
		Nickname:         "normal",
		ForbiddenEndTime: 0,
	})
	mustCreateUser(t, db, &models.User{
		Model:            models.Model{Id: 2},
		Nickname:         "forbidden",
		ForbiddenEndTime: now + 3600_000,
	})
	mustCreateUser(t, db, &models.User{
		Model:            models.Model{Id: 3},
		Nickname:         "forever",
		ForbiddenEndTime: -1,
	})
	mustCreateUser(t, db, &models.User{
		Model:            models.Model{Id: 4},
		Nickname:         "expired",
		ForbiddenEndTime: now - 3600_000,
	})

	users := postUserList(t, "forbidden=false")

	gotIDs := make([]int64, 0, len(users))
	for _, user := range users {
		gotIDs = append(gotIDs, int64(user["id"].(float64)))
		if forbidden, ok := user["forbidden"].(bool); !ok || forbidden {
			t.Fatalf("expected only non-forbidden users, got %#v", user)
		}
	}
	wantIDs := []int64{4, 1}
	if len(gotIDs) != len(wantIDs) {
		t.Fatalf("expected ids %v, got %v", wantIDs, gotIDs)
	}
	for i := range wantIDs {
		if gotIDs[i] != wantIDs[i] {
			t.Fatalf("expected ids %v, got %v", wantIDs, gotIDs)
		}
	}
}

func setupAdminUserTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	idcodec.Init(1)

	dsn := fmt.Sprintf("file:admin_user_test_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
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
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("auto migrate users: %v", err)
	}
	return db
}

func mustCreateUser(t *testing.T, db *gorm.DB, user *models.User) {
	t.Helper()

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
}

func postUserList(t *testing.T, body string) []map[string]interface{} {
	t.Helper()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/user/list", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx.Request = req

	UserList(ctx)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			Results []map[string]interface{} `json:"results"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response %q: %v", w.Body.String(), err)
	}
	if !result.Success {
		t.Fatalf("expected success response, got %s", w.Body.String())
	}
	return result.Data.Results
}
