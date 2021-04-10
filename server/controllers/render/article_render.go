package render

import (
	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/package/markdown"
	"github.com/mlogclub/simple"
)

func BuildArticle(article *model.Article) *model.ArticleResponse {
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

	rsp.User = BuildUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.Id)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = BuildTags(tags)

	if article.ContentType == constants.ContentTypeMarkdown {
		content := markdown.ToHTML(article.Content)
		rsp.Content = handleHtmlContent(content)
	} else if article.ContentType == constants.ContentTypeHtml {
		rsp.Content = handleHtmlContent(article.Content)
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
	rsp.CreateTime = article.CreateTime
	rsp.Status = article.Status

	rsp.User = BuildUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.Id)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = BuildTags(tags)

	if article.ContentType == constants.ContentTypeMarkdown {
		if len(rsp.Summary) == 0 {
			rsp.Summary = markdown.GetSummary(article.Content, constants.SummaryLen)
		}
	} else if article.ContentType == constants.ContentTypeHtml {
		if len(rsp.Summary) == 0 {
			rsp.Summary = simple.GetSummary(simple.GetHtmlText(article.Content), constants.SummaryLen)
		}
	}

	return rsp
}

func BuildSimpleArticles(articles []model.Article) []model.ArticleSimpleResponse {
	if articles == nil || len(articles) == 0 {
		return nil
	}
	var responses []model.ArticleSimpleResponse
	for _, article := range articles {
		responses = append(responses, *BuildSimpleArticle(&article))
	}
	return responses
}
