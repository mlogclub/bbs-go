package render

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/urls"
	"github.com/tidwall/gjson"
)

func BuildMessage(msg *model.Message) *model.MessageResponse {
	if msg == nil {
		return nil
	}

	from := BuildUserInfoDefaultIfNull(msg.FromId)
	if msg.FromId <= 0 {
		from.Nickname = "系统通知"
	}
	detailUrl := getMessageDetailUrl(msg)
	resp := &model.MessageResponse{
		MessageId:    msg.Id,
		From:         from,
		UserId:       msg.UserId,
		Title:        msg.Title,
		Content:      msg.Content,
		QuoteContent: msg.QuoteContent,
		Type:         msg.Type,
		DetailUrl:    detailUrl,
		ExtraData:    msg.ExtraData,
		Status:       msg.Status,
		CreateTime:   msg.CreateTime,
	}
	return resp
}

// BuildMessages 渲染消息列表
func BuildMessages(messages []model.Message) []model.MessageResponse {
	if len(messages) == 0 {
		return nil
	}
	var responses []model.MessageResponse
	for _, message := range messages {
		responses = append(responses, *BuildMessage(&message))
	}
	return responses
}

// getMessageDetailUrl 查看消息详情链接地址
func getMessageDetailUrl(msg *model.Message) string {
	msgType := constants.MsgType(msg.Type)
	if msgType == constants.MsgTypeTopicComment ||
		msgType == constants.MsgTypeCommentReply {
		entityType := gjson.Get(msg.ExtraData, "entityType")
		entityId := gjson.Get(msg.ExtraData, "entityId")

		if entityType.String() == constants.EntityArticle {
			return urls.ArticleUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTopic {
			return urls.TopicUrl(entityId.Int())
		}
	} else if msgType == constants.MsgTypeTopicLike ||
		msgType == constants.MsgTypeTopicFavorite ||
		msgType == constants.MsgTypeTopicRecommend {
		topicId := gjson.Get(msg.ExtraData, "topicId")
		if topicId.Exists() && topicId.Int() > 0 {
			return urls.TopicUrl(topicId.Int())
		}
	}
	return urls.AbsUrl("/user/messages")
}
