package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
)

type CategoryController struct {
	Ctx iris.Context
}

func (this *CategoryController) GetBy(id int64) *simple.JsonResult {
	t := services.CategoryService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *CategoryController) AnyList() *simple.JsonResult {
	list, paging := services.CategoryService.FindPageByParams(simple.NewQueryParams(this.Ctx).
		LikeByReq("name").
		EqByReq("status").
		PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *CategoryController) PostCreate() *simple.JsonResult {
	t := &model.Category{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	if len(t.Name) == 0 {
		return simple.JsonErrorMsg("name is required")
	}

	if services.CategoryService.FindByName(t.Name) != nil {
		return simple.JsonErrorMsg("分类「" + t.Name + "」已存在")
	}

	t.Status = model.CategoryStatusOk
	t.CreateTime = simple.NowTimestamp()
	t.UpdateTime = simple.NowTimestamp()

	err = services.CategoryService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *CategoryController) PostUpdate() *simple.JsonResult {
	id := this.Ctx.PostValueInt64Default("id", 0)
	if id <= 0 {
		return simple.JsonErrorMsg("id is required")
	}
	t := services.CategoryService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	if len(t.Name) == 0 {
		return simple.JsonErrorMsg("name is required")
	}

	if tmp := services.CategoryService.FindByName(t.Name); tmp != nil && tmp.Id != id {
		return simple.JsonErrorMsg("分类「" + t.Name + "」已存在")
	}

	t.UpdateTime = simple.NowTimestamp()

	err = services.CategoryService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

// options选项
func (this *CategoryController) AnyOptions() *simple.JsonResult {
	categories := services.CategoryService.GetCategories()

	var results []map[string]interface{}
	for _, cat := range categories {
		option := make(map[string]interface{})
		option["value"] = cat.Id
		option["label"] = cat.Name

		results = append(results, option)
	}
	return simple.JsonData(results)
}
