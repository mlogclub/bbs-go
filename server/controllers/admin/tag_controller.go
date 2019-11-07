package admin

import (
	"strconv"
	"strings"

	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
)

type TagController struct {
	Ctx iris.Context
}

func (this *TagController) GetBy(id int64) *simple.JsonResult {
	t := services.TagService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *TagController) AnyList() *simple.JsonResult {
	list, paging := services.TagService.FindPageByParams(simple.NewQueryParams(this.Ctx).
		LikeByReq("name").
		EqByReq("status").
		PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *TagController) PostCreate() *simple.JsonResult {
	t := &model.Tag{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	if len(t.Name) == 0 {
		return simple.JsonErrorMsg("name is required")
	}
	if services.TagService.GetByName(t.Name) != nil {
		return simple.JsonErrorMsg("标签「" + t.Name + "」已存在")
	}

	t.Status = model.TagStatusOk
	t.CreateTime = simple.NowTimestamp()
	t.UpdateTime = simple.NowTimestamp()

	err = services.TagService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *TagController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.TagService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	if len(t.Name) == 0 {
		return simple.JsonErrorMsg("name is required")
	}
	if tmp := services.TagService.GetByName(t.Name); tmp != nil && tmp.Id != id {
		return simple.JsonErrorMsg("标签「" + t.Name + "」已存在")
	}

	t.UpdateTime = simple.NowTimestamp()
	err = services.TagService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

// 自动完成
func (this *TagController) GetAutocomplete() *simple.JsonResult {
	keyword := strings.TrimSpace(this.Ctx.URLParam("keyword"))
	var tags []model.Tag
	if len(keyword) > 0 {
		tags = services.TagService.Find(simple.NewSqlCnd().Like("name", keyword).Desc("id").Limit(10))
	} else {
		tags = services.TagService.Find(simple.NewSqlCnd().Desc("id").Limit(10))
	}
	return simple.JsonData(render.BuildTags(tags))
}
