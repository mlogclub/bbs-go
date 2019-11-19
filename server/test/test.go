package main

import (
	"strconv"
	"time"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"

	"github.com/mlogclub/bbs-go/common/config"
)

func init() {
	config.InitConfig("./bbs-go.yaml")
}

func main() {
	sm := stm.NewSitemap(1)
	sm.SetDefaultHost("https://mlog.club")
	// sm.SetPublicPath("/Users/gaoyoubo/Downloads/sitemap")
	sm.SetSitemapsPath("sitemap")
	sm.SetVerbose(false)
	sm.SetCompress(false)
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
