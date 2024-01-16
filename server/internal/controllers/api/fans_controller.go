package api

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"
	"strconv"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type FansController struct {
	Ctx iris.Context
}

func (c *FansController) PostFollow() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}

	otherId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	if otherId <= 0 {
		return web.JsonErrorMsg("param: userId required")
	}

	err := services.UserFollowService.Follow(user.Id, otherId)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *FansController) PostUnfollow() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}

	otherId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	if otherId <= 0 {
		return web.JsonErrorMsg("param: userId required")
	}

	err := services.UserFollowService.UnFollow(user.Id, otherId)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *FansController) GetIsfollowed() *web.JsonResult {
	userId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	current := services.UserTokenService.GetCurrent(c.Ctx)
	var followed = false
	if current != nil && current.Id != userId {
		followed = services.UserFollowService.IsFollowed(current.Id, userId)
	}
	return web.JsonData(followed)
}

func (c *FansController) GetFans() *web.JsonResult {
	userId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFans(userId, cursor, 10)

	current := services.UserTokenService.GetCurrent(c.Ctx)
	var followedSet hashset.Set
	if current != nil {
		followedSet = services.UserFollowService.IsFollowedUsers(current.Id, userIds...)
	}

	var itemList []*models.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	return web.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}

func (c *FansController) GetFollowed() *web.JsonResult {
	userId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFollows(userId, cursor, 10)

	current := services.UserTokenService.GetCurrent(c.Ctx)
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

	var itemList []*models.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	return web.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}

func (c *FansController) GetRecentFans() *web.JsonResult {
	userId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFans(userId, 0, 10)

	current := services.UserTokenService.GetCurrent(c.Ctx)
	var followedSet hashset.Set
	if current != nil {
		followedSet = services.UserFollowService.IsFollowedUsers(current.Id, userIds...)
	}

	var itemList []*models.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	return web.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}

func (c *FansController) GetRecentFollow() *web.JsonResult {
	userId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFollows(userId, 0, 10)

	current := services.UserTokenService.GetCurrent(c.Ctx)
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

	var itemList []*models.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	return web.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}
