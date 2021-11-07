package api

import (
	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
	"strconv"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
)

type FansController struct {
	Ctx iris.Context
}

func (c *FansController) PostFollow() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	otherId := simple.FormValueInt64Default(c.Ctx, "userId", 0)
	if otherId <= 0 {
		return simple.JsonErrorMsg("param: userId required")
	}

	err := services.UserFollowService.Follow(user.Id, otherId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

func (c *FansController) PostUnfollow() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}

	otherId := simple.FormValueInt64Default(c.Ctx, "userId", 0)
	if otherId <= 0 {
		return simple.JsonErrorMsg("param: userId required")
	}

	err := services.UserFollowService.UnFollow(user.Id, otherId)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

func (c *FansController) GetFans() *simple.JsonResult {
	userId := simple.FormValueInt64Default(c.Ctx, "userId", 0)
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFans(userId, cursor, 20)

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
	return simple.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}

func (c *FansController) GetFollows() *simple.JsonResult {
	userId := simple.FormValueInt64Default(c.Ctx, "userId", 0)
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFollows(userId, cursor, 20)

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
	return simple.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}

func (c *FansController) GetRecentFans() *simple.JsonResult {
	userId := simple.FormValueInt64Default(c.Ctx, "userId", 0)
	userIds, cursor, hasMore := services.UserFollowService.GetFans(userId, 10, 20)

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
	return simple.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}

func (c *FansController) GetRecentFollow() *simple.JsonResult {
	userId := simple.FormValueInt64Default(c.Ctx, "userId", 0)
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
	return simple.JsonCursorData(itemList, strconv.FormatInt(cursor, 10), hasMore)
}
