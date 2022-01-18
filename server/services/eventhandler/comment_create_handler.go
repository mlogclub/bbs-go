package eventhandler

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/common"
	"bbs-go/pkg/event"
	"bbs-go/pkg/msg"
	"reflect"
)

import (
	"bbs-go/services"
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
func handleMsg(comment *model.Comment) {
	commentMsg := getCommentMsg(comment)

	handleEntityMsg(comment, commentMsg)
	handleQuoteMsg(comment, commentMsg)
	handleReplyMsg(comment, commentMsg)
}

// 给被回复的实体对象作者发送消息
func handleEntityMsg(comment *model.Comment, commentMsg *CommentMsg) {
	var (
		from = comment.UserId
		to   = commentMsg.EntityUserId
	)
	if from == to {
		return
	}

	// 如果回复的评论作者就是帖子作者，那么只给回复作者发消息即可，这里就不再给帖子作者发消息了
	if commentMsg.FirstComment != nil && commentMsg.FirstComment.UserId == to {
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
			EntityType: commentMsg.EntityType,
			EntityId:   commentMsg.EntityId,
			CommentId:  comment.Id,
			QuoteId:    comment.QuoteId,
		})
}

func handleReplyMsg(comment *model.Comment, commentMsg *CommentMsg) {
	if commentMsg.FirstComment == nil {
		return
	}

	var (
		from = comment.UserId
		to   = commentMsg.FirstComment.UserId
	)

	if from == to {
		return
	}

	// 如果回复的评论作者就是被引用消息的作者作者，那么只发引用消息即可
	if commentMsg.QuoteComment != nil && commentMsg.QuoteComment.UserId == to {
		return
	}

	var (
		title          = "回复了你的评论"
		content        = commentMsg.msgContent()
		repliedContent = common.GetSummary(commentMsg.FirstComment.ContentType, commentMsg.FirstComment.Content)
	)

	services.MessageService.SendMsg(from, to, msg.TypeCommentReply, title, content, repliedContent,
		&msg.CommentExtraData{
			EntityType: comment.EntityType,
			EntityId:   comment.EntityId,
			CommentId:  comment.Id,
			QuoteId:    comment.QuoteId,
		})
}

// handleQuoteMsg 给被引用人发送消息
func handleQuoteMsg(comment *model.Comment, commentMsg *CommentMsg) {
	if commentMsg.QuoteComment == nil {
		return
	}

	var (
		from           = comment.UserId
		to             = commentMsg.QuoteComment.UserId
		title          = "回复了你的评论"
		content        = commentMsg.msgContent()
		repliedContent = common.GetSummary(commentMsg.QuoteComment.ContentType, commentMsg.QuoteComment.Content)
	)

	if from == to {
		return
	}

	services.MessageService.SendMsg(from, to, msg.TypeCommentReply, title, content, repliedContent,
		&msg.CommentExtraData{
			EntityType: comment.EntityType,
			EntityId:   comment.EntityId,
			CommentId:  comment.Id,
			QuoteId:    comment.QuoteId,
		})
}

func getCommentMsg(comment *model.Comment) *CommentMsg {
	if comment.EntityType == constants.EntityTopic { // 帖子
		topic := services.TopicService.Get(comment.EntityId)
		if topic != nil && topic.Status == constants.StatusOk {
			return &CommentMsg{
				Comment:      comment,
				EntityType:   comment.EntityType,
				EntityId:     comment.EntityId,
				EntityUserId: topic.UserId,
				Entity:       topic,
			}
		}
	} else if comment.EntityType == constants.EntityArticle { // 文章
		article := services.ArticleService.Get(comment.EntityId)
		if article != nil && article.Status == constants.StatusOk {
			return &CommentMsg{
				Comment:      comment,
				EntityType:   comment.EntityType,
				EntityId:     comment.EntityId,
				EntityUserId: article.UserId,
				Entity:       article,
			}
		}
	} else if comment.EntityType == constants.EntityComment { // 二级评论
		firstComment := services.CommentService.Get(comment.EntityId)
		if firstComment == nil || firstComment.Status != constants.StatusOk {
			return nil
		}

		ret := &CommentMsg{
			Comment:      comment,
			EntityType:   firstComment.EntityType, // 二级评论时，取一级评论的
			EntityId:     firstComment.EntityId,   // 二级评论时，取一级评论的
			FirstComment: firstComment,            // 一级评论
		}

		if firstComment.EntityType == constants.EntityTopic {
			topic := services.TopicService.Get(firstComment.EntityId)
			if topic != nil && topic.Status == constants.StatusOk {
				ret.EntityUserId = topic.UserId
				ret.Entity = topic
			}
		} else if firstComment.EntityType == constants.EntityArticle {
			article := services.ArticleService.Get(firstComment.EntityId)
			if article != nil && article.Status == constants.StatusOk {
				ret.EntityUserId = article.UserId
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
	EntityType   string         // 实体类型
	EntityId     int64          // 实体ID
	EntityUserId int64          // 实体所属人
	Entity       interface{}    // 被评论实体
	Comment      *model.Comment // 当前评论
	FirstComment *model.Comment // 根评论
	QuoteComment *model.Comment // 引用评论
}

// msgType 消息类型
func (c *CommentMsg) msgType() msg.Type {
	if c.EntityType == constants.EntityTopic {
		return msg.TypeTopicComment
	} else if c.EntityType == constants.EntityArticle {
		return msg.TypeArticleComment
	}
	return msg.TypeTopicComment
}

// msgTitle 消息标题
func (c *CommentMsg) msgTitle() string {
	if c.EntityType == constants.EntityTopic {
		return "回复了你的话题"
	} else if c.EntityType == constants.EntityArticle {
		return "回复了你的文章"
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
		article := c.Entity.(*model.Article)
		return "《" + article.Title + "》"
	} else if c.EntityType == constants.EntityTopic {
		topic := c.Entity.(*model.Topic)
		return "《" + topic.GetTitle() + "》"
	}
	return ""
}
