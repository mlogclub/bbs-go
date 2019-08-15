package services

import (
	"strconv"
	"time"

	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/mlog/utils/github"
)

var ProjectService = newProjectService()

type ProjectScanCallback func(project model.Project)

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

func (this *projectService) GetByFullName(fullname string) *model.Project {
	return this.Take("full_name = ?", fullname)
}

func (this *projectService) Scan(callback ProjectScanCallback) {
	var cursor int64
	for {
		list, err := repositories.ProjectRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("id > ?",
			cursor).Order("id asc").Size(300))
		if err != nil {
			break
		}
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		for _, project := range list {
			callback(project)
		}
	}
}

func (this *projectService) updateOrCreate(repo *github.Repo) *model.Project {
	project := this.GetByFullName(repo.FullName)
	if project == nil {
		project = &model.Project{}
		project.CreateTime = simple.NowTimestamp()
	}

	project.Name = repo.Name
	project.FullName = repo.FullName
	project.Url = repo.Url
	project.Description = repo.Description
	project.Readme = repo.Readme

	if project.Id > 0 {
		_ = this.Update(project)
	} else {
		_ = this.Create(project)
	}
	return project
}

// 开始采集
func (this *projectService) StartCollect() {
	// github.CollectRepo(func(repo *github.Repo) {
	// 	logrus.Info("采集项目：" + repo.Url)
	// 	project := this.updateOrCreate(repo)
	// 	this.syncToTopic(project)
	// })

	github.Collect(func(path string) {
		fullName := github.GetFullnameByPath(path)
		project := this.GetByFullName(fullName)
		if project != nil {
			logrus.Info("已采集仓库：" + path)
			return
		}
		repo, err := github.GetGithubRepo(path)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info("采集项目：" + repo.Url)
		project = this.updateOrCreate(repo)
		this.syncToTopic(project)

		time.Sleep(time.Minute * 3) // 睡一下，否则会限制api访问
	})
}

// 从帖子中同步过来
func (this *projectService) SyncFromTopic() {
	TopicService.Scan(func(topics []model.Topic) {
		for _, topic := range topics {
			logrus.Info("同步项目，topicId=" + strconv.FormatInt(topic.Id, 10))
			fullNameRet := gjson.Get(topic.ExtraData, "full_name")
			if !fullNameRet.Exists() {
				continue
			}
			repo, err := github.GetGithubRepo("/" + fullNameRet.String())
			if err != nil {
				continue
			} else {
				project := this.updateOrCreate(repo)

				{
					content := ""
					if len(project.Name) > 0 {
						content += "- 项目名称：" + project.Name + "\n"
					}
					if len(project.Url) > 0 {
						content += "- 项目地址：" + project.Url + "\n"
					}
					if len(project.Description) > 0 {
						content += "- 项目简介：" + project.Description + "\n\n\n"
					}
					content += project.Readme

					title := project.Name
					if len(project.Description) > 0 {
						title += ": " + project.Description
					}

					topic.Content = content
					topic.Title = title
					_ = TopicService.Update(&topic)
				}
			}
		}
	})
}

// 同步到帖子
func (this *projectService) SyncToTopic(userId int64) {
	this.Scan(func(project model.Project) {
		this.syncToTopic(&project)
	})
}

func (this *projectService) syncToTopic(project *model.Project) {
	logrus.Info("同步项目到帖子:" + project.Url)

	content := ""
	if len(project.Name) > 0 {
		content += "- 项目名称：" + project.Name + "\n"
	}
	if len(project.Url) > 0 {
		content += "- 项目地址：" + project.Url + "\n"
	}
	if len(project.Description) > 0 {
		content += "- 项目简介：" + project.Description + "\n\n\n"
	}
	content += project.Readme

	title := project.Name
	if len(project.Description) > 0 {
		title += ": " + project.Description
	}
	topic, err := TopicService.Publish(1, []string{"开源项目", "Go语言"}, title, content,
		simple.NewEmptyRspBuilder().Put("url", project.Url).Put("full_name", project.FullName).Build())
	if err != nil {
		logrus.Error(err)
	} else {
		err2 := this.UpdateColumn(project.Id, "topic_id", topic.Id)
		if err2 != nil {
			logrus.Error(err2)
		}
	}
}
