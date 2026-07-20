package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/msg"
	"bbs-go/internal/repositories"
	"log/slog"
	"strings"
	"time"

	"bbs-go/internal/pkg/params"

	cachelib "github.com/goburrow/cache"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/tidwall/gjson"
)

var MessageService = newMessageService()
var emailNoticeLimitCache = cachelib.New(
	cachelib.WithMaximumSize(10000),
	cachelib.WithExpireAfterAccess(30*time.Minute),
)

func newMessageService() *messageService {
	return &messageService{}
}

type messageService struct {
}

func (s *messageService) Get(id int64) *models.Message {
	return repositories.MessageRepository.Get(sqls.DB(), id)
}

func (s *messageService) Take(where ...interface{}) *models.Message {
	return repositories.MessageRepository.Take(sqls.DB(), where...)
}

func (s *messageService) Find(cnd *sqls.Cnd) []models.Message {
	return repositories.MessageRepository.Find(sqls.DB(), cnd)
}

func (s *messageService) FindOne(cnd *sqls.Cnd) *models.Message {
	return repositories.MessageRepository.FindOne(sqls.DB(), cnd)
}

func (s *messageService) FindPageByParams(params *params.QueryParams) (list []models.Message, paging *sqls.Paging) {
	return repositories.MessageRepository.FindPageByParams(sqls.DB(), params)
}

func (s *messageService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Message, paging *sqls.Paging) {
	return repositories.MessageRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *messageService) Create(t *models.Message) error {
	return repositories.MessageRepository.Create(sqls.DB(), t)
}

func (s *messageService) Update(t *models.Message) error {
	return repositories.MessageRepository.Update(sqls.DB(), t)
}

func (s *messageService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.MessageRepository.Updates(sqls.DB(), id, columns)
}

func (s *messageService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.MessageRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *messageService) Delete(id int64) {
	repositories.MessageRepository.Delete(sqls.DB(), id)
}

// GetUnReadCount 获取未读消息数量
func (s *messageService) GetUnReadCount(userId int64) (count int64) {
	sqls.DB().Where("user_id = ? and status = ?", userId, msg.StatusUnread).Model(&models.Message{}).Count(&count)
	return
}

// MarkRead 将所有消息标记为已读
func (s *messageService) MarkRead(userId int64) {
	sqls.DB().Exec("update t_message set status = ? where user_id = ? and status = ?", msg.StatusHaveRead,
		userId, msg.StatusUnread)
}

// SendMsg 发送消息（站内信和/或邮件由通知配置分别控制）
func (s *messageService) SendMsg(from, to int64, msgType msg.Type,
	title, content, quoteContent string, extraData interface{}) {
	siteOn := SysConfigService.IsSiteNoticeEnabled(msgType)
	emailOn := SysConfigService.IsEmailNoticeEnabled(msgType)
	if !siteOn && !emailOn {
		return
	}
	t := &models.Message{
		FromId:       from,
		UserId:       to,
		Title:        title,
		Content:      content,
		QuoteContent: quoteContent,
		Type:         int(msgType),
		ExtraData:    jsons.ToJsonStr(extraData),
		Status:       msg.StatusUnread,
		CreateTime:   dates.NowTimestamp(),
	}
	if siteOn {
		if err := s.Create(t); err != nil {
			slog.Error(err.Error(), slog.Any("err", err))
			return
		}
	}
	if emailOn {
		s.SendEmailNotice(t)
	}
}

// SendEmailNotice 发送邮件通知
func (s *messageService) SendEmailNotice(t *models.Message) {
	msgType := msg.Type(t.Type)
	if !SysConfigService.IsEmailNoticeEnabled(msgType) {
		return
	}
	user := cache.UserCache.Get(t.UserId)
	if user == nil || strs.IsBlank(user.Email.String) {
		return
	}
	emailKey := strings.ToLower(user.Email.String)
	intervalSeconds := cache.SysConfigCache.GetInt(constants.SysConfigEmailNoticeIntervalSeconds)
	if intervalSeconds > 0 {
		if lastSend, found := emailNoticeLimitCache.GetIfPresent(emailKey); found {
			if now := time.Now().Unix(); now-lastSend.(int64) < int64(intervalSeconds) {
				return
			}
		}
	}
	var (
		siteTitle        = cache.SysConfigCache.GetStr(constants.SysConfigSiteTitle)
		noticeTitle      = s.buildEmailNoticeFallbackTitle(t)
		emailSubject     string
		emailContent     string
		emailDetailURL   string
		emailActionTitle = locales.Get("email.view_details")
	)

	if title := strings.TrimSpace(t.Title); title != "" {
		noticeTitle = title
	}
	emailSubject = s.buildEmailNoticeSubject(siteTitle, noticeTitle)
	emailContent = s.buildEmailNoticeContent(t.Content, noticeTitle)
	emailDetailURL = s.buildEmailNoticeDetailURL(t)

	var from *models.User
	if t.FromId > 0 {
		from = cache.UserCache.Get(t.FromId)
	}
	err := EmailService.SendTemplateEmail(from, user.Email.String, emailSubject, noticeTitle, emailContent,
		t.QuoteContent, &dto.ActionLink{
			Title: emailActionTitle,
			Url:   emailDetailURL,
		}, constants.EmailLogBizTypeMessageNotice)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return
	}
	if intervalSeconds > 0 {
		emailNoticeLimitCache.Put(emailKey, time.Now().Unix())
	}
}

