package render

import (
	"server/model"
	"server/pkg/es"
	"server/services"
)

func BuildSearchArticles(docs []es.ArticleDocument) []model.SearchArticleResponse {
	var items []model.SearchArticleResponse
	for _, doc := range docs {
		items = append(items, BuildSearchArticle(doc))
	}
	return items
}

func BuildSearchArticle(doc es.ArticleDocument) model.SearchArticleResponse {
	rsp := model.SearchArticleResponse{
		ArticleId:  doc.Id,
		Tags:       nil,
		Title:      doc.Title,
		Summary:    doc.Content,
		CreateTime: doc.CreateTime,
		User:       BuildUserInfoDefaultIfNull(doc.UserId),
	}

	tags := services.TopicService.GetTopicTags(doc.Id)
	rsp.Tags = BuildTags(tags)
	return rsp
}
