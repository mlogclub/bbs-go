package admin

import (
	"strconv"

	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"

	"bbs-go/internal/pkg/simple/web"
	"bbs-go/internal/pkg/simple/web/params"
)

type OperateLogController struct {
	Ctx iris.Context
}

func (c *OperateLogController) GetBy(id int64) *web.JsonResult {
	t := services.OperateLogService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *OperateLogController) AnyList() *web.JsonResult {
	list, paging := services.OperateLogService.FindPageByParams(params.NewQueryParams(c.Ctx).
		EqByReq("user_id").EqByReq("op_type").EqByReq("data_type").EqByReq("data_id").
		PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}
