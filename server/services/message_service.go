package services

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/common/email"
	"github.com/mlogclub/bbs-go/common/urls"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/bbs-go/services/cache"
)

var MessageService = newMessageService()

var messageLog = logrus.WithFields(logrus.Fields{
	"type": "message",
})

func newMessageService() *messageService {
	return &messageService{
		messagesChan: make(chan *model.Message),
	}
}

type messageService struct {
	messagesChan        chan *model.Message
	messagesConsumeOnce sync.Once
}

func (this *messageService) Get(id int64) *model.Message {
	return repositories.MessageRepository.Get(simple.DB(), id)
}

func (this *messageService) Take(where ...interface{}) *model.Message {
	return repositories.MessageRepository.Take(simple.DB(), where...)
}

func (this *messageService) Find(cnd *simple.SqlCnd) []model.Message {
	return repositories.MessageRepository.Find(simple.DB(), cnd)
}

func (this *messageService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Message) {
	cnd.FindOne(db, &ret)
	return
}

func (this *messageService) FindPageByParams(params *simple.QueryParams) (list []model.Message, paging *simple.Paging) {
	return repositories.MessageRepository.FindPageByParams(simple.DB(), params)
}

func (this *messageService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Message, paging *simple.Paging) {
	return repositories.MessageRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *messageService) Create(t *model.Message) error {
	return repositories.MessageRepository.Create(simple.DB(), t)
}

func (this *messageService) Update(t *model.Message) error {
	return repositories.MessageRepository.Update(simple.DB(), t)
}

func (this *messageService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.MessageRepository.Updates(simple.DB(), id, columns)
}

func (this *messageService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.MessageRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *messageService) Delete(id int64) {
	repositories.MessageRepository.Delete(simple.DB(), id)
}

// 获取未读消息数量
func (this *messageService) GetUnReadCount(userId int64) (count int64) {
	simple.DB().Where("user_id = ? and status = ?", userId, model.MsgStatusUnread).Model(&model.Message{}).Count(&count)
	return
}

// 将所有消息标记为已读
func (this *messageService) MarkRead(userId int64) {
	simple.DB().Exec("update t_message set status = ? where user_id = ? and status = ?", model.MsgStatusReaded,
		userId, model.MsgStatusUnread)
}

// 评论被回复消息
func (this *messageService) SendCommentMsg(comment *model.Comment) {
	user := cache.UserCache.Get(comment.UserId)
	summary := common.GetMarkdownSummary(comment.Content)

	// 文章评论消息
	var (
		userId          int64
		msgContent      string
		msgQuoteContent string
		quoteContent    string
	)
	if comment.EntityType == model.EntityTypeArticle {
		article := repositories.ArticleRepository.Get(simple.DB(), comment.EntityId)
		if article != nil && article.UserId != comment.UserId {
			userId = article.UserId
			msgContent = user.Nickname + " 回复了你的文章：" + summary
			msgQuoteContent = "《" + article.Title + "》"
		}
	} else if comment.EntityType == model.EntityTypeTopic {
		topic := repositories.TopicRepository.Get(simple.DB(), comment.EntityId)
		if topic != nil && topic.UserId != comment.UserId {
			userId = topic.UserId
			msgContent = user.Nickname + " 回复了你的主题：" + summary
			msgQuoteContent = "《" + topic.Title + "》"
		}
	}
	if userId > 0 {
		this.Produce(comment.UserId, userId, msgContent, msgQuoteContent, model.MsgTypeComment, map[string]interface{}{
			"entityType": comment.EntityType,
			"entityId":   comment.EntityId,
			"commentId":  comment.Id,
			"quoteId":    comment.QuoteId,
		})
	}

	// 评论被引用的时候，给被引用的人发送消息
	if comment.QuoteId > 0 {
		quote := repositories.CommentRepository.Get(simple.DB(), comment.QuoteId)
		if quote != nil && quote.UserId != comment.UserId {
			msgContent = user.Nickname + " 回复了你的评论：" + summary
			quoteContent = common.GetMarkdownSummary(quote.Content)
			this.Produce(comment.UserId, quote.UserId, msgContent, quoteContent, model.MsgTypeComment, map[string]interface{}{
				"entityType": comment.EntityType,
				"entityId":   comment.EntityId,
				"commentId":  comment.Id,
				"quoteId":    comment.QuoteId,
			})
		}
	}
}

// 生产，将消息数据放入chan
func (this *messageService) Produce(fromId, toId int64, content, quoteContent string, msgType int, extraDataMap map[string]interface{}) {
	this.Consume()

	var (
		extraData string
		err       error
	)
	if extraData, err = simple.FormatJson(extraDataMap); err != nil {
		messageLog.Error("格式化extraData错误", err)
	}
	this.messagesChan <- &model.Message{
		FromId:       fromId,
		UserId:       toId,
		Content:      content,
		QuoteContent: quoteContent,
		Type:         msgType,
		ExtraData:    extraData,
		Status:       model.MsgStatusUnread,
		CreateTime:   simple.NowTimestamp(),
	}
}

// 消费，消费chan中的消息
func (this *messageService) Consume() {
	this.messagesConsumeOnce.Do(func() {
		go func() {
			messageLog.Info("开始消费系统消息...")
			for {
				msg := <-this.messagesChan
				messageLog.Info("处理消息：from=", msg.FromId, " to=", msg.UserId)

				if err := this.Create(msg); err != nil {
					messageLog.Info("创建消息发生异常...", err)
				} else {
					this.SendEmailNotice(msg)
				}
			}
		}()
	})
}

// 发送邮件通知
func (this *messageService) SendEmailNotice(message *model.Message) {
	user := cache.UserCache.Get(message.UserId)
	if user != nil && len(user.Email.String) > 0 {
		email.SendTemplateEmail(user.Email.String, "M-LOG新消息提醒", "M-LOG新消息提醒", message.Content,
			message.QuoteContent, urls.AbsUrl("/user/messages"))
		messageLog.Info("发送邮件...email=", user.Email)
	} else {
		messageLog.Info("邮件未发送，没设置邮箱...")
	}
}
