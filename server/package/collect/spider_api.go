package collect

import (
	"bbs-go/model/constants"
	"bbs-go/package/baiduai"
	"errors"

	"bbs-go/model"
	"bbs-go/package/common"
	"bbs-go/services"
)

type SpiderApi struct {
}

func NewSpiderApi() *SpiderApi {
	return &SpiderApi{}
}

func (api *SpiderApi) Publish(article *Article) (articleId int64, err error) {
	if article.Summary == "" {
		article.Summary = common.GetSummary(article.ContentType, article.Content)
	}

	if len(article.Tags) == 0 {
		article.Tags = api.AnalyzeTags(article)
	}

	t, err := services.ArticleService.Publish(article.UserId, article.Title, article.Summary, article.Content,
		article.ContentType, article.Tags, article.SourceUrl)
	if err == nil {
		articleId = t.Id

		if article.PublishTime > 0 {
			_ = services.ArticleService.UpdateColumn(articleId, "create_time", article.PublishTime)
		}
	}
	return
}

func (api *SpiderApi) PublishComment(comment *Comment) (commentId int64, err error) {
	if len(comment.Content) == 0 {
		err = errors.New("评论内容不能为空")
		return
	}

	c, err := services.CommentService.Publish(comment.UserId, &model.CreateCommentForm{
		EntityType:  comment.EntityType,
		EntityId:    comment.EntityId,
		Content:     comment.Content,
		ContentType: constants.ContentTypeHtml,
	})
	if err == nil {
		commentId = c.Id

		if comment.PublishTime > 0 {
			_ = services.CommentService.UpdateColumn(commentId, "create_time", comment.PublishTime)
		}
	}
	return
}

func (api *SpiderApi) AnalyzeTags(article *Article) []string {
	var analyzeRet *baiduai.AiAnalyzeRet
	if article.ContentType == constants.ContentTypeMarkdown {
		analyzeRet, _ = baiduai.GetAi().AnalyzeMarkdown(article.Title, article.Content)
	} else if article.ContentType == constants.ContentTypeHtml {
		analyzeRet, _ = baiduai.GetAi().AnalyzeHtml(article.Title, article.Content)
	}
	var tags []string
	if analyzeRet != nil {
		tags = analyzeRet.Tags
		if article.Summary == "" {
			article.Summary = analyzeRet.Summary
		}
	}
	return tags
}

type Article struct {
	UserId      int64    `json:"userId" form:"userId"` // 发布用户编号
	Title       string   `json:"title" form:"title"`
	Summary     string   `json:"summary" form:"summary"`
	Content     string   `json:"content" form:"content"`
	ContentType string   `json:"contentType" form:"contentType"`
	Tags        []string `json:"tags" form:"tags"`
	SourceUrl   string   `json:"sourceUrl" form:"sourceUrl"`
	PublishTime int64    `json:"publishTime" form:"publishTime"`
}

type Comment struct {
	UserId      int64  `json:"userId" form:"userId"`
	Content     string `json:"content" form:"content"`
	EntityType  string `json:"entityType" form:"entityType"`
	EntityId    int64  `json:"entityId" form:"entityId"`
	PublishTime int64  `json:"publishTime" form:"publishTime"`
}
