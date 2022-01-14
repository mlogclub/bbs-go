package eventhandler

import (
	"bbs-go/pkg/event"
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
	services.MessageService.SendCommentMsg(comment)
}
