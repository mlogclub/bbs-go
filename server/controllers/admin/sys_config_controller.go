package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/services"
)

type SysConfigController struct {
	Ctx iris.Context
}

func (this *SysConfigController) GetBy(id int64) *simple.JsonResult {
	t := services.SysConfigService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *SysConfigController) AnyList() *simple.JsonResult {
	list, paging := services.SysConfigService.FindPageByParams(simple.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *SysConfigController) GetAll() *simple.JsonResult {
	config := services.SysConfigService.GetConfigResponse()
	return simple.JsonData(config)
}

func (this *SysConfigController) PostSave() *simple.JsonResult {
	config := this.Ctx.FormValue("config")
	if err := services.SysConfigService.SetAll(config); err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}
