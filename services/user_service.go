package services

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/core/errors"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/mlog/utils"
	"github.com/mlogclub/mlog/utils/validate"
	"github.com/mlogclub/simple"
	"strings"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
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
	return repositories.UserRepository.Create(simple.GetDB(), t)
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

// 登录
func (this *userService) SignIn(username, password string) (*model.User, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名/邮箱不能为空")
	}
	if len(password) == 0 {
		return nil, errors.New("密码不能为空")
	}
	var user *model.User = nil
	if validate.IsEmail(username) { // 如果用户输入的是邮箱
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

// 注册
func (this *userService) SignUp(username, email, password, rePassword, nickname, avatar string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if len(username) == 0 {
		return nil, errors.New("请输入用户名")
	}
	if simple.RuneLen(username) < 5 {
		return nil, errors.New("用户名长度不能少于5个字符")
	}
	if !validate.IsEmail(email) {
		return nil, errors.New("请输入合法的邮箱")
	}
	if len(password) == 0 {
		return nil, errors.New("请输入密码")
	}
	if simple.RuneLen(password) < 6 {
		return nil, errors.New("密码过于简单")
	}
	if len(nickname) == 0 {
		return nil, errors.New("昵称不能为空")
	}
	if password != rePassword {
		return nil, errors.New("两次输入密码不匹配")
	}

	if this.GetByUsername(username) != nil {
		return nil, errors.New("用户名：" + username + " 已被占用")
	}

	if this.GetByEmail(email) != nil {
		return nil, errors.New("邮箱：" + email + " 已被占用")
	}

	password = simple.EncodePassword(password)

	user := &model.User{
		Username:   username,
		Email:      email,
		Nickname:   nickname,
		Password:   password,
		Avatar:     avatar,
		Status:     model.UserStatusOk,
		CreateTime: simple.NowTimestamp(),
		UpdateTime: simple.NowTimestamp(),
	}

	err := this.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// 绑定账号
func (this *userService) Bind(githubId int64, bindType, username, email, password, rePassword, nickname string) (user *model.User, err error) {
	githubUser := repositories.GithubUserRepository.Get(simple.GetDB(), githubId)
	if githubUser == nil {
		err = errors.New("Github账号未找到")
		return
	}
	if githubUser.UserId > 0 {
		err = errors.New("Github账号已绑定了用户")
		return
	}

	if bindType == "login" { // 登录绑定
		user, err = this.SignIn(username, password)
		if err != nil {
			return
		}
	} else { // 注册绑定
		user, err = this.SignUp(username, email, password, rePassword, nickname, githubUser.AvatarUrl)
		if err != nil {
			return
		}
	}

	if user == nil {
		err = errors.New("未知异常")
		return
	}

	// 执行绑定
	githubUser.UserId = user.Id
	githubUser.UpdateTime = simple.NowTimestamp()
	err = repositories.GithubUserRepository.Update(simple.GetDB(), githubUser)
	return
}

// Github账号登录
func (this *userService) SignInByGithub(githubUser *model.GithubUser) (*model.User, *simple.CodeError) {
	user := this.Get(githubUser.UserId)
	if user != nil {
		return user, nil
	}

	if this.GetByUsername(githubUser.Login) == nil { // 用户名没有被占用，直接注册
		user = &model.User{
			Username:   githubUser.Login,
			Email:      githubUser.Email,
			Nickname:   githubUser.Name,
			Avatar:     githubUser.AvatarUrl,
			Status:     model.UserStatusOk,
			CreateTime: simple.NowTimestamp(),
			UpdateTime: simple.NowTimestamp(),
		}
		err := simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
			err := repositories.UserRepository.Create(tx, user)
			if err != nil {
				return err
			}
			err = repositories.GithubUserRepository.UpdateColumn(tx, githubUser.Id, "user_id", user.Id)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, simple.NewError2(err)
		}
		return user, nil
	} else { // 用户名称已经存在，那么要让用户修改
		return nil, simple.NewError(utils.ErrorCodeUserNameExists, "用户名已存在")
	}
}
