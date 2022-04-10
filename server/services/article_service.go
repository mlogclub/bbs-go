package services

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/bbsurls"
	"bbs-go/pkg/seo"
	"errors"
	"math"
	"path"
	"strings"
	"time"

	"github.com/emirpasic/gods/sets/hashset"

	"bbs-go/cache"
	"bbs-go/pkg/config"
	"bbs-go/repositories"

	"github.com/gorilla/feeds"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/files"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"bbs-go/model"
	"bbs-go/pkg/common"
)

var ArticleService = newArticleService()

func newArticleService() *articleService {
	return &articleService{}
}

type articleService struct {
}

func (s *articleService) Get(id int64) *model.Article {
	return repositories.ArticleRepository.Get(sqls.DB(), id)
}

func (s *articleService) Take(where ...interface{}) *model.Article {
	return repositories.ArticleRepository.Take(sqls.DB(), where...)
}

func (s *articleService) Find(cnd *sqls.Cnd) []model.Article {
	return repositories.ArticleRepository.Find(sqls.DB(), cnd)
}

func (s *articleService) FindOne(cnd *sqls.Cnd) *model.Article {
	return repositories.ArticleRepository.FindOne(sqls.DB(), cnd)
}

func (s *articleService) FindPageByParams(params *params.QueryParams) (list []model.Article, paging *sqls.Paging) {
	return repositories.ArticleRepository.FindPageByParams(sqls.DB(), params)
}

func (s *articleService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Article, paging *sqls.Paging) {
	return repositories.ArticleRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *articleService) Update(t *model.Article) error {
	err := repositories.ArticleRepository.Update(sqls.DB(), t)
	return err
}

func (s *articleService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.ArticleRepository.Updates(sqls.DB(), id, columns)
	return err
}

func (s *articleService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.ArticleRepository.UpdateColumn(sqls.DB(), id, name, value)
	return err
}

func (s *articleService) Delete(id int64) error {
	err := repositories.ArticleRepository.UpdateColumn(sqls.DB(), id, "status", constants.StatusDeleted)
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
	sqls.DB().Where("id in (?)", articleIds).Order("id desc").Find(&articles)
	return articles
}

// 获取文章对应的标签
func (s *articleService) GetArticleTags(articleId int64) []model.Tag {
	articleTags := repositories.ArticleTagRepository.Find(sqls.DB(), sqls.NewCnd().Where("article_id = ?", articleId))
	var tagIds []int64
	for _, articleTag := range articleTags {
		tagIds = append(tagIds, articleTag.TagId)
	}
	return cache.TagCache.GetList(tagIds)
}

// 文章列表
func (s *articleService) GetArticles(cursor int64) (articles []model.Article, nextCursor int64, hasMore bool) {
	limit := 20
	cnd := sqls.NewCnd().Eq("status", constants.StatusOk).Desc("id").Limit(limit)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	articles = repositories.ArticleRepository.Find(sqls.DB(), cnd)
	if len(articles) > 0 {
		nextCursor = articles[len(articles)-1].Id
		hasMore = len(articles) >= limit
	} else {
		nextCursor = cursor
	}
	return
}

// 标签文章列表
func (s *articleService) GetTagArticles(tagId int64, cursor int64) (articles []model.Article, nextCursor int64, hasMore bool) {
	limit := 20
	cnd := sqls.NewCnd().Eq("tag_id", tagId).Eq("status", constants.StatusOk).Desc("id").Limit(limit)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	nextCursor = cursor
	articleTags := repositories.ArticleTagRepository.Find(sqls.DB(), cnd)
	if len(articleTags) > 0 {
		var articleIds []int64
		for _, articleTag := range articleTags {
			articleIds = append(articleIds, articleTag.ArticleId)
			nextCursor = articleTag.Id
		}
		articles = s.GetArticleInIds(articleIds)
	}
	hasMore = len(articleTags) >= limit
	return
}

// 发布文章
func (s *articleService) Publish(userId int64, form model.CreateArticleForm) (article *model.Article, err error) {
	form.Title = strings.TrimSpace(form.Title)
	form.Summary = strings.TrimSpace(form.Summary)
	form.Content = strings.TrimSpace(form.Content)

	if strs.IsBlank(form.Title) {
		return nil, errors.New("标题不能为空")
	}
	if strs.IsBlank(form.Content) {
		return nil, errors.New("内容不能为空")
	}

	// 获取后台配置 否是开启发表文章审核
	status := constants.StatusOk
	sysConfigArticlePending := cache.SysConfigCache.GetValue(constants.SysConfigArticlePending)
	if strings.ToLower(sysConfigArticlePending) == "true" {
		status = constants.StatusPending
	}

	article = &model.Article{
		UserId:      userId,
		Title:       form.Title,
		Summary:     form.Summary,
		Content:     form.Content,
		ContentType: form.ContentType,
		Status:      status,
		SourceUrl:   form.SourceUrl,
		CreateTime:  dates.NowTimestamp(),
		UpdateTime:  dates.NowTimestamp(),
	}

	err = sqls.DB().Transaction(func(tx *gorm.DB) error {
		tagIds := repositories.TagRepository.GetOrCreates(tx, form.Tags)
		err := repositories.ArticleRepository.Create(tx, article)
		if err != nil {
			return err
		}
		repositories.ArticleTagRepository.AddArticleTags(tx, article.Id, tagIds)
		return nil
	})

	if err == nil {
		seo.Push(bbsurls.ArticleUrl(article.Id))
	}
	return
}

