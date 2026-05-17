package admin

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/locales"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

type TagController struct {
	Ctx iris.Context
}

func (c *TagController) GetBy(id int64) *web.JsonResult {
	t := services.TagService.Get(id)
	if t == nil {
		return web.JsonErrorMsg(locales.Get("admin.entity_not_found") + ", id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *TagController) AnyList() *web.JsonResult {
	list, paging := services.TagService.FindPageByParams(params.NewQueryParams(c.Ctx).
		LikeByReq("id").
		LikeByReq("name").
		EqByReq("status").
		PageByReq().Desc("id"))
	return web.JsonData(&web.PageResult{Results: list, Page: paging})
}

func (c *TagController) PostCreate() *web.JsonResult {
	t := &models.Tag{}
	err := params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	if len(t.Name) == 0 {
		return web.JsonErrorMsg(locales.Get("admin.tag_name_required"))
	}
	if services.TagService.GetByName(t.Name) != nil {
		return web.JsonErrorMsg(locales.Getf("admin.tag_name_exists", t.Name))
	}

	t.Status = constants.StatusOk
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()

	err = services.TagService.Create(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *TagController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := services.TagService.Get(id)
	if t == nil {
		return web.JsonErrorMsg(locales.Get("admin.entity_not_found"))
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}

	if len(t.Name) == 0 {
		return web.JsonErrorMsg(locales.Get("admin.tag_name_required"))
	}
	if tmp := services.TagService.GetByName(t.Name); tmp != nil && tmp.Id != id {
		return web.JsonErrorMsg(locales.Getf("admin.tag_name_exists", t.Name))
	}

	t.UpdateTime = dates.NowTimestamp()
	err = services.TagService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

// 自动完成
func (c *TagController) GetAutocomplete() *web.JsonResult {
	keyword := strings.TrimSpace(c.Ctx.URLParam("keyword"))
	var tags []models.Tag
	if len(keyword) > 0 {
		tags = services.TagService.Find(sqls.NewCnd().Starting("name", keyword).Desc("id"))
	} else {
		tags = services.TagService.Find(sqls.NewCnd().Desc("id").Limit(10))
	}
	return web.JsonData(render.BuildTags(tags))
}

// 根据标签编号批量获取
func (c *TagController) GetTags() *web.JsonResult {
	tagIds := params.FormValueInt64Array(c.Ctx, "tagIds")
	var tags *[]resp.TagResponse
	if len(tagIds) > 0 {
		tagArr := services.TagService.Find(sqls.NewCnd().In("id", tagIds))
		if len(tagArr) > 0 {
			tags = render.BuildTags(tagArr)
		}
	}
	return web.JsonData(tags)
}
