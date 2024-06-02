package sitemap

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/uploader"
	"bytes"
	"compress/gzip"
	"log/slog"
	"strings"
	"time"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/mlogclub/simple/common/dates"

	"bbs-go/internal/models"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/services"
)

const (
	changefreqAlways  = "always"
	changefreqHourly  = "hourly"
	changefreqDaily   = "daily"
	changefreqWeekly  = "weekly"
	changefreqMonthly = "monthly"
	changefreqYearly  = "yearly"
	changefreqNever   = "never"
)

var building = false

// Generate
func Generate() {
	if !config.Instance.IsProd() {
		return
	}
	if building {
		slog.Info("Sitemap in building...")
		return
	}
	building = true
	defer func() {
		building = false
	}()

	sm := stm.NewSitemap(0)
	sm.SetDefaultHost(config.Instance.BaseUrl) // 网站host
	if uploader.IsEnabledOss() {
		sm.SetSitemapsHost(config.Instance.Uploader.AliyunOss.Host) // 上传到阿里云所以host设置为阿里云
	} else {
		sm.SetPublicPath(config.Instance.Uploader.Local.Host)
	}
	sm.SetSitemapsPath("sitemap") // sitemap存放目录
	sm.SetVerbose(false)
	sm.SetPretty(false)
	sm.SetCompress(true)
	sm.SetAdapter(&myAdapter{})
	sm.Create()

	sm.Add(stm.URL{
		{"loc", bbsurls.AbsUrl("/")},
		{"lastmod", time.Now()},
		{"changefreq", changefreqHourly},
	})

	sm.Add(stm.URL{
		{"loc", bbsurls.AbsUrl("/topics")},
		{"lastmod", time.Now()},
		{"changefreq", changefreqHourly},
	})

	sm.Add(stm.URL{
		{"loc", bbsurls.AbsUrl("/articles")},
		{"lastmod", time.Now()},
		{"changefreq", changefreqAlways},
	})

	// var (
	// 	dateFrom = dates.Timestamp(time.Now().AddDate(0, -1, 0))
	// 	dateTo   = dates.NowTimestamp()
	// )

	services.ArticleService.ScanDesc(func(articles []models.Article) {
		for _, article := range articles {
			if article.Status == constants.StatusOk {
				articleUrl := bbsurls.ArticleUrl(article.Id)
				sm.Add(stm.URL{
					{"loc", articleUrl},
					{"lastmod", dates.FromTimestamp(article.UpdateTime)},
					{"changefreq", changefreqMonthly},
					{"priority", 0.6},
				})
			}
		}
	})

	services.TopicService.ScanDesc(func(topics []models.Topic) {
		for _, topic := range topics {
			if topic.Status == constants.StatusOk {
				topicUrl := bbsurls.TopicUrl(topic.Id)
				sm.Add(stm.URL{
					{"loc", topicUrl},
					{"lastmod", dates.FromTimestamp(topic.LastCommentTime)},
					{"changefreq", changefreqDaily},
					{"priority", 0.6},
				})
			}
		}
	})

	services.TagService.Scan(func(tags []models.Tag) {
		for _, tag := range tags {
			tagUrl := bbsurls.TagArticlesUrl(tag.Id)
			sm.Add(stm.URL{
				{"loc", tagUrl},
				{"lastmod", time.Now()},
				{"changefreq", changefreqMonthly},
			})
		}
	})

	services.UserService.Scan(func(users []models.User) {
		for _, user := range users {
			sm.Add(stm.URL{
				{"loc", bbsurls.UserUrl(user.Id)},
				{"lastmod", time.Now()},
				{"changefreq", changefreqWeekly},
			})
		}
	})

	sm.Finalize().PingSearchEngines()
	// sm.Finalize().PingSearchEngines("http://www.google.cn/webmasters/tools/ping?sitemap=%s")
}

// My Adapter
type myAdapter struct {
}

// Bytes gets written content.
func (adp *myAdapter) Bytes() [][]byte {
	return nil
}

// Write will create sitemap xml file into the file systems.
func (adp *myAdapter) Write(loc *stm.Location, data []byte) {
	if stm.GzipPtn.MatchString(loc.Filename()) { // gzip
		var out []byte
		var in bytes.Buffer
		w := gzip.NewWriter(&in)
		_, _ = w.Write(data)
		_ = w.Close()
		out = in.Bytes()

		// 写入gzip格式数据
		adp.ossWrite(loc.PathInPublic(), out)

		// 写入原始数据
		adp.ossWrite(strings.ReplaceAll(loc.PathInPublic(), ".gz", ""), data)
	} else { // 非gzip
		adp.ossWrite(loc.PathInPublic(), data)
	}
}

// oss写入
func (adp *myAdapter) ossWrite(fileKey string, out []byte) {
	if _url, err := uploader.PutObject(fileKey, out, ""); err != nil {
		slog.Error("Upload sitemap error:", slog.Any("err", err))
	} else {
		slog.Info("Upload sitemap:", slog.String("url", _url))
	}
}
