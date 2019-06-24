package services

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

type OauthClientService struct {
	OauthClientRepository *repositories.OauthClientRepository
}

func NewOauthClientService() *OauthClientService {
	return &OauthClientService{
		OauthClientRepository: repositories.NewOauthClientRepository(),
	}
}

func (this *OauthClientService) Get(id int64) *model.OauthClient {
	return this.OauthClientRepository.Get(simple.GetDB(), id)
}

func (this *OauthClientService) Take(where ...interface{}) *model.OauthClient {
	return this.OauthClientRepository.Take(simple.GetDB(), where...)
}

func (this *OauthClientService) QueryCnd(cnd *simple.QueryCnd) (list []model.OauthClient, err error) {
	return this.OauthClientRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *OauthClientService) Query(queries *simple.ParamQueries) (list []model.OauthClient, paging *simple.Paging) {
	return this.OauthClientRepository.Query(simple.GetDB(), queries)
}

func (this *OauthClientService) Create(t *model.OauthClient) error {
	return this.OauthClientRepository.Create(simple.GetDB(), t)
}

func (this *OauthClientService) Update(t *model.OauthClient) error {
	return this.OauthClientRepository.Update(simple.GetDB(), t)
}

func (this *OauthClientService) Updates(id int64, columns map[string]interface{}) error {
	return this.OauthClientRepository.Updates(simple.GetDB(), id, columns)
}

func (this *OauthClientService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.OauthClientRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *OauthClientService) Delete(id int64) {
	this.OauthClientRepository.Delete(simple.GetDB(), id)
}

func (this *OauthClientService) GetByClientId(clientId string) *model.OauthClient {
	return this.OauthClientRepository.GetByClientId(simple.GetDB(), clientId)
}
