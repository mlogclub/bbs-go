package services

import (
	"time"

	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/common/config"
	"bbs-go/common/sitemap"
	"bbs-go/common/urls"
	"bbs-go/model"
	"bbs-go/repositories"
)

const (
	sitemapPath = "sitemap"
)

var SitemapService = newSitemapService()

func newSitemapService() *sitemapService {
	return &sitemapService{}
}

type sitemapService struct {
	building bool // is in building
}

func (s *sitemapService) Get(id int64) *model.Sitemap {
	return repositories.SitemapRepository.Get(simple.DB(), id)
}

func (s *sitemapService) Take(where ...interface{}) *model.Sitemap {
	return repositories.SitemapRepository.Take(simple.DB(), where...)
}

func (s *sitemapService) Find(cnd *simple.SqlCnd) []model.Sitemap {
	return repositories.SitemapRepository.Find(simple.DB(), cnd)
}

func (s *sitemapService) FindOne(cnd *simple.SqlCnd) *model.Sitemap {
	return repositories.SitemapRepository.FindOne(simple.DB(), cnd)
}

func (s *sitemapService) FindPageByParams(params *simple.QueryParams) (list []model.Sitemap, paging *simple.Paging) {
	return repositories.SitemapRepository.FindPageByParams(simple.DB(), params)
}

func (s *sitemapService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Sitemap, paging *simple.Paging) {
	return repositories.SitemapRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *sitemapService) Create(t *model.Sitemap) error {
	return repositories.SitemapRepository.Create(simple.DB(), t)
}

func (s *sitemapService) Update(t *model.Sitemap) error {
	return repositories.SitemapRepository.Update(simple.DB(), t)
}

func (s *sitemapService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.SitemapRepository.Updates(simple.DB(), id, columns)
}

func (s *sitemapService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.SitemapRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *sitemapService) Delete(id int64) {
	repositories.SitemapRepository.Delete(simple.DB(), id)
}

func (s *sitemapService) GenerateToday() {
	if s.building {
		logrus.Info("sitemap is in building")
		return
	}

	s.building = true
	defer func() {
		s.building = false
	}()

	dateFrom := simple.WithTimeAsStartOfDay(time.Now())
	dateTo := dateFrom.Add(time.Hour * 24)

	s.GenerateMisc()
	s.GenerateUser()
	s.Generate(simple.Timestamp(dateFrom), simple.Timestamp(dateTo))
}

func (s *sitemapService) Generate(dateFrom, dateTo int64) {
	sitemapName := "sitemap-" + simple.TimeFormat(simple.TimeFromTimestamp(dateFrom), simple.FMT_DATE)
	sm := sitemap.NewGenerator(config.Conf.AliyunOss.Host, sitemapPath, sitemapName, func(sm *sitemap.Generator, sitemapLoc string) {
		s.AddSitemapIndex(sm, sitemapLoc)
	})

	// topics
	TopicService.ScanDesc(dateFrom, dateTo, func(topics []model.Topic) {
		for _, topic := range topics {
			if topic.Status == model.StatusOk {
				sm.AddURL(sitemap.URL{
					Loc:        urls.TopicUrl(topic.Id),
					Lastmod:    simple.TimeFromTimestamp(topic.LastCommentTime),
					Changefreq: sitemap.ChangefreqDaily,
					Priority:   "1.0",
				})
			}
		}
	})

	// articles
	ArticleService.ScanDesc(dateFrom, dateTo, func(articles []model.Article) {
		for _, article := range articles {
			if article.Status == model.StatusOk {
				sm.AddURL(sitemap.URL{
					Loc:        urls.ArticleUrl(article.Id),
					Lastmod:    simple.TimeFromTimestamp(article.UpdateTime),
					Changefreq: sitemap.ChangefreqWeekly,
					Priority:   "1.0",
				})
			}
		}
	})

	// projects
	ProjectService.ScanDesc(dateFrom, dateTo, func(projects []model.Project) {
		for _, project := range projects {
			sm.AddURL(sitemap.URL{
				Loc:        urls.ProjectUrl(project.Id),
				Lastmod:    simple.TimeFromTimestamp(project.CreateTime),
				Changefreq: sitemap.ChangefreqMonthly,
				Priority:   "1.0",
			})
		}
	})

	sm.Finalize()
}

