package admin

import (
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/locales"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/common/dates"

	"bbs-go/internal/models"
	"bbs-go/internal/services"
)

// ensureAncestors 为过滤结果补充父节点，保证树形结构完整
func ensureAncestors(list []models.Category) []models.Category {
	idSet := make(map[int64]bool)
	for _, n := range list {
		idSet[n.Id] = true
	}
	for {
		added := false
		for _, n := range list {
			if n.ParentId > 0 && !idSet[n.ParentId] {
				parent := services.CategoryService.Get(n.ParentId)
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

func filterCategoryListByNodeID(list []models.Category, nodeID int64) []models.Category {
	if nodeID <= 0 {
		return list
	}

	childrenByParentID := make(map[int64][]int64)
	byID := make(map[int64]models.Category)
	for _, node := range list {
		byID[node.Id] = node
		childrenByParentID[node.ParentId] = append(childrenByParentID[node.ParentId], node.Id)
	}
	if _, ok := byID[nodeID]; !ok {
		return nil
	}

	keep := make(map[int64]bool)
	for id := nodeID; id > 0; {
		node, ok := byID[id]
		if !ok {
			break
		}
		keep[id] = true
		id = node.ParentId
	}

	var collectDescendants func(id int64)
	collectDescendants = func(id int64) {
		for _, childID := range childrenByParentID[id] {
			if keep[childID] {
				continue
			}
			keep[childID] = true
			collectDescendants(childID)
		}
	}
	collectDescendants(nodeID)

	filtered := make([]models.Category, 0, len(keep))
	for _, node := range list {
		if keep[node.Id] {
			filtered = append(filtered, node)
		}
	}
	return filtered
}

// PostDelete 删除节点（一级有子节点时禁止）
func CategoryDetail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	t := services.CategoryService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("Not found, id="+strconv.FormatInt(id, 10)))
		return
	}
	ginx.WriteJSON(ctx, t)

}

func CategoryList(ctx *gin.Context) {
	nodeID, _ := params.GetInt64(ctx, "categoryId")
	if nodeID <= 0 {
		nodeID, _ = params.GetInt64(ctx, "parentId")
	}

	list := services.CategoryService.Find(params.NewSqlCnd(ctx,
		params.QueryFilter{
			ParamName: "name",
			Op:        params.Like,
		},
		params.QueryFilter{
			ParamName: "type",
			Op:        params.Eq,
		},
		params.QueryFilter{
			ParamName: "status",
			Op:        params.Eq,
		},
	).Asc("sort_no").Desc("id"))
	list = filterCategoryListByNodeID(list, nodeID)
	// 确保父节点在列表中，以便正确构建树
	list = ensureAncestors(list)
	ginx.WriteJSON(ctx, render.BuildCategoryTree(0, list))

}

func CategoryCreate(ctx *gin.Context) {
	t := &models.Category{}
	if err := ginx.Bind(ctx, t); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	t.SortNo = services.CategoryService.GetNextSortNo()
	if t.ParentId < 0 {
		t.ParentId = 0
	}
	// 子节点类型必须与父节点一致，直接取父节点的 type
	if t.ParentId > 0 {
		parent := services.CategoryService.Get(t.ParentId)
		if parent == nil {
			ginx.WriteJSON(ctx, ginx.ErrorMessage("parent node not found"))
			return
		}
		t.Type = parent.Type
	} else {
		if t.Type == "" {
			t.Type = constants.CategoryTypeNormal
		}
	}
	t.CreateTime = dates.NowTimestamp()
	if err := services.CategoryService.Create(t); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, t)

}

func CategoryUpdate(ctx *gin.Context) {
	id, err := params.FormValueInt64(ctx, "id")
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	t := services.CategoryService.Get(id)
	if t == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("entity not found"))
		return
	}

	err = ginx.Bind(ctx, t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	if strings.TrimSpace(t.Description) == "" {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("param: description required"))
		return
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
		parent := services.CategoryService.Get(t.ParentId)
		if parent == nil {
			ginx.WriteJSON(ctx, ginx.ErrorMessage("parent node not found"))
			return
		}
		if t.Type != parent.Type {
			ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("topic.category.child_type_must_match_parent")))
			return
		}
		t.Type = parent.Type
	} else {
		// 一级节点：校验 type 必填，且编辑时联动更新所有子节点类型
		if strings.TrimSpace(string(t.Type)) == "" {
			ginx.WriteJSON(ctx, ginx.ErrorMessage("param: type required"))
			return
		}
		if err := services.CategoryService.UpdateChildrenType(id, t.Type); err != nil {
			ginx.WriteJSON(ctx, err)
			return
		}
	}

	err = services.CategoryService.Update(t)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, t)

}

func CategoryOptions(ctx *gin.Context) {

	list := services.CategoryService.GetCategories()
	ginx.WriteJSON(ctx, list)

}

func CategoryUpdateSort(ctx *gin.Context) {
	var ids []int64
	if err := ginx.BindJSON(ctx, &ids); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	if err := services.CategoryService.UpdateSort(ids); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func CategoryRemove(ctx *gin.Context) {
	ids := params.GetInt64Arr(ctx, "ids")
	if len(ids) == 0 {
		ginx.WriteJSON(ctx, ginx.ErrorMessage("delete ids is empty"))
		return
	}
	for _, id := range ids {
		if err := services.CategoryService.DeleteWithCheck(id); err != nil {
			ginx.WriteJSON(ctx, ginx.ErrorMessage(err.Error()))
			return
		}
	}
	ginx.WriteJSON(ctx, nil)

}
