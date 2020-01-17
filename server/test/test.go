package main

import (
	"strconv"
	"time"

	"bbs-go/common/config"
	"bbs-go/common/sitemap"
)

func main() {
	// captchaId := captcha.NewLen(3)
	// captcha.Verify(captchaId, []byte())
	// captcha.WriteImage()

	config.InitConfig("./bbs-go.yaml")

	sm := sitemap.NewSitemap("https://file.mlog.club/", "sitemap-test", "sitemap-1")
	for i := 1; i <= 10000; i++ {
		sm.Add(sitemap.URL{
			Loc:        "/article/" + strconv.Itoa(i),
			Lastmod:    time.Now(),
			Changefreq: sitemap.ChangefreqDaily,
			Priority:   "1.0",
		})
	}
	sm.Write()
}
