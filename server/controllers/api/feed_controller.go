package api

import (
	"bbs-go/controllers/render"
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"strconv"
)

type FeedController struct {
	Ctx iris.Context
}

func (c *FeedController) GetTopics() *simple.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return simple.JsonError(simple.ErrorNotLogin)
	}
	cursor := simple.FormValueInt64Default(c.Ctx, "cursor", 0)
	topics, cursor, hasMore := services.UserFeedService.GetTopics(user.Id, cursor)
	return simple.JsonCursorData(render.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
}
