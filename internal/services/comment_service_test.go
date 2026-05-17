package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/permissions"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/repositories"
	"testing"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
)

func setupCommentServiceTestDB(t *testing.T) {
	t.Helper()
	config.Instance = &config.Config{Language: config.DefaultLanguage}
	db := setupTestDB(t)
	if err := db.AutoMigrate(&models.Comment{}, &models.Role{}, &models.UserRole{}, &models.Permission{}, &models.RolePermission{}); err != nil {
		t.Fatalf("auto migrate comment: %v", err)
	}
	PermissionService.ClearCache()
}

func mustCreateComment(t *testing.T, comment *models.Comment) *models.Comment {
	t.Helper()
	if comment.Status == 0 {
		comment.Status = constants.StatusOk
	}
	if err := repositories.CommentRepository.Create(sqls.DB(), comment); err != nil {
		t.Fatalf("create comment: %v", err)
	}
	return comment
}

func TestCommentService_DeleteByUserRejectsNonAuthorWithoutPermission(t *testing.T) {
	setupCommentServiceTestDB(t)
	comment := mustCreateComment(t, &models.Comment{
		UserId:      10,
		EntityType:  constants.EntityTopic,
		EntityId:    20,
		Content:     "hello",
		ContentType: constants.ContentTypeText,
	})

	regularUser := &models.User{Roles: ""}
	if err := CommentService.DeleteByUser(regularUser, comment.Id); err == nil {
		t.Fatalf("expected permission error for regular user")
	}

	got := CommentService.Get(comment.Id)
	if got == nil {
		t.Fatalf("expected comment to still exist")
	}
	if got.Status != constants.StatusOk {
		t.Fatalf("expected comment status ok, got %d", got.Status)
	}
}

func TestCommentService_DeleteByUserAllowsAuthor(t *testing.T) {
	setupCommentServiceTestDB(t)
	comment := mustCreateComment(t, &models.Comment{
		UserId:      10,
		EntityType:  constants.EntityTopic,
		EntityId:    20,
		Content:     "hello",
		ContentType: constants.ContentTypeText,
	})

	author := &models.User{Model: models.Model{Id: 10}}
	if err := CommentService.DeleteByUser(author, comment.Id); err != nil {
		t.Fatalf("delete by author: %v", err)
	}

	got := CommentService.Get(comment.Id)
	if got == nil {
		t.Fatalf("expected comment to still exist")
	}
	if got.Status != constants.StatusDeleted {
		t.Fatalf("expected comment status deleted, got %d", got.Status)
	}
}

func TestCommentService_DeleteByUserAllowsCommentDeletePermission(t *testing.T) {
	setupCommentServiceTestDB(t)
	comment := mustCreateComment(t, &models.Comment{
		UserId:      10,
		EntityType:  constants.EntityTopic,
		EntityId:    20,
		Content:     "hello",
		ContentType: constants.ContentTypeText,
	})
	now := dates.NowTimestamp()
	moderator := mustCreateUser(t, now)
	role := mustCreateRole(t, "comment-moderator", constants.StatusOk)
	permission := mustCreatePermission(t, permissions.PermissionCommentDelete.Code, constants.StatusOk)
	mustAssignRole(t, moderator, role)
	mustGrantPermission(t, role, permission)

	if err := CommentService.DeleteByUser(moderator, comment.Id); err != nil {
		t.Fatalf("delete with comment permission: %v", err)
	}

	got := CommentService.Get(comment.Id)
	if got == nil {
		t.Fatalf("expected comment to still exist")
	}
	if got.Status != constants.StatusDeleted {
		t.Fatalf("expected comment status deleted, got %d", got.Status)
	}
}

func TestCommentService_DeleteByUserAllowsOwner(t *testing.T) {
	setupCommentServiceTestDB(t)
	comment := mustCreateComment(t, &models.Comment{
		UserId:      10,
		EntityType:  constants.EntityComment,
		EntityId:    20,
		Content:     "reply",
		ContentType: constants.ContentTypeText,
	})

	ownerUser := &models.User{Roles: constants.RoleOwner}
	if err := CommentService.DeleteByUser(ownerUser, comment.Id); err != nil {
		t.Fatalf("delete by owner: %v", err)
	}

	got := CommentService.Get(comment.Id)
	if got == nil {
		t.Fatalf("expected comment to still exist")
	}
	if got.Status != constants.StatusDeleted {
		t.Fatalf("expected comment status deleted, got %d", got.Status)
	}
}
