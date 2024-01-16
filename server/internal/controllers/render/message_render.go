package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/msg"

	"github.com/tidwall/gjson"
)

func BuildMessage(msg *models.Message) *models.MessageResponse {
	if msg == nil {
		return nil
	}

	from := BuildUserInfoDefaultIfNull(msg.FromId)
	if msg.FromId <= 0 {
		from.Nickname = "系统通知"
	}
	detailUrl := getMessageDetailUrl(msg)
	resp := &models.MessageResponse{
		Id:           msg.Id,
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
func BuildMessages(messages []models.Message) []models.MessageResponse {
	if len(messages) == 0 {
		return nil
	}
	var responses []models.MessageResponse
	for _, message := range messages {
		responses = append(responses, *BuildMessage(&message))
	}
	return responses
}

// getMessageDetailUrl 查看消息详情链接地址
func getMessageDetailUrl(t *models.Message) string {
	msgType := msg.Type(t.Type)
	if msgType == msg.TypeTopicComment || msgType == msg.TypeArticleComment {
		entityType := gjson.Get(t.ExtraData, "entityType")
		entityId := gjson.Get(t.ExtraData, "entityId")
		if entityType.String() == constants.EntityArticle {
			return bbsurls.ArticleUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTopic {
			return bbsurls.TopicUrl(entityId.Int())
		}
	} else if msgType == msg.TypeCommentReply {
		entityType := gjson.Get(t.ExtraData, "rootEntityType")
		entityId := gjson.Get(t.ExtraData, "rootEntityId")

		if entityType.String() == constants.EntityArticle {
			return bbsurls.ArticleUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTopic {
			return bbsurls.TopicUrl(entityId.Int())
		}
	} else if msgType == msg.TypeTopicLike ||
		msgType == msg.TypeTopicFavorite ||
		msgType == msg.TypeTopicRecommend {
		topicId := gjson.Get(t.ExtraData, "topicId")
		if topicId.Exists() && topicId.Int() > 0 {
			return bbsurls.TopicUrl(topicId.Int())
		}
	}
	return bbsurls.AbsUrl("/user/messages")
}
