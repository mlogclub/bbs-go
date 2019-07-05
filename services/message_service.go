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

var MessageService = newMessageService()

func newMessageService() *messageService {
	return &messageService{}
}

type messageService struct {
}

func (this *messageService) Get(id int64) *model.Message {
	return repositories.MessageRepository.Get(simple.GetDB(), id)
}

func (this *messageService) Take(where ...interface{}) *model.Message {
	return repositories.MessageRepository.Take(simple.GetDB(), where...)
}

func (this *messageService) QueryCnd(cnd *simple.QueryCnd) (list []model.Message, err error) {
	return repositories.MessageRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *messageService) Query(queries *simple.ParamQueries) (list []model.Message, paging *simple.Paging) {
	return repositories.MessageRepository.Query(simple.GetDB(), queries)
}

func (this *messageService) Create(t *model.Message) error {
	return repositories.MessageRepository.Create(simple.GetDB(), t)
}

func (this *messageService) Update(t *model.Message) error {
	return repositories.MessageRepository.Update(simple.GetDB(), t)
}

func (this *messageService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.MessageRepository.Updates(simple.GetDB(), id, columns)
}

func (this *messageService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.MessageRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *messageService) Delete(id int64) {
	repositories.MessageRepository.Delete(simple.GetDB(), id)
}

func (this *messageService) GetUnReadCount(userId int64) (count int64) {
	simple.GetDB().Where("user_id = ? and status = ?", userId, model.MsgStatusUnread).Model(&model.Message{}).Count(&count)
	return
}

// 读消息
func (this *messageService) Read(id int64) *model.Message {
	msg := this.Get(id)
	if msg != nil && msg.Status == model.MsgStatusUnread {
		_ = this.UpdateColumn(id, "status", model.MsgStatusReaded) // 标记为已读
	}
	return msg
}

// 将所有消息标记为已读
func (this *messageService) MarkReadAll(userId int64) {
	simple.GetDB().Exec("update t_message set status = ? where user_id = ? and status = ?", model.MsgStatusReaded,
		userId, model.MsgStatusUnread)
}

func (this *messageService) Send(userId int64, content, quoteContent string, msgType int, extraData map[string]interface{}) {
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

func (this *messageService) sendEmailNotice(message *model.Message) {
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

func (this *messageService) SendCommentMsg(comment *model.Comment) {
	commentUser := repositories.UserRepository.Get(simple.GetDB(), comment.UserId)
	commentSummary := utils.GetMarkdownSummary(comment.Content)
	// 引用消息
	if comment.QuoteId > 0 {
		quote := repositories.CommentRepository.Get(simple.GetDB(), comment.QuoteId)
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
			article := repositories.ArticleRepository.Get(simple.GetDB(), comment.EntityId)
			if article != nil && article.UserId != comment.UserId {
				userId = article.UserId
				msgContent = commentUser.Nickname + " 回复了你的文章：" + commentSummary
				msgQuoteContent = "《" + article.Title + "》"
			}
		} else if comment.EntityType == model.EntityTypeTopic {
			topic := repositories.TopicRepository.Get(simple.GetDB(), comment.EntityId)
			if topic != nil && topic.UserId != comment.UserId {
				userId = topic.UserId
				msgContent = commentUser.Nickname + " 回复了你的主题：" + commentSummary
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
