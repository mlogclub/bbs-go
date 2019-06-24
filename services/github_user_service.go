package services

import (
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
)

type GithubUserService struct {
	GithubUserRepository *repositories.GithubUserRepository
}

func NewGithubUserService() *GithubUserService {
	return &GithubUserService{
		GithubUserRepository: repositories.NewGithubUserRepository(),
	}
}

func (this *GithubUserService) Get(id int64) *model.GithubUser {
	return this.GithubUserRepository.Get(simple.GetDB(), id)
}

func (this *GithubUserService) Take(where ...interface{}) *model.GithubUser {
	return this.GithubUserRepository.Take(simple.GetDB(), where...)
}

func (this *GithubUserService) QueryCnd(cnd *simple.QueryCnd) (list []model.GithubUser, err error) {
	return this.GithubUserRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *GithubUserService) Query(queries *simple.ParamQueries) (list []model.GithubUser, paging *simple.Paging) {
	return this.GithubUserRepository.Query(simple.GetDB(), queries)
}

func (this *GithubUserService) Create(t *model.GithubUser) error {
	return this.GithubUserRepository.Create(simple.GetDB(), t)
}

func (this *GithubUserService) Update(t *model.GithubUser) error {
	return this.GithubUserRepository.Update(simple.GetDB(), t)
}

func (this *GithubUserService) Updates(id int64, columns map[string]interface{}) error {
	return this.GithubUserRepository.Updates(simple.GetDB(), id, columns)
}

func (this *GithubUserService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.GithubUserRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *GithubUserService) Delete(id int64) {
	this.GithubUserRepository.Delete(simple.GetDB(), id)
}

func (this *GithubUserService) GetByGithubId(githubId int64) *model.GithubUser {
	return this.Take("github_id = ?", githubId)
}
