package api

import (
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/emirpasic/gods/sets/hashset"
)

func FansFollow(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}

	otherId := common.GetID(ctx, "userId")
	if otherId <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("param: userId required"))
		return
	}

	err := services.UserFollowService.Follow(user.Id, otherId)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func FansUnfollow(ctx *gin.Context) {
	user := common.GetCurrentUser(ctx)
	if user == nil {
		ginx.WriteJSON(ctx, errs.NotLogin())
		return
	}

	otherId := common.GetID(ctx, "userId")
	if otherId <= 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("param: userId required"))
		return
	}

	err := services.UserFollowService.UnFollow(user.Id, otherId)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func FansIsFollowed(ctx *gin.Context) {
	userId := common.GetID(ctx, "userId")
	current := common.GetCurrentUser(ctx)
	var followed = false
	if current != nil && current.Id != userId {
		followed = services.UserFollowService.IsFollowed(current.Id, userId)
	}
	ginx.WriteJSON(ctx, followed)

}

func FansFans(ctx *gin.Context) {
	userId := common.GetID(ctx, "userId")
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFans(userId, cursor, 10)

	current := common.GetCurrentUser(ctx)
	var followedSet hashset.Set
	if current != nil {
		followedSet = services.UserFollowService.IsFollowedUsers(current.Id, userIds...)
	}

	var itemList []*resp.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	ginx.WriteJSON(ctx, ginx.CursorData(itemList, strconv.FormatInt(cursor, 10), hasMore))

}

func FansFollowed(ctx *gin.Context) {
	userId := common.GetID(ctx, "userId")
	cursor := params.FormValueInt64Default(ctx, "cursor", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFollows(userId, cursor, 10)

	current := common.GetCurrentUser(ctx)
	var followedSet hashset.Set
	if current != nil {
		if current.Id == userId {
			followedSet = *hashset.New()
			for _, id := range userIds {
				followedSet.Add(id)
			}
		} else {
			followedSet = services.UserFollowService.IsFollowedUsers(current.Id, userIds...)
		}
	}

	var itemList []*resp.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	ginx.WriteJSON(ctx, ginx.CursorData(itemList, strconv.FormatInt(cursor, 10), hasMore))

}

func FansRecentFans(ctx *gin.Context) {
	userId := common.GetID(ctx, "userId")
	userIds, cursor, hasMore := services.UserFollowService.GetFans(userId, 0, 10)

	current := common.GetCurrentUser(ctx)
	var followedSet hashset.Set
	if current != nil {
		followedSet = services.UserFollowService.IsFollowedUsers(current.Id, userIds...)
	}

	var itemList []*resp.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	ginx.WriteJSON(ctx, ginx.CursorData(itemList, strconv.FormatInt(cursor, 10), hasMore))

}

func FansRecentFollow(ctx *gin.Context) {
	userId := common.GetID(ctx, "userId")
	userIds, cursor, hasMore := services.UserFollowService.GetFollows(userId, 0, 10)

	current := common.GetCurrentUser(ctx)
	var followedSet hashset.Set
	if current != nil {
		if current.Id == userId {
			followedSet = *hashset.New()
			for _, id := range userIds {
				followedSet.Add(id)
			}
		} else {
			followedSet = services.UserFollowService.IsFollowedUsers(current.Id, userIds...)
		}
	}

	var itemList []*resp.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	ginx.WriteJSON(ctx, ginx.CursorData(itemList, strconv.FormatInt(cursor, 10), hasMore))

}
