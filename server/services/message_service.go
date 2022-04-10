package services

import (
	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/bbsurls"
	"bbs-go/pkg/email"
	"bbs-go/pkg/msg"
	"bbs-go/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"github.com/sirupsen/logrus"
)

var MessageService = newMessageService()

func newMessageService() *messageService {
	return &messageService{}
}

type messageService struct {
}

func (s *messageService) Get(id int64) *model.Message {
	return repositories.MessageRepository.Get(sqls.DB(), id)
}

func (s *messageService) Take(where ...interface{}) *model.Message {
	return repositories.MessageRepository.Take(sqls.DB(), where...)
}

func (s *messageService) Find(cnd *sqls.Cnd) []model.Message {
	return repositories.MessageRepository.Find(sqls.DB(), cnd)
}

func (s *messageService) FindOne(cnd *sqls.Cnd) *model.Message {
	return repositories.MessageRepository.FindOne(sqls.DB(), cnd)
}

func (s *messageService) FindPageByParams(params *params.QueryParams) (list []model.Message, paging *sqls.Paging) {
	return repositories.MessageRepository.FindPageByParams(sqls.DB(), params)
}

func (s *messageService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Message, paging *sqls.Paging) {
	return repositories.MessageRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *messageService) Create(t *model.Message) error {
	return repositories.MessageRepository.Create(sqls.DB(), t)
}

func (s *messageService) Update(t *model.Message) error {
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
	sqls.DB().Where("user_id = ? and status = ?", userId, msg.StatusUnread).Model(&model.Message{}).Count(&count)
	return
}

// MarkRead 将所有消息标记为已读
func (s *messageService) MarkRead(userId int64) {
	sqls.DB().Exec("update t_message set status = ? where user_id = ? and status = ?", msg.StatusHaveRead,
		userId, msg.StatusUnread)
}

// SendMsg 发送消息
func (s *messageService) SendMsg(from, to int64, msgType msg.Type,
	title, content, quoteContent string, extraData interface{}) {

	t := &model.Message{
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
	if err := s.Create(t); err != nil {
		logrus.Error(err)
	} else {
		s.SendEmailNotice(t)
	}
}

// SendEmailNotice 发送邮件通知
func (s *messageService) SendEmailNotice(t *model.Message) {
	msgType := msg.Type(t.Type)

	// 话题被删除不发送邮件提醒
	if msgType == msg.TypeTopicDelete {
		return
	}
	user := cache.UserCache.Get(t.UserId)
	if user == nil || len(user.Email.String) == 0 {
		return
	}
	var (
		siteTitle  = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
		emailTitle = siteTitle + " - 新消息提醒"
	)

	if msgType == msg.TypeTopicComment {
		emailTitle = siteTitle + " - 收到话题评论"
	} else if msgType == msg.TypeCommentReply {
		emailTitle = siteTitle + " - 收到他人回复"
	} else if msgType == msg.TypeTopicLike {
		emailTitle = siteTitle + " - 收到点赞"
	} else if msgType == msg.TypeTopicFavorite {
		emailTitle = siteTitle + " - 话题被收藏"
	} else if msgType == msg.TypeTopicRecommend {
		emailTitle = siteTitle + " - 话题被设为推荐"
	} else if msgType == msg.TypeTopicDelete {
		emailTitle = siteTitle + " - 话题被删除"
	} else if msgType == msg.TypeArticleComment {
		emailTitle = siteTitle + " - 收到文章评论"
	}

	var from *model.User
	if t.FromId > 0 {
		from = cache.UserCache.Get(t.FromId)
	}
	err := email.SendTemplateEmail(from, user.Email.String, emailTitle, emailTitle, t.Content,
		t.QuoteContent, &model.ActionLink{
			Title: "点击查看详情",
			Url:   bbsurls.AbsUrl("/user/messages"),
		})
	if err != nil {
		logrus.Error(err)
	}
}
