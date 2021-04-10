package baiduai

import (
	"bbs-go/model/constants"
	"bbs-go/package/config"
	"bbs-go/package/markdown"
	"encoding/json"
	"errors"
	"github.com/mlogclub/simple/date"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-resty/resty/v2"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type ai struct {
	ApiKey    string
	SecretKey string

	accessToken           string // accessToken
	accessTokenCreateTime int64  // accessToken创建时间
}

var once sync.Once
var instance *ai

func GetAi() *ai {
	once.Do(func() {
		instance = &ai{
			ApiKey:    config.Instance.BaiduAi.ApiKey,
			SecretKey: config.Instance.BaiduAi.SecretKey,
		}
	})
	return instance
}

// 获取baidu api token 临时用
func (a *ai) GetToken() string {
	durationMillis := date.NowTimestamp() - a.accessTokenCreateTime
	if len(a.accessToken) == 0 || durationMillis > (86400*1000) { // accessToken为空或者生成时间超过一天
		c := NewClient(a.ApiKey, a.SecretKey)
		err := c.Auth()
		if err != nil {
			logrus.Error(err)
		}
		a.accessToken = c.AccessToken
		a.accessTokenCreateTime = date.NowTimestamp()
	}
	return a.accessToken
}

func (a *ai) GetTags(title, content string) *AiTags {
	if title == "" || content == "" {
		return nil
	}
	data := make(map[string]interface{})
	data["title"] = title
	data["content"] = simple.Substr(content, 0, 10000)

	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	url := "https://aip.baidubce.com/rpc/2.0/nlp/v1/keyword?charset=UTF-8&access_token=" + a.GetToken()
	response, err := resty.New().R().SetBody(string(bytesData)).Post(url)
	if err != nil {
		return nil
	}

	tags := &AiTags{}
	err = json.Unmarshal(response.Body(), tags)
	if err != nil {
		return nil
	}
	return tags
}

func (a *ai) GetCategories(title, content string) *AiCategories {
	if title == "" || content == "" {
		return nil
	}

	data := make(map[string]interface{})
	data["title"] = title
	data["content"] = simple.Substr(content, 0, 10000)

	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	url := "https://aip.baidubce.com/rpc/2.0/nlp/v1/topic?charset=UTF-8&access_token=" + a.GetToken()
	response, err := resty.New().R().SetBody(string(bytesData)).Post(url)
	if err != nil {
		return nil
	}

	categories := &AiCategories{}
	err = json.Unmarshal(response.Body(), categories)
	if err != nil {
		return nil
	}
	return categories
}

func (a *ai) GetNewsSummary(title, content string, maxSummaryLen int) (string, error) {
	if title == "" || content == "" {
		return "", errors.New("标题或内容为空")
	}
	if maxSummaryLen <= 0 {
		maxSummaryLen = constants.SummaryLen
	}

	data := make(map[string]interface{})
	data["title"] = title
	data["content"] = simple.Substr(content, 0, 3000)
	data["max_summary_len"] = maxSummaryLen

	bytesData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	url := "https://aip.baidubce.com/rpc/2.0/nlp/v1/news_summary?charset=UTF-8&access_token=" + a.GetToken()
	response, err := resty.New().R().SetBody(string(bytesData)).Post(url)
	if err != nil {
		return "", err
	}
	ret := gjson.Get(string(response.Body()), "summary")
	return ret.String(), nil
}

func (a *ai) AnalyzeMarkdown(title, markdownStr string) (*AiAnalyzeRet, error) {
	content := markdown.ToHTML(markdownStr)
	return a.AnalyzeHtml(title, content)
}

func (a *ai) AnalyzeHtml(title, html string) (*AiAnalyzeRet, error) {
	if title == "" || html == "" {
		return nil, errors.New("内容为空")
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	text := doc.Text()
	return a.AnalyzeText(title, text)
}

func (a *ai) AnalyzeText(title, text string) (*AiAnalyzeRet, error) {
	if title == "" || text == "" {
		return nil, errors.New("内容为空")
	}
	aiCategories := a.GetCategories(title, text)
	aiTags := a.GetTags(title, text)
	summary, _ := a.GetNewsSummary(title, text, constants.SummaryLen)

	set := hashset.New()
	if aiCategories != nil {
		if len(aiCategories.Item.TopCategory) > 0 {
			for _, v := range aiCategories.Item.TopCategory {
				set.Add(v.Tag)
			}
		}
		if len(aiCategories.Item.SecondCatrgory) > 0 {
			for _, v := range aiCategories.Item.SecondCatrgory {
				set.Add(v.Tag)
			}
		}
	}
	if aiTags != nil && len(aiTags.Items) > 0 {
		for _, v := range aiTags.Items {
			set.Add(v.Tag)
		}
	}

	var tags []string
	for _, v := range set.Values() {
		tags = append(tags, v.(string))
	}
	return &AiAnalyzeRet{
		Tags:    tags,
		Summary: summary,
	}, nil
}
