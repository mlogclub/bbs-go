package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/model"
	"bbs-go/services"
)

type FavoriteController struct {
	Ctx iris.Context
}

func (c *FavoriteController) GetBy(id int64) *mvc.JsonResult {
	t := services.FavoriteService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *FavoriteController) AnyList() *mvc.JsonResult {
	list, paging := services.FavoriteService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return mvc.JsonData(&sqls.PageResult{Results: list, Page: paging})
}

func (c *FavoriteController) PostCreate() *mvc.JsonResult {
	t := &model.Favorite{}
	params.ReadForm(c.Ctx, t)

	err := services.FavoriteService.Create(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}

func (c *FavoriteController) PostUpdate() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	t := services.FavoriteService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("entity not found")
	}

	params.ReadForm(c.Ctx, t)

	err = services.FavoriteService.Update(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}
