package services

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

var OauthClientService = newOauthClientService()

func newOauthClientService() *oauthClientService {
	return &oauthClientService{
		OauthClientRepository: repositories.NewOauthClientRepository(),
	}
}

type oauthClientService struct {
	OauthClientRepository *repositories.OauthClientRepository
}

func (this *oauthClientService) Get(id int64) *model.OauthClient {
	return this.OauthClientRepository.Get(simple.GetDB(), id)
}

func (this *oauthClientService) Take(where ...interface{}) *model.OauthClient {
	return this.OauthClientRepository.Take(simple.GetDB(), where...)
}

func (this *oauthClientService) QueryCnd(cnd *simple.QueryCnd) (list []model.OauthClient, err error) {
	return this.OauthClientRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *oauthClientService) Query(queries *simple.ParamQueries) (list []model.OauthClient, paging *simple.Paging) {
	return this.OauthClientRepository.Query(simple.GetDB(), queries)
}

func (this *oauthClientService) Create(t *model.OauthClient) error {
	return this.OauthClientRepository.Create(simple.GetDB(), t)
}

func (this *oauthClientService) Update(t *model.OauthClient) error {
	return this.OauthClientRepository.Update(simple.GetDB(), t)
}

func (this *oauthClientService) Updates(id int64, columns map[string]interface{}) error {
	return this.OauthClientRepository.Updates(simple.GetDB(), id, columns)
}

func (this *oauthClientService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.OauthClientRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *oauthClientService) Delete(id int64) {
	this.OauthClientRepository.Delete(simple.GetDB(), id)
}

func (this *oauthClientService) GetByClientId(clientId string) *model.OauthClient {
	return this.OauthClientRepository.GetByClientId(simple.GetDB(), clientId)
}
