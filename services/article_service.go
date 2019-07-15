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
	return repositories.ArticleRepository.Update(simple.GetDB(), t)
}

func (this *articleService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ArticleRepository.Updates(simple.GetDB(), id, columns)
}

func (this *articleService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ArticleRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *articleService) Delete(id int64) error {
	return repositories.ArticleRepository.UpdateColumn(simple.GetDB(), id, "status", model.ArticleStatusDeleted)
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
	tagIds []int64, sourceUrl string, share bool) (article *model.Article, err error) {

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

	// 标签滤重
	var tagIdsUnique []int64
	if tagIds != nil && len(tagIds) > 0 {
		for _, tagId := range tagIds {
			if !simple.Contains(tagId, tagIdsUnique) {
				tagIdsUnique = append(tagIdsUnique, tagId)
			}
		}
	}

	err = simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		err := repositories.ArticleRepository.Create(tx, article)
		if err != nil {
			return err
		}

		if tagIdsUnique != nil && len(tagIdsUnique) > 0 {
			for _, tagId := range tagIdsUnique {
				if tagId <= 0 {
					continue
				}
				err := repositories.ArticleTagRepository.Create(tx, &model.ArticleTag{
					ArticleId:  article.Id,
					TagId:      tagId,
					CreateTime: simple.NowTimestamp(),
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err == nil {
		// 清理首页文章列表缓存
		cache.ArticleCache.InvalidateIndexList()
		utils.BaiduUrlPush([]string{utils.BuildArticleUrl(article.Id)})
	}
	return
}

// 添加文章标签
func (this *articleService) AddArticleTag(articleId, tagId int64) {
	articleTag := repositories.ArticleTagRepository.GetUnique(simple.GetDB(), articleId, tagId)
	if articleTag != nil {
		return
	}
	_ = repositories.ArticleTagRepository.Create(simple.GetDB(), &model.ArticleTag{
		ArticleId:  articleId,
		TagId:      tagId,
		CreateTime: simple.NowTimestamp(),
	})
}

// 删除文章标签
func (this *articleService) DelArticleTag(articleId, tagId int64) {
	articleTag := repositories.ArticleTagRepository.GetUnique(simple.GetDB(), articleId, tagId)
	if articleTag == nil {
		return
	}
	repositories.ArticleTagRepository.Delete(simple.GetDB(), articleTag.Id)
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

	feed := &feeds.Feed{
		Title:       config.Conf.SiteTitle,
		Link:        &feeds.Link{Href: config.Conf.BaseUrl},
		Description: "分享生活",
		Author:      &feeds.Author{Name: config.Conf.SiteTitle},
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
