package services

import (
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var GithubUserService = newGithubUserService()

func newGithubUserService() *githubUserService {
	return &githubUserService{}
}

type githubUserService struct {
}

func (this *githubUserService) Get(id int64) *model.GithubUser {
	return repositories.GithubUserRepository.Get(simple.GetDB(), id)
}

func (this *githubUserService) Take(where ...interface{}) *model.GithubUser {
	return repositories.GithubUserRepository.Take(simple.GetDB(), where...)
}

func (this *githubUserService) QueryCnd(cnd *simple.QueryCnd) (list []model.GithubUser, err error) {
	return repositories.GithubUserRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *githubUserService) Query(queries *simple.ParamQueries) (list []model.GithubUser, paging *simple.Paging) {
	return repositories.GithubUserRepository.Query(simple.GetDB(), queries)
}

func (this *githubUserService) Create(t *model.GithubUser) error {
	return repositories.GithubUserRepository.Create(simple.GetDB(), t)
}

func (this *githubUserService) Update(t *model.GithubUser) error {
	return repositories.GithubUserRepository.Update(simple.GetDB(), t)
}

func (this *githubUserService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.GithubUserRepository.Updates(simple.GetDB(), id, columns)
}

func (this *githubUserService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.GithubUserRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *githubUserService) Delete(id int64) {
	repositories.GithubUserRepository.Delete(simple.GetDB(), id)
}

func (this *githubUserService) GetByGithubId(githubId int64) *model.GithubUser {
	return repositories.GithubUserRepository.GetByGithubId(simple.GetDB(), githubId)
}
