package collect

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/mattn/godown"
	"github.com/mlogclub/simple"
	"github.com/sundy-li/html2article"

	"github.com/mlogclub/mlog/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36"
)

type GetImageAbsoluteUrlFunc func(inputSrc string) string

// 自动采集
func Collect(articleUrl string, toMarkdown bool) (title, content string, err error) {
	response, err := resty.SetRedirectPolicy(resty.FlexibleRedirectPolicy(3)).R().Get(articleUrl)
	if err != nil {
		return "", "", err
	}
	if response.StatusCode() != 200 {
		return "", "", errors.New("Http Error " + strconv.Itoa(response.StatusCode()))
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body()))

	// 智能分析文章
	html, err := doc.Html()
	if err != nil {
		return "", "", err
	}
	ext, err := html2article.NewFromHtml(html)
	if err != nil {
		return "", "", err
	}
	article, err := ext.ToArticle()
	if err != nil {
		return "", "", err
	}

	content, err = handleImageSrc(article.Html, "src", func(inputSrc string) string {
		outputSrc, err := simple.AbsoluteURL(inputSrc, "", articleUrl)
		if err != nil {
			return ""
		}
		return outputSrc
	})
	if err != nil {
		return "", "", err
	}

	if toMarkdown {
		content, err = html2md(content)
		if err != nil {
			return "", "", err
		}
	}
	return article.Title, content, nil
}

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

func handleImageSrc(html string, srcAttr string, urlFunc GetImageAbsoluteUrlFunc) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	if len(srcAttr) == 0 {
		srcAttr = "src"
	}

	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		src, exist := selection.Attr(srcAttr)
		if !exist || src == "" {
			return
		}
		src = urlFunc(src)
		if len(src) > 0 {
			output, err := utils.AliyunOss.CopyImage(src)
			if err == nil {
				selection.SetAttr("src", output)
			} else {
				logrus.Error(err)
			}
		}
	})
	return doc.Html()
}
