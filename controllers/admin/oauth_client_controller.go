package admin

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/simple"
	"strconv"
)

type OauthClientController struct {
	Ctx iris.Context
}

func (this *OauthClientController) GetBy(id int64) *simple.JsonResult {
	t := services.OauthClientService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *OauthClientController) AnyList() *simple.JsonResult {
	list, paging := services.OauthClientService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *OauthClientController) PostCreate() *simple.JsonResult {
	t := &model.OauthClient{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	if len(t.ClientId) == 0 {
		return simple.ErrorMsg("clientId 不能为空")
	}
	if len(t.ClientSecret) == 0 {
		return simple.ErrorMsg("clientSecret 不能为空")
	}
	if len(t.Domain) == 0 {
		return simple.ErrorMsg("domain 不能为空")
	}
	if len(t.CallbackUrl) == 0 {
		return simple.ErrorMsg("callbackUrl 不能为空")
	}

	err = services.OauthClientService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *OauthClientController) PostUpdate() *simple.JsonResult {
	id := this.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return simple.ErrorMsg("id is required")
	}
	t := services.OauthClientService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}

	if len(t.ClientId) == 0 {
		return simple.ErrorMsg("clientId 不能为空")
	}
	if len(t.ClientSecret) == 0 {
		return simple.ErrorMsg("clientSecret 不能为空")
	}
	if len(t.Domain) == 0 {
		return simple.ErrorMsg("domain 不能为空")
	}
	if len(t.CallbackUrl) == 0 {
		return simple.ErrorMsg("callbackUrl 不能为空")
	}

	err = services.OauthClientService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}
