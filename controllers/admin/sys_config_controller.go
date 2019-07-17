package admin

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/simple"
	"strconv"
)

type SysConfigController struct {
	Ctx iris.Context
}

func (this *SysConfigController) GetBy(id int64) *simple.JsonResult {
	t := services.SysConfigService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *SysConfigController) AnyList() *simple.JsonResult {
	list, paging := services.SysConfigService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *SysConfigController) PostCreate() *simple.JsonResult {
	t := &model.SysConfig{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	err = services.SysConfigService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *SysConfigController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := services.SysConfigService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	err = services.SysConfigService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
