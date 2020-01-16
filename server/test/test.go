package main

import (
	"fmt"
	"time"

	"bbs-go/common/config"
	"bbs-go/common/sitemap"
)

func main() {
	// captchaId := captcha.NewLen(3)
	// captcha.Verify(captchaId, []byte())
	// captcha.WriteImage()

	config.InitConfig("./bbs-go.yaml")

	sm := sitemap.NewSitemap("https://file.mlog.club/", "sitemap-test", "fuck")
	sm.Add(sitemap.URL{
		Loc:        "/shit",
		Lastmod:    time.Now(),
		Changefreq: sitemap.ChangefreqDaily,
		Priority:   "1.0",
	})
	fmt.Println(sm.SitemapXml())
	sm.Write()
}
