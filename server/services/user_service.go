package services

import (
	"database/sql"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/core/errors"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/services/cache"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var UserService = newUserService()

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (this *userService) Get(id int64) *model.User {
	return repositories.UserRepository.Get(simple.GetDB(), id)
}

func (this *userService) Take(where ...interface{}) *model.User {
	return repositories.UserRepository.Take(simple.GetDB(), where...)
}

func (this *userService) QueryCnd(cnd *simple.QueryCnd) (list []model.User, err error) {
	return repositories.UserRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *userService) Query(queries *simple.ParamQueries) (list []model.User, paging *simple.Paging) {
	return repositories.UserRepository.Query(simple.GetDB(), queries)
}

func (this *userService) Create(t *model.User) error {
	err := repositories.UserRepository.Create(simple.GetDB(), t)
	if err == nil {
		cache.UserCache.Invalidate(t.Id)
	}
	return nil
}

func (this *userService) Update(t *model.User) error {
	err := repositories.UserRepository.Update(simple.GetDB(), t)
	cache.UserCache.Invalidate(t.Id)
	return err
}

func (this *userService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.UserRepository.Updates(simple.GetDB(), id, columns)
	cache.UserCache.Invalidate(id)
	return err
}

func (this *userService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.UserRepository.UpdateColumn(simple.GetDB(), id, name, value)
	cache.UserCache.Invalidate(id)
	return err
}

func (this *userService) Delete(id int64) {
	repositories.UserRepository.Delete(simple.GetDB(), id)
	cache.UserCache.Invalidate(id)
}

func (this *userService) GetByEmail(email string) *model.User {
	return repositories.UserRepository.GetByEmail(simple.GetDB(), email)
}

func (this *userService) GetByUsername(username string) *model.User {
	return repositories.UserRepository.GetByUsername(simple.GetDB(), username)
}

// 注册
func (this *userService) SignUp(username, email, nickname, avatar, password, rePassword string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	nickname = strings.TrimSpace(nickname)

	if len(nickname) == 0 {
		return nil, errors.New("昵称不能为空")
	}

	if err := common.IsValidateUsername(username); err != nil {
		return nil, err
	}

	// 验证密码
	err := common.IsValidatePassword(password, rePassword)
	if err != nil {
		return nil, err
	}

	// 如果设置了邮箱，那么需要验证邮箱
	if len(email) > 0 {
		if err := common.IsValidateEmail(email); err != nil {
			return nil, err
		}
		if this.GetByEmail(email) != nil {
			return nil, errors.New("邮箱：" + email + " 已被占用")
		}
	}

	// 验证用户名是否存在
	if this.isUsernameExists(username) {
		return nil, errors.New("用户名：" + username + " 已被占用")
	}

	user := &model.User{
		Username:   simple.SqlNullString(username),
		Email:      simple.SqlNullString(email),
		Nickname:   nickname,
		Password:   simple.EncodePassword(password),
		Avatar:     avatar,
		Status:     model.UserStatusOk,
		CreateTime: simple.NowTimestamp(),
		UpdateTime: simple.NowTimestamp(),
	}

	if err := this.Create(user); err != nil {
		return nil, err
	}
	cache.UserCache.Invalidate(user.Id)
	return user, nil
}

// 登录
func (this *userService) SignIn(username, password string) (*model.User, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名/邮箱不能为空")
	}
	if len(password) == 0 {
		return nil, errors.New("密码不能为空")
	}
	var user *model.User = nil
	if err := common.IsValidateEmail(username); err == nil { // 如果用户输入的是邮箱
		user = this.GetByEmail(username)
	} else {
		user = this.GetByUsername(username)
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	if !simple.ValidatePassword(user.Password, password) {
		return nil, errors.New("密码错误")
	}
	return user, nil
}

// Github账号登录
func (this *userService) SignInByThirdAccount(thirdAccount *model.ThirdAccount) (*model.User, *simple.CodeError) {
	user := this.Get(thirdAccount.UserId)
	if user != nil {
		return user, nil
	}

	user = &model.User{
		Username:   sql.NullString{},
		Nickname:   thirdAccount.Nickname,
		Avatar:     thirdAccount.Avatar,
		Status:     model.UserStatusOk,
		CreateTime: simple.NowTimestamp(),
		UpdateTime: simple.NowTimestamp(),
	}
	err := simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		err := repositories.UserRepository.Create(tx, user)
		if err != nil {
			return err
		}
		err = repositories.ThirdAccountRepository.UpdateColumn(tx, thirdAccount.Id, "user_id", user.Id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, simple.FromError(err)
	}
	cache.UserCache.Invalidate(user.Id)
	return user, nil
}

// 邮箱是否存在
func (this *userService) isEmailExists(email string) bool {
	if len(email) == 0 { // 如果邮箱为空，那么就认为是不存在
		return false
	}
	return this.GetByEmail(email) != nil
}

// 用户名是否存在
func (this *userService) isUsernameExists(username string) bool {
	return this.GetByUsername(username) != nil
}

// 设置用户名
func (this *userService) SetUsername(userId int64, username string) error {
	username = strings.TrimSpace(username)
	if err := common.IsValidateUsername(username); err != nil {
		return err
	}

	user := this.Get(userId)
	if len(user.Username.String) > 0 {
		return errors.New("你已设置了用户名，无法重复设置。")
	}
	if this.isUsernameExists(username) {
		return errors.New("用户名：" + username + " 已被占用")
	}
	return this.UpdateColumn(userId, "username", username)
}

// 设置密码
func (this *userService) SetEmail(userId int64, email string) error {
	email = strings.TrimSpace(email)
	if err := common.IsValidateEmail(email); err != nil {
		return err
	}
	if this.isEmailExists(email) {
		return errors.New("邮箱：" + email + " 已被占用")
	}
	return this.UpdateColumn(userId, "email", email)
}

// 设置密码
func (this *userService) SetPassword(userId int64, password, rePassword string) error {
	if err := common.IsValidatePassword(password, rePassword); err != nil {
		return err
	}
	user := this.Get(userId)
	if len(user.Password) > 0 {
		return errors.New("你已设置了密码，如需修改请前往修改页面。")
	}
	password = simple.EncodePassword(password)
	return this.UpdateColumn(userId, "password", password)
}

// 修改密码
func (this *userService) UpdatePassword(userId int64, oldPassword, password, rePassword string) error {
	if err := common.IsValidatePassword(password, rePassword); err != nil {
		return err
	}
	user := this.Get(userId)

	if len(user.Password) == 0 {
		return errors.New("你没设置密码，请先设置密码")
	}

	if !simple.ValidatePassword(user.Password, oldPassword) {
		return errors.New("旧密码验证失败")
	}

	return this.UpdateColumn(userId, "password", simple.EncodePassword(password))
}
