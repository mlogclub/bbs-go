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
	event.RegHandler(reflect.TypeFor[event.TopicCreateEvent](), handleTaskTopicCreateEvent)
	event.RegHandler(reflect.TypeFor[event.CommentCreateEvent](), handleTaskCommentCreateEvent)
	event.RegHandler(reflect.TypeFor[event.FollowEvent](), handleTaskFollowEvent)
	event.RegHandler(reflect.TypeFor[event.UserFavoriteEvent](), handleTaskFavoriteEvent)
	event.RegHandler(reflect.TypeFor[event.UserLikeEvent](), handleTaskLikeEvent)
	event.RegHandler(reflect.TypeFor[event.CheckInEvent](), handleTaskCheckInEvent)
	event.RegHandler(reflect.TypeFor[event.UserLoginEvent](), handleTaskUserLoginEvent)
	event.RegHandler(reflect.TypeFor[event.LevelUpEvent](), handleTaskLevelUpEvent)
	event.RegHandler(reflect.TypeFor[event.BadgeGrantEvent](), handleBadgeGrantEvent)
}

func handleTaskTopicCreateEvent(i any) {
	e := i.(event.TopicCreateEvent)
	services.TaskEngineService.HandleUserEvent(e.UserId, constants.TaskEventTypeTopicCreate, e.CreateTime)
}

func handleTaskCommentCreateEvent(i any) {
	e := i.(event.CommentCreateEvent)
	services.TaskEngineService.HandleUserEvent(e.UserId, constants.TaskEventTypeCommentCreate, 0)
}

func handleTaskFollowEvent(i any) {
	e := i.(event.FollowEvent)
	services.TaskEngineService.HandleUserEvent(e.UserId, constants.TaskEventTypeFollowCreate, 0)
}

func handleTaskFavoriteEvent(i any) {
	e := i.(event.UserFavoriteEvent)
	services.TaskEngineService.HandleUserEvent(e.UserId, constants.TaskEventTypeFavoriteCreate, 0)
}

func handleTaskLikeEvent(i any) {
	e := i.(event.UserLikeEvent)
	services.TaskEngineService.HandleUserEvent(e.UserId, constants.TaskEventTypeLikeCreate, 0)
}

func handleTaskCheckInEvent(i any) {
	e := i.(event.CheckInEvent)
	services.TaskEngineService.HandleUserEvent(e.UserId, constants.TaskEventTypeCheckIn, 0)
}

func handleTaskUserLoginEvent(i any) {
	e := i.(event.UserLoginEvent)
	services.TaskEngineService.HandleUserEvent(e.UserId, constants.TaskEventTypeUserLogin, e.LoginTime)
}

func handleTaskLevelUpEvent(i any) {
	e := i.(event.LevelUpEvent)
	if e.OldLevel < 10 && e.NewLevel >= 10 {
		services.TaskEngineService.HandleUserEvent(e.UserId, constants.TaskEventTypeLevel10, e.UpdateTime)
	}

	services.MessageService.SendMsg(0, e.UserId, msg.TypeUserLevelUp,
		locales.Get("message.user_level_up_msg_title"),
		locales.Getf("message.user_level_up_msg_content", e.NewLevel),
		"", &msg.UserLevelUpExtraData{
			OldLevel: e.OldLevel,
			NewLevel: e.NewLevel,
		})
}

func handleBadgeGrantEvent(i any) {
	e := i.(event.BadgeGrantEvent)
	badgeTitle := ""
	if badge := services.BadgeService.Get(e.BadgeId); badge != nil {
		badgeTitle = badge.Title
	}
	if badgeTitle == "" {
		badgeTitle = locales.Get("message.user_badge_default_name")
	}
	services.MessageService.SendMsg(0, e.UserId, msg.TypeUserBadgeGrant,
		locales.Get("message.user_badge_grant_msg_title"),
		locales.Getf("message.user_badge_grant_msg_content", badgeTitle),
		"", &msg.UserBadgeGrantExtraData{
			BadgeId: e.BadgeId,
		})
}
