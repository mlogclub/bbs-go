package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/services"
)

type SysConfigController struct {
	Ctx iris.Context
}

func (c *SysConfigController) GetBy(id int64) *simple.JsonResult {
	t := services.SysConfigService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *SysConfigController) AnyList() *simple.JsonResult {
	list, paging := services.SysConfigService.FindPageByParams(simple.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *SysConfigController) GetAll() *simple.JsonResult {
	config := services.SysConfigService.GetConfig()
	return simple.JsonData(config)
}

func (c *SysConfigController) PostSave() *simple.JsonResult {
	config := c.Ctx.FormValue("config")
	if err := services.SysConfigService.SetAll(config); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}
