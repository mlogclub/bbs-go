package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/resp"
)

func BuildDict(element models.Dict) resp.DictResponse {
	item := resp.DictResponse{
		Id:         element.Id,
		TypeId:     element.TypeId,
		Name:       element.Name,
		Label:      element.Label,
		Value:      element.Value,
		SortNo:     element.SortNo,
		Status:     element.Status,
		CreateTime: element.CreateTime,
		UpdateTime: element.UpdateTime,
	}
	if element.ParentId > 0 {
		item.ParentId = &element.ParentId
	}
	return item
}

func BuildDictTree(parentId int64, list []models.Dict) (ret []resp.DictListResponse) {
	for _, element := range list {
		if element.ParentId == parentId {
			item := resp.DictListResponse{
				DictResponse: BuildDict(element),
				Children:     BuildDictTree(element.Id, list),
			}
			ret = append(ret, item)
		}
	}
	if len(ret) == 0 {
		return
	}
	return
}
