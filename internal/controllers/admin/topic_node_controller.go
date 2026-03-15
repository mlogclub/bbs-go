package admin

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/locales"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

type TopicNodeController struct {
	Ctx iris.Context
}

func (c *TopicNodeController) GetBy(id int64) *web.JsonResult {
	t := services.TopicNodeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(t)
}

func (c *TopicNodeController) AnyList() *web.JsonResult {
	list := services.TopicNodeService.Find(params.NewSqlCnd(c.Ctx,
		params.QueryFilter{
			ParamName: "name",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "type",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "parentId",
			Op:        params.Eq,
		},
	).Asc("sort_no").Desc("id"))
	// 确保父节点在列表中，以便正确构建树
	list = ensureAncestors(list)
	return web.JsonData(render.BuildTopicNodeTree(0, list))
}

// ensureAncestors 为过滤结果补充父节点，保证树形结构完整
func ensureAncestors(list []models.TopicNode) []models.TopicNode {
	idSet := make(map[int64]bool)
	for _, n := range list {
		idSet[n.Id] = true
	}
	for {
		added := false
		for _, n := range list {
			if n.ParentId > 0 && !idSet[n.ParentId] {
				parent := services.TopicNodeService.Get(n.ParentId)
				if parent != nil {
					list = append(list, *parent)
					idSet[parent.Id] = true
					added = true
				}
			}
		}
		if !added {
			break
		}
	}
	return list
}

func (c *TopicNodeController) PostCreate() *web.JsonResult {
	t := &models.TopicNode{}
	if err := params.ReadForm(c.Ctx, t); err != nil {
		return web.JsonError(err)
	}
	t.SortNo = services.TopicNodeService.GetNextSortNo()
	if t.ParentId < 0 {
		t.ParentId = 0
	}
	// 子节点类型必须与父节点一致，直接取父节点的 type
	if t.ParentId > 0 {
		parent := services.TopicNodeService.Get(t.ParentId)
		if parent == nil {
			return web.JsonErrorMsg("parent node not found")
		}
		t.Type = parent.Type
	} else {
		if t.Type == "" {
			t.Type = constants.TopicNodeTypeNormal
		}
	}
	t.CreateTime = dates.NowTimestamp()
	if err := services.TopicNodeService.Create(t); err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *TopicNodeController) PostUpdate() *web.JsonResult {
	id, err := params.FormValueInt64(c.Ctx, "id")
	if err != nil {
		return web.JsonError(err)
	}
	t := services.TopicNodeService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("entity not found")
	}

	err = params.ReadForm(c.Ctx, t)
	if err != nil {
		return web.JsonError(err)
	}
	if strings.TrimSpace(t.Description) == "" {
		return web.JsonErrorMsg("param: description required")
	}
	if t.ParentId < 0 {
		t.ParentId = 0
	}
	// 禁止将父节点设为自己
	if t.ParentId == id {
		t.ParentId = 0
	}

	// 子节点类型必须与父节点一致，直接取父节点的 type，忽略表单传入
	if t.ParentId > 0 {
		parent := services.TopicNodeService.Get(t.ParentId)
		if parent == nil {
			return web.JsonErrorMsg("parent node not found")
		}
		if t.Type != parent.Type {
			return web.JsonErrorMsg(locales.Get("topic.node.child_type_must_match_parent"))
		}
		t.Type = parent.Type
	} else {
		// 一级节点：校验 type 必填，且编辑时联动更新所有子节点类型
		if strings.TrimSpace(string(t.Type)) == "" {
			return web.JsonErrorMsg("param: type required")
		}
		if err := services.TopicNodeService.UpdateChildrenType(id, t.Type); err != nil {
			return web.JsonError(err)
		}
	}

	err = services.TopicNodeService.Update(t)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(t)
}

func (c *TopicNodeController) GetNodes() *web.JsonResult {
	list := services.TopicNodeService.GetNodes()
	return web.JsonData(list)
}

func (c *TopicNodeController) PostUpdate_sort() *web.JsonResult {
	var ids []int64
	if err := c.Ctx.ReadJSON(&ids); err != nil {
		return web.JsonError(err)
	}
	if err := services.TopicNodeService.UpdateSort(ids); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// PostDelete 删除节点（一级有子节点时禁止）
func (c *TopicNodeController) PostDelete() *web.JsonResult {
	ids := params.GetInt64Arr(c.Ctx, "ids")
	if len(ids) == 0 {
		return web.JsonErrorMsg("delete ids is empty")
	}
	for _, id := range ids {
		if err := services.TopicNodeService.DeleteWithCheck(id); err != nil {
			return web.JsonErrorMsg(err.Error())
		}
	}
	return web.JsonSuccess()
}