func (s *sitemapService) GenerateMisc() {
	sm := sitemap.NewGenerator(config.Conf.AliyunOss.Host, sitemapPath, "sitemap-misc", func(sm *sitemap.Generator, sitemapLoc string) {
		s.AddSitemapIndex(sm, sitemapLoc)
	})
	sm.AddURL(sitemap.URL{
		Loc:        urls.AbsUrl("/"),
		Lastmod:    time.Now(),
		Changefreq: sitemap.ChangefreqDaily,
		Priority:   "1.0",
	})
	sm.AddURL(sitemap.URL{
		Loc:        urls.AbsUrl("/topics"),
		Lastmod:    time.Now(),
		Changefreq: sitemap.ChangefreqDaily,
		Priority:   "1.0",
	})
	sm.AddURL(sitemap.URL{
		Loc:        urls.AbsUrl("/articles"),
		Lastmod:    time.Now(),
		Changefreq: sitemap.ChangefreqDaily,
		Priority:   "1.0",
	})
	sm.AddURL(sitemap.URL{
		Loc:        urls.AbsUrl("/projects"),
		Lastmod:    time.Now(),
		Changefreq: sitemap.ChangefreqDaily,
		Priority:   "1.0",
	})

	TagService.Scan(func(tags []model.Tag) bool {
		for _, tag := range tags {
			tagUrl := urls.TagArticlesUrl(tag.Id)

			sm.AddURL(sitemap.URL{
				Loc:        tagUrl,
				Lastmod:    time.Now(),
				Changefreq: sitemap.ChangefreqDaily,
				Priority:   "0.6",
			})
		}
		return true
	})

	sm.Finalize()
}

func (s *sitemapService) GenerateUser() {
	sm := sitemap.NewGenerator(config.Conf.AliyunOss.Host, sitemapPath, "sitemap-user", func(sm *sitemap.Generator, sitemapLoc string) {
		s.AddSitemapIndex(sm, sitemapLoc)
	})
	UserService.Scan(func(users []model.User) {
		for _, user := range users {
			sm.AddURL(sitemap.URL{
				Loc:        urls.UserUrl(user.Id),
				Lastmod:    time.Now(),
				Changefreq: sitemap.ChangefreqWeekly,
				Priority:   "0.6",
			})
		}
	})

	sm.Finalize()
}

func (s *sitemapService) AddSitemapIndex(sm *sitemap.Generator, sitemapLoc string) {
	locName := simple.MD5(sitemapLoc)
	t := s.FindOne(simple.NewSqlCnd().Eq("loc_name", locName))
	if t == nil {
		_ = s.Create(&model.Sitemap{
			Model:      model.Model{},
			Loc:        sitemapLoc,
			Lastmod:    simple.NowTimestamp(),
			LocName:    locName,
			CreateTime: simple.NowTimestamp(),
		})
	} else {
		t.Lastmod = simple.NowTimestamp()
		_ = s.Update(t)
	}

	go func() {
		s.GenerateSitemapIndex(sm)
	}()
}

func (s *sitemapService) GenerateSitemapIndex(sm *sitemap.Generator) {
	sitemaps := s.Find(simple.NewSqlCnd().Desc("id"))

	if len(sitemaps) == 0 {
		return
	}

	var sitemapLocs []sitemap.IndexURL
	for _, s := range sitemaps {
		sitemapLocs = append(sitemapLocs, sitemap.IndexURL{
			Loc:     s.Loc,
			Lastmod: simple.TimeFromTimestamp(s.Lastmod),
		})
	}
	sm.WriteIndex(sitemapLocs)
}
