package services

import (
	"errors"
	"math"
	"path"
	"strings"
	"time"

	"github.com/emirpasic/gods/sets/hashset"

	"github.com/mlogclub/bbs-go/common/baiduseo"
	"github.com/mlogclub/bbs-go/common/config"
	"github.com/mlogclub/bbs-go/common/urls"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/bbs-go/services/cache"

	"github.com/gorilla/feeds"
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/model"
)

type ScanArticleCallback func(articles []model.Article) bool

var ArticleService = newArticleService()

func newArticleService() *articleService {
	return &articleService{

	}
}

type articleService struct {
}

func (this *articleService) Get(id int64) *model.Article {
	return repositories.ArticleRepository.Get(simple.DB(), id)
}

func (this *articleService) Take(where ...interface{}) *model.Article {
	return repositories.ArticleRepository.Take(simple.DB(), where...)
}

func (this *articleService) Find(cnd *simple.SqlCnd) []model.Article {
	return repositories.ArticleRepository.Find(simple.DB(), cnd)
}

func (this *articleService) FindOne(cnd *simple.SqlCnd) *model.Article {
	return repositories.ArticleRepository.FindOne(simple.DB(), cnd)
}

func (this *articleService) FindPageByParams(params *simple.QueryParams) (list []model.Article, paging *simple.Paging) {
	return repositories.ArticleRepository.FindPageByParams(simple.DB(), params)
}

func (this *articleService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Article, paging *simple.Paging) {
	return repositories.ArticleRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *articleService) Update(t *model.Article) error {
	err := repositories.ArticleRepository.Update(simple.DB(), t)
	return err
}

func (this *articleService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.ArticleRepository.Updates(simple.DB(), id, columns)
	return err
}

func (this *articleService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.ArticleRepository.UpdateColumn(simple.DB(), id, name, value)
	return err
}

func (this *articleService) Delete(id int64) error {
	err := repositories.ArticleRepository.UpdateColumn(simple.DB(), id, "status", model.ArticleStatusDeleted)
	if err == nil {
		// 删掉专栏文章
		SubjectContentService.DeleteByEntity(model.EntityTypeArticle, id)
		// 删掉标签文章
		ArticleTagService.DeleteByArticleId(id)
	}
	return err
}

// 根据文章编号批量获取文章
func (this *articleService) GetArticleInIds(articleIds []int64) []model.Article {
	if len(articleIds) == 0 {
		return nil
	}
	var articles []model.Article
	simple.DB().Where("id in (?)", articleIds).Find(&articles)
	return articles
}

// 获取文章对应的标签
func (this *articleService) GetArticleTags(articleId int64) []model.Tag {
	articleTags := repositories.ArticleTagRepository.Find(simple.DB(), simple.NewSqlCnd().Where("article_id = ?", articleId))
	var tagIds []int64
	for _, articleTag := range articleTags {
		tagIds = append(tagIds, articleTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// 文章列表
func (this *articleService) GetArticles(cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := simple.NewSqlCnd().Eq("status", model.ArticleStatusPublished).Desc("id").Limit(20)
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
func (this *articleService) GetTagArticles(tagId int64, cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := simple.NewSqlCnd().Eq("tag_id", tagId).Eq("status", model.ArticleTagStatusOk).Desc("id").Limit(20)
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
		articles = this.GetArticleInIds(articleIds)
	}
	return
}

// 分类文章列表
func (this *articleService) GetCategoryArticles(categoryId int64, cursor int64) (articles []model.Article, nextCursor int64) {
	cnd := simple.NewSqlCnd().Eq("category_id", categoryId).Eq("status", model.ArticleStatusPublished).Limit(20).Desc("id")
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

// 发布文章
func (this *articleService) Publish(userId int64, title, summary, content, contentType string, categoryId int64,
	tags []string, sourceUrl string, share bool) (article *model.Article, err error) {

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
		CategoryId:  categoryId,
		Status:      model.ArticleStatusPublished,
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
func (this *articleService) Edit(articleId int64, tags []string, title, content string) *simple.CodeError {
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
func (this *articleService) GetRelatedArticles(articleId int64) []model.Article {
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

	return this.GetArticleInIds(articleIds)
}

// 最新文章
func (this *articleService) GetUserNewestArticles(userId int64) []model.Article {
	return repositories.ArticleRepository.Find(simple.DB(), simple.NewSqlCnd().Where("user_id = ? and status = ?",
		userId, model.ArticleStatusPublished).Desc("id").Limit(10))
}

// 扫描
func (this *articleService) Scan(cb ScanArticleCallback) {
	var cursor int64
	for {
		list := repositories.ArticleRepository.Find(simple.DB(), simple.NewSqlCnd().Where("id > ? ",
			cursor).Asc("id").Limit(100))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		if !cb(list) {
			break
		}
	}
}

// 从新往旧扫描
func (this *articleService) ScanDesc(cb ScanArticleCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.ArticleRepository.Find(simple.DB(), simple.NewSqlCnd("id", "status", "update_time").Where("id < ? ",
			cursor).Desc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		if !cb(list) {
			break
		}
	}
}

// 扫描
func (this *articleService) ScanWithDate(dateFrom, dateTo int64, cb ScanArticleCallback) {
	var cursor int64
	for {
		list := repositories.ArticleRepository.Find(simple.DB(), simple.NewSqlCnd().Where("id > ? and status = ? and create_time >= ? and create_time < ?",
			cursor, model.ArticleStatusPublished, dateFrom, dateTo).Asc("id").Limit(300))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		cb(list)
	}
}

// rss
func (this *articleService) GenerateRss() {
	articles := repositories.ArticleRepository.Find(simple.DB(),
		simple.NewSqlCnd().Where("status = ?", model.ArticleStatusPublished).Desc("id").Limit(1000))

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

// 生成码农日报内容
func (this *articleService) GetDailyContent(userIds []int64) string {
	if userIds == nil || len(userIds) == 0 {
		return ""
	}

	content := "\n"

	dateFromTemp := time.Now().Add(-time.Hour * 24)
	dateToTemp := time.Now()
	dateFrom := time.Date(dateFromTemp.Year(), dateFromTemp.Month(), dateFromTemp.Day(), 0, 0, 0, 0, dateFromTemp.Location())
	dateTo := time.Date(dateToTemp.Year(), dateToTemp.Month(), dateToTemp.Day(), 0, 0, 0, 0, dateToTemp.Location())

	this.ScanWithDate(simple.Timestamp(dateFrom), simple.Timestamp(dateTo), func(articles []model.Article) bool {
		for _, article := range articles {
			if common.IndexOf(userIds, article.UserId) != -1 {
				content += "## " + article.Title + "\n\n"
				if len(strings.TrimSpace(article.Summary)) > 0 {
					content += strings.TrimSpace(article.Summary) + "\n\n"
				}
				content += "[点击查看原文>>](" + urls.ArticleUrl(article.Id) + ")\n\n"
			}
		}
		return true
	})
	return content
}
