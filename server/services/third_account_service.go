package services

import (
	"strconv"
	"strings"

	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/common/github"
	"github.com/mlogclub/bbs-go/common/qq"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var ThirdAccountService = newThirdAccountService()

func newThirdAccountService() *thirdAccountService {
	return &thirdAccountService{}
}

type thirdAccountService struct {
}

func (this *thirdAccountService) Get(id int64) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Get(simple.GetDB(), id)
}

func (this *thirdAccountService) Take(where ...interface{}) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Take(simple.GetDB(), where...)
}

func (this *thirdAccountService) QueryCnd(cnd *simple.QueryCnd) (list []model.ThirdAccount, err error) {
	return repositories.ThirdAccountRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *thirdAccountService) Query(queries *simple.ParamQueries) (list []model.ThirdAccount, paging *simple.Paging) {
	return repositories.ThirdAccountRepository.Query(simple.GetDB(), queries)
}

func (this *thirdAccountService) Create(t *model.ThirdAccount) error {
	return repositories.ThirdAccountRepository.Create(simple.GetDB(), t)
}

func (this *thirdAccountService) Update(t *model.ThirdAccount) error {
	return repositories.ThirdAccountRepository.Update(simple.GetDB(), t)
}

func (this *thirdAccountService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ThirdAccountRepository.Updates(simple.GetDB(), id, columns)
}

func (this *thirdAccountService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ThirdAccountRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *thirdAccountService) Delete(id int64) {
	repositories.ThirdAccountRepository.Delete(simple.GetDB(), id)
}

func (this *thirdAccountService) GetThirdAccount(thirdType string, thirdId string) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Take(simple.GetDB(), "third_type = ? and third_id = ?", thirdType, thirdId)
}

func (this *thirdAccountService) GetOrCreateByGithub(code string) (*model.ThirdAccount, error) {
	userInfo, err := github.GetUserInfoByCode(code)
	if err != nil {
		return nil, err
	}

	account := this.GetThirdAccount(model.ThirdAccountTypeGithub, strconv.FormatInt(userInfo.Id, 10))
	if account != nil {
		return account, nil
	}

	nickname := userInfo.Login
	if len(userInfo.Name) > 0 {
		nickname = strings.TrimSpace(userInfo.Name)
	}

	userInfoJson, _ := simple.FormatJson(userInfo)
	account = &model.ThirdAccount{
		Avatar:     userInfo.AvatarUrl,
		Nickname:   nickname,
		ThirdType:  model.ThirdAccountTypeGithub,
		ThirdId:    strconv.FormatInt(userInfo.Id, 10),
		ExtraData:  userInfoJson,
		CreateTime: simple.NowTimestamp(),
		UpdateTime: simple.NowTimestamp(),
	}
	err = this.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (this *thirdAccountService) GetOrCreateByQQ(code, state string) (*model.ThirdAccount, error) {
	userInfo, err := qq.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	thirdId := userInfo.UnionId
	if len(thirdId) == 0 {
		thirdId = userInfo.OpenId
	}
	account := this.GetThirdAccount(model.ThirdAccountTypeQQ, thirdId)
	if account != nil {
		return account, nil
	}

	userInfoJson, _ := simple.FormatJson(userInfo)
	account = &model.ThirdAccount{
		Avatar:     userInfo.FigureurlQQ1,
		Nickname:   strings.TrimSpace(userInfo.Nickname),
		ThirdType:  model.ThirdAccountTypeQQ,
		ThirdId:    thirdId,
		ExtraData:  userInfoJson,
		CreateTime: simple.NowTimestamp(),
		UpdateTime: simple.NowTimestamp(),
	}
	err = this.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
