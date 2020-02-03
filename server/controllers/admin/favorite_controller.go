package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/services"
)

type FavoriteController struct {
	Ctx iris.Context
}

func (c *FavoriteController) GetBy(id int64) *simple.JsonResult {
	t := services.FavoriteService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *FavoriteController) AnyList() *simple.JsonResult {
	list, paging := services.FavoriteService.FindPageByParams(simple.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *FavoriteController) PostCreate() *simple.JsonResult {
	t := &model.Favorite{}
	simple.ReadForm(c.Ctx, t)

	err := services.FavoriteService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *FavoriteController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.FavoriteService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	simple.ReadForm(c.Ctx, t)

	err = services.FavoriteService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
