package render

import "bbs-go/model"

func BuildTag(tag *model.Tag) *model.TagResponse {
	if tag == nil {
		return nil
	}
	return &model.TagResponse{TagId: tag.Id, TagName: tag.Name}
}

func BuildTags(tags []model.Tag) *[]model.TagResponse {
	if len(tags) == 0 {
		return nil
	}
	var responses []model.TagResponse
	for _, tag := range tags {
		responses = append(responses, *BuildTag(&tag))
	}
	return &responses
}
