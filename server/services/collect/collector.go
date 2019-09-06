package collect

import (
	"github.com/gocolly/colly"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
)

type Rule struct {
	EnterUrls                []string `json:"enterUrls" form:"enterUrls"`                               // 入口地址
	LinkCssQueries           []string `json:"linkCssQueries" form:"linkCssQueries"`                     // 链接选择器
	ArticleCssQuery          string   `json:"articleCssQuery" form:"articleCssQuery"`                   // 文章选择器
	ArticleTitleCssQuery     string   `json:"articleTitleCssQuery" form:"articleTitleCssQuery"`         // 文章标题选择器
	ArticleSummaryCssQuery   string   `json:"articleSummaryCssQuery" form:"articleSummaryCssQuery"`     // 文章摘要选择器
	ArticleContentCssQuery   string   `json:"articleContentCssQuery" form:"articleContentCssQuery"`     // 文章内容选择器
	ArticleContentImgSrcAttr string   `json:"articleContentImgSrcAttr" form:"articleContentImgSrcAttr"` // 图片属性，默认：src
	RemoveCssQueries         []string `json:"removeCssQueries" form:"removeCssQueries"`                 // 需要移除的元素选择器
	AutoPublish              bool     `json:"autoPublish" form:"autoPublish"`                           // 自动发布
	UserId                   int64    `json:"userId" form:"userId"`                                     // 使用哪个用户去发表文章
}

type MyCollector interface {
	Start(maxDepth int)
	ClawLink(c *colly.Collector)
	ClawArticle(c *colly.Collector)
}

func ParseRule(ruleStr string) (*Rule, error) {
	rule := &Rule{}
	err := simple.ParseJson(ruleStr, rule)
	return rule, err
}

func Start(ruleId int64, maxDepth int) {
	collector, err := NewDefaultCollector(ruleId)
	if err != nil {
		logrus.Error(err)
	}
	collector.Start(maxDepth)
}
