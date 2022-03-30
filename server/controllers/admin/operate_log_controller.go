package admin

import (
	"bbs-go/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"
)

type OperateLogController struct {
	Ctx iris.Context
}

func (c *OperateLogController) GetBy(id int64) *mvc.JsonResult {
	t := services.OperateLogService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *OperateLogController) AnyList() *mvc.JsonResult {
	list, paging := services.OperateLogService.FindPageByParams(params.NewQueryParams(c.Ctx).
		EqByReq("user_id").EqByReq("op_type").EqByReq("data_type").EqByReq("data_id").
		PageByReq().Desc("id"))
	return mvc.JsonData(&sqls.PageResult{Results: list, Page: paging})
}
