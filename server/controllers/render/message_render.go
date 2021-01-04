package render

import (
	"bbs-go/model"
	"bbs-go/services"
)

func BuildMessage(msg *model.Message) *model.MessageResponse {
	if msg == nil {
		return nil
	}

	from := BuildUserDefaultIfNull(msg.FromId)
	if msg.FromId <= 0 {
		from.Nickname = "系统通知"
	}
	detailUrl := services.MessageService.GetMessageDetailUrl(msg)
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
