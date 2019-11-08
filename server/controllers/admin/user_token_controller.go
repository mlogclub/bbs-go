package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/services"
)

type UserTokenController struct {
	Ctx iris.Context
}

func (this *UserTokenController) GetBy(id int64) *simple.JsonResult {
	t := services.UserTokenService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *UserTokenController) AnyList() *simple.JsonResult {
	list, paging := services.UserTokenService.FindPageByParams(simple.NewQueryParams(this.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}
