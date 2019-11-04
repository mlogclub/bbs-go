package services

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
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
	return repositories.ThirdAccountRepository.Get(simple.DB(), id)
}

func (this *thirdAccountService) Take(where ...interface{}) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Take(simple.DB(), where...)
}

func (this *thirdAccountService) Find(cnd *simple.SqlCnd) []model.ThirdAccount {
	return repositories.ThirdAccountRepository.Find(simple.DB(), cnd)
}

func (this *thirdAccountService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.ThirdAccount) {
	cnd.FindOne(db, &ret)
	return
}

func (this *thirdAccountService) FindPageByParams(params *simple.QueryParams) (list []model.ThirdAccount, paging *simple.Paging) {
	return repositories.ThirdAccountRepository.FindPageByParams(simple.DB(), params)
}

func (this *thirdAccountService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.ThirdAccount, paging *simple.Paging) {
	return repositories.ThirdAccountRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *thirdAccountService) Create(t *model.ThirdAccount) error {
	return repositories.ThirdAccountRepository.Create(simple.DB(), t)
}

func (this *thirdAccountService) Update(t *model.ThirdAccount) error {
	return repositories.ThirdAccountRepository.Update(simple.DB(), t)
}

func (this *thirdAccountService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ThirdAccountRepository.Updates(simple.DB(), id, columns)
}

func (this *thirdAccountService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ThirdAccountRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *thirdAccountService) Delete(id int64) {
	repositories.ThirdAccountRepository.Delete(simple.DB(), id)
}

func (this *thirdAccountService) GetThirdAccount(thirdType string, thirdId string) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Take(simple.DB(), "third_type = ? and third_id = ?", thirdType, thirdId)
}

func (this *thirdAccountService) GetOrCreateByGithub(code, state string) (*model.ThirdAccount, error) {
	userInfo, err := github.GetUserInfoByCode(code, state)
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
		UserId:     sql.NullInt64{},
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

	account := this.GetThirdAccount(model.ThirdAccountTypeQQ, userInfo.Unionid)
	if account != nil {
		return account, nil
	}

	userInfoJson, _ := simple.FormatJson(userInfo)
	account = &model.ThirdAccount{
		UserId:     sql.NullInt64{},
		Avatar:     userInfo.FigureurlQQ1,
		Nickname:   strings.TrimSpace(userInfo.Nickname),
		ThirdType:  model.ThirdAccountTypeQQ,
		ThirdId:    userInfo.Unionid,
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
