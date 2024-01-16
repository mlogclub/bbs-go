package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/controllers/render"
	"bbs-go/internal/services"
)

type UserScoreLogController struct {
	Ctx iris.Context
}

func (c *UserScoreLogController) GetBy(id int64) *web.JsonResult {
	t := services.UserScoreLogService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *UserScoreLogController) AnyList() *web.JsonResult {
	list, paging := services.UserScoreLogService.FindPageByParams(params.NewQueryParams(c.Ctx).
		EqByReq("user_id").EqByReq("source_type").EqByReq("source_id").EqByReq("type").PageByReq().Desc("id"))

	var results []map[string]interface{}
	for _, userScoreLog := range list {
		user := render.BuildUserInfoDefaultIfNull(userScoreLog.UserId)
		item := web.NewRspBuilder(userScoreLog).Put("user", user).Build()
		results = append(results, item)
	}

	return web.JsonData(&web.PageResult{Results: results, Page: paging})
}
