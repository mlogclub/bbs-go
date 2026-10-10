package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func setupUserTokenTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:user_token_test_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
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

	if err := db.AutoMigrate(
		&models.User{},
		&models.UserToken{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	return db
}

func TestUserTokenService_SignoutAndDisableCacheInvalidation(t *testing.T) {
	db := setupUserTokenTestDB(t)

	// 1. 创建测试用户
	user := &models.User{
		Model:  models.Model{Id: 1001},
		Status: constants.StatusOk,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	// 2. 测试 Signout 退出登录时必须失效缓存
	t.Run("SignoutInvalidatesCache", func(t *testing.T) {
		token, err := UserTokenService.Generate(user.Id)
		if err != nil {
			t.Fatalf("generate token: %v", err)
		}

		// 预热缓存：先读取一次，确保 token 进入 UserTokenCache
		cachedToken := cache.UserTokenCache.Get(token)
		if cachedToken == nil || cachedToken.Status != constants.StatusOk {
			t.Fatalf("expected active token in cache, got: %+v", cachedToken)
		}

		// 构造 Gin 上下文并设置 Cookie
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest(http.MethodPost, "/api/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  constants.CookieTokenKey,
			Value: token,
		})
		ctx.Request = req

		// 执行退出登录
		if err := UserTokenService.Signout(ctx); err != nil {
			t.Fatalf("signout failed: %v", err)
		}

		// 等待 LoadingCache 异步事件循环完成 Invalidate 处理
		time.Sleep(20 * time.Millisecond)

		// 验证缓存已被失效：再次通过 UserTokenCache.Get 获取，将触发重载并拿到 StatusDeleted
		afterSignoutToken := cache.UserTokenCache.Get(token)
		if afterSignoutToken == nil {
			t.Fatalf("expected token record reloaded, got nil")
		}
		if afterSignoutToken.Status != constants.StatusDeleted {
			t.Fatalf("expected token status to be deleted (1) after signout, got: %d", afterSignoutToken.Status)
		}

		// 验证 GetCurrent 判定为未登录 (返回 nil)
		currentUser := UserTokenService.GetCurrent(ctx)
		if currentUser != nil {
			t.Fatalf("expected current user to be nil after signout, got: %+v", currentUser)
		}
	})

	// 3. 测试 Disable 禁用凭证时必须失效缓存
	t.Run("DisableInvalidatesCache", func(t *testing.T) {
		token, err := UserTokenService.Generate(user.Id)
		if err != nil {
			t.Fatalf("generate token: %v", err)
		}

		// 预热缓存
		cachedToken := cache.UserTokenCache.Get(token)
		if cachedToken == nil || cachedToken.Status != constants.StatusOk {
			t.Fatalf("expected active token in cache, got: %+v", cachedToken)
		}

		// 执行管理员禁用
		if err := UserTokenService.Disable(token); err != nil {
			t.Fatalf("disable token failed: %v", err)
		}

		// 等待 LoadingCache 异步事件循环完成 Invalidate 处理
		time.Sleep(20 * time.Millisecond)

		// 验证缓存已失效并更新为 StatusDeleted
		reloadedToken := cache.UserTokenCache.Get(token)
		if reloadedToken == nil {
			t.Fatalf("expected token record reloaded, got nil")
		}
		if reloadedToken.Status != constants.StatusDeleted {
			t.Fatalf("expected token status to be deleted (1) after disable, got: %d", reloadedToken.Status)
		}
	})
}