func (s *messageService) buildEmailNoticeFallbackTitle(t *models.Message) string {
	msgType := msg.Type(t.Type)
	switch msgType {
	case msg.TypeTopicComment:
		return locales.Get("email.topic_comment")
	case msg.TypeCommentReply:
		return locales.Get("email.comment_reply")
	case msg.TypeTopicLike:
		return locales.Get("email.topic_like")
	case msg.TypeTopicFavorite:
		return locales.Get("email.topic_favorite")
	case msg.TypeTopicRecommend:
		return locales.Get("email.topic_recommend")
	case msg.TypeTopicDelete:
		return locales.Get("email.topic_delete")
	case msg.TypeArticleComment:
		return locales.Get("email.article_comment")
	case msg.TypeUserLevelUp:
		return locales.Get("email.user_level_up")
	case msg.TypeUserBadgeGrant:
		return locales.Get("email.user_badge_grant")
	case msg.TypeQaAnswerAccepted:
		bountyScore := gjson.Get(t.ExtraData, "bountyScore").Int()
		if bountyScore > 0 {
			return locales.Getf("email.qa_answer_accepted_with_bounty", int(bountyScore))
		}
		return locales.Get("email.qa_answer_accepted")
	case msg.TypeMention:
		return locales.Get("email.mention")
	}
	return locales.Get("email.new_message")
}

func (s *messageService) buildEmailNoticeSubject(siteTitle, noticeTitle string) string {
	siteTitle = strings.TrimSpace(siteTitle)
	noticeTitle = strings.TrimSpace(noticeTitle)
	if siteTitle == "" {
		return noticeTitle
	}
	if noticeTitle == "" {
		return siteTitle
	}
	return siteTitle + " - " + noticeTitle
}

func (s *messageService) buildEmailNoticeContent(content, noticeTitle string) string {
	content = strings.TrimSpace(content)
	if content != "" {
		return content
	}
	return strings.TrimSpace(noticeTitle)
}

func (s *messageService) buildEmailNoticeDetailURL(t *models.Message) string {
	msgType := msg.Type(t.Type)
	switch msgType {
	case msg.TypeTopicComment, msg.TypeArticleComment:
		entityType := gjson.Get(t.ExtraData, "entityType")
		entityId := gjson.Get(t.ExtraData, "entityId")
		if entityType.String() == constants.EntityArticle {
			return bbsurls.ArticleUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTopic {
			return bbsurls.TopicUrl(entityId.Int())
		}
	case msg.TypeCommentReply:
		entityType := gjson.Get(t.ExtraData, "rootEntityType")
		entityId := gjson.Get(t.ExtraData, "rootEntityId")
		if entityType.String() == constants.EntityArticle {
			return bbsurls.ArticleUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTopic {
			return bbsurls.TopicUrl(entityId.Int())
		}
	case msg.TypeTopicLike, msg.TypeTopicFavorite, msg.TypeTopicRecommend, msg.TypeQaAnswerAccepted:
		topicId := gjson.Get(t.ExtraData, "topicId")
		if topicId.Exists() && topicId.Int() > 0 {
			return bbsurls.TopicUrl(topicId.Int())
		}
	case msg.TypeUserLevelUp:
		return bbsurls.AbsUrl("/tasks")
	case msg.TypeUserBadgeGrant:
		return bbsurls.UserUrl(t.UserId) + "/badges"
	case msg.TypeMention:
		entityType := gjson.Get(t.ExtraData, "entityType")
		entityId := gjson.Get(t.ExtraData, "entityId")
		if entityType.String() == constants.EntityTopic {
			return bbsurls.TopicUrl(entityId.Int())
		}
	}
	return bbsurls.AbsUrl("/user/messages")
}
