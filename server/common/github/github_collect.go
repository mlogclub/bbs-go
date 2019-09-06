package github

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type CollectCallback func(repo *Repo)
type CollectPathCallback func(path string)

func Collect(callback CollectPathCallback) {
	for i := 1; i <= 100; i++ {
		paths, err := GetGithubRepos(i)
		if err != nil {
			logrus.Error(err)
		} else {
			for _, path := range paths {
				callback(path)
			}
		}
	}
}

func CollectRepo(callback CollectCallback) {
	Collect(func(path string) {
		repo, err := GetGithubRepo(path)
		if err != nil {
			logrus.Error(err)
		} else {
			callback(repo)
		}
	})
}

func GetGithubRepos(page int) ([]string, error) {
	rsp, err := resty.New().R().SetQueryParams(map[string]string{
		"p":    strconv.Itoa(page),
		"q":    "stars:>200 language:Go",
		"ref":  "advsearch",
		"type": "Repositories",
		"utf8": "✓",
	}).Get("https://github.com/search")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rsp.Body()))
	if err != nil {
		return nil, err
	}
	var ret []string
	doc.Find(".repo-list .repo-list-item h3 > a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")
		if err == nil {
			ret = append(ret, href)
		}
	})
	return ret, nil
}

func GetGithubRepo(path string) (*Repo, error) {
	repoJson, err := getGithubRepoByApi(path)
	if err != nil {
		return nil, err
	}
	messageRet := gjson.Get(repoJson, "message")
	if messageRet.Exists() {
		return nil, errors.New(messageRet.String())
	}
	branch := gjson.Get(repoJson, "default_branch").String()
	readme, err := getGithubRepoReadme(path, branch)
	if err != nil {
		return nil, err
	}
	return &Repo{
		Url:         gjson.Get(repoJson, "html_url").String(),
		Name:        gjson.Get(repoJson, "name").String(),
		FullName:    gjson.Get(repoJson, "full_name").String(),
		Description: gjson.Get(repoJson, "description").String(),
		Readme:      readme,
	}, nil
}

func getGithubRepoByApi(path string) (string, error) {
	url := "https://api.github.com/repos" + path
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return "", err
	}
	return string(rsp.Body()), nil
}

// 根据Path获取fullName
func GetFullnameByPath(path string) string {
	return path[0:]
}

// README
func getGithubRepoReadme(path, branch string) (string, error) {
	url := "https://raw.githubusercontent.com" + path + "/" + branch + "/README.md"
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return "", err
	}
	return string(rsp.Body()), nil
}

type Repo struct {
	Url         string
	Name        string
	FullName    string
	Description string
	Readme      string
}
