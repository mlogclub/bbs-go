package services

import (
	"math"
	"path"
	"time"

	"github.com/gorilla/feeds"
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/common/config"
	"github.com/mlogclub/bbs-go/common/urls"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/bbs-go/services/cache"
)

var ProjectService = newProjectService()

type ProjectScanCallback func(projects []model.Project) bool

func newProjectService() *projectService {
	return &projectService{}
}

type projectService struct {
}

func (this *projectService) Get(id int64) *model.Project {
	return repositories.ProjectRepository.Get(simple.DB(), id)
}

func (this *projectService) Take(where ...interface{}) *model.Project {
	return repositories.ProjectRepository.Take(simple.DB(), where...)
}

func (this *projectService) Find(cnd *simple.SqlCnd) []model.Project {
	return repositories.ProjectRepository.Find(simple.DB(), cnd)
}

func (this *projectService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Project) {
	cnd.FindOne(db, &ret)
	return
}

func (this *projectService) FindPageByParams(params *simple.QueryParams) (list []model.Project, paging *simple.Paging) {
	return repositories.ProjectRepository.FindPageByParams(simple.DB(), params)
}

func (this *projectService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Project, paging *simple.Paging) {
	return repositories.ProjectRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *projectService) Create(t *model.Project) error {
	return repositories.ProjectRepository.Create(simple.DB(), t)
}

func (this *projectService) Update(t *model.Project) error {
	return repositories.ProjectRepository.Update(simple.DB(), t)
}

func (this *projectService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ProjectRepository.Updates(simple.DB(), id, columns)
}

func (this *projectService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ProjectRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *projectService) Delete(id int64) {
	repositories.ProjectRepository.Delete(simple.DB(), id)
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
	err := repositories.ProjectRepository.Create(simple.DB(), project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (this *projectService) Scan(callback ProjectScanCallback) {
	var cursor int64
	for {
		list := repositories.ProjectRepository.Find(simple.DB(), simple.NewSqlCnd().Where("id > ?",
			cursor).Asc("id").Limit(100))
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
		list := repositories.ProjectRepository.Find(simple.DB(), simple.NewSqlCnd().Where("id < ?",
			cursor).Desc("id").Limit(100))
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
	projects := repositories.ProjectRepository.Find(simple.DB(),
		simple.NewSqlCnd().Where("1 = 1").Desc("id").Limit(2000))

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
			Author:      &feeds.Author{Name: user.Avatar, Email: user.Email.String},
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
