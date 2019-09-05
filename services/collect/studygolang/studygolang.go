package studygolang

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common/oss"
)

type PageFunc func(url string)
type ProjectFunc func()

type Project struct {
	Name        string // 名称
	Title       string // 标题
	Content     string // 内容
	Logo        string
	Url         string // 项目主页
	DocUrl      string // 项目文档
	DownloadUrl string // 下载地址
}

// 采集项目分页
func GetStudyGoLangPage(page int, pageFunc PageFunc) {
	resp, err := resty.New().SetRedirectPolicy(resty.FlexibleRedirectPolicy(10)).R().Get("https://studygolang.com/projects?p=" + strconv.Itoa(page))
	if err != nil {
		logrus.Error(err)
		return
	}
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		logrus.Error(err)
		return
	}
	document.Find(".article .row h2 a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")
		if len(href) == 0 {
			return
		}
		url := "https://studygolang.com" + href
		pageFunc(url)
	})
}

// 采集项目
func GetStudyGolangProject(url string) *Project {
	resp, err := resty.New().SetRedirectPolicy(resty.FlexibleRedirectPolicy(10)).R().Get(url)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		logrus.Error(err)
		return nil
	}

	p := &Project{}

	p.Name = strings.TrimSpace(document.Find(".page .title h1 u").Remove().Text()) // name
	p.Content = document.Find(".project .desc").Text()                             // content
	p.Title = strings.TrimSpace(document.Find(".page .title h1").Text())           // title
	logo := document.Find(".page .title h1 img").AttrOr("src", "")                 // LOGO
	if logo != "" {
		p.Logo, _ = oss.CopyImage(logo)
	}
	document.Find("ul.urls li a").Each(func(i int, selection *goquery.Selection) {
		txt := selection.Text()
		if "项目首页" == txt {
			p.Url = selection.AttrOr("href", "")
		} else if "项目文档" == txt {
			p.DocUrl = selection.AttrOr("href", "")
		} else if "软件下载" == txt {
			p.DownloadUrl = selection.AttrOr("href", "")
		}
	})
	return p
}
