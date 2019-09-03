package services

import (
	"math"
	"path"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/common"
	"github.com/mlogclub/mlog/common/config"
	"github.com/mlogclub/mlog/common/urls"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/mlog/services/cache"
)

var ProjectService = newProjectService()

type ProjectScanCallback func(projects []model.Project) bool

func newProjectService() *projectService {
	return &projectService{}
}

type projectService struct {
}

func (this *projectService) Get(id int64) *model.Project {
	return repositories.ProjectRepository.Get(simple.GetDB(), id)
}

func (this *projectService) Take(where ...interface{}) *model.Project {
	return repositories.ProjectRepository.Take(simple.GetDB(), where...)
}

func (this *projectService) QueryCnd(cnd *simple.QueryCnd) (list []model.Project, err error) {
	return repositories.ProjectRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *projectService) Query(queries *simple.ParamQueries) (list []model.Project, paging *simple.Paging) {
	return repositories.ProjectRepository.Query(simple.GetDB(), queries)
}

func (this *projectService) Create(t *model.Project) error {
	return repositories.ProjectRepository.Create(simple.GetDB(), t)
}

func (this *projectService) Update(t *model.Project) error {
	return repositories.ProjectRepository.Update(simple.GetDB(), t)
}

func (this *projectService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ProjectRepository.Updates(simple.GetDB(), id, columns)
}

func (this *projectService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ProjectRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *projectService) Delete(id int64) {
	repositories.ProjectRepository.Delete(simple.GetDB(), id)
}

// 发布
func (this *projectService) Publish(userId int64, name, title, logo, url, docUrl, downloadUrl, contentType,
	content string) (*model.Project, error) {
	project := &model.Project{
		UserId:      userId,
		Name:        name,
		Title:       title,
		Logo:        logo,
		Url:         url,
		DocUrl:      docUrl,
		DownloadUrl: downloadUrl,
		ContentType: contentType,
		Content:     content,
		CreateTime:  simple.NowTimestamp(),
	}
	err := repositories.ProjectRepository.Create(simple.GetDB(), project)
	if err != nil {
		return nil, err
	}
	common.BaiduUrlPush([]string{urls.ProjectUrl(project.Id)})
	return project, nil
}

func (this *projectService) Scan(callback ProjectScanCallback) {
	var cursor int64
	for {
		list, err := repositories.ProjectRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ?",
			cursor).Order("id asc").Size(100))
		if err != nil {
			break
		}
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		if !callback(list) {
			break
		}
	}
}

func (this *projectService) ScanDesc(callback ProjectScanCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list, err := repositories.ProjectRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id < ?",
			cursor).Order("id desc").Size(100))
		if err != nil {
			break
		}
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		if !callback(list) {
			break
		}
	}
}

// rss
func (this *projectService) GenerateRss() {
	projects, err := repositories.ProjectRepository.QueryCnd(simple.GetDB(),
		simple.NewQueryCnd("1 = 1").Order("id desc").Size(2000))
	if err != nil {
		logrus.Error(err)
		return
	}

	var items []*feeds.Item

	for _, project := range projects {
		projectUrl := urls.ProjectUrl(project.Id)
		user := cache.UserCache.Get(project.UserId)
		if user == nil {
			continue
		}
		description := ""
		if project.ContentType == model.ContentTypeMarkdown {
			description = common.GetMarkdownSummary(project.Content)
		} else {
			description = common.GetHtmlSummary(project.Content)
		}
		item := &feeds.Item{
			Title:       project.Name + " - " + project.Title,
			Link:        &feeds.Link{Href: projectUrl},
			Description: description,
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email},
			Created:     simple.TimeFromTimestamp(project.CreateTime),
		}
		items = append(items, item)
	}
	siteTitle := cache.SysConfigCache.GetValue(model.SysConfigSiteTitle)
	siteDescription := cache.SysConfigCache.GetValue(model.SysConfigSiteDescription)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: config.Conf.BaseUrl},
		Description: siteDescription,
		Author:      &feeds.Author{Name: siteTitle},
		Created:     time.Now(),
		Items:       items,
	}
	atom, err := feed.ToAtom()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Conf.StaticPath, "project_atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Conf.StaticPath, "project_rss.xml"), rss, false)
	}
}
