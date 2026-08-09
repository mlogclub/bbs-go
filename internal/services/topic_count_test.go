package services

import (
	"os"
	"path/filepath"
	"testing"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

func setupTopicCountTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dir, err := os.MkdirTemp("", "bbs-go-test-search-")
	if err != nil {
		t.Fatalf("mk temp dir: %v", err)
	}
	config.Instance = &config.Config{
		Language: config.DefaultLanguage,
		Search:   config.SearchConfig{IndexPath: filepath.Join(dir, "index")},
	}
	search.Init()

	db := setupTestDB(t)
	if err := db.AutoMigrate(
		&models.Topic{},
		&models.TopicTag{},
		&models.Attachment{},
		&models.Comment{},
	); err != nil {
		t.Fatalf("auto migrate topic count: %v", err)
	}
	return db
}

func mustCreateTopicWithStatus(t *testing.T, userId int64, status int) *models.Topic {
	t.Helper()
	now := dates.NowTimestamp()
	topic := &models.Topic{
		Type:            constants.TopicTypeTopic,
		CategoryId:      1,
		QaStatus:        constants.QaStatusUnsolved,
		UserId:          userId,
		Title:           "test topic",
		ContentType:     constants.ContentTypeText,
		Content:         "content",
		Status:          status,
		LastCommentTime: now,
		CreateTime:      now,
	}
	if err := repositories.TopicRepository.Create(sqls.DB(), topic); err != nil {
		t.Fatalf("create topic: %v", err)
	}
	return topic
}

func mustSetUserCount(t *testing.T, userId int64, column string, count int) {
	t.Helper()
	if err := repositories.UserRepository.UpdateColumn(sqls.DB(), userId, column, count); err != nil {
		t.Fatalf("set %s: %v", column, err)
	}
}

func getUserTopicCount(t *testing.T, userId int64) int {
	t.Helper()
	user := UserService.Get(userId)
	if user == nil {
		t.Fatalf("user not found")
	}
	return user.TopicCount
}

func TestTopicService_Delete_DecrementsTopicCountForPublished(t *testing.T) {
	setupTopicCountTestDB(t)
	user := mustCreateUser(t, dates.NowTimestamp())
	topic := mustCreateTopicWithStatus(t, user.Id, constants.StatusOk)
	mustSetUserCount(t, user.Id, "topic_count", 1)

	if err := TopicService.Delete(topic.Id, user.Id, nil); err != nil {
		t.Fatalf("delete topic: %v", err)
	}
	if got := getUserTopicCount(t, user.Id); got != 0 {
		t.Fatalf("expected topic_count 0 after delete, got %d", got)
	}
}

func TestTopicService_Delete_DoesNotDecrementReviewTopic(t *testing.T) {
	setupTopicCountTestDB(t)
	user := mustCreateUser(t, dates.NowTimestamp())
	topic := mustCreateTopicWithStatus(t, user.Id, constants.StatusReview)
	mustSetUserCount(t, user.Id, "topic_count", 0)

	if err := TopicService.Delete(topic.Id, user.Id, nil); err != nil {
		t.Fatalf("delete review topic: %v", err)
	}
	if got := getUserTopicCount(t, user.Id); got != 0 {
		t.Fatalf("expected topic_count still 0 after deleting review topic, got %d", got)
	}
}

func TestTopicService_Undelete_IncrementsTopicCount(t *testing.T) {
	setupTopicCountTestDB(t)
	user := mustCreateUser(t, dates.NowTimestamp())
	topic := mustCreateTopicWithStatus(t, user.Id, constants.StatusDeleted)
	mustSetUserCount(t, user.Id, "topic_count", 0)

	if err := TopicService.Undelete(topic.Id); err != nil {
		t.Fatalf("undelete topic: %v", err)
	}
	if got := getUserTopicCount(t, user.Id); got != 1 {
		t.Fatalf("expected topic_count 1 after undelete, got %d", got)
	}
}

func TestTopicService_Undelete_ReviewTopicCountsAfterRestore(t *testing.T) {
	setupTopicCountTestDB(t)
	user := mustCreateUser(t, dates.NowTimestamp())
	// 待审核话题发布时未计数，恢复（Undelete）为已发布后应计入
	topic := mustCreateTopicWithStatus(t, user.Id, constants.StatusReview)
	mustSetUserCount(t, user.Id, "topic_count", 0)

	if err := TopicService.Undelete(topic.Id); err != nil {
		t.Fatalf("undelete review topic: %v", err)
	}
	if got := getUserTopicCount(t, user.Id); got != 1 {
		t.Fatalf("expected topic_count 1 after undelete review topic, got %d", got)
	}
}

func TestTopicService_Audit_IncrementsTopicCountOnce(t *testing.T) {
	setupTopicCountTestDB(t)
	user := mustCreateUser(t, dates.NowTimestamp())
	topic := mustCreateTopicWithStatus(t, user.Id, constants.StatusReview)
	mustSetUserCount(t, user.Id, "topic_count", 0)

	if err := TopicService.Audit(topic.Id); err != nil {
		t.Fatalf("audit topic: %v", err)
	}
	if got := getUserTopicCount(t, user.Id); got != 1 {
		t.Fatalf("expected topic_count 1 after audit, got %d", got)
	}

	// 二次审核（已是正常状态）不应再次计数
	if err := TopicService.Audit(topic.Id); err != nil {
		t.Fatalf("re-audit topic: %v", err)
	}
	if got := getUserTopicCount(t, user.Id); got != 1 {
		t.Fatalf("expected topic_count still 1 after re-audit, got %d", got)
	}
}

func TestCommentService_Delete_DecrementsCommentCount(t *testing.T) {
	setupTopicCountTestDB(t)
	user := mustCreateUser(t, dates.NowTimestamp())
	comment := &models.Comment{
		UserId:      user.Id,
		EntityType:  constants.EntityTopic,
		EntityId:    1,
		Content:     "hello",
		ContentType: constants.ContentTypeText,
		Status:      constants.StatusOk,
		CreateTime:  dates.NowTimestamp(),
	}
	if err := repositories.CommentRepository.Create(sqls.DB(), comment); err != nil {
		t.Fatalf("create comment: %v", err)
	}
	mustSetUserCount(t, user.Id, "comment_count", 1)

	if err := CommentService.Delete(comment.Id); err != nil {
		t.Fatalf("delete comment: %v", err)
	}
	user = UserService.Get(user.Id)
	if user == nil || user.CommentCount != 0 {
		t.Fatalf("expected comment_count 0 after delete, got %+v", user)
	}
}
