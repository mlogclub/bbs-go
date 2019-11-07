package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
)

type FavoriteController struct {
	Ctx iris.Context
}

func (this *FavoriteController) GetBy(id int64) *simple.JsonResult {
	t := services.FavoriteService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *FavoriteController) AnyList() *simple.JsonResult {
	list, paging := services.FavoriteService.FindPageByParams(simple.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *FavoriteController) PostCreate() *simple.JsonResult {
	t := &model.Favorite{}
	this.Ctx.ReadForm(t)

	err := services.FavoriteService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *FavoriteController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.FavoriteService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	err = services.FavoriteService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
