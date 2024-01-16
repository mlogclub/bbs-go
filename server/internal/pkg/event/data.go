package event

// FollowEvent 关注
type FollowEvent struct {
	UserId  int64 `json:"userId"`
	OtherId int64 `json:"otherId"`
}

// UnFollowEvent 取消关注
type UnFollowEvent struct {
	UserId  int64 `json:"userId"`
	OtherId int64 `json:"otherId"`
}

type TopicCreateEvent struct {
	UserId     int64 `json:"userId"`
	TopicId    int64 `json:"topicId"`
	CreateTime int64 `json:"createTime"`
}

type TopicDeleteEvent struct {
	UserId       int64 `json:"userId"`
	TopicId      int64 `json:"topicId"`
	DeleteUserId int64 `json:"deleteUserId"`
}

type UserLikeEvent struct {
	UserId     int64  `json:"userId"`
	EntityId   int64  `json:"entityId"`
	EntityType string `json:"entityType"`
}

type UserUnLikeEvent struct {
	UserId     int64  `json:"userId"`
	EntityId   int64  `json:"entityId"`
	EntityType string `json:"entityType"`
}

type UserFavoriteEvent struct {
	UserId     int64  `json:"userId"`
	EntityId   int64  `json:"entityId"`
	EntityType string `json:"entityType"`
}

type CommentCreateEvent struct {
	UserId    int64 `json:"userId"`
	CommentId int64 `json:"commentId"`
}

type TopicRecommendEvent struct {
	TopicId   int64 `json:"topicId"`
	Recommend bool  `json:"recommend"`
}
