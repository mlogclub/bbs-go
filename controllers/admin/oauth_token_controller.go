package admin

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/simple"
	"strconv"
)

type OauthTokenController struct {
	Ctx               iris.Context
	OauthTokenService *services.OauthTokenService
}

func (this *OauthTokenController) GetBy(id int64) *simple.JsonResult {
	t := this.OauthTokenService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *OauthTokenController) AnyList() *simple.JsonResult {
	list, paging := this.OauthTokenService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}
