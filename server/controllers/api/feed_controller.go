package api

import (
	"bbs-go/controllers/render"
	"bbs-go/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
)

type FeedController struct {
	Ctx iris.Context
}

func (c *FeedController) GetTopics() *mvc.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return mvc.JsonError(simple.ErrorNotLogin)
	}
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	topics, cursor, hasMore := services.UserFeedService.GetTopics(user.Id, cursor)
	return mvc.JsonCursorData(render.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
}
