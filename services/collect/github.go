package collect

import (
	"bytes"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
)

type GithubCollectCallback func(repo *GithubRepo)

func CollectGithub(callback GithubCollectCallback) {
	for i := 1; i <= 100; i++ {
		paths, err := GetGithubRepos(i)
		if err != nil {
			logrus.Error(err)
			continue
		}
		for _, path := range paths {
			repo, err := GetGithubRepo(path)
			if err != nil {
				logrus.Error(err)
				continue
			}
			callback(repo)
		}
	}
}

func GetGithubRepos(page int) ([]string, error) {
	rsp, err := resty.R().SetQueryParams(map[string]string{
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

func GetGithubRepo(path string) (*GithubRepo, error) {
	repoJson, err := getGithubRepoByApi(path)
	if err != nil {
		return nil, err
	}
	branch := gjson.Get(repoJson, "default_branch").String()
	readme, err := getGithubRepoReadme(path, branch)
	if err != nil {
		return nil, err
	}
	return &GithubRepo{
		Url:         gjson.Get(repoJson, "html_url").String(),
		Name:        gjson.Get(repoJson, "name").String(),
		FullName:    gjson.Get(repoJson, "full_name").String(),
		Description: gjson.Get(repoJson, "description").String(),
		Readme:      readme,
	}, nil
}

func getGithubRepoByApi(path string) (string, error) {
	url := "https://api.github.com/repos" + path
	rsp, err := resty.R().Get(url)
	if err != nil {
		return "", err
	}
	return string(rsp.Body()), nil
}

// README
func getGithubRepoReadme(path, branch string) (string, error) {
	url := "https://raw.githubusercontent.com" + path + "/" + branch + "/README.md"
	rsp, err := resty.R().Get(url)
	if err != nil {
		return "", err
	}
	return string(rsp.Body()), nil
}

type GithubRepo struct {
	Url         string
	Name        string
	FullName    string
	Description string
	Readme      string
}