// 修改文章
func (s *articleService) Edit(articleId int64, tags []string, title, content string) *web.CodeError {
	if len(title) == 0 {
		return web.NewErrorMsg("请输入标题")
	}
	if len(content) == 0 {
		return web.NewErrorMsg("请填写文章内容")
	}

	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		err := repositories.ArticleRepository.Updates(sqls.DB(), articleId, map[string]interface{}{
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
	return web.FromError(err)
}

func (s *articleService) PutTags(articleId int64, tags []string) {
	tagIds := repositories.TagRepository.GetOrCreates(sqls.DB(), tags)             // 创建文章对应标签
	repositories.ArticleTagRepository.DeleteArticleTags(sqls.DB(), articleId)      // 先删掉所有的标签
	repositories.ArticleTagRepository.AddArticleTags(sqls.DB(), articleId, tagIds) // 然后重新添加标签
	cache.ArticleTagCache.Invalidate(articleId)
}

// 相关文章
func (s *articleService) GetRelatedArticles(articleId int64) []model.Article {
	tagIds := cache.ArticleTagCache.Get(articleId)
	if len(tagIds) == 0 {
		return nil
	}
	var articleTags []model.ArticleTag
	sqls.DB().Where("tag_id in (?)", tagIds).Limit(30).Find(&articleTags)

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

// 近期文章
func (s *articleService) GetNearlyArticles(articleId int64) []model.Article {
	articles := repositories.ArticleRepository.Find(sqls.DB(), sqls.NewCnd().Where("id < ?", articleId).Desc("id").Limit(10))
	var ret []model.Article
	for _, article := range articles {
		if article.Status == constants.StatusOk {
			ret = append(ret, article)
		}
	}
	return ret
}

// 倒序扫描
func (s *articleService) ScanDesc(callback func(articles []model.Article)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.ArticleRepository.Find(sqls.DB(), sqls.NewCnd().
			Cols("id", "status", "create_time", "update_time").
			Lt("id", cursor).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

func (s *articleService) ScanByUser(userId int64, callback func(articles []model.Article)) {
	var cursor int64 = 0
	for {
		list := repositories.ArticleRepository.Find(sqls.DB(), sqls.NewCnd().
			Eq("user_id", userId).Gt("id", cursor).Asc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// 倒序扫描
func (s *articleService) ScanDescWithDate(dateFrom, dateTo int64, callback func(articles []model.Article)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.ArticleRepository.Find(sqls.DB(), sqls.NewCnd().
			Cols("id", "status", "create_time", "update_time").
			Lt("id", cursor).Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// rss
func (s *articleService) GenerateRss() {
	articles := repositories.ArticleRepository.Find(sqls.DB(),
		sqls.NewCnd().Where("status = ?", constants.StatusOk).Desc("id").Limit(200))

	var items []*feeds.Item
	for _, article := range articles {
		articleUrl := bbsurls.ArticleUrl(article.Id)
		user := cache.UserCache.Get(article.UserId)
		if user == nil {
			continue
		}
		description := common.GetSummary(article.ContentType, article.Content)
		item := &feeds.Item{
			Title:       article.Title,
			Link:        &feeds.Link{Href: articleUrl},
			Description: description,
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email.String},
			Created:     dates.FromTimestamp(article.CreateTime),
		}
		items = append(items, item)
	}

	siteTitle := cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
	siteDescription := cache.SysConfigCache.GetValue(constants.SysConfigSiteDescription)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: config.Instance.BaseUrl},
		Description: siteDescription,
		Author:      &feeds.Author{Name: siteTitle},
		Created:     time.Now(),
		Items:       items,
	}
	atom, err := feed.ToAtom()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = files.WriteString(path.Join(config.Instance.StaticPath, "atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = files.WriteString(path.Join(config.Instance.StaticPath, "rss.xml"), rss, false)
	}
}

// 浏览数+1
func (s *articleService) IncrViewCount(articleId int64) {
	sqls.DB().Exec("update t_article set view_count = view_count + 1 where id = ?", articleId)
}

func (s *articleService) GetUserArticles(userId, cursor int64) (articles []model.Article, nextCursor int64, hasMore bool) {
	limit := 20
	cnd := sqls.NewCnd()
	if userId > 0 {
		cnd.Eq("user_id", userId)
	}
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	cnd.Eq("status", constants.StatusOk).Desc("id").Limit(limit)
	articles = repositories.ArticleRepository.Find(sqls.DB(), cnd)
	if len(articles) > 0 {
		nextCursor = articles[len(articles)-1].Id
		hasMore = len(articles) >= limit
	} else {
		nextCursor = cursor
	}
	return
}
