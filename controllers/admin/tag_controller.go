package admin

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/simple"
	"strconv"
)

type TagController struct {
	Ctx             iris.Context
	TagService      *services.TagService
	CategoryService *services.CategoryService
}

func (this *TagController) GetBy(id int64) *simple.JsonResult {
	t := this.TagService.Get(id)
	if t == nil {
		return simple.ErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *TagController) AnyList() *simple.JsonResult {
	list, paging := this.TagService.Query(simple.NewParamQueries(this.Ctx).
		EqAuto("category_id").
		EqAuto("status").
		PageAuto().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *TagController) PostCreate() *simple.JsonResult {
	t := &model.Tag{}
	this.Ctx.ReadForm(t)

	if len(t.Name) == 0 {
		return simple.ErrorMsg("name is required")
	}
	if this.TagService.GetByName(t.Name) != nil {
		return simple.ErrorMsg("标签「" + t.Name + "」已存在")
	}
	if t.CategoryId <= 0 {
		return simple.ErrorMsg("请选择标签分类")
	}

	t.Status = model.TagStatusOk
	t.CreateTime = simple.NowTimestamp()
	t.UpdateTime = simple.NowTimestamp()

	err := this.TagService.Create(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *TagController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	t := this.TagService.Get(id)
	if t == nil {
		return simple.ErrorMsg("entity not found")
	}

	this.Ctx.ReadForm(t)

	if len(t.Name) == 0 {
		return simple.ErrorMsg("name is required")
	}
	if tmp := this.TagService.GetByName(t.Name); tmp != nil && tmp.Id != id {
		return simple.ErrorMsg("标签「" + t.Name + "」已存在")
	}
	if t.CategoryId <= 0 {
		return simple.ErrorMsg("请选择标签分类")
	}

	t.UpdateTime = simple.NowTimestamp()
	err = this.TagService.Update(t)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *TagController) AnyListAll() *simple.JsonResult {
	categoryId, err := strconv.ParseInt(this.Ctx.FormValue("categoryId"), 10, 64)
	if err != nil {
		return simple.ErrorMsg(err.Error())
	}
	if categoryId < 0 {
		return simple.ErrorMsg("请指定categoryId")
	}
	list, err := this.TagService.ListAll(categoryId)
	if err != nil {
		return simple.JsonData([]interface{}{})
	}
	return simple.JsonData(list)
}

// 标签数据级联选择器
func (this *TagController) GetCascader() *simple.JsonResult {
	categories, err := this.CategoryService.GetCategories()
	if err != nil {
		return simple.ErrorMsg("数据加载失败")
	}

	var results []map[string]interface{}

	for _, cat := range categories {
		tags, err := this.TagService.ListAll(cat.Id)
		if err != nil || len(tags) == 0 {
			continue
		}

		var tagOptions []map[string]interface{}
		for _, tag := range tags {
			tagOption := make(map[string]interface{})
			tagOption["value"] = tag.Id
			tagOption["label"] = tag.Name
			tagOptions = append(tagOptions, tagOption)
		}

		option := make(map[string]interface{})
		option["value"] = cat.Id
		option["label"] = cat.Name
		option["children"] = tagOptions

		results = append(results, option)
	}

	return simple.JsonData(results)

}
