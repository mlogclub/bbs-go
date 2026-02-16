package admin

import (
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type UserTaskLogController struct {
	Ctx iris.Context
}

func (c *UserTaskLogController) GetBy(id int64) *web.JsonResult {
	t := services.UserTaskLogService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *UserTaskLogController) AnyList() *web.JsonResult {
	list, paging := services.UserTaskLogService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "userId",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "taskId",
			Op:        params.Eq,
		},
	).Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}
