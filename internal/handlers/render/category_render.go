package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/services"

	"github.com/mlogclub/simple/common/strs"
)

func BuildCategory(category *models.Category) *resp.CategoryResponse {
	if category == nil {
		return nil
	}
	if strs.IsBlank(category.Logo) {
		category.Logo = "/res/images/category_default.svg"
	}
	return &resp.CategoryResponse{
		Id:          category.Id,
		ParentId:    category.ParentId,
		Name:        category.Name,
		Type:        category.Type,
		Logo:        category.Logo,
		Description: category.Description,
	}
}

func BuildCategoryWithChildren(category *models.Category) *resp.CategoryResponse {
	r := BuildCategory(category)
	if r == nil {
		return nil
	}
	if category.ParentId == 0 {
		children := services.CategoryService.GetChildren(category.Id)
		if len(children) > 0 {
			r.Children = BuildCategoryResponses(children)
		}
	}
	return r
}

func BuildCategoryResponses(categories []models.Category) []resp.CategoryResponse {
	if len(categories) == 0 {
		return nil
	}
	var ret []resp.CategoryResponse
	for _, category := range categories {
		ret = append(ret, *BuildCategory(&category))
	}
	return ret
}

func BuildCategoryResponseTree(parentId int64, list []models.Category) []resp.CategoryResponse {
	var ret []resp.CategoryResponse
	for _, category := range list {
		if category.ParentId != parentId {
			continue
		}
		item := BuildCategory(&category)
		if item == nil {
			continue
		}
		children := BuildCategoryResponseTree(category.Id, list)
		if len(children) > 0 {
			item.Children = children
		}
		ret = append(ret, *item)
	}
	return ret
}

func BuildCategoryTree(parentId int64, list []models.Category) []resp.CategoryTreeItem {
	var ret []resp.CategoryTreeItem
	for _, category := range list {
		if category.ParentId == parentId {
			children := BuildCategoryTree(category.Id, list)
			logo := category.Logo
			if strs.IsBlank(logo) {
				logo = "/res/images/category_default.svg"
			}
			ret = append(ret, resp.CategoryTreeItem{
				Id:          category.Id,
				ParentId:    category.ParentId,
				Name:        category.Name,
				Type:        category.Type,
				Logo:        logo,
				Description: category.Description,
				SortNo:      category.SortNo,
				Status:      category.Status,
				CreateTime:  category.CreateTime,
				Children:    children,
			})
		}
	}
	return ret
}
