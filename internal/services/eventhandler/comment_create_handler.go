package eventhandler

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/msg"
	"bbs-go/internal/services"
	"log/slog"
	"reflect"

	"github.com/spf13/cast"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.CommentCreateEvent{}), handleCommentCreate)
}

func handleCommentCreate(i interface{}) {
	e := i.(event.CommentCreateEvent)

	comment := services.CommentService.Get(e.CommentId)

	// 发送消息
	handleMsg(comment)
}

// 处理评论消息
func handleMsg(comment *models.Comment) {
	commentMsg := getCommentMsg(comment)

	handleEntityMsg(comment, commentMsg)
	handleQuoteMsg(comment, commentMsg)
	handleReplyMsg(comment, commentMsg)
}

// 给被回复的实体对象作者发送消息
func handleEntityMsg(comment *models.Comment, commentMsg *CommentMsg) {
	var (
		from = comment.UserId
		to   = commentMsg.rootEntityUserId()
	)
	if from == to {
		return
	}
	if to <= 0 {
		slog.Warn("消息发送失败", slog.Any("to", to))
		return
	}

	// 如果回复的评论作者就是帖子作者，那么只给回复作者发消息即可，这里就不再给帖子作者发消息了
	if commentMsg.ParentComment != nil && commentMsg.ParentComment.UserId == to {
		return
	}
	// 同上
	if commentMsg.QuoteComment != nil && commentMsg.QuoteComment.UserId == to {
		return
	}

	services.MessageService.SendMsg(from, to,
		commentMsg.msgType(),
		commentMsg.msgTitle(),
		commentMsg.msgContent(),
		commentMsg.msgRepliedContent(),
		&msg.CommentExtraData{
			EntityType:     commentMsg.EntityType,
			EntityId:       commentMsg.EntityId,
			QuoteId:        comment.QuoteId,
			RootEntityType: commentMsg.rootEntityType(),
			RootEntityId:   cast.ToString(commentMsg.rootEntityId()),
		})
}

func handleReplyMsg(comment *models.Comment, commentMsg *CommentMsg) {
	if commentMsg.ParentComment == nil {
		return
	}

	var (
		from = comment.UserId
		to   = commentMsg.ParentComment.UserId
	)

	if from == to {
		return
	}

	// 如果回复的评论作者就是被引用消息的作者作者，那么只发引用消息即可
	if commentMsg.QuoteComment != nil && commentMsg.QuoteComment.UserId == to {
		return
	}

	var (
		title          = commentMsg.msgTitle()
		content        = commentMsg.msgContent()
		repliedContent = common.GetSummary(commentMsg.ParentComment.ContentType, commentMsg.ParentComment.Content)
	)

	services.MessageService.SendMsg(from, to, msg.TypeCommentReply, title, content, repliedContent,
		&msg.CommentExtraData{
			EntityType:     comment.EntityType,
			EntityId:       comment.EntityId,
			QuoteId:        comment.QuoteId,
			RootEntityType: commentMsg.rootEntityType(),
			RootEntityId:   cast.ToString(commentMsg.rootEntityId()),
		})
}

// handleQuoteMsg 给被引用人发送消息
func handleQuoteMsg(comment *models.Comment, commentMsg *CommentMsg) {
	if commentMsg.QuoteComment == nil {
		return
	}

	var (
		from           = comment.UserId
		to             = commentMsg.QuoteComment.UserId
		title          = commentMsg.msgTitle()
		content        = commentMsg.msgContent()
		repliedContent = common.GetSummary(commentMsg.QuoteComment.ContentType, commentMsg.QuoteComment.Content)
	)

	if from == to {
		return
	}

	services.MessageService.SendMsg(from, to, msg.TypeCommentReply, title, content, repliedContent,
		&msg.CommentExtraData{
			EntityType:     comment.EntityType,
			EntityId:       comment.EntityId,
			QuoteId:        comment.QuoteId,
			RootEntityType: commentMsg.rootEntityType(),
			RootEntityId:   cast.ToString(commentMsg.rootEntityId()),
		})
}

