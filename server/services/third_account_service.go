package services

import (
	"bbs-go/model/constants"
	"bbs-go/package/github"
	"bbs-go/package/qq"
	"database/sql"
	"github.com/mlogclub/simple/date"
	"github.com/mlogclub/simple/json"
	"strconv"
	"strings"

	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/repositories"
)

var ThirdAccountService = newThirdAccountService()

func newThirdAccountService() *thirdAccountService {
	return &thirdAccountService{}
}

type thirdAccountService struct {
}

func (s *thirdAccountService) Get(id int64) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Get(simple.DB(), id)
}

func (s *thirdAccountService) Take(where ...interface{}) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Take(simple.DB(), where...)
}

func (s *thirdAccountService) Find(cnd *simple.SqlCnd) []model.ThirdAccount {
	return repositories.ThirdAccountRepository.Find(simple.DB(), cnd)
}

func (s *thirdAccountService) FindOne(cnd *simple.SqlCnd) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.FindOne(simple.DB(), cnd)
}

func (s *thirdAccountService) FindPageByParams(params *simple.QueryParams) (list []model.ThirdAccount, paging *simple.Paging) {
	return repositories.ThirdAccountRepository.FindPageByParams(simple.DB(), params)
}

func (s *thirdAccountService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.ThirdAccount, paging *simple.Paging) {
	return repositories.ThirdAccountRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *thirdAccountService) Create(t *model.ThirdAccount) error {
	return repositories.ThirdAccountRepository.Create(simple.DB(), t)
}

func (s *thirdAccountService) Update(t *model.ThirdAccount) error {
	return repositories.ThirdAccountRepository.Update(simple.DB(), t)
}

func (s *thirdAccountService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ThirdAccountRepository.Updates(simple.DB(), id, columns)
}

func (s *thirdAccountService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ThirdAccountRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *thirdAccountService) Delete(id int64) {
	repositories.ThirdAccountRepository.Delete(simple.DB(), id)
}

func (s *thirdAccountService) GetThirdAccount(thirdType string, thirdId string) *model.ThirdAccount {
	return repositories.ThirdAccountRepository.Take(simple.DB(), "third_type = ? and third_id = ?", thirdType, thirdId)
}

func (s *thirdAccountService) GetOrCreateByGithub(code, state string) (*model.ThirdAccount, error) {
	userInfo, err := github.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetThirdAccount(constants.ThirdAccountTypeGithub, strconv.FormatInt(userInfo.Id, 10))
	if account != nil {
		return account, nil
	}

	nickname := userInfo.Login
	if len(userInfo.Name) > 0 {
		nickname = strings.TrimSpace(userInfo.Name)
	}

	userInfoJson, _ := json.ToStr(userInfo)
	account = &model.ThirdAccount{
		UserId:     sql.NullInt64{},
		Avatar:     userInfo.AvatarUrl,
		Nickname:   nickname,
		ThirdType:  constants.ThirdAccountTypeGithub,
		ThirdId:    strconv.FormatInt(userInfo.Id, 10),
		ExtraData:  userInfoJson,
		CreateTime: date.NowTimestamp(),
		UpdateTime: date.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *thirdAccountService) GetOrCreateByQQ(code, state string) (*model.ThirdAccount, error) {
	userInfo, err := qq.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetThirdAccount(constants.ThirdAccountTypeQQ, userInfo.Unionid)
	if account != nil {
		return account, nil
	}

	userInfoJson, _ := json.ToStr(userInfo)
	account = &model.ThirdAccount{
		UserId:     sql.NullInt64{},
		Avatar:     userInfo.FigureurlQQ1,
		Nickname:   strings.TrimSpace(userInfo.Nickname),
		ThirdType:  constants.ThirdAccountTypeQQ,
		ThirdId:    userInfo.Unionid,
		ExtraData:  userInfoJson,
		CreateTime: date.NowTimestamp(),
		UpdateTime: date.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
