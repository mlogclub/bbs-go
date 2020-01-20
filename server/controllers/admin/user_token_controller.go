package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/services"
)

type UserTokenController struct {
	Ctx iris.Context
}

func (c *UserTokenController) GetBy(id int64) *simple.JsonResult {
	t := services.UserTokenService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *UserTokenController) AnyList() *simple.JsonResult {
	list, paging := services.UserTokenService.FindPageByParams(simple.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}
