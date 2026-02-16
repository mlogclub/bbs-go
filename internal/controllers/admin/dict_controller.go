package admin

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

type DictController struct {
	Ctx iris.Context
}

func (c *DictController) GetBy(id int64) *web.JsonResult {
	t := services.DictService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(render.BuildDict(*t))
}

func (c *DictController) GetList() *web.JsonResult {
	typeId, _ := params.GetInt64(c.Ctx, "typeId")
	list := services.DictService.Find(sqls.NewCnd().Eq("type_id", typeId).Asc("sort_no").Desc("id"))
	return web.JsonData(render.BuildDictTree(0, list))
}

func (c *DictController) PostCreate() *web.JsonResult {
	t := &models.Dict{}
	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	t.SortNo = services.DictService.GetNextSortNo()
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()
	if err := services.DictService.Create(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *DictController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	t := services.DictService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	t.UpdateTime = dates.NowTimestamp()
	if err := services.DictService.Update(t); err != nil {
		return web.JsonErrorMsg(err.Error())
	}
	return web.JsonData(t)
}

func (c *DictController) PostDelete() *web.JsonResult {
	ids := params.GetInt64Arr(c.Ctx, "ids")
	if len(ids) == 0 {
		return web.JsonErrorMsg("delete ids is empty")
	}
	for _, id := range ids {
		services.DictService.Delete(id)
	}
	return web.JsonSuccess()
}

func (c *DictController) PostUpdate_sort() *web.JsonResult {
	var ids []int64
	if err := c.Ctx.ReadJSON(&ids); err != nil {
		return web.JsonError(err)
	}
	if err := services.DictService.UpdateSort(ids); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

func (c *DictController) GetDicts() *web.JsonResult {
	var (
		typeId, _ = params.GetInt64(c.Ctx, "typeId")
		code, _   = params.Get(c.Ctx, "code")
	)
	if typeId <= 0 {
		if strs.IsNotBlank(code) {
			dictType := services.DictTypeService.GetByCode(code)
			if dictType != nil {
				typeId = dictType.Id
			}
		}
	}
	list := services.DictService.FindByTypeId(typeId)
	return web.JsonData(render.BuildDictTree(0, list))
}
