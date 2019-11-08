package baiduai

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-resty/resty/v2"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// AiData 要获取描述的数据
type AiData struct {
	Title string
	Desc  string
}

type AiTags struct {
	LogID     int64   `json:"log_id"`
	Items     []AiTag `json:"items"`
	ErrorMSG  string  `json:"error_msg"`
	ErrorCode int     `json:"error_code"`
}

type AiTag struct {
	Score float64 `json:"score"`
	Tag   string  `json:"tag"`
}

type AiTagParam struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type AiCategories struct {
	LogID     int64          `json:"log_id"`
	Item      AiCategoryItem `json:"item"`
	ErrorMSG  string         `json:"error_msg"`
	ErrorCode int            `json:"error_code"`
}

type AiCategoryItem struct {
	TopCategory    []AiTag `json:"lv1_tag_list"`
	SecondCatrgory []AiTag `json:"lv2_tag_list"`
}

type AiAnalyzeRet struct {
	Tags    []string
	Summary string
}

type Ai struct {
	ApiKey    string
	SecretKey string
}

var accessToken = ""                // accessToken
var accessTokenCreateTime int64 = 0 // accessToken创建时间

// 获取baidu api token 临时用
func (a *Ai) GetToken() string {
	durationMillis := simple.NowTimestamp() - accessTokenCreateTime
	if len(accessToken) == 0 || durationMillis > (86400*1000) { // accessToken为空或者生成时间超过一天
		c := NewClient(a.ApiKey, a.SecretKey)
		err := c.Auth()
		if err != nil {
			logrus.Error(err)
		}
		accessToken = c.AccessToken
		accessTokenCreateTime = simple.NowTimestamp()
	}
	return accessToken
}

func (a *Ai) GetTags(title, content string) *AiTags {
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

func (a *Ai) GetCategories(title, content string) *AiCategories {
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

func (a *Ai) GetNewsSummary(title, content string, maxSummaryLen int) (string, error) {
	if title == "" || content == "" {
		return "", errors.New("标题或内容为空")
	}
	if maxSummaryLen <= 0 {
		maxSummaryLen = 256
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

func (a *Ai) AnalyzeMarkdown(title, markdown string) (*AiAnalyzeRet, error) {
	mdResult := simple.NewMd().Run(markdown)
	return a.AnalyzeHtml(title, mdResult.ContentHtml)
}

func (a *Ai) AnalyzeHtml(title, html string) (*AiAnalyzeRet, error) {
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

func (a *Ai) AnalyzeText(title, text string) (*AiAnalyzeRet, error) {
	if title == "" || text == "" {
		return nil, errors.New("内容为空")
	}
	aiCategories := a.GetCategories(title, text)
	aiTags := a.GetTags(title, text)
	summary, _ := a.GetNewsSummary(title, text, 256)

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
