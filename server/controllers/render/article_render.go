package render

import (
	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/html"
	"bbs-go/pkg/markdown"
	"bbs-go/pkg/text"
	"bbs-go/services"
)

func BuildArticle(article *model.Article, currentUser *model.User) *model.ArticleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleResponse{}
	rsp.ArticleId = article.Id
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

	rsp.Cover = buildImage(article.Cover)

	if currentUser != nil {
		rsp.Favorited = services.FavoriteService.IsFavorited(currentUser.Id, constants.EntityArticle, article.Id)
	}

	return rsp
}

func BuildSimpleArticle(article *model.Article) *model.ArticleSimpleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleSimpleResponse{}
	rsp.ArticleId = article.Id
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

	rsp.Cover = buildImage(article.Cover)

	return rsp
}

func BuildSimpleArticles(articles []model.Article) []model.ArticleSimpleResponse {
	if len(articles) == 0 {
		return nil
	}
	var responses []model.ArticleSimpleResponse
	for _, article := range articles {
		responses = append(responses, *BuildSimpleArticle(&article))
	}
	return responses
}
