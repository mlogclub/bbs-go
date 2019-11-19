package main

import (
	"strconv"
	"time"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common/config"
	"github.com/mlogclub/bbs-go/services/task"
)

func init() {
	config.InitConfig("./bbs-go.yaml")
}

func main() {
	go func() {
		buildSitemap()
	}()
	go func() {
		buildSitemap()
	}()
}

var sitemapBuilding = false

func buildSitemap() {
	if sitemapBuilding {
		logrus.Info("Sitemap in building...")
		return
	}
	sitemapBuilding = true
	defer func() {
		sitemapBuilding = false
	}()

	sm := stm.NewSitemap(1)
	sm.SetDefaultHost("https://mlog.club")
	// sm.SetPublicPath("/Users/gaoyoubo/Downloads/sitemap")
	sm.SetSitemapsPath("sitemap1")
	sm.SetFilename("1")
	sm.SetVerbose(false)
	sm.SetCompress(true)
	sm.SetAdapter(&task.AliyunOssAdapter{})
	sm.Create()

	sm.Add(stm.URL{
		{"loc", "/topics"},
		{"lastmod", time.Now()},
		{"changefreq", "daily"},
		{"priority", 1.0},
	})

	sm.Add(stm.URL{
		{"loc", "/articles"},
		{"lastmod", time.Now()},
		{"changefreq", "daily"},
		{"priority", 1.0},
	})

	for i := 0; i < 1; i++ {
		url := "https://mlog.club/article/" + strconv.Itoa(i)
		sm.Add(stm.URL{
			{"loc", url},
			{"lastmod", time.Now()},
			{"priority", 0.5},
		})
	}

	sm.Finalize()
}
