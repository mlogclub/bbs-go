package render

import "bbs-go/internal/models"

func BuildTag(tag *models.Tag) *models.TagResponse {
	if tag == nil {
		return nil
	}
	return &models.TagResponse{Id: tag.Id, Name: tag.Name}
}

func BuildTags(tags []models.Tag) *[]models.TagResponse {
	if len(tags) == 0 {
		return nil
	}
	var responses []models.TagResponse
	for _, tag := range tags {
		responses = append(responses, *BuildTag(&tag))
	}
	return &responses
}
