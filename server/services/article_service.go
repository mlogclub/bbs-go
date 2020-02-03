package services

import (
	"errors"
	"math"
	"path"
	"strings"
	"time"

	"github.com/emirpasic/gods/sets/hashset"

	"bbs-go/common/baiduseo"
	"bbs-go/common/config"
	"bbs-go/common/urls"
	"bbs-go/repositories"
	"bbs-go/services/cache"

	"github.com/gorilla/feeds"
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/common"
	"bbs-go/model"
)

type ScanArticleCallback func(articles []model.Article)

var ArticleService = newArticleService()

func newArticleService() *articleService {
	return &articleService{}
}

type articleService struct {
}

func (s *articleService) Get(id int64) *model.Article {
	return repositories.ArticleRepository.Get(simple.DB(), id)
}

func (s *articleService) Take(where ...interface{}) *model.Article {
	return repositories.ArticleRepository.Take(simple.DB(), where...)
}

func (s *articleService) Find(cnd *simple.SqlCnd) []model.Article {
	return repositories.ArticleRepository.Find(simple.DB(), cnd)
}

func (s *articleService) FindOne(cnd *simple.SqlCnd) *model.Article {
	return repositories.ArticleRepository.FindOne(simple.DB(), cnd)
}

func (s *articleService) FindPageByParams(params *simple.QueryParams) (list []model.Article, paging *simple.Paging) {
	return repositories.ArticleRepository.FindPageByParams(simple.DB(), params)
}

func (s *articleService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Article, paging *simple.Paging) {
	return repositories.ArticleRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *articleService) Update(t *model.Article) error {
	err := repositories.ArticleRepository.Update(simple.DB(), t)
	return err
}

func (s *articleService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.ArticleRepository.Updates(simple.DB(), id, columns)
	return err
}

func (s *articleService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.ArticleRepository.UpdateColumn(simple.DB(), id, name, value)
	return err
}

func (s *articleService) Delete(id int64) error {
	err := repositories.ArticleRepository.UpdateColumn(simple.DB(), id, "status", model.StatusDeleted)
	if err == nil {
		// 删掉标签文章
		ArticleTagService.DeleteByArticleId(id)
	}
	return err
}

// 根据文章编号批量获取文章
func (s *articleService) GetArticleInIds(articleIds []int64) []model.Article {
	if len(articleIds) == 0 {
		return nil
	}
	var articles []model.Article
	simple.DB().Where("id in (?)", articleIds).Find(&articles)
	return articles
}

