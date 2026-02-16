package admin

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type VoteRecordController struct {
	Ctx iris.Context
}

func (c *VoteRecordController) GetBy(id int64) *web.JsonResult {
	t := services.VoteRecordService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *VoteRecordController) AnyList() *web.JsonResult {
	list, paging := services.VoteRecordService.FindPageByCnd(params.NewPagedSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "id",
		},
	).Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *VoteRecordController) PostCreate() *web.JsonResult {
	t := &models.VoteRecord{}
	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	if err := services.VoteRecordService.Create(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *VoteRecordController) PostUpdate() *web.JsonResult {
	id, _ := params.GetInt64(c.Ctx, "id")
	t := services.VoteRecordService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	if err := services.VoteRecordService.Update(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *VoteRecordController) PostDelete() *web.JsonResult {
	ids := params.GetInt64Arr(c.Ctx, "ids")
	if len(ids) == 0 {
		return web.JsonErrorMsg("delete ids is empty")
	}
	for _, id := range ids {
		services.VoteRecordService.Delete(id)
	}
	return web.JsonSuccess()
}

