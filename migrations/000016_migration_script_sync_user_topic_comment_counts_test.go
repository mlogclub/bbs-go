package migrations

import (
	"fmt"
	"testing"
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"

	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func setupSyncCountTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:sync_count_test_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
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

	if err := db.AutoMigrate(&models.User{}, &models.Topic{}, &models.Comment{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func TestMigrateSyncUserTopicCommentCounts(t *testing.T) {
	setupSyncCountTestDB(t)
	now := dates.NowTimestamp()

	user := &models.User{Nickname: "u", Status: constants.StatusOk, CreateTime: now, UpdateTime: now}
	if err := repositories.UserRepository.Create(sqls.DB(), user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	// 1 个已发布 + 1 个已删除 + 1 个待审核
	mustCreateSyncTopic(t, user.Id, constants.StatusOk)
	mustCreateSyncTopic(t, user.Id, constants.StatusDeleted)
	mustCreateSyncTopic(t, user.Id, constants.StatusReview)
	// 1 个正常 + 1 个已删除
	mustCreateSyncComment(t, user.Id, constants.StatusOk)
	mustCreateSyncComment(t, user.Id, constants.StatusDeleted)

	// 人为写入错误计数
	if err := repositories.UserRepository.UpdateColumn(sqls.DB(), user.Id, "topic_count", 3); err != nil {
		t.Fatalf("set topic_count: %v", err)
	}
	if err := repositories.UserRepository.UpdateColumn(sqls.DB(), user.Id, "comment_count", 2); err != nil {
		t.Fatalf("set comment_count: %v", err)
	}

	if err := migrate_sync_user_topic_comment_counts(); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	got := repositories.UserRepository.Get(sqls.DB(), user.Id)
	if got == nil {
		t.Fatal("user not found")
	}
	if got.TopicCount != 1 {
		t.Fatalf("expected topic_count 1, got %d", got.TopicCount)
	}
	if got.CommentCount != 1 {
		t.Fatalf("expected comment_count 1, got %d", got.CommentCount)
	}
}

func mustCreateSyncTopic(t *testing.T, userId int64, status int) {
	t.Helper()
	now := dates.NowTimestamp()
	topic := &models.Topic{
		Type:            constants.TopicTypeTopic,
		CategoryId:      1,
		QaStatus:        constants.QaStatusUnsolved,
		UserId:          userId,
		Title:           "t",
		ContentType:     constants.ContentTypeText,
		Content:         "c",
		Status:          status,
		LastCommentTime: now,
		CreateTime:      now,
	}
	if err := repositories.TopicRepository.Create(sqls.DB(), topic); err != nil {
		t.Fatalf("create topic: %v", err)
	}
}

func mustCreateSyncComment(t *testing.T, userId int64, status int) {
	t.Helper()
	comment := &models.Comment{
		UserId:      userId,
		EntityType:  constants.EntityTopic,
		EntityId:    1,
		Content:     "c",
		ContentType: constants.ContentTypeText,
		Status:      status,
		CreateTime:  dates.NowTimestamp(),
	}
	if err := repositories.CommentRepository.Create(sqls.DB(), comment); err != nil {
		t.Fatalf("create comment: %v", err)
	}
}