// 获取文章对应的标签
func (s *articleService) GetArticleTags(articleId int64) []model.Tag {
	articleTags := repositories.ArticleTagRepository.Find(simple.DB(), simple.NewSqlCnd().Where("article_id = ?", articleId))
	var tagIds []int64
	for _, articleTag := range articleTags {
		tagIds = append(tagIds, articleTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// 文章列表
func (s *articleService) GetArticles(cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := simple.NewSqlCnd().Eq("status", model.StatusOk).Desc("id").Limit(20)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	articles = repositories.ArticleRepository.Find(simple.DB(), cnd)
	if len(articles) > 0 {
		nextCursor = articles[len(articles)-1].Id
	} else {
		nextCursor = cursor
	}
	return
}

// 标签文章列表
func (s *articleService) GetTagArticles(tagId int64, cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := simple.NewSqlCnd().Eq("tag_id", tagId).Eq("status", model.StatusOk).Desc("id").Limit(20)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	nextCursor = cursor
	articleTags := repositories.ArticleTagRepository.Find(simple.DB(), cnd)
	if len(articleTags) > 0 {
		var articleIds []int64
		for _, articleTag := range articleTags {
			articleIds = append(articleIds, articleTag.ArticleId)
			nextCursor = articleTag.Id
		}
		articles = s.GetArticleInIds(articleIds)
	}
	return
}

// 发布文章
func (s *articleService) Publish(userId int64, title, summary, content, contentType string, tags []string,
	sourceUrl string, share bool) (article *model.Article, err error) {
	title = strings.TrimSpace(title)
	summary = strings.TrimSpace(summary)
	content = strings.TrimSpace(content)

	if len(title) == 0 {
		return nil, errors.New("标题不能为空")
	}
	if share { // 如果是分享的内容，必须有Summary和SourceUrl
		if len(summary) == 0 {
			return nil, errors.New("分享内容摘要不能为空")
		}
		if len(sourceUrl) == 0 {
			return nil, errors.New("分享内容原文链接不能为空")
		}
	} else {
		if len(content) == 0 {
			return nil, errors.New("内容不能为空")
		}
	}
	article = &model.Article{
		UserId:      userId,
		Title:       title,
		Summary:     summary,
		Content:     content,
		ContentType: contentType,
		Status:      model.StatusOk,
		Share:       share,
		SourceUrl:   sourceUrl,
		CreateTime:  simple.NowTimestamp(),
		UpdateTime:  simple.NowTimestamp(),
	}

	err = simple.Tx(simple.DB(), func(tx *gorm.DB) error {
		tagIds := repositories.TagRepository.GetOrCreates(tx, tags)
		err := repositories.ArticleRepository.Create(tx, article)
		if err != nil {
			return err
		}
		repositories.ArticleTagRepository.AddArticleTags(tx, article.Id, tagIds)
		return nil
	})

	if err == nil {
		baiduseo.PushUrl(urls.ArticleUrl(article.Id))
	}
	return
}

// 修改文章
func (s *articleService) Edit(articleId int64, tags []string, title, content string) *simple.CodeError {
	if len(title) == 0 {
		return simple.NewErrorMsg("请输入标题")
	}
	if len(content) == 0 {
		return simple.NewErrorMsg("请填写文章内容")
	}

	err := simple.Tx(simple.DB(), func(tx *gorm.DB) error {
		err := repositories.ArticleRepository.Updates(simple.DB(), articleId, map[string]interface{}{
			"title":   title,
			"content": content,
		})
		if err != nil {
			return err
		}
		tagIds := repositories.TagRepository.GetOrCreates(tx, tags)             // 创建文章对应标签
		repositories.ArticleTagRepository.DeleteArticleTags(tx, articleId)      // 先删掉所有的标签
		repositories.ArticleTagRepository.AddArticleTags(tx, articleId, tagIds) // 然后重新添加标签
		return nil
	})
	cache.ArticleTagCache.Invalidate(articleId)
	return simple.FromError(err)
}

// 相关文章
func (s *articleService) GetRelatedArticles(articleId int64) []model.Article {
	tagIds := cache.ArticleTagCache.Get(articleId)
	if len(tagIds) == 0 {
		return nil
	}
	var articleTags []model.ArticleTag
	simple.DB().Where("tag_id in (?)", tagIds).Limit(30).Find(&articleTags)

	set := hashset.New()
	if len(articleTags) > 0 {
		for _, articleTag := range articleTags {
			set.Add(articleTag.ArticleId)
		}
	}

	var articleIds []int64
	for i, articleId := range set.Values() {
		if i < 10 {
			articleIds = append(articleIds, articleId.(int64))
		}
	}

	return s.GetArticleInIds(articleIds)
}

// 最新文章
func (s *articleService) GetUserNewestArticles(userId int64) []model.Article {
	return repositories.ArticleRepository.Find(simple.DB(), simple.NewSqlCnd().Where("user_id = ? and status = ?",
		userId, model.StatusOk).Desc("id").Limit(10))
}

// 倒序扫描
func (s *articleService) ScanDesc(dateFrom, dateTo int64, cb ScanArticleCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.ArticleRepository.Find(simple.DB(), simple.NewSqlCnd("id", "status", "create_time", "update_time").
			Lt("id", cursor).Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		cb(list)
	}
}

// rss
func (s *articleService) GenerateRss() {
	articles := repositories.ArticleRepository.Find(simple.DB(),
		simple.NewSqlCnd().Where("status = ?", model.StatusOk).Desc("id").Limit(1000))

	var items []*feeds.Item
	for _, article := range articles {
		articleUrl := urls.ArticleUrl(article.Id)
		user := cache.UserCache.Get(article.UserId)
		if user == nil {
			continue
		}
		description := ""
		if article.ContentType == model.ContentTypeMarkdown {
			description = common.GetMarkdownSummary(article.Content)
		} else {
			description = common.GetHtmlSummary(article.Content)
		}
		item := &feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: articleUrl},
			Description: description,
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email.String},
			Created:     simple.TimeFromTimestamp(article.CreateTime),
		}
		items = append(items, item)
	}

	siteTitle := cache.SysConfigCache.GetValue(model.SysConfigSiteTitle)
	siteDescription := cache.SysConfigCache.GetValue(model.SysConfigSiteDescription)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: config.Conf.BaseUrl},
		Description: siteDescription,
		Author:      &feeds.Author{Name: siteTitle},
		Created:     time.Now(),
		Items:       items,
	}
	atom, err := feed.ToAtom()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Conf.StaticPath, "atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Conf.StaticPath, "rss.xml"), rss, false)
	}
}

// 浏览数+1
func (s *articleService) IncrViewCount(articleId int64) {
	simple.DB().Exec("update t_article set view_count = view_count + 1 where id = ?", articleId)
}
