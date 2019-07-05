package services

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

var OauthTokenService = newOauthTokenService()

func newOauthTokenService() *oauthTokenService {
	return &oauthTokenService{
		OauthTokenRepository: repositories.NewOauthTokenRepository(),
	}
}

type oauthTokenService struct {
	OauthTokenRepository *repositories.OauthTokenRepository
}

func (this *oauthTokenService) Get(id int64) *model.OauthToken {
	return this.OauthTokenRepository.Get(simple.GetDB(), id)
}

func (this *oauthTokenService) Take(where ...interface{}) *model.OauthToken {
	return this.OauthTokenRepository.Take(simple.GetDB(), where...)
}

func (this *oauthTokenService) QueryCnd(cnd *simple.QueryCnd) (list []model.OauthToken, err error) {
	return this.OauthTokenRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *oauthTokenService) Query(queries *simple.ParamQueries) (list []model.OauthToken, paging *simple.Paging) {
	return this.OauthTokenRepository.Query(simple.GetDB(), queries)
}

func (this *oauthTokenService) Create(t *model.OauthToken) error {
	return this.OauthTokenRepository.Create(simple.GetDB(), t)
}

func (this *oauthTokenService) Update(t *model.OauthToken) error {
	return this.OauthTokenRepository.Update(simple.GetDB(), t)
}

func (this *oauthTokenService) Updates(id int64, columns map[string]interface{}) error {
	return this.OauthTokenRepository.Updates(simple.GetDB(), id, columns)
}

func (this *oauthTokenService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.OauthTokenRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *oauthTokenService) Delete(id int64) {
	this.OauthTokenRepository.Delete(simple.GetDB(), id)
}
