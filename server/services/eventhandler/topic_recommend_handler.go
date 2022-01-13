package eventhandler

import (
	"bbs-go/pkg/event"
	"bbs-go/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicRecommendEvent{}), handleTopicRecommend)
}

func handleTopicRecommend(i interface{}) {
	e := i.(event.TopicRecommendEvent)

	if e.Recommend {
		services.MessageService.SendTopicRecommendMsg(e.TopicId)
	}
}
