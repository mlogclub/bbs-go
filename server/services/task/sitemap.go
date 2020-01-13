package task

import (
	"bytes"
	"compress/gzip"
	"time"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/common/config"
	"bbs-go/common/oss"
	"bbs-go/common/urls"
	"bbs-go/model"
	"bbs-go/services"
)

var sitemapBuilding = false

// 生成sitemap
func SitemapTask() {
	if sitemapBuilding {
		logrus.Info("Sitemap in building...")
		return
	}
	sitemapBuilding = true
	defer func() {
		sitemapBuilding = false
	}()


	// gzip 格式sitemap
	smGzip := stm.NewSitemap(1)
	smGzip.SetDefaultHost(config.Conf.BaseUrl)         // 网站host
	smGzip.SetSitemapsHost(config.Conf.AliyunOss.Host) // 上传到阿里云所以host设置为阿里云
	smGzip.SetSitemapsPath("sitemap")                  // sitemap存放目录
	smGzip.SetVerbose(false)
	smGzip.SetPretty(false)
	smGzip.SetCompress(true)
	smGzip.SetAdapter(&AliyunOssAdapter{})
	smGzip.Create()

	// xml 格式sitemap
	smXml := stm.NewSitemap(1)
	smXml.SetDefaultHost(config.Conf.BaseUrl)         // 网站host
	smXml.SetSitemapsHost(config.Conf.AliyunOss.Host) // 上传到阿里云所以host设置为阿里云
	smXml.SetSitemapsPath("sitemap")                  // sitemap存放目录
	smXml.SetVerbose(false)
	smXml.SetPretty(false)
	smXml.SetCompress(false)
	smXml.SetAdapter(&AliyunOssAdapter{})
	smXml.Create()

	generate(smGzip, smXml)
}

func generate(smArr ...*stm.Sitemap) {
	addStmUrl(stm.URL{
		{"loc", "/"},
		{"lastmod", time.Now()},
		{"changefreq", "daily"},
		{"priority", 1.0},
	}, smArr...)

	addStmUrl(stm.URL{
		{"loc", "/topics"},
		{"lastmod", time.Now()},
		{"changefreq", "daily"},
		{"priority", 1.0},
	}, smArr...)

	addStmUrl(stm.URL{
		{"loc", "/articles"},
		{"lastmod", time.Now()},
		{"changefreq", "daily"},
		{"priority", 1.0},
	}, smArr...)

	addStmUrl(stm.URL{
		{"loc", "/projects"},
		{"lastmod", time.Now()},
		{"changefreq", "daily"},
		{"priority", 1.0},
	}, smArr...)

	services.TopicService.ScanDesc(func(topics []model.Topic) bool {
		for _, topic := range topics {
			if topic.Status == model.StatusOk {
				topicUrl := urls.TopicUrl(topic.Id)
				addStmUrl(stm.URL{
					{"loc", topicUrl},
					{"lastmod", simple.TimeFromTimestamp(topic.CreateTime)},
				}, smArr...)
			}
		}
		return true
	})

	services.ArticleService.ScanDesc(func(articles []model.Article) bool {
		for _, article := range articles {
			if article.Status == model.StatusOk {
				articleUrl := urls.ArticleUrl(article.Id)
				addStmUrl(stm.URL{
					{"loc", articleUrl},
					{"lastmod", simple.TimeFromTimestamp(article.UpdateTime)},
				}, smArr...)
			}
		}
		return true
	})

	services.UserService.Scan(func(users []model.User) {
		for _, user := range users {
			userUrl := urls.UserUrl(user.Id)
			addStmUrl(stm.URL{
				{"loc", userUrl},
				{"lastmod", time.Now()},
			}, smArr...)
		}
	})

	services.ProjectService.Scan(func(projects []model.Project) bool {
		for _, project := range projects {
			projectUrl := urls.ProjectUrl(project.Id)
			addStmUrl(stm.URL{
				{"loc", projectUrl},
				{"lastmod", simple.TimeFromTimestamp(project.CreateTime)},
			}, smArr...)
		}
		return true
	})

	services.TagService.Scan(func(tags []model.Tag) bool {
		for _, tag := range tags {
			tagUrl := urls.TagArticlesUrl(tag.Id)
			addStmUrl(stm.URL{
				{"loc", tagUrl},
				{"lastmod", time.Now()},
				{"changefreq", "daily"},
				{"priority", 0.6},
			}, smArr...)
		}
		return true
	})

	for _, sm := range smArr {
		sm.Finalize().PingSearchEngines()
	}
}

func addStmUrl(stmUrl stm.URL, smArr ...*stm.Sitemap) {
	for _, sm := range smArr {
		sm.Add(stmUrl)
	}
}

// sitemap上传到aliyun
type AliyunOssAdapter struct {
}

// Bytes gets written content.
func (adp *AliyunOssAdapter) Bytes() [][]byte {
	return nil
}

// Write will create sitemap xml file into the file systems.
func (adp *AliyunOssAdapter) Write(loc *stm.Location, data []byte) {
	var out []byte
	if stm.GzipPtn.MatchString(loc.Filename()) { // 如果需要压缩
		var in bytes.Buffer
		w := gzip.NewWriter(&in)
		_, _ = w.Write(data)
		_ = w.Close()
		out = in.Bytes()
	} else { // 如果不需要压缩
		out = data
	}

	sitemapUrl, err := oss.PutObject(loc.PathInPublic(), out)
	if err != nil {
		logrus.Error("Upload sitemap to aliyun oss error:", err)
	} else {
		logrus.Info("Upload sitemap:", sitemapUrl)
	}
}
