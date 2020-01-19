package main

import (
	"fmt"
	"strconv"
	"time"

	"bbs-go/common/config"
	"bbs-go/common/sitemap"
	"bbs-go/common/urls"
)

func main() {
	// captchaId := captcha.NewLen(3)
	// captcha.Verify(captchaId, []byte())
	// captcha.WriteImage()

	config.InitConfig("./bbs-go.yaml")

	sm := sitemap.NewGenerator("https://file.mlog.club/", "sitemap-test", "sitemap-2020-01-17", func(sitemapLoc string) {
		fmt.Println("sitemap hahaha  ", sitemapLoc)
	})
	for i := 1; i <= 7; i++ {
		sm.AddURL(sitemap.URL{
			Loc:        urls.AbsUrl("/article/" + strconv.Itoa(i)),
			Lastmod:    time.Now(),
			Changefreq: sitemap.ChangefreqDaily,
			Priority:   "1.0",
		})
	}
	sm.Finalize()
}
