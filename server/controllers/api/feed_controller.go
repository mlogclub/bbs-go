package api

import (
	"bbs-go/controllers/render"
	"bbs-go/pkg/common"
	"bbs-go/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type FeedController struct {
	Ctx iris.Context
}

func (c *FeedController) GetTopics() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(common.ErrorNotLogin)
	}
	cursor := params.FormValueInt64Default(c.Ctx, "cursor", 0)
	topics, cursor, hasMore := services.UserFeedService.GetTopics(user.Id, cursor)
	return web.JsonCursorData(render.BuildSimpleTopics(topics, user), strconv.FormatInt(cursor, 10), hasMore)
}
