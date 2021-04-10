package services

import (
	"bbs-go/model/constants"
	"bbs-go/package/config"
	"bbs-go/package/urls"
	"github.com/mlogclub/simple/date"
	"math"
	"path"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/package/common"
	"bbs-go/repositories"
)

var ProjectService = newProjectService()

func newProjectService() *projectService {
	return &projectService{}
}

type projectService struct {
}

func (s *projectService) Get(id int64) *model.Project {
	return repositories.ProjectRepository.Get(simple.DB(), id)
}

func (s *projectService) Take(where ...interface{}) *model.Project {
	return repositories.ProjectRepository.Take(simple.DB(), where...)
}

func (s *projectService) Find(cnd *simple.SqlCnd) []model.Project {
	return repositories.ProjectRepository.Find(simple.DB(), cnd)
}

func (s *projectService) FindOne(cnd *simple.SqlCnd) *model.Project {
	return repositories.ProjectRepository.FindOne(simple.DB(), cnd)
}

func (s *projectService) FindPageByParams(params *simple.QueryParams) (list []model.Project, paging *simple.Paging) {
	return repositories.ProjectRepository.FindPageByParams(simple.DB(), params)
}

func (s *projectService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Project, paging *simple.Paging) {
	return repositories.ProjectRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *projectService) Create(t *model.Project) error {
	return repositories.ProjectRepository.Create(simple.DB(), t)
}

func (s *projectService) Update(t *model.Project) error {
	return repositories.ProjectRepository.Update(simple.DB(), t)
}

func (s *projectService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ProjectRepository.Updates(simple.DB(), id, columns)
}

func (s *projectService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ProjectRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *projectService) Delete(id int64) {
	repositories.ProjectRepository.Delete(simple.DB(), id)
}

// 发布
func (s *projectService) Publish(userId int64, name, title, logo, url, docUrl, downloadUrl, contentType,
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
		CreateTime:  date.NowTimestamp(),
	}
	err := repositories.ProjectRepository.Create(simple.DB(), project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *projectService) ScanDesc(callback func(projects []model.Project)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.ProjectRepository.Find(simple.DB(), simple.NewSqlCnd().Lt("id", cursor).
			Desc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

func (s *projectService) ScanDescWithDate(dateFrom, dateTo int64, callback func(projects []model.Project)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.ProjectRepository.Find(simple.DB(), simple.NewSqlCnd().Lt("id", cursor).
			Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// rss
func (s *projectService) GenerateRss() {
	projects := repositories.ProjectRepository.Find(simple.DB(),
		simple.NewSqlCnd().Where("1 = 1").Desc("id").Limit(200))

	var items []*feeds.Item
	for _, project := range projects {
		projectUrl := urls.ProjectUrl(project.Id)
		user := cache.UserCache.Get(project.UserId)
		if user == nil {
			continue
		}
		description := common.GetSummary(project.ContentType, project.Content)
		item := &feeds.Item{
			Title:       project.Name + " - " + project.Title,
			Link:        &feeds.Link{Href: projectUrl},
			Description: description,
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email.String},
			Created:     date.FromTimestamp(project.CreateTime),
		}
		items = append(items, item)
	}
	siteTitle := cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
	siteDescription := cache.SysConfigCache.GetValue(constants.SysConfigSiteDescription)
	feed := &feeds.Feed{
		Title:       siteTitle,
		Link:        &feeds.Link{Href: config.Instance.BaseUrl},
		Description: siteDescription,
		Author:      &feeds.Author{Name: siteTitle},
		Created:     time.Now(),
		Items:       items,
	}
	atom, err := feed.ToAtom()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Instance.StaticPath, "project_atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = simple.WriteString(path.Join(config.Instance.StaticPath, "project_rss.xml"), rss, false)
	}
}
