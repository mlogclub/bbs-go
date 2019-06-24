package admin

import (
	"strconv"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
)

type FavoriteController struct {
	Ctx             iris.Context
	FavoriteService *services.FavoriteService
}

func (this *FavoriteController) GetBy(id int64) *simple.JsonResult {
	t := this.FavoriteService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *FavoriteController) AnyList() *simple.JsonResult {
	list, paging := this.FavoriteService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *FavoriteController) PostCreate() *simple.JsonResult {
	t := &model.Favorite{}
	this.Ctx.ReadForm(t)

	err := this.FavoriteService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *FavoriteController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := this.FavoriteService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = this.FavoriteService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
