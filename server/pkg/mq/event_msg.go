package mq

import (
	"errors"

	"github.com/mlogclub/simple/json"
	"github.com/tidwall/gjson"
)

type EventMsg struct {
	Type EventType `json:"type"`
	Body string    `json:"body"`
}

func NewEventMsg(eventType EventType, event interface{}) (string, error) {
	body, err := json.ToStr(event)
	if err != nil {
		return "", err
	}
	eventMsg := &EventMsg{
		Type: eventType,
		Body: body,
	}
	return json.ToStr(eventMsg)
}

func ParseEvent(msg string) (EventType, interface{}, error) {
	typ := gjson.Get(msg, "type").String()
	if len(typ) == 0 {
		return EventTypeNone, nil, errors.New("invalid type")
	}
	body := gjson.Get(msg, "body").String()
	var (
		evtType = EventTypeNone
		evt     interface{}
		err     error
	)
	if typ == string(EventTypeFollow) {
		evtType = EventTypeFollow
		evt = new(FollowEvent)
		err = json.Parse(body, evt)
	} else if typ == string(EventTypeUnFollow) {
		evtType = EventTypeUnFollow
		evt = new(UnFollowEvent)
		err = json.Parse(body, evt)
	} else if typ == string(EventTypeTopicCreate) {
		evtType = EventTypeTopicCreate
		evt = new(TopicCreateEvent)
		err = json.Parse(body, evt)
	} else if typ == string(EventTypeTopicDelete) {
		evtType = EventTypeTopicDelete
		evt = new(TopicDeleteEvent)
		err = json.Parse(body, evt)
	} else if typ == string(EventTypeTopicLike) {
		evtType = EventTypeTopicLike
		evt = new(TopicLikeEvent)
		err = json.Parse(body, evt)
	} else if typ == string(EventTypeTopicFavorite) {
		evtType = EventTypeTopicFavorite
		evt = new(TopicFavoriteEvent)
		err = json.Parse(body, evt)
	} else if typ == string(EventTypeTopicRecommend) {
		evtType = EventTypeTopicRecommend
		evt = new(TopicRecommendEvent)
		err = json.Parse(body, evt)
	} else {
		err = errors.New("未处理的事件:" + msg)
	}
	return evtType, evt, err

}
