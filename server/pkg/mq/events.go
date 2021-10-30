package mq

type EventType string

const (
	EventTypeNone           EventType = "none"
	EventTypeFollow         EventType = "follow"
	EventTypeUnFollow       EventType = "unfollow"
	EventTypeTopicCreate    EventType = "topicCreate"
	EventTypeTopicDelete    EventType = "topicDelete"
	EventTypeTopicLike      EventType = "topicLike"
	EventTypeTopicFavorite  EventType = "topicFavorite"
	EventTypeTopicRecommend EventType = "topicRecommend"
)

// 关注
type FollowEvent struct {
	UserId  int64 `json:"userId"`
	OtherId int64 `json:"otherId"`
}

// 取消关注
type UnFollowEvent struct {
	UserId  int64 `json:"userId"`
	OtherId int64 `json:"otherId"`
}

type TopicCreateEvent struct {
	UserId  int64 `json:"userId"`
	TopicId int64 `json:"topicId"`
}

type TopicDeleteEvent struct {
	UserId       int64 `json:"userId"`
	TopicId      int64 `json:"topicId"`
	DeleteUserId int64 `json:"deleteUserId"`
}

type TopicLikeEvent struct {
	UserId  int64 `json:"userId"`
	TopicId int64 `json:"topicId"`
}

type TopicFavoriteEvent struct {
	UserId  int64 `json:"userId"`
	TopicId int64 `json:"topicId"`
}

type TopicRecommendEvent struct {
	UserId  int64 `json:"userId"`
	TopicId int64 `json:"topicId"`
}
