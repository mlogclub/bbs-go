package render

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/html"
	"bbs-go/internal/pkg/markdown"
	"bbs-go/internal/pkg/text"
	"bbs-go/internal/services"
)

func BuildArticle(article *models.Article, currentUser *models.User) *models.ArticleResponse {
	if article == nil {
		return nil
	}

	rsp := &models.ArticleResponse{}
	rsp.Id = article.Id
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.SourceUrl = article.SourceUrl
	rsp.ViewCount = article.ViewCount
	rsp.CreateTime = article.CreateTime
	rsp.Status = article.Status

	rsp.User = BuildUserInfoDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.Id)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = BuildTags(tags)

	if article.ContentType == constants.ContentTypeMarkdown {
		content := markdown.ToHTML(article.Content)
		rsp.Content = handleHtmlContent(content)
	} else if article.ContentType == constants.ContentTypeHtml {
		rsp.Content = handleHtmlContent(article.Content)
	}

	rsp.Cover = BuildImage(article.Cover)

	if currentUser != nil {
		rsp.Favorited = services.FavoriteService.IsFavorited(currentUser.Id, constants.EntityArticle, article.Id)
	}

	return rsp
}

func BuildSimpleArticle(article *models.Article) *models.ArticleSimpleResponse {
	if article == nil {
		return nil
	}

	rsp := &models.ArticleSimpleResponse{}
	rsp.Id = article.Id
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.SourceUrl = article.SourceUrl
	rsp.ViewCount = article.ViewCount
	rsp.CommentCount = article.CommentCount
	rsp.LikeCount = article.LikeCount
	rsp.CreateTime = article.CreateTime
	rsp.Status = article.Status

	rsp.User = BuildUserInfoDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.Id)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = BuildTags(tags)

	if article.ContentType == constants.ContentTypeMarkdown {
		if len(rsp.Summary) == 0 {
			rsp.Summary = markdown.GetSummary(article.Content, constants.SummaryLen)
		}
	} else if article.ContentType == constants.ContentTypeHtml {
		if len(rsp.Summary) == 0 {
			rsp.Summary = text.GetSummary(html.GetHtmlText(article.Content), constants.SummaryLen)
		}
	}

	rsp.Cover = BuildImage(article.Cover)

	return rsp
}

func BuildSimpleArticles(articles []models.Article) []models.ArticleSimpleResponse {
	if len(articles) == 0 {
		return nil
	}
	var responses []models.ArticleSimpleResponse
	for _, article := range articles {
		responses = append(responses, *BuildSimpleArticle(&article))
	}
	return responses
}
