// 采集相关任务

package task

import (
	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/common/config"
	"github.com/mlogclub/bbs-go/common/urls"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/bbs-go/services/collect/oschina"
	"github.com/mlogclub/bbs-go/services/collect/studygolang"
)

const (
	collectUserId int64 = 2 // 采集用户
)

// 采集studygolang项目
func CollectStudyGoLangProjectTask() {
	// 只采集第一页
	studygolang.GetStudyGoLangPage(1, func(url string) {
		p := studygolang.GetStudyGolangProject(url)
		if p == nil {
			logrus.Warn("项目采集失败：" + url)
			return
		}
		temp := services.ProjectService.Take("name = ?", p.Name)
		if temp != nil {
			logrus.Warn("项目已经存在：" + temp.Name)
			return
		}
		logrus.Info("采集项目：" + p.Name + ", " + url)
		_, _ = services.ProjectService.Publish(collectUserId, p.Name, p.Title, p.Logo, p.Url, p.DocUrl, p.DownloadUrl,
			model.ContentTypeMarkdown, p.Content)
	})
}

// 采集开源中国项目
func CollectOschinaProjectTask() {
	urls := oschina.GetPage(1)
	if len(urls) > 0 {
		for _, url := range urls {
			p := oschina.GetProject(url)
			temp := services.ProjectService.Take("name = ?", p.Name)
			if temp != nil {
				logrus.Warn("项目已经存在：" + temp.Name)
				continue
			}
			logrus.Info("采集项目：" + p.Name + ", " + url)
			_, _ = services.ProjectService.Publish(2, p.Name, p.Title, p.Logo, p.Url, p.DocUrl, p.DownloadUrl,
				model.ContentTypeHtml, p.Content)
		}
	}
}

// 生成sitemap
func SitemapTask() {
	sm := stm.NewSitemap(1)
	sm.SetDefaultHost(config.Conf.BaseUrl)
	sm.SetPublicPath(config.Conf.StaticPath)
	sm.SetSitemapsPath("")
	sm.SetVerbose(false)
	sm.SetCompress(false)
	sm.SetPretty(false)
	sm.Create()

	services.ArticleService.ScanDesc(func(articles []model.Article) bool {
		for _, article := range articles {
			if article.Status == model.ArticleStatusPublished {
				articleUrl := urls.ArticleUrl(article.Id)
				sm.Add(stm.URL{{"loc", articleUrl}, {"lastmod", simple.TimeFromTimestamp(article.UpdateTime)}})
			}
		}
		return true
	})

	services.TopicService.ScanDesc(func(topics []model.Topic) bool {
		for _, topic := range topics {
			if topic.Status == model.TopicStatusOk {
				topicUrl := urls.TopicUrl(topic.Id)
				sm.Add(stm.URL{{"loc", topicUrl}, {"lastmod", simple.TimeFromTimestamp(topic.CreateTime)}})
			}
		}
		return true
	})

	services.ProjectService.ScanDesc(func(projects []model.Project) bool {
		for _, project := range projects {
			projectUrl := urls.ProjectUrl(project.Id)
			sm.Add(stm.URL{{"loc", projectUrl}, {"lastmod", simple.TimeFromTimestamp(project.CreateTime)}})
		}
		return true
	})

	sm.Finalize().PingSearchEngines()
}

// 生成rss
func RssTask() {
	services.ArticleService.GenerateRss()

	services.TopicService.GenerateRss()

	services.ProjectService.GenerateRss()
}

// Ping百度
func BaiduPing() {
	services.ArticleService.Scan(func(articles []model.Article) bool {
		var pushUrls []string
		for _, article := range articles {
			if article.Status == model.ArticleStatusPublished {
				pushUrls = append(pushUrls, urls.ArticleUrl(article.Id))
			}
		}
		common.BaiduUrlPush(pushUrls)
		return true
	})

	services.TopicService.Scan(func(topics []model.Topic) bool {
		var pushUrls []string
		for _, topic := range topics {
			if topic.Status == model.TopicStatusOk {
				pushUrls = append(pushUrls, urls.TopicUrl(topic.Id))
			}
		}
		common.BaiduUrlPush(pushUrls)
		return true
	})

	services.ProjectService.Scan(func(projects []model.Project) bool {
		var pushUrls []string
		for _, project := range projects {
			pushUrls = append(pushUrls, urls.ProjectUrl(project.Id))
		}
		common.BaiduUrlPush(pushUrls)
		return true
	})
}
