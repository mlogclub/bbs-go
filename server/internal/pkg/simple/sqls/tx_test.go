package sqls_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"bbs-go/internal/pkg/simple/sqls"
)

// 测试用的模型
type TestUser struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

// 设置测试数据库
func setupTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 创建测试表
	err = db.AutoMigrate(&TestUser{})
	assert.NoError(t, err)

	sqls.SetDB(db)
}

// 测试成功场景：事务成功，回调执行
func TestWithTransaction_Success(t *testing.T) {
	callbackExecuted := false

	setupTestDB(t)

	err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		// 注册回调
		ctx.RegisterCallback(func() {
			callbackExecuted = true
		})

		return nil
	})

	assert.NoError(t, err)
	assert.True(t, callbackExecuted, "回调应该被执行")
}

// 测试失败场景：事务失败，回调不执行
func TestWithTransaction_TransactionFailed(t *testing.T) {
	callbackExecuted := false

	setupTestDB(t)

	err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		// 注册回调
		ctx.RegisterCallback(func() {
			callbackExecuted = true
		})

		// 故意返回错误
		return errors.New("transaction failed")
	})

	assert.Error(t, err)
	assert.False(t, callbackExecuted, "事务失败时回调不应该被执行")
}

// 测试多个回调场景：注册多个回调，全部执行
func TestWithTransaction_MultipleCallbacks(t *testing.T) {
	callback1Executed := false
	callback2Executed := false

	setupTestDB(t)

	err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		// 注册多个回调
		ctx.RegisterCallback(func() {
			callback1Executed = true
		})

		ctx.RegisterCallback(func() {
			callback2Executed = true
		})

		return nil
	})

	assert.NoError(t, err)
	assert.True(t, callback1Executed, "第一个回调应该被执行")
	assert.True(t, callback2Executed, "第二个回调应该被执行")
}

// 测试空回调场景：注册空回调，不执行
func TestWithTransaction_NilCallback(t *testing.T) {
	setupTestDB(t)

	err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		// 注册空回调
		ctx.RegisterCallback(nil)
		return nil
	})

	assert.NoError(t, err)
}

// 测试条件回调场景：根据条件注册回调
func TestWithTransaction_ConditionalCallbacks(t *testing.T) {
	callbackExecuted := false
	condition := true

	setupTestDB(t)

	err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		// 根据条件注册回调
		if condition {
			ctx.RegisterCallback(func() {
				callbackExecuted = true
			})
		}

		return nil
	})

	assert.NoError(t, err)
	assert.True(t, callbackExecuted, "条件满足时回调应该被执行")
}

// 测试嵌套事务场景：在回调中执行数据库操作
func TestWithTransaction_NestedOperations(t *testing.T) {
	setupTestDB(t)
	callbackExecuted := false

	err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		// 创建用户
		user := &TestUser{Name: "Test User"}
		if err := ctx.Tx.Create(user).Error; err != nil {
			return err
		}

		// 注册回调，在回调中执行数据库操作
		ctx.RegisterCallback(func() {
			// 注意：这里的操作不会在事务中执行
			ctx.Tx.Create(&TestUser{Name: "Callback User"})
			callbackExecuted = true
		})

		return nil
	})

	assert.NoError(t, err)
	assert.True(t, callbackExecuted, "回调应该被执行")
}

// 测试并发场景：多个事务同时执行
func TestWithTransaction_Concurrent(t *testing.T) {
	done := make(chan bool)
	const numTransactions = 10

	setupTestDB(t)

	for i := 0; i < numTransactions; i++ {
		go func() {
			// 为每个goroutine创建一个新的数据库实例
			callbackExecuted := false

			err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
				// 注册回调
				ctx.RegisterCallback(func() {
					callbackExecuted = true
				})

				return nil
			})

			assert.NoError(t, err)
			assert.True(t, callbackExecuted, "回调应该被执行")
			done <- true
		}()
	}

	// 等待所有事务完成
	for i := 0; i < numTransactions; i++ {
		<-done
	}
}
