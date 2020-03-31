package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type UserScoreLogController struct {
	Ctx iris.Context
}

func (c *UserScoreLogController) GetBy(id int64) *simple.JsonResult {
	t := services.UserScoreLogService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *UserScoreLogController) AnyList() *simple.JsonResult {
	list, paging := services.UserScoreLogService.FindPageByParams(simple.NewQueryParams(c.Ctx).
		EqByReq("user_id").EqByReq("source_type").EqByReq("source_id").EqByReq("type").PageByReq().Desc("id"))

	var results []map[string]interface{}
	for _, userScoreLog := range list {
		user := render.BuildUserDefaultIfNull(userScoreLog.UserId)
		item := simple.NewRspBuilder(userScoreLog).Put("user", user).Build()
		results = append(results, item)
	}

	return simple.JsonData(&simple.PageResult{Results: results, Page: paging})
}
