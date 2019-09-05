package oschina

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common/oss"
	"github.com/mlogclub/bbs-go/model"
)

func GetPage(page int) (urls []string) {
	document, err := getDocument("https://www.oschina.net/project/widgets/_project_list?company=0&tag=0&lang=358&os=0&sort=time&recommend=false&cn=false&weekly=false&type=ajax&p=" + strconv.Itoa(page))
	if err != nil {
		logrus.Error(err)
		return
	}
	document.Find(".items .item .header a").Each(func(i int, selection *goquery.Selection) {
		url := selection.AttrOr("href", "")
		if url != "" {
			urls = append(urls, url)
		}
	})
	return
}

func GetProject(url string) *model.Project {
	document, err := getDocument(url)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	document.Find(".ad-wrap").Remove()

	project := &model.Project{
		ContentType: model.ContentTypeHtml,
	}
	project.Name = strings.TrimSpace(document.Find(".row .article-detail h1.header .project-name").Text())
	project.Title = strings.TrimSpace(document.Find(".row .article-detail h1.header .project-title").Text())

	// 处理链接地址
	document.Find(".row .project-detail-container .article-detail .related-links a").Each(func(i int, selection *goquery.Selection) {
		text := strings.TrimSpace(selection.Text())
		if text == "软件首页" {
			project.Url = selection.AttrOr("href", "")
		} else if text == "软件文档" {
			project.DocUrl = selection.AttrOr("href", "")
		}
	})

	// 处理内容
	contentDom := document.Find(".row .project-detail-container .article-detail .content")
	contentDom.Find("img").Each(func(i int, selection *goquery.Selection) {
		src, exists := selection.Attr("src")
		if !exists || src == "" {
			return
		}

		output, err := oss.CopyImage(src)
		if err == nil {
			selection.SetAttr("src", output)
		} else {
			logrus.Error(err)
		}
	})
	tempContent, err := contentDom.Html()
	if err == nil {
		project.Content = strings.TrimSpace(tempContent)
	} else {
		logrus.Error(err)
		return nil
	}

	// 处理Logo
	logoDom := document.Find(".row .article-detail .logo-wrap img")
	if logoDom.Size() > 0 {
		logo := logoDom.AttrOr("src", "")
		if logo != "" {
			project.Logo, _ = oss.CopyImage(logo)
		}
	}

	return project
}

func getDocument(url string) (*goquery.Document, error) {
	resp, err := resty.New().R().SetHeaders(map[string]string{
		"Cookie": "Hm_lvt_a411c4d1664dd70048ee98afe7b28f0b=1565313534,1565608243,1566022640; _user_behavior_=6059fc2a-40ec-446e-b492-6d3c9eeb04dd; oscid=p%2BjynSY2nk%2Fs%2FA%2FnkvANsEIXBXLR92Dx%2BnZPLIVxeWDfmQy%2BhYqwffw%2FlmcP7mspcZBARBNn9kmcqq2I0pV2nZdcH7QtN%2FvvAQPiY%2BJEGB3qRtga6Q18iubUfBRGFE6kvO8LoBs8HJcahAwSwXtBXw%3D%3D; banner_osc_scr_0729=1; aliyungf_tc=AQAAAGmzSQiz7AUATby9r1NGxLJQOrwz; Hm_lpvt_a411c4d1664dd70048ee98afe7b28f0b=1566022640",
	}).Get(url)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
}
