package eventhandler

import (
	"bbs-go/internal/pkg/event"
	"bbs-go/internal/services"
	"reflect"
)

func init() {
	event.RegHandler(reflect.TypeOf(event.MentionCreateEvent{}), handleMentionCreate)
}

func handleMentionCreate(i interface{}) {
	// Mention notifications are already sent directly in MentionService.SendMentionNotifications
	// This handler can be used for future extensibility (e.g., email digests, analytics)
	_ = i.(event.MentionCreateEvent)
	_ = services.MentionService
}