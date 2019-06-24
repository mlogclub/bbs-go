package services

import (
	"bytes"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/utils"
	"github.com/mlogclub/simple"
	"html/template"
)

type MessageService struct {
	MessageRepository *repositories.MessageRepository
	UserRepository    *repositories.UserRepository
	CommentRepository *repositories.CommentRepository
	ArticleRepository *repositories.ArticleRepository
	TopicRepository   *repositories.TopicRepository
}

func NewMessageService() *MessageService {
	return &MessageService{
		MessageRepository: repositories.NewMessageRepository(),
		CommentRepository: repositories.NewCommentRepository(),
		ArticleRepository: repositories.NewArticleRepository(),
		TopicRepository:   repositories.NewTopicRepository(),
	}
}

func (this *MessageService) Get(id int64) *model.Message {
	return this.MessageRepository.Get(simple.GetDB(), id)
}

func (this *MessageService) Take(where ...interface{}) *model.Message {
	return this.MessageRepository.Take(simple.GetDB(), where...)
}

func (this *MessageService) QueryCnd(cnd *simple.QueryCnd) (list []model.Message, err error) {
	return this.MessageRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *MessageService) Query(queries *simple.ParamQueries) (list []model.Message, paging *simple.Paging) {
	return this.MessageRepository.Query(simple.GetDB(), queries)
}

func (this *MessageService) Create(t *model.Message) error {
	return this.MessageRepository.Create(simple.GetDB(), t)
}

func (this *MessageService) Update(t *model.Message) error {
	return this.MessageRepository.Update(simple.GetDB(), t)
}

func (this *MessageService) Updates(id int64, columns map[string]interface{}) error {
	return this.MessageRepository.Updates(simple.GetDB(), id, columns)
}

func (this *MessageService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.MessageRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *MessageService) Delete(id int64) {
	this.MessageRepository.Delete(simple.GetDB(), id)
}

func (this *MessageService) GetUnReadCount(userId int64) (count int64) {
	simple.GetDB().Where("user_id = ? and status = ?", userId, model.MsgStatusUnread).Model(&model.Message{}).Count(&count)
	return
}

// 读消息
func (this *MessageService) Read(id int64) *model.Message {
	msg := this.Get(id)
	if msg != nil && msg.Status == model.MsgStatusUnread {
		_ = this.UpdateColumn(id, "status", model.MsgStatusReaded) // 标记为已读
	}
	return msg
}

// 将所有消息标记为已读
func (this *MessageService) MarkReadAll(userId int64) {
	simple.GetDB().Exec("update t_message set status = ? where user_id = ? and status = ?", model.MsgStatusReaded,
		userId, model.MsgStatusUnread)
}

func (this *MessageService) Send(userId int64, content, quoteContent string, msgType int, extraData map[string]interface{}) {
	extraDataStr, _ := simple.FormatJson(extraData)
	message := &model.Message{
		UserId:       userId,
		Content:      content,
		QuoteContent: quoteContent,
		Type:         msgType,
		ExtraData:    extraDataStr,
		Status:       model.MsgStatusUnread,
		CreateTime:   simple.NowTimestamp(),
	}
	err := this.Create(message)
	if err == nil {
		go func() {
			this.sendEmailNotice(message)
		}()
	}
}

func (this *MessageService) sendEmailNotice(message *model.Message) {
	user := cache.UserCache.Get(message.UserId)
	if len(user.Email) == 0 {
		return
	}

	tpl, err := template.New("msg").Parse(`
<div style="font-size: 14px;">
	<div style="margin-bottom: 10px;">{{.Content}}</div>
	{{if .QuoteContent}}
		<blockquote style="font-size:12px; padding: 10px 15px; margin: 0 0 20px; border: 1px dotted #eeeeee; border-left: 3px solid #eeeeee; background-color: #fbfbfb;">
		{{.QuoteContent}}
		</blockquote>
	{{end}}
	点击查看详情：<a href="{{.Url}}" target="_blank">{{.Url}}</a>
</div>
`)
	if err != nil {
		return
	}

	var b bytes.Buffer
	err = tpl.Execute(&b, map[string]interface{}{
		"Url":          utils.BuildUserUrl(message.UserId) + "/messages",
		"Content":      message.Content,
		"QuoteContent": message.QuoteContent,
	})
	if err != nil {
		return
	}

	_ = utils.SendEmail(user.Email, "M-LOG：消息提醒", b.String())
}

func (this *MessageService) SendCommentMsg(comment *model.Comment) {
	commentUser := this.UserRepository.Get(simple.GetDB(), comment.UserId)
	commentSummary := utils.GetMarkdownSummary(comment.Content)
	// 引用消息
	if comment.QuoteId > 0 {
		quote := this.CommentRepository.Get(simple.GetDB(), comment.QuoteId)
		if quote != nil && quote.UserId != comment.UserId {
			msgContent := commentUser.Nickname + " 回复了你的评论：" + commentSummary
			quoteContent := utils.GetMarkdownSummary(quote.Content)
			this.Send(quote.UserId, msgContent, quoteContent, model.MsgTypeComment, map[string]interface{}{
				"entityType": comment.EntityType,
				"entityId":   comment.EntityId,
				"commentId":  comment.Id,
				"quoteId":    comment.QuoteId,
			})
		}
	}

	// 文章评论消息
	{
		var userId int64 = 0
		var msgContent = ""
		var msgQuoteContent = ""
		if comment.EntityType == model.EntityTypeArticle {
			article := this.ArticleRepository.Get(simple.GetDB(), comment.EntityId)
			if article != nil && article.UserId != comment.UserId {
				userId = article.UserId
				msgContent = commentUser.Nickname + " 回复了你的文章：" + commentSummary
				msgQuoteContent = "《" + article.Title + "》"
			}
		} else if comment.EntityType == model.EntityTypeTopic {
			topic := this.TopicRepository.Get(simple.GetDB(), comment.EntityId)
			if topic != nil && topic.UserId != comment.UserId {
				userId = topic.UserId
				msgContent = commentUser.Nickname + " 回复了你的讨论：" + commentSummary
				msgQuoteContent = "《" + topic.Title + "》"
			}
		}
		if userId > 0 {
			this.Send(userId, msgContent, msgQuoteContent, model.MsgTypeComment, map[string]interface{}{
				"entityType": comment.EntityType,
				"entityId":   comment.EntityId,
				"commentId":  comment.Id,
			})
		}
	}
}
