package collect

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
)

func NewDefaultCollector(ruleId int64) (MyCollector, error) {
	collectRule := services.CollectRuleService.Get(ruleId)
	if collectRule == nil {
		return nil, errors.New("没找到采集规则:" + strconv.FormatInt(ruleId, 10))
	}

	if collectRule.Status != model.CollectRuleStatusOk {
		return nil, errors.New("采集规则被禁用...ruleId=" + strconv.FormatInt(ruleId, 10))
	}

	rule, err := ParseRule(collectRule.Rule)
	if err != nil {
		return nil, err
	}

	return &DefaultCollector{
		collectRule: collectRule,
		rule:        rule,
	}, nil
}

type DefaultCollector struct {
	collectRule *model.CollectRule
	rule        *Rule
}

func (this *DefaultCollector) Start(maxDepth int) {
	c := this.newCollyCollector(maxDepth)

	this.ClawLink(c)
	this.ClawArticle(c)

	if len(this.rule.EnterUrls) > 0 {
		for _, url := range this.rule.EnterUrls {
			err := c.Visit(url)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

func (this *DefaultCollector) newCollyCollector(maxDepth int) *colly.Collector {
	c := colly.NewCollector(colly.DetectCharset(), colly.MaxDepth(maxDepth))

	// 随机UA
	extensions.RandomUserAgent(c)
	// 保留Referer
	extensions.Referer(c)

	c.OnRequest(func(request *colly.Request) {
		logrus.Infof("Collector request:" + request.URL.String())
	})
	c.OnError(func(response *colly.Response, e error) {
		logrus.Warn("Collector error:", e)
	})
	return c
}

// 采集链接
func (this *DefaultCollector) ClawLink(c *colly.Collector) {
	if len(this.rule.LinkCssQueries) > 0 {
		for _, query := range this.rule.LinkCssQueries {
			c.OnHTML(query, func(element *colly.HTMLElement) {
				href := element.Attr("href")
				href = element.Request.AbsoluteURL(href)
				err := c.Visit(href)
				if err != nil {
					logrus.Error(err)
				}
			})
		}
	}
}

// 采集文章
func (this *DefaultCollector) ClawArticle(c *colly.Collector) {
	c.OnHTML(this.rule.ArticleCssQuery, func(element *colly.HTMLElement) {
		var (
			sourceUrl string
			title     string
			summary   string
			content   string
		)

		// 原文链接
		sourceUrl = element.Request.URL.String()

		// 文章标题
		if len(this.rule.ArticleTitleCssQuery) > 0 {
			ele := element.DOM.Find(this.rule.ArticleTitleCssQuery).First()
			title = strings.TrimSpace(ele.Text())
		}

		if title == "" {
			logrus.Warn("未获取到文章标题，跳过..." + sourceUrl)
			return
		}

		// 判断文章是否存在
		if services.CollectArticleService.IsExists(sourceUrl, title) {
			logrus.Warn("文章已采集，跳过..." + sourceUrl)
			return
		}

		// 文章摘要
		if len(this.rule.ArticleSummaryCssQuery) > 0 {
			ele := element.DOM.Find(this.rule.ArticleSummaryCssQuery).First()
			summary = strings.TrimSpace(ele.Text())
		}

		// 文章内容
		if len(this.rule.ArticleContentCssQuery) > 0 {
			// 清理不必要的数据
			if len(this.rule.RemoveCssQueries) > 0 {
				for _, cssQuery := range this.rule.RemoveCssQueries {
					element.DOM.Find(cssQuery).Remove()
				}
			}

			ele := element.DOM.Find(this.rule.ArticleContentCssQuery).First()

			htmlContentClean(ele)
			htmlSelectionImageReplace(ele, this.rule.ArticleContentImgSrcAttr, func(inputSrc string) string {
				return element.Request.AbsoluteURL(inputSrc)
			})

			content, _ = ele.Html()
			content = strings.TrimSpace(content)
		}

		if content == "" {
			logrus.Warn("未获取到文章内容，跳过..." + sourceUrl)
			return
		}

		logrus.Info("采集文章：" + element.Request.URL.String())

		ca, err := services.CollectArticleService.Create(this.collectRule.Id, this.rule.UserId, sourceUrl, title, summary, content)
		if err != nil {
			logrus.Error(err)
		} else if this.rule.AutoPublish {
			if this.rule.UserId <= 0 {
				logrus.Warn("没配置文章发布作者，无法发布...ruleId=" + strconv.FormatInt(this.collectRule.Id, 10))
			} else {
				logrus.Info("自动发布..." + strconv.FormatInt(ca.Id, 10))
				if err := services.CollectArticleService.Publish(ca.Id); err != nil {
					logrus.Error(err)
				}
			}
		}
	})
}
