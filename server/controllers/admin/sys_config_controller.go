package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/services"
)

type SysConfigController struct {
	Ctx iris.Context
}

func (c *SysConfigController) GetBy(id int64) *mvc.JsonResult {
	t := services.SysConfigService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *SysConfigController) AnyList() *mvc.JsonResult {
	list, paging := services.SysConfigService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return mvc.JsonData(&sqls.PageResult{Results: list, Page: paging})
}

func (c *SysConfigController) GetAll() *mvc.JsonResult {
	config := services.SysConfigService.GetConfig()
	return mvc.JsonData(config)
}

func (c *SysConfigController) PostSave() *mvc.JsonResult {
	config := c.Ctx.FormValue("config")
	if err := services.SysConfigService.SetAll(config); err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}
