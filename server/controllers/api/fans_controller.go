package api

import (
	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
	"strconv"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
)

type FansController struct {
	Ctx iris.Context
}

func (c *FansController) PostFollow() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}

	otherId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	if otherId <= 0 {
		return mvc.JsonErrorMsg("param: userId required")
	}

	err := services.UserFollowService.Follow(user.Id, otherId)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

func (c *FansController) PostUnfollow() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}

	otherId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	if otherId <= 0 {
		return mvc.JsonErrorMsg("param: userId required")
	}

	err := services.UserFollowService.UnFollow(user.Id, otherId)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}

func (c *FansController) GetIsfollowed() *mvc.JsonResult {
	userId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	current := services.UserTokenService.GetCurrent(c.Ctx)
	var followed = false
	if current != nil && current.Id != userId {
		followed = services.UserFollowService.IsFollowed(current.Id, userId)
	}
	return mvc.JsonData(followed)
}

func (c *FansController) GetFans() *mvc.JsonResult {
	userId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFans(userId, cursor, 10)

	current := services.UserTokenService.GetCurrent(c.Ctx)
	var followedSet hashset.Set
	if current != nil {
		followedSet = services.UserFollowService.IsFollowedUsers(current.Id, userIds...)
	}

	var itemList []*model.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	return mvc.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}

func (c *FansController) GetFollows() *mvc.JsonResult {
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

	var itemList []*model.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	return mvc.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}

func (c *FansController) GetRecentFans() *mvc.JsonResult {
	userId := params.FormValueInt64Default(c.Ctx, "userId", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFans(userId, 0, 10)

	current := services.UserTokenService.GetCurrent(c.Ctx)
	var followedSet hashset.Set
	if current != nil {
		followedSet = services.UserFollowService.IsFollowedUsers(current.Id, userIds...)
	}

	var itemList []*model.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	return mvc.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}

func (c *FansController) GetRecentFollow() *mvc.JsonResult {
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

	var itemList []*model.UserInfo
	for _, id := range userIds {
		item := render.BuildUserInfoDefaultIfNull(id)
		item.Followed = followedSet.Contains(id)
		itemList = append(itemList, item)
	}
	return mvc.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}
