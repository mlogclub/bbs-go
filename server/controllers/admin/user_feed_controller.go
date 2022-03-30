package admin

import (
	"bbs-go/model"
	"bbs-go/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"
)

type UserFeedController struct {
	Ctx iris.Context
}

func (c *UserFeedController) GetBy(id int64) *mvc.JsonResult {
	t := services.UserFeedService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *UserFeedController) AnyList() *mvc.JsonResult {
	list, paging := services.UserFeedService.FindPageByParams(params.NewQueryParams(c.Ctx).PageByReq().Desc("id"))
	return mvc.JsonData(&sqls.PageResult{Results: list, Page: paging})
}

func (c *UserFeedController) PostCreate() *mvc.JsonResult {
	t := &model.UserFeed{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	err = services.UserFeedService.Create(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}

func (c *UserFeedController) PostUpdate() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	t := services.UserFeedService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	err = services.UserFeedService.Update(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}
