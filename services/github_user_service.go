package services

import (
	context2 "context"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/mlog/utils/github"
	"github.com/mlogclub/simple"
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

func (this *githubUserService) GetGithubUser(code string) (*model.GithubUser, error) {
	token, err := github.GetOauthConfig(nil).Exchange(context2.TODO(), code)
	if err != nil {
		return nil, err
	}

	third, err := github.GetUserInfo(token.AccessToken)
	if err != nil {
		return nil, err
	}

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
		Bio:        third.Bio,
		Company:    third.Company,
		Blog:       third.Blog,
		Location:   third.Location,
		CreateTime: simple.NowTimestamp(),
		UpdateTime: simple.NowTimestamp(),
	}

	err = this.Create(githubUser)
	if err != nil {
		return nil, err
	}
	return githubUser, nil
}
