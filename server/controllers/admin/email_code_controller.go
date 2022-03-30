package admin

import (
	"bbs-go/model"
	"bbs-go/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type EmailCodeController struct {
	Ctx iris.Context
}

func (c *EmailCodeController) GetBy(id int64) *web.JsonResult {
	t := services.EmailCodeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *EmailCodeController) AnyList() *web.JsonResult {
	list, paging := services.EmailCodeService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *EmailCodeController) PostCreate() *web.JsonResult {
	t := &model.EmailCode{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	err = services.EmailCodeService.Create(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *EmailCodeController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := services.EmailCodeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	err = services.EmailCodeService.Update(t)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}
