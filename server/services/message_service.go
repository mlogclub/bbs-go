package services

import (
	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/common"
	"bbs-go/pkg/email"
	"bbs-go/pkg/msg"
	"bbs-go/pkg/urls"
	"bbs-go/repositories"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/date"
	"github.com/mlogclub/simple/json"
	"github.com/sirupsen/logrus"
)

var MessageService = newMessageService()

func newMessageService() *messageService {
	return &messageService{}
}

type messageService struct {
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

// GetUnReadCount 获取未读消息数量
func (s *messageService) GetUnReadCount(userId int64) (count int64) {
	simple.DB().Where("user_id = ? and status = ?", userId, msg.StatusUnread).Model(&model.Message{}).Count(&count)
	return
}

// MarkRead 将所有消息标记为已读
func (s *messageService) MarkRead(userId int64) {
	simple.DB().Exec("update t_message set status = ? where user_id = ? and status = ?", msg.StatusHaveRead,
		userId, msg.StatusUnread)
}

// SendTopicLikeMsg 话题收到点赞
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
	s.SendMsg(likeUserId, topic.UserId, title, "", quoteContent, msg.TypeTopicLike, &msg.TopicLikeExtraData{
		TopicId:    topicId,
		LikeUserId: likeUserId,
	})
}

// SendTopicFavoriteMsg 话题被收藏
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
	s.SendMsg(favoriteUserId, topic.UserId, title, "", quoteContent, msg.TypeTopicFavorite, &msg.TopicFavoriteExtraData{
		TopicId:        topicId,
		FavoriteUserId: favoriteUserId,
	})
}

// SendTopicRecommendMsg 话题被设为推荐
func (s *messageService) SendTopicRecommendMsg(topicId int64) {
	topic := repositories.TopicRepository.Get(simple.DB(), topicId)
	if topic == nil {
		return
	}
	var (
		title        = "你的话题被设为推荐"
		quoteContent = "《" + topic.GetTitle() + "》"
	)
	s.SendMsg(0, topic.UserId, title, "", quoteContent, msg.TypeTopicRecommend, &msg.TopicRecommendExtraData{
		TopicId: topicId,
	})
}

// SendTopicDeleteMsg 话题被删除消息
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
	s.SendMsg(0, topic.UserId, title, "", quoteContent, msg.TypeTopicDelete,
		&msg.TopicDeleteExtraData{
			TopicId:      topicId,
			DeleteUserId: deleteUserId,
		},
	)
}

// SendCommentMsg 评论被回复消息
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
			s.SendMsg(fromId, toId, title, content, quoteContent, msg.TypeTopicComment,
				&msg.CommentExtraData{
					EntityType: comment.EntityType,
					EntityId:   comment.EntityId,
					CommentId:  comment.Id,
					QuoteId:    comment.QuoteId,
				},
			)
		}

		// 给被引用的人发消息（收到他人回复）
		if fromId != quote.UserId {
			s.SendMsg(fromId, quote.UserId, "回复了你的评论", content, common.GetMarkdownSummary(quote.Content),
				msg.TypeCommentReply, &msg.CommentExtraData{
					EntityType: comment.EntityType,
					EntityId:   comment.EntityId,
					CommentId:  comment.Id,
					QuoteId:    comment.QuoteId,
				},
			)
		}
	} else if fromId != toId { // 回复主贴，并且不是自己回复自己
		// 给帖子作者发消息（收到话题评论）
		s.SendMsg(fromId, toId, title, content, quoteContent, msg.TypeTopicComment,
			&msg.CommentExtraData{
				EntityType: comment.EntityType,
				EntityId:   comment.EntityId,
				CommentId:  comment.Id,
				QuoteId:    comment.QuoteId,
			},
		)
	}
}

func (s *messageService) getQuoteComment(quoteId int64) *model.Comment {
	if quoteId <= 0 {
		return nil
	}
	return repositories.CommentRepository.Get(simple.DB(), quoteId)
}

// SendMsg 发送消息
func (s *messageService) SendMsg(fromId, toId int64, title, content, quoteContent string, msgType msg.Type, extraData interface{}) {
	to := cache.UserCache.Get(toId)
	if to == nil || to.Type != constants.UserTypeNormal {
		return
	}

	t := &model.Message{
		FromId:       fromId,
		UserId:       toId,
		Title:        title,
		Content:      content,
		QuoteContent: quoteContent,
		Type:         int(msgType),
		ExtraData:    json.ToJsonStr(extraData),
		Status:       msg.StatusUnread,
		CreateTime:   date.NowTimestamp(),
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
	}

	var from *model.User
	if t.FromId > 0 {
		from = cache.UserCache.Get(t.FromId)
	}
	err := email.SendTemplateEmail(from, user.Email.String, emailTitle, emailTitle, t.Content,
		t.QuoteContent, &model.ActionLink{
			Title: "点击查看详情",
			Url:   urls.AbsUrl("/user/messages"),
		})
	if err != nil {
		logrus.Error(err)
	}
}
