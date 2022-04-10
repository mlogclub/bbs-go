package services

import (
	"bbs-go/model/constants"
	"bbs-go/pkg/bbsurls"
	"bbs-go/pkg/config"
	"math"
	"path"
	"time"

	"github.com/gorilla/feeds"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/files"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"github.com/sirupsen/logrus"

	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/pkg/common"
	"bbs-go/repositories"
)

var ProjectService = newProjectService()

func newProjectService() *projectService {
	return &projectService{}
}

type projectService struct {
}

func (s *projectService) Get(id int64) *model.Project {
	return repositories.ProjectRepository.Get(sqls.DB(), id)
}

func (s *projectService) Take(where ...interface{}) *model.Project {
	return repositories.ProjectRepository.Take(sqls.DB(), where...)
}

func (s *projectService) Find(cnd *sqls.Cnd) []model.Project {
	return repositories.ProjectRepository.Find(sqls.DB(), cnd)
}

func (s *projectService) FindOne(cnd *sqls.Cnd) *model.Project {
	return repositories.ProjectRepository.FindOne(sqls.DB(), cnd)
}

func (s *projectService) FindPageByParams(params *params.QueryParams) (list []model.Project, paging *sqls.Paging) {
	return repositories.ProjectRepository.FindPageByParams(sqls.DB(), params)
}

func (s *projectService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Project, paging *sqls.Paging) {
	return repositories.ProjectRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *projectService) Create(t *model.Project) error {
	return repositories.ProjectRepository.Create(sqls.DB(), t)
}

func (s *projectService) Update(t *model.Project) error {
	return repositories.ProjectRepository.Update(sqls.DB(), t)
}

func (s *projectService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ProjectRepository.Updates(sqls.DB(), id, columns)
}

func (s *projectService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ProjectRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *projectService) Delete(id int64) {
	repositories.ProjectRepository.Delete(sqls.DB(), id)
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
		CreateTime:  dates.NowTimestamp(),
	}
	err := repositories.ProjectRepository.Create(sqls.DB(), project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *projectService) ScanDesc(callback func(projects []model.Project)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.ProjectRepository.Find(sqls.DB(), sqls.NewCnd().Lt("id", cursor).
			Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

func (s *projectService) ScanDescWithDate(dateFrom, dateTo int64, callback func(projects []model.Project)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.ProjectRepository.Find(sqls.DB(), sqls.NewCnd().Lt("id", cursor).
			Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// rss
func (s *projectService) GenerateRss() {
	projects := repositories.ProjectRepository.Find(sqls.DB(),
		sqls.NewCnd().Where("1 = 1").Desc("id").Limit(200))

	var items []*feeds.Item
	for _, project := range projects {
		projectUrl := bbsurls.ProjectUrl(project.Id)
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
			Created:     dates.FromTimestamp(project.CreateTime),
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
		_ = files.WriteString(path.Join(config.Instance.StaticPath, "project_atom.xml"), atom, false)
	}

	rss, err := feed.ToRss()
	if err != nil {
		logrus.Error(err)
	} else {
		_ = files.WriteString(path.Join(config.Instance.StaticPath, "project_rss.xml"), rss, false)
	}
}
