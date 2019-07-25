package admin

import (
	"strconv"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/services"
)

type UserTokenController struct {
	Ctx iris.Context
}

func (this *UserTokenController) GetBy(id int64) *simple.JsonResult {
	t := services.UserTokenService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *UserTokenController) AnyList() *simple.JsonResult {
	list, paging := services.UserTokenService.Query(simple.NewParamQueries(this.Ctx).PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}
