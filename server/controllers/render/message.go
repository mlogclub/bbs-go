package render

import (
	"bbs-go/common/avatar"
	"bbs-go/common/urls"
	"bbs-go/model"
	"bbs-go/model/constants"
	"github.com/tidwall/gjson"
)

func BuildMessage(message *model.Message) *model.MessageResponse {
	if message == nil {
		return nil
	}

	resp := &model.MessageResponse{
		MessageId:    message.Id,
		UserId:       message.UserId,
		Content:      message.Content,
		QuoteContent: message.QuoteContent,
		Type:         message.Type,
		ExtraData:    message.ExtraData,
		Status:       message.Status,
		CreateTime:   message.CreateTime,
	}

	// 消息发送人
	resp.From = BuildUserDefaultIfNull(message.FromId)
	if message.FromId <= 0 {
		resp.From.Nickname = "系统通知"
		resp.From.Avatar = avatar.DefaultAvatar
	}

	// 详情链接地址
	if message.Type == constants.MsgTypeComment {
		entityType := gjson.Get(message.ExtraData, "entityType")
		entityId := gjson.Get(message.ExtraData, "entityId")
		if entityType.String() == constants.EntityArticle {
			resp.DetailUrl = urls.ArticleUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTopic {
			resp.DetailUrl = urls.TopicUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTweet {
			resp.DetailUrl = urls.TweetUrl(entityId.Int())
		}
	}

	if message.Type == constants.MsgTypeComment {
		// TODO 渲染消息对应的文章、话题，样式参见开源中国的动态消息
	}

	return resp
}

// 渲染消息列表
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
