package admin

import (
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type EmailLogController struct {
	Ctx iris.Context
}

func (c *EmailLogController) GetBy(id int64) *web.JsonResult {
	t := services.EmailLogService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *EmailLogController) AnyList() *web.JsonResult {
	list, paging := services.EmailLogService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "toEmail",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "bizType",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
	).Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}
