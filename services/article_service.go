package services

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/mlogclub/mlog/services/cache"
	"path"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/mlog/utils"
)

type ScanArticleCallback func(articles []model.Article) bool

type ArticleService struct {
	ArticleRepository    *repositories.ArticleRepository
	ArticleTagRepository *repositories.ArticleTagRepository
	TagRepository        *repositories.TagRepository
}

func NewArticleService() *ArticleService {
	return &ArticleService{
		ArticleRepository:    repositories.NewArticleRepository(),
		ArticleTagRepository: repositories.NewArticleTagRepository(),
		TagRepository:        repositories.NewTagRepository(),
	}
}

func (this *ArticleService) Get(id int64) *model.Article {
	return this.ArticleRepository.Get(simple.GetDB(), id)
}

func (this *ArticleService) Take(where ...interface{}) *model.Article {
	return this.ArticleRepository.Take(simple.GetDB(), where...)
}

func (this *ArticleService) QueryCnd(cnd *simple.QueryCnd) (list []model.Article, err error) {
	return this.ArticleRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *ArticleService) Query(queries *simple.ParamQueries) (list []model.Article, paging *simple.Paging) {
	return this.ArticleRepository.Query(simple.GetDB(), queries)
}

func (this *ArticleService) Update(t *model.Article) error {
	return this.ArticleRepository.Update(simple.GetDB(), t)
}

func (this *ArticleService) Updates(id int64, columns map[string]interface{}) error {
	return this.ArticleRepository.Updates(simple.GetDB(), id, columns)
}

func (this *ArticleService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.ArticleRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *ArticleService) Delete(id int64) error {
	return this.ArticleRepository.UpdateColumn(simple.GetDB(), id, "status", model.ArticleStatusDeleted)
}

// 根据文章编号批量获取文章
func (this *ArticleService) GetArticleInIds(articleIds []int64) []model.Article {
	if len(articleIds) == 0 {
		return nil
	}
	var articles []model.Article
	simple.GetDB().Where("id in (?)", articleIds).Find(&articles)
	return articles
}

// 标签文章列表
func (this *ArticleService) GetTagArticles(tagId int64, page int) (articles []model.Article, paging *simple.Paging) {
	articleTags, paging := this.ArticleTagRepository.Query(simple.GetDB(), simple.NewParamQueries(nil).
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
func (this *ArticleService) Publish(userId int64, title, summary, content, contentType string, categoryId int64,
	tagIds []int64, sourceUrl string) (article *model.Article, err error) {

	article = &model.Article{
		UserId:      userId,
		Title:       title,
		Summary:     summary,
		Content:     content,
		ContentType: contentType,
		CategoryId:  categoryId,
		Status:      model.ArticleShareStatusOk,
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
		err := this.ArticleRepository.Create(tx, article)
		if err != nil {
			return err
		}

		if tagIdsUnique != nil && len(tagIdsUnique) > 0 {
			for _, tagId := range tagIdsUnique {
				if tagId <= 0 {
					continue
				}
				err := this.ArticleTagRepository.Create(tx, &model.ArticleTag{
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
		utils.BaiduUrlPush([]string{utils.BuildArticleUrl(article.Id)})
	}
	return
}

// 添加文章标签
func (this *ArticleService) AddArticleTag(articleId, tagId int64) {
	articleTag := this.ArticleTagRepository.GetUnique(simple.GetDB(), articleId, tagId)
	if articleTag != nil {
		return
	}
	_ = this.ArticleTagRepository.Create(simple.GetDB(), &model.ArticleTag{
		ArticleId:  articleId,
		TagId:      tagId,
		CreateTime: simple.NowTimestamp(),
	})
}

// 删除文章标签
func (this *ArticleService) DelArticleTag(articleId, tagId int64) {
	articleTag := this.ArticleTagRepository.GetUnique(simple.GetDB(), articleId, tagId)
	if articleTag == nil {
		return
	}
	this.ArticleTagRepository.Delete(simple.GetDB(), articleTag.Id)
}

// 相关文章
func (this *ArticleService) GetRelatedArticles(articleId int64) []model.Article {
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
func (this *ArticleService) GetUserNewestArticles(userId int64) []model.Article {
	articles, err := this.ArticleRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("user_id = ? and status = ?",
		userId, model.ArticleStatusPublished).Order("id desc").Size(10))
	if err != nil {
		return nil
	}
	return articles
}

// 扫描
func (this *ArticleService) Scan(cb ScanArticleCallback) {
	var cursor int64
	for {
		list, err := this.ArticleRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ? and status = ? ",
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
func (this *ArticleService) ScanWithDate(dateFrom, dateTo int64, cb ScanArticleCallback) {
	var cursor int64
	for {
		list, err := this.ArticleRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ? and status = ? and create_time >= ? and create_time < ?",
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
func (this *ArticleService) GenerateSitemap() {
	articles, err := this.ArticleRepository.QueryCnd(simple.GetDB(),
		simple.NewQueryCnd("status = ?", model.ArticleStatusPublished).Order("id desc").Size(1000))
	if err != nil {
		logrus.Error(err)
		return
	}

	sm := stm.NewSitemap(0)
	sm.SetDefaultHost(utils.Conf.BaseUrl)
	sm.Create()

	for _, article := range articles {
		articleUrl := utils.BuildArticleUrl(article.Id)
		sm.Add(stm.URL{{"loc", articleUrl}, {"lastmod", simple.TimeFromTimestamp(article.UpdateTime)}})
	}

	data := sm.XMLContent()
	_ = simple.WriteString(path.Join(utils.Conf.StaticPath, "sitemap.xml"), string(data), false)
}

// rss
func (this *ArticleService) GenerateRss() {
	articles, err := this.ArticleRepository.QueryCnd(simple.GetDB(),
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
		Title:       utils.Conf.SiteTitle,
		Link:        &feeds.Link{Href: utils.Conf.BaseUrl},
		Description: "分享生活",
		Author:      &feeds.Author{Name: utils.Conf.SiteTitle},
		Created:     time.Now(),
		Items:       items,
	}
	atom, err := feed.ToAtom()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(utils.Conf.StaticPath, "atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(utils.Conf.StaticPath, "rss.xml"), rss, false)
	}
}

// 每日分享
// title: 标题，最终文章标题为：title + yyyy-MM-dd
// summary: 分享专题的描述
// userIds: 用户编号
func (this *ArticleService) CreateDailyShare(title, summary string, userIds []int64) {
	if userIds == nil || len(userIds) == 0 {
		return
	}

	content := "\n"
	if len(strings.TrimSpace(summary)) > 0 {
		content += "> " + strings.TrimSpace(summary) + "\n\n"
	}

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

	title = title + "（" + simple.TimeFormat(dateFrom, simple.FMT_DATE) + "）"
	_, _ = this.Publish(199, title, "", content, model.ArticleContentTypeMarkdown, 0, nil, "")
}

func (this *ArticleService) GetDailyContent(userIds []int64) string {
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
