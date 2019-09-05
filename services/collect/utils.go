package collect

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/mattn/godown"
	"github.com/mlogclub/simple"
	"github.com/sundy-li/html2article"

	"github.com/mlogclub/bbs-go/common/oss"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36"
)

type GetImageAbsoluteUrlFunc func(inputSrc string) string

// 自动采集
func Collect(articleUrl string, toMarkdown bool) (title, content string, err error) {
	response, err := resty.New().SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).R().Get(articleUrl)
	if err != nil {
		return
	}
	if response.StatusCode() != 200 {
		err = errors.New("Http Error " + strconv.Itoa(response.StatusCode()))
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body()))
	if err != nil {
		return
	}

	// 智能分析文章
	html, err := doc.Html()
	if err != nil {
		return
	}
	ext, err := html2article.NewFromHtml(html)
	if err != nil {
		return
	}
	article, err := ext.ToArticle()
	if err != nil {
		return
	}

	content, err = htmlContentImageReplace(article.Html, "src", func(inputSrc string) string {
		outputSrc, err := simple.AbsoluteURL(inputSrc, "", articleUrl)
		if err != nil {
			return ""
		}
		return outputSrc
	})
	if err != nil {
		return
	}

	if toMarkdown {
		content, err = html2md(content)
		if err != nil {
			return
		}
	}
	return article.Title, content, nil
}

// html2markdown
func html2md(html string) (string, error) {
	var buf bytes.Buffer
	err := godown.Convert(&buf, strings.NewReader(html), &godown.Option{
		Style: true,
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// html内容清洗
func htmlContentClean(ele *goquery.Selection) {
	if ele == nil {
		return
	}

	ele.Find("script").Remove()
	ele.Find("link").Remove()
	ele.Find("style").Remove()
	ele.Find("ins").Remove()

	// 生成的行号
	ele.Find("figure.highlight table td.gutter").Remove()
}

// 处理图片，并将图片转存到我们自己的存储服务
func htmlContentImageReplace(html string, srcAttr string, urlFunc GetImageAbsoluteUrlFunc) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html, err
	}

	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		handleImageSrc(selection, srcAttr, urlFunc)
	})
	return doc.Html()
}

// 处理图片，并将图片转存到我们自己的存储服务
func htmlSelectionImageReplace(ele *goquery.Selection, srcAttr string, urlFunc GetImageAbsoluteUrlFunc) {
	if ele == nil {
		return
	}

	ele.Find("img").Each(func(i int, selection *goquery.Selection) {
		handleImageSrc(selection, srcAttr, urlFunc)
	})
}

// 处理图片
func handleImageSrc(selection *goquery.Selection, srcAttr string, urlFunc GetImageAbsoluteUrlFunc) {
	if len(srcAttr) == 0 {
		srcAttr = "src"
	}

	src, exist := selection.Attr(srcAttr)
	if !exist || src == "" {
		return
	}
	src = urlFunc(src)
	output, err := oss.CopyImage(src)
	if err == nil {
		selection.SetAttr("src", output)
	} else {
		logrus.Error(err)
	}
}
