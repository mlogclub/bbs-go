package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/msg"
	"bbs-go/internal/repositories"
	"log/slog"
	"regexp"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
)

// mentionRegex matches @username patterns
// Matches @ followed by alphanumeric chars, Chinese chars, hyphens, underscores
var mentionRegex = regexp.MustCompile(`@([\p{L}\p{N}_\-]+)`)

var MentionService = newMentionService()

func newMentionService() *mentionService {
	return &mentionService{}
}

type mentionService struct {
}

// ParseMentions extracts mentioned usernames from content
func (s *mentionService) ParseMentions(content string) []string {
	if content == "" {
		return nil
	}
	matches := mentionRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	var usernames []string
	for _, match := range matches {
		if len(match) >= 2 {
			username := match[1]
			if !seen[username] {
				seen[username] = true
				usernames = append(usernames, username)
			}
		}
	}
	return usernames
}

// SendMentionNotifications creates mention records and sends notifications
// entityType: "topic" or "comment"
func (s *mentionService) SendMentionNotifications(mentionByUserId int64, entityType string, entityId int64, content string) {
	usernames := s.ParseMentions(content)
	if len(usernames) == 0 {
		return
	}

	now := dates.NowTimestamp()

	for _, username := range usernames {
		user := repositories.UserRepository.GetByUsername(sqls.DB(), username)
		if user == nil {
			continue
		}
		// Don't mention yourself
		if user.Id == mentionByUserId {
			continue
		}

		// Create mention record
		mention := &models.Mention{
			UserId:          user.Id,
			MentionByUserId: mentionByUserId,
			EntityType:      entityType,
			EntityId:        entityId,
			Status:          0,
			CreateTime:      now,
			UpdateTime:      now,
		}
		if err := repositories.MentionRepository.Create(sqls.DB(), mention); err != nil {
			slog.Error("create mention record failed",
				slog.String("err", err.Error()),
				slog.Int64("userId", user.Id),
				slog.Int64("mentionByUserId", mentionByUserId),
			)
			continue
		}

		// Send notification message
		s.sendNotification(mentionByUserId, user.Id, entityType, entityId)
	}
}

func (s *mentionService) sendNotification(fromId, toUserId int64, entityType string, entityId int64) {
	var title string
	var content string

	if entityType == constants.EntityTopic {
		title = locales.Get("message.mention_in_topic_msg_title")
	} else {
		title = locales.Get("message.mention_msg_title")
	}

	content = title

	MessageService.SendMsg(fromId, toUserId,
		msg.TypeMention,
		title,
		content,
		"",
		&msg.MentionExtraData{
			EntityType: entityType,
			EntityId:   entityId,
		})

	// Send mention event
	event.Send(event.MentionCreateEvent{
		UserId:     toUserId,
		EntityType: entityType,
		EntityId:   entityId,
	})
}

// RemoveHTMLTags strips HTML tags to extract plain text for mention parsing
func (s *mentionService) RemoveHTMLTags(html string) string {
	// Simple tag removal - removes everything between < and >
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(html, "")
}

// ParseContentForMentions extracts mentions based on content type
func (s *mentionService) ParseContentForMentions(contentType constants.ContentType, content string) string {
	switch contentType {
	case constants.ContentTypeHtml:
		return s.RemoveHTMLTags(content)
	case constants.ContentTypeMarkdown:
		// For markdown, we can search directly since @username is plain text
		return content
	default:
		return content
	}
}

// SendMentionNotificationsForContent handles both html and markdown content
func (s *mentionService) SendMentionNotificationsForContent(mentionByUserId int64, entityType string, entityId int64, contentType constants.ContentType, content string) {
	text := s.ParseContentForMentions(contentType, content)
	s.SendMentionNotifications(mentionByUserId, entityType, entityId, text)
}

// GetUnreadCount returns unread mention count for a user
func (s *mentionService) GetUnreadCount(userId int64) int64 {
	return repositories.MentionRepository.Count(sqls.DB(),
		sqls.NewCnd().Eq("user_id", userId).Eq("status", 0))
}

// MarkRead marks mentions as read for a user
func (s *mentionService) MarkRead(userId int64) {
	sqls.DB().Model(&models.Mention{}).Where("user_id = ? and status = ?", userId, 0).UpdateColumn("status", 1)
}

// GetMentions returns mentions for a user
func (s *mentionService) GetMentions(userId int64, cursor int64) (mentions []models.Mention, nextCursor int64, hasMore bool) {
	limit := 20
	cnd := sqls.NewCnd().Eq("user_id", userId).Desc("id").Limit(limit)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	mentions = repositories.MentionRepository.Find(sqls.DB(), cnd)
	if len(mentions) > 0 {
		nextCursor = mentions[len(mentions)-1].Id
		hasMore = len(mentions) >= limit
	}
	return
}

// CleanInvalidMentions removes mentions that reference non-existent content
// This should be called when topics/comments are deleted
func (s *mentionService) CleanEntityMentions(entityType string, entityId int64) {
	sqls.DB().Delete(&models.Mention{}, "entity_type = ? and entity_id = ?", entityType, entityId)
}

// GetMentionSummary returns a summary of who was mentioned (for display)
func (s *mentionService) GetMentionSummary(entityType string, entityId int64) []int64 {
	var userIds []int64
	mentions := repositories.MentionRepository.Find(sqls.DB(),
		sqls.NewCnd().Eq("entity_type", entityType).Eq("entity_id", entityId))
	for _, m := range mentions {
		userIds = append(userIds, m.UserId)
	}
	return userIds
}

