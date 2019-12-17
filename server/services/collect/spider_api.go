package collect

import (
	"errors"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/common/baiduai"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
)

type SpiderApi struct {
}

func NewSpiderApi() *SpiderApi {
	return &SpiderApi{}
}

func (this *SpiderApi) Publish(article *SpiderArticle) (articleId int64, err error) {
	if article.Summary == "" {
		article.Summary = common.GetSummary(article.ContentType, article.Content)
	}

	if len(article.Tags) == 0 {
		article.Tags = this.AnalyzeTags(article)
	}

	t, err := services.ArticleService.Publish(article.UserId, article.Title, article.Summary, article.Content,
		article.ContentType, article.Tags, article.SourceUrl, true)
	if err == nil {
		articleId = t.Id

		if article.PublishTime > 0 {
			services.ArticleService.UpdateColumn(articleId, "create_time", article.PublishTime)
		}
	}
	return
}

func (this *SpiderApi) PublishComment(comment *SpiderComment) (commentId int64, err error) {
	if len(comment.Content) == 0 {
		err = errors.New("评论内容不能为空")
		return
	}

	c, err := services.CommentService.Publish(comment.UserId, &model.CreateCommentForm{
		EntityType: comment.EntityType,
		EntityId:   comment.EntityId,
		Content:    comment.Content,
	})
	if err == nil {
		commentId = c.Id

		if comment.PublishTime > 0 {
			services.CommentService.UpdateColumn(commentId, "create_time", comment.PublishTime)
		}
	}
	return
}

func (this *SpiderApi) AnalyzeTags(article *SpiderArticle) []string {
	var analyzeRet *baiduai.AiAnalyzeRet
	if article.ContentType == model.ContentTypeMarkdown {
		analyzeRet, _ = baiduai.GetAi().AnalyzeMarkdown(article.Title, article.Content)
	} else {
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

type SpiderArticle struct {
	UserId      int64    `json:"userId" form:"userId"` // 发布用户编号
	Title       string   `json:"title" form:"title"`
	Summary     string   `json:"summary" form:"summary"`
	Content     string   `json:"content" form:"content"`
	ContentType string   `json:"contentType" form:"contentType"`
	Tags        []string `json:"tags" form:"tags"`
	SourceUrl   string   `json:"sourceUrl" form:"sourceUrl"`
	PublishTime int64    `json:"publishTime" form:"publishTime"`
}

type SpiderComment struct {
	UserId      int64  `json:"userId" form:"userId"`
	Content     string `json:"content" form:"content"`
	EntityType  string `json:"entityType" form:"entityType"`
	EntityId    int64  `json:"entityId" form:"entityId"`
	PublishTime int64  `json:"publishTime" form:"publishTime"`
}
