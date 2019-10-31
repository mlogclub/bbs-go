package services

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/mlogclub/bbs-go/common/baiduai"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"

	"github.com/mlogclub/simple"
)

var CollectArticleService = newCollectArticleService()

func newCollectArticleService() *collectArticleService {
	return &collectArticleService{}
}

type collectArticleService struct {
}

func (this *collectArticleService) Get(id int64) *model.CollectArticle {
	return repositories.CollectArticleRepository.Get(simple.GetDB(), id)
}

func (this *collectArticleService) Take(where ...interface{}) *model.CollectArticle {
	return repositories.CollectArticleRepository.Take(simple.GetDB(), where...)
}

func (this *collectArticleService) QueryCnd(cnd *simple.SqlCnd) (list []model.CollectArticle, err error) {
	return repositories.CollectArticleRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *collectArticleService) Query(params *simple.QueryParams) (list []model.CollectArticle, paging *simple.Paging) {
	return repositories.CollectArticleRepository.Query(simple.GetDB(), queries)
}

func (this *collectArticleService) Update(t *model.CollectArticle) error {
	return repositories.CollectArticleRepository.Update(simple.GetDB(), t)
}

func (this *collectArticleService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CollectArticleRepository.Updates(simple.GetDB(), id, columns)
}

func (this *collectArticleService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CollectArticleRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *collectArticleService) Delete(id int64) {
	repositories.CollectArticleRepository.Delete(simple.GetDB(), id)
}

// 文章是否存在
func (this *collectArticleService) IsExists(sourceUrl, title string) bool {
	if sourceUrl != "" {
		if tmp := repositories.CollectArticleRepository.Take(simple.GetDB(), "source_url_md5 = ?", simple.MD5(sourceUrl)); tmp != nil {
			return true
		}
	}
	if title != "" {
		if tmp := repositories.CollectArticleRepository.Take(simple.GetDB(), "source_title_md5 = ?", simple.MD5(title)); tmp != nil {
			return true
		}
	}
	return false
}

// 创建采集文章
func (this *collectArticleService) Create(ruleId, userId int64, sourceUrl, title, summary, content string) (*model.CollectArticle, error) {
	title = strings.TrimSpace(title)
	summary = strings.TrimSpace(summary)
	content = strings.TrimSpace(content)

	if len(title) == 0 || len(content) == 0 {
		return nil, errors.New("文章标题或内容为空")
	}

	collectArticle := &model.CollectArticle{
		UserId:         userId,
		RuleId:         ruleId,
		Title:          title,
		Summary:        summary,
		Content:        content,
		Status:         model.CollectArticleStatusPending,
		ContentType:    model.ContentTypeHtml,
		SourceUrl:      sourceUrl,
		SourceUrlMd5:   simple.MD5(sourceUrl),
		SourceTitleMd5: simple.MD5(title),
		CreateTime:     simple.NowTimestamp(),
	}

	err := simple.Tx(simple.GetDB(), func(tx *gorm.DB) (err error) {
		err = repositories.CollectArticleRepository.Create(tx, collectArticle)
		return
	})
	if err != nil {
		return nil, err
	}
	return collectArticle, err
}

// 发布采集文章
func (this *collectArticleService) Publish(collectArticleId int64) error {
	return simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		ca := repositories.CollectArticleRepository.Get(tx, collectArticleId)
		if ca == nil {
			return errors.New("没找到该采集文章")
		}

		var analyzeRet *baiduai.AiAnalyzeRet
		if ca.ContentType == model.ContentTypeHtml {
			analyzeRet, _ = baiduai.AnalyzeHtml(ca.Title, ca.Content)
		} else if ca.ContentType == model.ContentTypeMarkdown {
			analyzeRet, _ = baiduai.AnalyzeMarkdown(ca.Title, ca.Content)
		}

		var tags []string
		var summary string
		if analyzeRet != nil {
			tags = analyzeRet.Tags
			summary = analyzeRet.Summary
		}

		article, err := ArticleService.Publish(ca.UserId, ca.Title, summary, ca.Content, ca.ContentType,
			0, tags, ca.SourceUrl, false)
		if err != nil {
			return err
		}
		return repositories.CollectArticleRepository.Updates(tx, collectArticleId, map[string]interface{}{
			"status":     model.CollectArticleStatusPublished,
			"article_id": article.Id,
		})
	})
}

type CollectArticleScanCallback func(*model.CollectArticle)

// 扫描
func (this *collectArticleService) Scan(callback CollectArticleScanCallback) {
	var cursor int64
	for {
		list, err := repositories.CollectArticleRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ?", cursor).Order("id asc").Size(100))
		if err != nil {
			break
		}
		if len(list) == 0 {
			break
		}
		for _, item := range list {
			cursor = item.Id
			callback(&item)
		}
	}
}
