package admin

import (
	"bbs-go/model/constants"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/mvc"
	"github.com/mlogclub/simple/mvc/params"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
)

type TagController struct {
	Ctx iris.Context
}

func (c *TagController) GetBy(id int64) *mvc.JsonResult {
	t := services.TagService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return mvc.JsonData(t)
}

func (c *TagController) AnyList() *mvc.JsonResult {
	list, paging := services.TagService.FindPageByParams(params.NewQueryParams(c.Ctx).
		LikeByReq("id").
		LikeByReq("name").
		EqByReq("status").
		PageByReq().Desc("id"))
	return mvc.JsonData(&sqls.PageResult{Results: list, Page: paging})
}

func (c *TagController) PostCreate() *mvc.JsonResult {
	t := &model.Tag{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	if len(t.Name) == 0 {
		return mvc.JsonErrorMsg("name is required")
	}
	if services.TagService.GetByName(t.Name) != nil {
		return mvc.JsonErrorMsg("标签「" + t.Name + "」已存在")
	}

	t.Status = constants.StatusOk
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()

	err = services.TagService.Create(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}

func (c *TagController) PostUpdate() *mvc.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	t := services.TagService.Get(id)
	if t == nil {
		return mvc.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	if len(t.Name) == 0 {
		return mvc.JsonErrorMsg("name is required")
	}
	if tmp := services.TagService.GetByName(t.Name); tmp != nil && tmp.Id != id {
		return mvc.JsonErrorMsg("标签「" + t.Name + "」已存在")
	}

	t.UpdateTime = dates.NowTimestamp()
	err = services.TagService.Update(t)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonData(t)
}

// 自动完成
func (c *TagController) GetAutocomplete() *mvc.JsonResult {
	keyword := strings.TrimSpace(c.Ctx.URLParam("keyword"))
	var tags []model.Tag
	if len(keyword) > 0 {
		tags = services.TagService.Find(sqls.NewSqlCnd().Starting("name", keyword).Desc("id"))
	} else {
		tags = services.TagService.Find(sqls.NewSqlCnd().Desc("id").Limit(10))
	}
	return mvc.JsonData(render.BuildTags(tags))
}

// 根据标签编号批量获取
func (c *TagController) GetTags() *mvc.JsonResult {
	tagIds := params.FormValueInt64Array(c.Ctx, "tagIds")
	var tags *[]model.TagResponse
	if len(tagIds) > 0 {
		tagArr := services.TagService.Find(sqls.NewSqlCnd().In("id", tagIds))
		if len(tagArr) > 0 {
			tags = render.BuildTags(tagArr)
		}
	}
	return mvc.JsonData(tags)
}
