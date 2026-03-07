package eventhandler

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/msg"
	"bbs-go/internal/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.QaAnswerAcceptedEvent{}), handleQaAnswerAcceptedEvent)
}

func handleQaAnswerAcceptedEvent(i interface{}) {
	e := i.(event.QaAnswerAcceptedEvent)
	sendQaAnswerAcceptedMsg(e.TopicId, e.UserId)
}

// sendQaAnswerAcceptedMsg 回答被采纳站内通知（通知被采纳者）
func sendQaAnswerAcceptedMsg(topicId, acceptedUserId int64) {
	topic := services.TopicService.Get(topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return
	}
	if topic.UserId == acceptedUserId {
		return
	}
	var (
		from         = topic.UserId
		to           = acceptedUserId
		bountyScore  = topic.BountyScore
		quoteContent = "《" + topic.GetTitle() + "》"
	)
	title := locales.Get("message.qa_answer_accepted_msg_title")
	if bountyScore > 0 {
		title = locales.Getf("message.qa_answer_accepted_msg_title_with_bounty", bountyScore)
	}
	services.MessageService.SendMsg(from, to, msg.TypeQaAnswerAccepted, title, "", quoteContent,
		&msg.QaAnswerAcceptedExtraData{
			TopicId:     topicId,
			BountyScore: bountyScore,
		})
}
