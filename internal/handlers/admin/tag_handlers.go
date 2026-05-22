package admin

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/locales"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/handlers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

// 自动完成
// 根据标签编号批量获取
func TagDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.TagService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.entity_not_found")+", id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func TagList(ctx *gin.Context) {
	list, paging := services.TagService.FindPageByParams(params.NewQueryParams(ctx).
		LikeByReq("id").
		LikeByReq("name").
		EqByReq("status").
		PageByReq().Desc("id"))
	ginx.WriteJSON(ctx, &web.PageResult{Results: list, Page: paging})

}

func TagCreate(ctx *gin.Context) {
	t := &models.Tag{}
	err := ginx.Bind(ctx, t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if len(t.Name) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.tag_name_required")))
		return
	}
	if services.TagService.GetByName(t.Name) != nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Getf("admin.tag_name_exists", t.Name)))
		return
	}

	t.Status = constants.StatusOk
	t.CreateTime = dates.NowTimestamp()
	t.UpdateTime = dates.NowTimestamp()

	err = services.TagService.Create(t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, t)

}

func TagUpdate(ctx *gin.Context) {
	id, err := params.FormValueInt64(ctx, "id")
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	t := services.TagService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.entity_not_found")))
		return
	}

	err = ginx.Bind(ctx, t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if len(t.Name) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("admin.tag_name_required")))
		return
	}
	if tmp := services.TagService.GetByName(t.Name); tmp != nil && tmp.Id != id {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Getf("admin.tag_name_exists", t.Name)))
		return
	}

	t.UpdateTime = dates.NowTimestamp()
	err = services.TagService.Update(t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, t)

}

func TagAutocomplete(ctx *gin.Context) {
	keyword := strings.TrimSpace(ctx.Query("keyword"))
	var tags []models.Tag
	if len(keyword) > 0 {
		tags = services.TagService.Find(sqls.NewCnd().Starting("name", keyword).Desc("id"))
	} else {
		tags = services.TagService.Find(sqls.NewCnd().Desc("id").Limit(10))
	}
	ginx.WriteJSON(ctx, render.BuildTags(tags))

}

func TagTags(ctx *gin.Context) {
	var tags *[]resp.TagResponse
	tagIds := params.FormValueInt64Array(ctx, "tagIds")
	if len(tagIds) > 0 {
		tagArr := services.TagService.Find(sqls.NewCnd().In("id", tagIds))
		if len(tagArr) > 0 {
			tags = render.BuildTags(tagArr)
		}
	}
	ginx.WriteJSON(ctx, tags)

}
