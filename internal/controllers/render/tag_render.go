package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/resp"
)

func BuildTag(tag *models.Tag) *resp.TagResponse {
	if tag == nil {
		return nil
	}
	return &resp.TagResponse{Id: tag.Id, Name: tag.Name}
}

func BuildTags(tags []models.Tag) *[]resp.TagResponse {
	if len(tags) == 0 {
		return nil
	}
	var responses []resp.TagResponse
	for _, tag := range tags {
		responses = append(responses, *BuildTag(&tag))
	}
	return &responses
}
