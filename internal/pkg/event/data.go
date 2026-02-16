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

type TopicUpdateEvent struct {
	UserId  int64 `json:"userId"`
	TopicId int64 `json:"topicId"`
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

// UserLoginEvent 登录成功
type UserLoginEvent struct {
	UserId     int64 `json:"userId"`
	LoginTime  int64 `json:"loginTime"`
	IsNewLogin bool  `json:"isNewLogin"` // 是否是新登录
}

// CheckInEvent 签到
type CheckInEvent struct {
	UserId  int64 `json:"userId"`
	DayName int   `json:"dayName"`
}

// LevelUpEvent 等级提升
type LevelUpEvent struct {
	UserId     int64 `json:"userId"`
	OldLevel   int   `json:"oldLevel"`
	NewLevel   int   `json:"newLevel"`
	UpdateTime int64 `json:"updateTime"`
}
