package render

import "bbs-go/internal/models"

func BuildMenuTree(parentId int64, list []models.Menu) (ret []models.MenuResponse) {
	return _BuildMenuTree(parentId, 1, list)
}

func _BuildMenuTree(parentId int64, level int, list []models.Menu) (ret []models.MenuResponse) {
	for _, element := range list {
		if element.ParentId == parentId {
			ret = append(ret, models.MenuResponse{
				Id:         element.Id,
				ParentId:   element.ParentId,
				Name:       element.Name,
				Title:      element.Title,
				Icon:       element.Icon,
				Path:       element.Path,
				SortNo:     element.SortNo,
				Status:     element.Status,
				CreateTime: element.CreateTime,
				UpdateTime: element.UpdateTime,
				Level:      level,
				Children:   _BuildMenuTree(element.Id, level+1, list),
			})
		}
	}
	return
}
