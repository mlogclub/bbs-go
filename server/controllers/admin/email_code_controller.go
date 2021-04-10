package admin

import (
	"bbs-go/model"
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"strconv"
)

type EmailCodeController struct {
	Ctx iris.Context
}

func (c *EmailCodeController) GetBy(id int64) *simple.JsonResult {
	t := services.EmailCodeService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *EmailCodeController) AnyList() *simple.JsonResult {
	list, paging := services.EmailCodeService.FindPageByParams(simple.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (c *EmailCodeController) PostCreate() *simple.JsonResult {
	t := &model.EmailCode{}
	err := simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.EmailCodeService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (c *EmailCodeController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.EmailCodeService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = simple.ReadForm(c.Ctx, t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.EmailCodeService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
