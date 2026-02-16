package eventhandler

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/pkg/seo"
	"bbs-go/internal/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.TopicUpdateEvent{}), handleTopicUpdateEvent)
}

func handleTopicUpdateEvent(i interface{}) {
	e := i.(event.TopicUpdateEvent)

	// 操作日志
	services.OperateLogService.AddOperateLog(e.UserId, constants.OpTypeUpdate, constants.EntityTopic, e.TopicId, "", nil)

	// 百度链接推送
	seo.Push(bbsurls.TopicUrl(e.TopicId))
}
