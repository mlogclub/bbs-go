package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type UserScoreLogController struct {
	Ctx iris.Context
}

func (c *UserScoreLogController) GetBy(id int64) *mvc.JsonResult {
	t := services.UserScoreLogService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *UserScoreLogController) AnyList() *mvc.JsonResult {
	list, paging := services.UserScoreLogService.FindPageByParams(params.NewQueryParams(c.Ctx).
		EqByReq("user_id").EqByReq("source_type").EqByReq("source_id").EqByReq("type").PageByReq().Desc("id"))

	var results []map[string]interface{}
	for _, userScoreLog := range list {
		user := render.BuildUserInfoDefaultIfNull(userScoreLog.UserId)
		item := mvc.NewRspBuilder(userScoreLog).Put("user", user).Build()
		results = append(results, item)
	}

	return mvc.JsonData(&sqls.PageResult{Results: results, Page: paging})
}
