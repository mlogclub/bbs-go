package render

import (
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/services"
)

func BuildSearchArticles(docs []search.ArticleDocument) []resp.SearchArticleResponse {
	var items []resp.SearchArticleResponse
	for _, doc := range docs {
		items = append(items, BuildSearchArticle(doc))
	}
	return items
}

func BuildSearchArticle(doc search.ArticleDocument) resp.SearchArticleResponse {
	rsp := resp.SearchArticleResponse{
		Id:         doc.Id,
		Title:      doc.Title,
		Summary:    doc.Summary,
		CreateTime: doc.CreateTime,
		User:       BuildUserInfoDefaultIfNull(doc.UserId),
	}
	if rsp.Summary == "" {
		rsp.Summary = doc.Content
	}
	tags := services.ArticleService.GetArticleTags(doc.Id)
	rsp.Tags = BuildTags(tags)
	return rsp
}
