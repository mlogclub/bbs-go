package services

import (
	context2 "context"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/mlog/utils/github"
	"github.com/mlogclub/simple"
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
	return this.GithubUserRepository.GetByGithubId(simple.GetDB(), githubId)
}

func (this *GithubUserService) GetGithubUser(code string) (*model.GithubUser, error) {
	token, err := github.OauthConfig.Exchange(context2.TODO(), code)
	if err != nil {
		return nil, err
	}

	third, _ := github.GetUserInfo(token.AccessToken)
	githubUser := this.GetByGithubId(third.Id)

	if githubUser != nil {
		return githubUser, nil
	}
	githubUser = &model.GithubUser{
		GithubId:   third.Id,
		Login:      third.Login,
		NodeId:     third.NodeId,
		AvatarUrl:  third.AvatarUrl,
		Url:        third.Url,
		HtmlUrl:    third.HtmlUrl,
		Email:      third.Email,
		Name:       third.Name,
		CreateTime: simple.NowTimestamp(),
		UpdateTime: simple.NowTimestamp(),
	}

	err = this.Create(githubUser)
	if err != nil {
		return nil, err
	}
	return githubUser, nil
}
