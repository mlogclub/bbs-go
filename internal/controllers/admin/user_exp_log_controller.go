package admin

import (
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type UserExpLogController struct {
	Ctx iris.Context
}

func (c *UserExpLogController) GetBy(id int64) *web.JsonResult {
	t := services.UserExpLogService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *UserExpLogController) AnyList() *web.JsonResult {
	list, paging := services.UserExpLogService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "userId",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "sourceType",
			Op:        params.Eq,
		},
	).Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}
