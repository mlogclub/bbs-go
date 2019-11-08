
package admin

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/simple"
	"github.com/kataras/iris/v12"
	"strconv"
)

type ThirdAccountController struct {
	Ctx             iris.Context
}

func (this *ThirdAccountController) GetBy(id int64) *simple.JsonResult {
	t := services.ThirdAccountService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *ThirdAccountController) AnyList() *simple.JsonResult {
	list, paging := services.ThirdAccountService.FindPageByParams(simple.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *ThirdAccountController) PostCreate() *simple.JsonResult {
	t := &model.ThirdAccount{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.ThirdAccountService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *ThirdAccountController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.ThirdAccountService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	err = services.ThirdAccountService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

