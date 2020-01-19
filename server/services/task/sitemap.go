package task

import "bbs-go/services"

// 生成sitemap
func SitemapTask() {
	services.SitemapService.GenerateToday()
}
