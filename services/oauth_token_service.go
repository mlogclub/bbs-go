package services

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

type OauthTokenService struct {
	OauthTokenRepository *repositories.OauthTokenRepository
}

func NewOauthTokenService() *OauthTokenService {
	return &OauthTokenService{
		OauthTokenRepository: repositories.NewOauthTokenRepository(),
	}
}

func (this *OauthTokenService) Get(id int64) *model.OauthToken {
	return this.OauthTokenRepository.Get(simple.GetDB(), id)
}

func (this *OauthTokenService) Take(where ...interface{}) *model.OauthToken {
	return this.OauthTokenRepository.Take(simple.GetDB(), where...)
}

func (this *OauthTokenService) QueryCnd(cnd *simple.QueryCnd) (list []model.OauthToken, err error) {
	return this.OauthTokenRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *OauthTokenService) Query(queries *simple.ParamQueries) (list []model.OauthToken, paging *simple.Paging) {
	return this.OauthTokenRepository.Query(simple.GetDB(), queries)
}

func (this *OauthTokenService) Create(t *model.OauthToken) error {
	return this.OauthTokenRepository.Create(simple.GetDB(), t)
}

func (this *OauthTokenService) Update(t *model.OauthToken) error {
	return this.OauthTokenRepository.Update(simple.GetDB(), t)
}

func (this *OauthTokenService) Updates(id int64, columns map[string]interface{}) error {
	return this.OauthTokenRepository.Updates(simple.GetDB(), id, columns)
}

func (this *OauthTokenService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.OauthTokenRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *OauthTokenService) Delete(id int64) {
	this.OauthTokenRepository.Delete(simple.GetDB(), id)
}
