package services

import (
	"bbs-go/model/constants"
	"bbs-go/package/email"
	"bbs-go/package/urls"
	"github.com/mlogclub/simple/date"
	"github.com/mlogclub/simple/json"
	"github.com/tidwall/gjson"
	"sync"

	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/package/common"
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
	simple.DB().Exec("update t_message set status = ? where user_id = ? and status = ?", constants.MsgStatusHaveRead,
		userId, constants.MsgStatusUnread)
}

// 话题收到点赞
func (s *messageService) SendTopicLikeMsg(topicId, likeUserId int64) {
	topic := repositories.TopicRepository.Get(simple.DB(), topicId)
	if topic == nil {
		return
	}
	if topic.UserId == likeUserId {
		return
	}
	var (
		title        = "点赞了你的话题"
		quoteContent = "《" + topic.GetTitle() + "》"
	)
	s.Produce(likeUserId, topic.UserId, title, "", quoteContent, constants.MsgTypeTopicLike, map[string]interface{}{
		"topicId":    topicId,
		"likeUserId": likeUserId,
	})
}

// 话题被收藏
func (s *messageService) SendTopicFavoriteMsg(topicId, favoriteUserId int64) {
	topic := repositories.TopicRepository.Get(simple.DB(), topicId)
	if topic == nil {
		return
	}
	if topic.UserId == favoriteUserId {
		return
	}
	var (
		title        = "收藏了你的话题"
		quoteContent = "《" + topic.GetTitle() + "》"
	)
	s.Produce(favoriteUserId, topic.UserId, title, "", quoteContent, constants.MsgTypeTopicFavorite, map[string]interface{}{
		"topicId":        topicId,
		"favoriteUserId": favoriteUserId,
	})
}

// 话题被设为推荐
func (s *messageService) SendTopicRecommendMsg(topicId int64) {
	topic := repositories.TopicRepository.Get(simple.DB(), topicId)
	if topic == nil {
		return
	}
	var (
		title        = "你的话题被设为推荐"
		quoteContent = "《" + topic.GetTitle() + "》"
	)
	s.Produce(0, topic.UserId, title, "", quoteContent, constants.MsgTypeTopicRecommend, map[string]interface{}{
		"topicId": topicId,
	})
}

// 话题被删除消息
func (s *messageService) SendTopicDeleteMsg(topicId, deleteUserId int64) {
	topic := repositories.TopicRepository.Get(simple.DB(), topicId)
	if topic == nil {
		return
	}
	if topic.UserId == deleteUserId {
		return
	}
	var (
		title        = "你的话题被删除"
		quoteContent = "《" + topic.GetTitle() + "》"
	)
	s.Produce(0, topic.UserId, title, "", quoteContent, constants.MsgTypeTopicDelete, map[string]interface{}{
		"topicId":      topicId,
		"deleteUserId": deleteUserId,
	})
}