func getCommentMsg(comment *models.Comment) *CommentMsg {
	if comment.EntityType == constants.EntityTopic { // 帖子
		topic := services.TopicService.Get(comment.EntityId)
		if topic != nil && topic.Status == constants.StatusOk {
			return &CommentMsg{
				Comment:    comment,
				EntityType: comment.EntityType,
				EntityId:   comment.EntityId,
				Entity:     topic,
			}
		}
	} else if comment.EntityType == constants.EntityArticle { // 文章
		article := services.ArticleService.Get(comment.EntityId)
		if article != nil && article.Status == constants.StatusOk {
			return &CommentMsg{
				Comment:    comment,
				EntityType: comment.EntityType,
				EntityId:   comment.EntityId,
				Entity:     article,
			}
		}
	} else if comment.EntityType == constants.EntityComment { // 二级评论
		parentComment := services.CommentService.Get(comment.EntityId)
		if parentComment == nil || parentComment.Status != constants.StatusOk {
			return nil
		}

		ret := &CommentMsg{
			Comment:       comment,
			EntityType:    parentComment.EntityType, // 二级评论时，取一级评论的
			EntityId:      parentComment.EntityId,   // 二级评论时，取一级评论的
			ParentComment: parentComment,            // 一级评论
		}

		if parentComment.EntityType == constants.EntityTopic {
			topic := services.TopicService.Get(parentComment.EntityId)
			if topic != nil && topic.Status == constants.StatusOk {
				ret.Entity = topic
			}
		} else if parentComment.EntityType == constants.EntityArticle {
			article := services.ArticleService.Get(parentComment.EntityId)
			if article != nil && article.Status == constants.StatusOk {
				ret.Entity = article
			}
		} else {
			return nil
		}

		if comment.QuoteId > 0 { // 三级评论
			quoteComment := services.CommentService.Get(comment.QuoteId)
			if quoteComment != nil && quoteComment.Status == constants.StatusOk {
				ret.QuoteComment = quoteComment
			}
		}

		return ret
	}
	return nil
}

type CommentMsg struct {
	EntityType    string          // 实体类型
	EntityId      int64           // 实体ID
	Entity        interface{}     // 被评论实体
	Comment       *models.Comment // 当前评论
	ParentComment *models.Comment // 上一级评论（二级评论的时候有值）
	QuoteComment  *models.Comment // 引用评论
}

// msgType 消息类型
func (c *CommentMsg) msgType() msg.Type {
	if c.EntityType == constants.EntityTopic {
		return msg.TypeTopicComment
	} else if c.EntityType == constants.EntityArticle {
		return msg.TypeArticleComment
	} else if c.EntityType == constants.EntityComment {
		return msg.TypeCommentReply
	}
	return msg.TypeTopicComment
}

// msgTitle 消息标题
func (c *CommentMsg) msgTitle() string {
	if c.EntityType == constants.EntityTopic {
		return "回复了你的话题"
	} else if c.EntityType == constants.EntityArticle {
		return "回复了你的文章"
	} else if c.EntityType == constants.EntityComment {
		return "回复了你的评论"
	}
	return ""
}

// msgContent 回复内容
func (c *CommentMsg) msgContent() string {
	return common.GetSummary(c.Comment.ContentType, c.Comment.Content)
}

// msgRepliedContent 被回复的内容
func (c *CommentMsg) msgRepliedContent() string {
	if c.EntityType == constants.EntityArticle {
		article := c.Entity.(*models.Article)
		return "《" + article.Title + "》"
	} else if c.EntityType == constants.EntityTopic {
		topic := c.Entity.(*models.Topic)
		return "《" + topic.GetTitle() + "》"
	}
	return ""
}

func (c *CommentMsg) rootEntityUserId() int64 {
	if c.ParentComment != nil { // 二级评论
		if c.ParentComment.EntityType == constants.EntityTopic {
			topic := c.Entity.(*models.Topic)
			return topic.UserId
		} else if c.ParentComment.EntityType == constants.EntityArticle {
			article := c.Entity.(*models.Article)
			return article.UserId
		}
	} else {
		if c.Comment.EntityType == constants.EntityTopic {
			topic := c.Entity.(*models.Topic)
			return topic.UserId
		} else if c.Comment.EntityType == constants.EntityArticle {
			article := c.Entity.(*models.Article)
			return article.UserId
		}
	}
	return 0
}

func (c *CommentMsg) rootEntityType() string {
	if c.ParentComment != nil { // 二级评论
		return c.ParentComment.EntityType
	} else {
		return c.Comment.EntityType
	}
}

func (c *CommentMsg) rootEntityId() int64 {
	if c.ParentComment != nil { // 二级评论
		return c.ParentComment.EntityId
	} else {
		return c.Comment.EntityId
	}
}
