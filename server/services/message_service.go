package services

import (
	"bbs-go/common/urls"
	"bbs-go/model/constants"
	"github.com/mlogclub/simple/json"
	"sync"

	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/cache"
	"bbs-go/common"
	"bbs-go/common/email"
	"bbs-go/model"
	"bbs-go/repositories"
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

func (s *messageService) Get(id int64) *model.Message {
	return repositories.MessageRepository.Get(simple.DB(), id)
}

func (s *messageService) Take(where ...interface{}) *model.Message {
	return repositories.MessageRepository.Take(simple.DB(), where...)
}

func (s *messageService) Find(cnd *simple.SqlCnd) []model.Message {
	return repositories.MessageRepository.Find(simple.DB(), cnd)
}

func (s *messageService) FindOne(cnd *simple.SqlCnd) *model.Message {
	return repositories.MessageRepository.FindOne(simple.DB(), cnd)
}

func (s *messageService) FindPageByParams(params *simple.QueryParams) (list []model.Message, paging *simple.Paging) {
	return repositories.MessageRepository.FindPageByParams(simple.DB(), params)
}

func (s *messageService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Message, paging *simple.Paging) {
	return repositories.MessageRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *messageService) Create(t *model.Message) error {
	return repositories.MessageRepository.Create(simple.DB(), t)
}

func (s *messageService) Update(t *model.Message) error {
	return repositories.MessageRepository.Update(simple.DB(), t)
}

func (s *messageService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.MessageRepository.Updates(simple.DB(), id, columns)
}

func (s *messageService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.MessageRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *messageService) Delete(id int64) {
	repositories.MessageRepository.Delete(simple.DB(), id)
}

// 获取未读消息数量
func (s *messageService) GetUnReadCount(userId int64) (count int64) {
	simple.DB().Where("user_id = ? and status = ?", userId, constants.MsgStatusUnread).Model(&model.Message{}).Count(&count)
	return
}

// 将所有消息标记为已读
func (s *messageService) MarkRead(userId int64) {
	simple.DB().Exec("update t_message set status = ? where user_id = ? and status = ?", constants.MsgStatusReaded,
		userId, constants.MsgStatusUnread)
}

// 评论被回复消息
func (s *messageService) SendCommentMsg(comment *model.Comment) {
	user := cache.UserCache.Get(comment.UserId)
	quote := s.getQuoteComment(comment.QuoteId)
	summary := common.GetSummary(comment.ContentType, comment.Content)

	var (
		fromId       = comment.UserId // 消息发送人
		authorId     int64            // 帖子作者编号
		content      string           // 消息内容
		quoteContent string           // 引用内容
	)

	if comment.EntityType == constants.EntityArticle { // 文章被评论
		article := repositories.ArticleRepository.Get(simple.DB(), comment.EntityId)
		if article != nil {
			authorId = article.UserId
			content = user.Nickname + " 回复了你的文章：" + summary
			quoteContent = "《" + article.Title + "》"
		}
	} else if comment.EntityType == constants.EntityTopic { // 话题被评论
		topic := repositories.TopicRepository.Get(simple.DB(), comment.EntityId)
		if topic != nil {
			authorId = topic.UserId
			content = user.Nickname + " 回复了你的话题：" + summary
			quoteContent = "《" + topic.Title + "》"
		}
	} else if comment.EntityType == constants.EntityTweet { // 动态被评论
		tweet := repositories.TweetRepository.Get(simple.DB(), comment.EntityId)
		if tweet != nil {
			authorId = tweet.UserId
			content = user.Nickname + " 回复了你的话题：" + summary
			quoteContent = tweet.Content
		}
	}

	if authorId <= 0 {
		return
	}

	if quote != nil { // 回复跟帖
		if comment.UserId != authorId && quote.UserId != authorId { // 回复人和帖子作者不是同一个人，并且引用的用户不是帖子作者，需要给帖子作者也发送一下消息
			// 给帖子作者发消息
			s.Produce(fromId, authorId, content, quoteContent, constants.MsgTypeComment, map[string]interface{}{
				"entityType": comment.EntityType,
				"entityId":   comment.EntityId,
				"commentId":  comment.Id,
				"quoteId":    comment.QuoteId,
			})
		}

		// 给被引用的人发消息
		s.Produce(fromId, quote.UserId, user.Nickname+" 回复了你的评论："+summary,
			common.GetMarkdownSummary(quote.Content), constants.MsgTypeComment, map[string]interface{}{
				"entityType": comment.EntityType,
				"entityId":   comment.EntityId,
				"commentId":  comment.Id,
				"quoteId":    comment.QuoteId,
			})
	} else if comment.UserId != authorId { // 回复主贴，并且不是自己回复自己
		// 给帖子作者发消息
		s.Produce(fromId, authorId, content, quoteContent, constants.MsgTypeComment, map[string]interface{}{
			"entityType": comment.EntityType,
			"entityId":   comment.EntityId,
			"commentId":  comment.Id,
			"quoteId":    comment.QuoteId,
		})
	}
}

func (s *messageService) getQuoteComment(quoteId int64) *model.Comment {
	if quoteId <= 0 {
		return nil
	}
	return repositories.CommentRepository.Get(simple.DB(), quoteId)
}

// 生产，将消息数据放入chan
func (s *messageService) Produce(fromId, toId int64, content, quoteContent string, msgType int,
	extraDataMap map[string]interface{}) {
	to := cache.UserCache.Get(toId)
	if to == nil || to.Type != constants.UserTypeNormal {
		return
	}

	s.Consume()

	var (
		extraData string
		err       error
	)
	if extraData, err = json.ToStr(extraDataMap); err != nil {
		messageLog.Error("格式化extraData错误", err)
	}
	s.messagesChan <- &model.Message{
		FromId:       fromId,
		UserId:       toId,
		Content:      content,
		QuoteContent: quoteContent,
		Type:         msgType,
		ExtraData:    extraData,
		Status:       constants.MsgStatusUnread,
		CreateTime:   simple.NowTimestamp(),
	}
}

// 消费，消费chan中的消息
func (s *messageService) Consume() {
	s.messagesConsumeOnce.Do(func() {
		go func() {
			messageLog.Info("开始消费系统消息...")
			for {
				msg := <-s.messagesChan
				messageLog.Info("处理消息：from=", msg.FromId, " to=", msg.UserId)

				if err := s.Create(msg); err != nil {
					messageLog.Info("创建消息发生异常...", err)
				} else {
					s.SendEmailNotice(msg)
				}
			}
		}()
	})
}

// 发送邮件通知
func (s *messageService) SendEmailNotice(message *model.Message) {
	user := cache.UserCache.Get(message.UserId)
	if user != nil && len(user.Email.String) > 0 {
		var (
			siteTitle  = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
			emailTitle = "新消息提醒 - " + siteTitle
		)

		_ = email.SendTemplateEmail(user.Email.String, emailTitle, emailTitle, message.Content,
			message.QuoteContent, &model.ActionLink{
				Title: "点击查看详情",
				Url:   urls.AbsUrl("/user/messages"),
			})
		messageLog.Info("发送邮件...email=", user.Email)
	} else {
		messageLog.Info("邮件未发送，没设置邮箱...")
	}
}
