package admin

import (
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/simple"
	"strconv"
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
	list := services.SysConfigService.GetAll()
	return simple.JsonData(list)
}

func (this *SysConfigController) PostSave() *simple.JsonResult {
	config := this.Ctx.FormValue("config")
	data := make(map[string]string)
	err := json.Unmarshal([]byte(config), &data)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	err = services.SysConfigService.SetAll(data)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}
