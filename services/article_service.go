package services

import (
	"errors"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/utils/config"
	"path"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/utils"
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
	return repositories.ArticleRepository.Get(simple.GetDB(), id)
}

func (this *articleService) Take(where ...interface{}) *model.Article {
	return repositories.ArticleRepository.Take(simple.GetDB(), where...)
}

func (this *articleService) QueryCnd(cnd *simple.QueryCnd) (list []model.Article, err error) {
	return repositories.ArticleRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *articleService) Query(queries *simple.ParamQueries) (list []model.Article, paging *simple.Paging) {
	return repositories.ArticleRepository.Query(simple.GetDB(), queries)
}

func (this *articleService) Update(t *model.Article) error {
	err := repositories.ArticleRepository.Update(simple.GetDB(), t)
	if err == nil {
		cache.ArticleCache.InvalidateIndexList()
	}
	return err
}

func (this *articleService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.ArticleRepository.Updates(simple.GetDB(), id, columns)
	if err == nil {
		cache.ArticleCache.InvalidateIndexList()
	}
	return err
}

func (this *articleService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.ArticleRepository.UpdateColumn(simple.GetDB(), id, name, value)
	if err == nil {
		cache.ArticleCache.InvalidateIndexList()
	}
	return err
}

func (this *articleService) Delete(id int64) error {
	err := repositories.ArticleRepository.UpdateColumn(simple.GetDB(), id, "status", model.ArticleStatusDeleted)
	if err == nil {
		cache.ArticleCache.InvalidateIndexList()
	}
	return err
}

// 根据文章编号批量获取文章
func (this *articleService) GetArticleInIds(articleIds []int64) []model.Article {
	if len(articleIds) == 0 {
		return nil
	}
	var articles []model.Article
	simple.GetDB().Where("id in (?)", articleIds).Find(&articles)
	return articles
}

// 获取文章对应的标签
func (this *articleService) GetArticleTags(articleId int64) []model.Tag {
	articleTags, err := repositories.ArticleTagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("article_id = ?", articleId))
	if err != nil {
		return nil
	}
	var tagIds []int64
	for _, articleTag := range articleTags {
		tagIds = append(tagIds, articleTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// 标签文章列表
func (this *articleService) GetTagArticles(tagId int64, page int) (articles []model.Article, paging *simple.Paging) {
	articleTags, paging := repositories.ArticleTagRepository.Query(simple.GetDB(), simple.NewParamQueries(nil).
		Eq("tag_id", tagId).
		Page(page, 20).Desc("id"))
	if len(articleTags) > 0 {
		var articleIds []int64
		for _, articleTag := range articleTags {
			articleIds = append(articleIds, articleTag.ArticleId)
		}
		articles = this.GetArticleInIds(articleIds)
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

	err = simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		tagIds := repositories.TagRepository.GetOrCreates(tx, tags)
		err := repositories.ArticleRepository.Create(tx, article)
		if err != nil {
			return err
		}
		repositories.ArticleTagRepository.AddArticleTags(tx, article.Id, tagIds)
		return nil
	})

	if err == nil {
		// 清理首页文章列表缓存
		cache.ArticleCache.InvalidateIndexList()
		utils.BaiduUrlPush([]string{utils.BuildArticleUrl(article.Id)})
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

	err := simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		tagIds := repositories.TagRepository.GetOrCreates(tx, tags)
		err := repositories.ArticleRepository.Updates(simple.GetDB(), articleId, map[string]interface{}{
			"title":   title,
			"content": content,
		})
		if err != nil {
			return err
		}
		repositories.ArticleTagRepository.RemoveArticleTags(tx, articleId)      // 先删掉所有的标签
		repositories.ArticleTagRepository.AddArticleTags(tx, articleId, tagIds) // 然后重新添加标签
		return nil
	})
	cache.ArticleTagCache.Invalidate(articleId)
	return simple.NewError2(err)
}

// 相关文章
func (this *articleService) GetRelatedArticles(articleId int64) []model.Article {
	tagIds := cache.ArticleTagCache.Get(articleId)
	if len(tagIds) == 0 {
		return nil
	}
	var articleTags []model.ArticleTag
	simple.GetDB().Where("tag_id in (?)", tagIds).Limit(30).Find(&articleTags)

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
	articles, err := repositories.ArticleRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("user_id = ? and status = ?",
		userId, model.ArticleStatusPublished).Order("id desc").Size(10))
	if err != nil {
		return nil
	}
	return articles
}

// 扫描
func (this *articleService) Scan(cb ScanArticleCallback) {
	var cursor int64
	for {
		list, err := repositories.ArticleRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ? and status = ? ",
			cursor, model.ArticleStatusPublished).Order("id asc").Size(300))
		if err != nil {
			break
		}
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		_continue := cb(list)
		if !_continue {
			break
		}
	}
}

// 扫描
func (this *articleService) ScanWithDate(dateFrom, dateTo int64, cb ScanArticleCallback) {
	var cursor int64
	for {
		list, err := repositories.ArticleRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ? and status = ? and create_time >= ? and create_time < ?",
			cursor, model.ArticleStatusPublished, dateFrom, dateTo).Order("id asc").Size(300))
		if err != nil {
			break
		}
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		cb(list)
	}
}

// sitemap
func (this *articleService) GenerateSitemap() {
	articles, err := repositories.ArticleRepository.QueryCnd(simple.GetDB(),
		simple.NewQueryCnd("status = ?", model.ArticleStatusPublished).Order("id desc").Size(1000))
	if err != nil {
		logrus.Error(err)
		return
	}

	sm := stm.NewSitemap(0)
	sm.SetDefaultHost(config.Conf.BaseUrl)
	sm.Create()

	for _, article := range articles {
		articleUrl := utils.BuildArticleUrl(article.Id)
		sm.Add(stm.URL{{"loc", articleUrl}, {"lastmod", simple.TimeFromTimestamp(article.UpdateTime)}})
	}

	data := sm.XMLContent()
	_ = simple.WriteString(path.Join(config.Conf.RootStaticPath, "sitemap.xml"), string(data), false)
}

// rss
func (this *articleService) GenerateRss() {
	articles, err := repositories.ArticleRepository.QueryCnd(simple.GetDB(),
		simple.NewQueryCnd("status = ?", model.ArticleStatusPublished).Order("id desc").Size(1000))
	if err != nil {
		logrus.Error(err)
		return
	}

	var items []*feeds.Item

	for _, article := range articles {
		articleUrl := utils.BuildArticleUrl(article.Id)
		user := cache.UserCache.Get(article.UserId)
		if user == nil {
			continue
		}
		item := &feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: articleUrl},
			Description: article.Summary,
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email},
			Created:     simple.TimeFromTimestamp(article.CreateTime),
		}
		items = append(items, item)
	}

	siteTitle := cache.SysConfigCache.GetValue(model.SysConfigSiteTitle)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: config.Conf.BaseUrl},
		Description: "分享生活",
		Author:      &feeds.Author{Name: siteTitle},
		Created:     time.Now(),
		Items:       items,
	}
	atom, err := feed.ToAtom()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Conf.RootStaticPath, "atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Conf.RootStaticPath, "rss.xml"), rss, false)
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
			if utils.IndexOf(userIds, article.UserId) != -1 {
				content += "## " + article.Title + "\n\n"
				if len(strings.TrimSpace(article.Summary)) > 0 {
					content += strings.TrimSpace(article.Summary) + "\n\n"
				}
				content += "[点击查看原文>>](" + utils.BuildArticleUrl(article.Id) + ")\n\n"
			}
		}
		return true
	})
	return content
}
