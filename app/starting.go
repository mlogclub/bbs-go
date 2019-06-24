package app

import (
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/utils"
)

func StartOn() {
	if utils.IsProd() {
		// 开启定时任务
		startSchedule()
		// 生成sitemap和rss
		generateSitemapAndRss()
	}
}

// 生成sitemap和rss
func generateSitemapAndRss() {
	go func() {
		articleService := services.NewArticleService()
		articleService.GenerateSitemap()
		articleService.GenerateRss()
	}()
}
