package collect

import (
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
	var analyzeRet *baiduai.AiAnalyzeRet
	if article.ContentType == model.ContentTypeMarkdown {
		analyzeRet, _ = baiduai.AnalyzeMarkdown(article.Title, article.Content)
	} else {
		analyzeRet, _ = baiduai.AnalyzeHtml(article.Title, article.Content)
	}

	var tags []string
	if analyzeRet != nil {
		tags = analyzeRet.Tags
		if article.Summary == "" {
			article.Summary = analyzeRet.Summary
		}
	}
	if article.Summary == "" {
		article.Summary = common.GetSummary(article.ContentType, article.Content)
	}

	t, err := services.ArticleService.Publish(article.UserId, article.Title, article.Summary, article.Content,
		article.ContentType, 0, tags, article.SourceUrl, true)
	if err == nil {
		articleId = t.Id
	}
	return
}

type SpiderArticle struct {
	UserId      int64  `json:"userId" form:"userId"` // 发布用户编号
	Title       string `json:"title" form:"title"`
	Summary     string `json:"summary" form:"summary"`
	Content     string `json:"content" form:"content"`
	ContentType string `json:"contentType" form:"contentType"`
	SourceUrl   string `json:"sourceUrl" form:"sourceUrl"`
	// SourceUserId     string `json:"sourceUserId" form:"sourceUserId"`
	// SourceUserName   string `json:"sourceUserName" form:"sourceUserName"`
	// SourceUserAvatar string `json:"sourceUserAvatar" form:"sourceUserAvatar"`
}