// 评论被回复消息
func (s *messageService) SendCommentMsg(comment *model.Comment) {
	var (
		fromId       = comment.UserId                                          // 消息发送人
		toId         int64                                                     // 消息接收人
		title        string                                                    // 消息的标题
		content      = common.GetSummary(comment.ContentType, comment.Content) // 消息内容
		quoteContent string                                                    // 引用内容
	)

	if comment.EntityType == constants.EntityArticle { // 文章被评论
		article := repositories.ArticleRepository.Get(simple.DB(), comment.EntityId)
		if article != nil {
			toId = article.UserId
			title = "回复了你的文章"
			quoteContent = "《" + article.Title + "》"
		}
	} else if comment.EntityType == constants.EntityTopic { // 话题被评论
		topic := repositories.TopicRepository.Get(simple.DB(), comment.EntityId)
		if topic != nil {
			toId = topic.UserId
			title = "回复了你的话题"
			quoteContent = "《" + topic.GetTitle() + "》"
		}
	}

	if toId <= 0 {
		return
	}

	quote := s.getQuoteComment(comment.QuoteId)
	if quote != nil { // 回复跟帖
		// 回复人和帖子作者不是同一个人，并且引用的用户不是帖子作者，需要给帖子作者也发送一下消息
		if fromId != toId && quote.UserId != toId {
			// 给帖子作者发消息（收到话题评论）
			s.Produce(fromId, toId, title, content, quoteContent, constants.MsgTypeTopicComment,
				map[string]interface{}{
					"entityType": comment.EntityType,
					"entityId":   comment.EntityId,
					"commentId":  comment.Id,
					"quoteId":    comment.QuoteId,
				})
		}

		// 给被引用的人发消息（收到他人回复）
		if fromId != quote.UserId {
			s.Produce(fromId, quote.UserId, "回复了你的评论", content, common.GetMarkdownSummary(quote.Content),
				constants.MsgTypeCommentReply, map[string]interface{}{
					"entityType": comment.EntityType,
					"entityId":   comment.EntityId,
					"commentId":  comment.Id,
					"quoteId":    comment.QuoteId,
				})
		}
	} else if fromId != toId { // 回复主贴，并且不是自己回复自己
		// 给帖子作者发消息（收到话题评论）
		s.Produce(fromId, toId, title, content, quoteContent, constants.MsgTypeTopicComment,
			map[string]interface{}{
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
func (s *messageService) Produce(fromId, toId int64, title, content, quoteContent string, msgType int, extraDataMap map[string]interface{}) {
	s.Consume()

	to := cache.UserCache.Get(toId)
	if to == nil || to.Type != constants.UserTypeNormal {
		return
	}

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
		Title:        title,
		Content:      content,
		QuoteContent: quoteContent,
		Type:         msgType,
		ExtraData:    extraData,
		Status:       constants.MsgStatusUnread,
		CreateTime:   date.NowTimestamp(),
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
func (s *messageService) SendEmailNotice(msg *model.Message) {
	user := cache.UserCache.Get(msg.UserId)
	if user != nil && len(user.Email.String) > 0 {
		var (
			siteTitle  = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
			emailTitle = siteTitle + " - 新消息提醒"
		)

		if msg.Type == constants.MsgTypeTopicComment {
			emailTitle = siteTitle + " - 收到话题评论"
		} else if msg.Type == constants.MsgTypeCommentReply {
			emailTitle = siteTitle + " - 收到他人回复"
		} else if msg.Type == constants.MsgTypeTopicLike {
			emailTitle = siteTitle + " - 收到点赞"
		} else if msg.Type == constants.MsgTypeTopicFavorite {
			emailTitle = siteTitle + " - 话题被收藏"
		} else if msg.Type == constants.MsgTypeTopicRecommend {
			emailTitle = siteTitle + " - 话题被设为推荐"
		} else if msg.Type == constants.MsgTypeTopicDelete {
			emailTitle = siteTitle + " - 话题被删除"
		}

		var from *model.User
		if msg.FromId > 0 {
			from = cache.UserCache.Get(msg.FromId)
		}
		err := email.SendTemplateEmail(from, user.Email.String, emailTitle, emailTitle, msg.Content,
			msg.QuoteContent, &model.ActionLink{
				Title: "点击查看详情",
				Url:   urls.AbsUrl("/user/messages"),
			})
		if err != nil {
			logrus.Error(err)
		}
	} else {
		messageLog.Info("邮件未发送，没设置邮箱...")
	}
}

// 查看消息详情链接地址
func (s *messageService) GetMessageDetailUrl(msg *model.Message) string {
	if msg.Type == constants.MsgTypeTopicComment ||
		msg.Type == constants.MsgTypeCommentReply {
		entityType := gjson.Get(msg.ExtraData, "entityType")
		entityId := gjson.Get(msg.ExtraData, "entityId")

		if entityType.String() == constants.EntityArticle {
			return urls.ArticleUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTopic {
			return urls.TopicUrl(entityId.Int())
		}
	} else if msg.Type == constants.MsgTypeTopicLike ||
		msg.Type == constants.MsgTypeTopicFavorite ||
		msg.Type == constants.MsgTypeTopicRecommend {
		topicId := gjson.Get(msg.ExtraData, "topicId")
		if topicId.Exists() && topicId.Int() > 0 {
			return urls.TopicUrl(topicId.Int())
		}
	}
	return urls.AbsUrl("/user/messages")
}
