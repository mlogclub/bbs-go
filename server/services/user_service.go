package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"server/model/constants"
	"server/pkg/bbsurls"
	"server/pkg/email"
	"server/pkg/errs"
	"server/pkg/event"
	"server/pkg/uploader"
	"server/pkg/validate"
	"strconv"
	"strings"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/passwd"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"

	"server/cache"

	"server/model"
	"server/repositories"
)

// 邮箱验证邮件有效期（小时）
const emailVerifyExpireHour = 24

var UserService = newUserService()

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (s *userService) Get(id int64) *model.User {
	return repositories.UserRepository.Get(sqls.DB(), id)
}

func (s *userService) Take(where ...interface{}) *model.User {
	return repositories.UserRepository.Take(sqls.DB(), where...)
}

func (s *userService) Find(cnd *sqls.Cnd) []model.User {
	return repositories.UserRepository.Find(sqls.DB(), cnd)
}

func (s *userService) FindOne(cnd *sqls.Cnd) *model.User {
	return repositories.UserRepository.FindOne(sqls.DB(), cnd)
}

func (s *userService) FindPageByParams(params *params.QueryParams) (list []model.User, paging *sqls.Paging) {
	return repositories.UserRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userService) FindPageByCnd(cnd *sqls.Cnd) (list []model.User, paging *sqls.Paging) {
	return repositories.UserRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userService) Create(t *model.User) error {
	err := repositories.UserRepository.Create(sqls.DB(), t)
	if err == nil {
		cache.UserCache.Invalidate(t.Id)
	}
	return nil
}

func (s *userService) Update(t *model.User) error {
	err := repositories.UserRepository.Update(sqls.DB(), t)
	cache.UserCache.Invalidate(t.Id)
	return err
}

func (s *userService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.UserRepository.Updates(sqls.DB(), id, columns)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.UserRepository.UpdateColumn(sqls.DB(), id, name, value)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) Delete(id int64) {
	repositories.UserRepository.Delete(sqls.DB(), id)
	cache.UserCache.Invalidate(id)
}

// Scan 扫描
func (s *userService) Scan(callback func(users []model.User)) {
	var cursor int64
	for {
		list := repositories.UserRepository.Find(sqls.DB(), sqls.NewCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// Forbidden 禁言
func (s *userService) Forbidden(operatorId, userId int64, days int, reason string, r *http.Request) error {
	var forbiddenEndTime int64
	if days == -1 { // 永久禁言
		forbiddenEndTime = -1
	} else if days > 0 {
		forbiddenEndTime = dates.Timestamp(time.Now().Add(time.Hour * 24 * time.Duration(days)))
	} else {
		return errors.New("禁言时间错误")
	}
	if repositories.UserRepository.UpdateColumn(sqls.DB(), userId, "forbidden_end_time", forbiddenEndTime) == nil {
		cache.UserCache.Invalidate(userId)
		description := ""
		if strs.IsNotBlank(reason) {
			description = "禁言原因：" + reason
		}
		OperateLogService.AddOperateLog(operatorId, constants.OpTypeForbidden, constants.EntityUser, userId,
			description, r)

		// 永久禁言
		if days == -1 {
			user := cache.UserCache.Get(userId)
			_ = s.DecrScore(userId, user.Score, constants.EntityUser, strconv.FormatInt(operatorId, 10), "永久禁言")
			go func() {
				// 删除话题
				TopicService.ScanByUser(userId, func(topics []model.Topic) {
					for _, topic := range topics {
						if topic.Status != constants.StatusDeleted {
							_ = TopicService.Delete(topic.Id, operatorId, nil)
						}
					}
				})

				// 删除文章
				ArticleService.ScanByUser(userId, func(articles []model.Article) {
					for _, article := range articles {
						if article.Status != constants.StatusDeleted {
							_ = ArticleService.Delete(article.Id)
						}
					}
				})

				// 删除评论
				CommentService.ScanByUser(userId, func(comments []model.Comment) {
					for _, comment := range comments {
						if comment.Status != constants.StatusDeleted {
							_ = CommentService.Delete(comment.Id)
						}
					}
				})

			}()
		}
	}
	return nil
}

// RemoveForbidden 移除禁言
func (s *userService) RemoveForbidden(operatorId, userId int64, r *http.Request) {
	user := s.Get(userId)
	if user == nil || !user.IsForbidden() {
		return
	}
	if repositories.UserRepository.UpdateColumn(sqls.DB(), userId, "forbidden_end_time", 0) == nil {
		cache.UserCache.Invalidate(user.Id)
		OperateLogService.AddOperateLog(operatorId, constants.OpTypeRemoveForbidden, constants.EntityUser, userId, "", r)
	}
}

// GetByEmail 根据邮箱查找
func (s *userService) GetByEmail(email string) *model.User {
	return repositories.UserRepository.GetByEmail(sqls.DB(), email)
}

// GetByUsername 根据用户名查找
func (s *userService) GetByUsername(username string) *model.User {
	return repositories.UserRepository.GetByUsername(sqls.DB(), username)
}

// SignUp 注册
func (s *userService) SignUp(username, email, nickname, password, rePassword, refereeCode string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	nickname = strings.TrimSpace(nickname)
	refereeCode = strings.TrimSpace(refereeCode)

	// 验证昵称
	if len(nickname) == 0 {
		return nil, errors.New("昵称不能为空")
	}

	// 验证密码
	err := validate.IsPassword(password, rePassword)
	if err != nil {
		return nil, err
	}

	// 验证邮箱
	if len(email) > 0 {
		if err := validate.IsEmail(email); err != nil {
			return nil, err
		}
		if s.GetByEmail(email) != nil {
			return nil, errors.New("邮箱：" + email + " 已被占用")
		}
	} else {
		return nil, errors.New("请输入邮箱")
	}

	// 验证用户名
	if len(username) > 0 {
		if err := validate.IsUsername(username); err != nil {
			return nil, err
		}
		if s.isUsernameExists(username) {
			return nil, errors.New("用户名：" + username + " 已被占用")
		}
	}

	user := &model.User{
		Username:   sqls.SqlNullString(username),
		Email:      sqls.SqlNullString(email),
		Nickname:   nickname,
		Password:   passwd.EncodePassword(password),
		Status:     constants.StatusOk,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}

	err = repositories.UserRepository.Create(sqls.DB(), user)
	if err != nil {
		return nil, err
	}

	if len(refereeCode) > 0 {
		userReferrer := &model.UserReferee{
			UserId:      user.Id,
			RefereeCode: refereeCode,
			Status:      0,
		}

		err = repositories.UserRefereeRepository.Create(sqls.DB(), userReferrer)
		if err != nil {
			logrus.Errorf("推荐人创建失败 referrerCode : [%s] - userId : [%d]", refereeCode, user.Id)
		}
	}

	return user, nil
}

// SignIn 登录
func (s *userService) SignIn(username, password string) (*model.User, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名/邮箱不能为空")
	}
	if len(password) == 0 {
		return nil, errors.New("密码不能为空")
	}
	var user *model.User = nil
	if err := validate.IsEmail(username); err == nil { // 如果用户输入的是邮箱
		user = s.GetByEmail(username)
	} else {
		user = s.GetByUsername(username)
	}
	user.QQStatus = ThirdAccountService.FindOne(sqls.NewCnd().Eq("user_id", user.Id).Eq("third_type", constants.ThirdAccountTypeQQ)) != nil
	if user == nil || user.Status != constants.StatusOk {
		return nil, errors.New("用户不存在或被禁用")
	}
	if !passwd.ValidatePassword(user.Password, password) {
		return nil, errors.New("密码错误")
	}
	return user, nil
}

// SignInByThirdAccount 第三方账号登录
func (s *userService) SignInByThirdAccount(thirdAccount *model.ThirdAccount) (*model.User, *web.CodeError) {
	user := s.Get(thirdAccount.UserId.Int64)
	if user != nil {
		if user.Status != constants.StatusOk {
			return nil, web.NewErrorMsg("用户已被禁用")
		}
		return user, nil
	}

	var homePage string
	var description string
	if thirdAccount.ThirdType == constants.ThirdAccountTypeGithub {
		if blog := gjson.Get(thirdAccount.ExtraData, "blog"); blog.Exists() && len(blog.String()) > 0 {
			homePage = blog.String()
		} else if htmlUrl := gjson.Get(thirdAccount.ExtraData, "html_url"); htmlUrl.Exists() && len(htmlUrl.String()) > 0 {
			homePage = htmlUrl.String()
		}

		description = gjson.Get(thirdAccount.ExtraData, "bio").String()
	}

	user = &model.User{
		Username:    sql.NullString{},
		Nickname:    thirdAccount.Nickname,
		Status:      constants.StatusOk,
		HomePage:    homePage,
		Description: description,
		CreateTime:  dates.NowTimestamp(),
		UpdateTime:  dates.NowTimestamp(),
	}
	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.UserRepository.Create(tx, user); err != nil {
			return err
		}

		if err := repositories.ThirdAccountRepository.UpdateColumn(tx, thirdAccount.Id, "user_id", user.Id); err != nil {
			return err
		}

		avatarUrl := s.HandleThirdAvatar(thirdAccount.Avatar)

		if err := repositories.UserRepository.UpdateColumn(tx, user.Id, "avatar", avatarUrl); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, web.FromError(err)
	}
	cache.UserCache.Invalidate(user.Id)
	return user, nil
}

// HandleThirdAvatar 处理第三方头像
func (s *userService) HandleThirdAvatar(thirdAvatar string) string {
	if strs.IsBlank(thirdAvatar) {
		return ""
	}
	avatar, err := uploader.CopyImage(thirdAvatar)
	if err != nil {
		return ""
	}
	return avatar
}

// isEmailExists 邮箱是否存在
func (s *userService) isEmailExists(email string) bool {
	if len(email) == 0 { // 如果邮箱为空，那么就认为是不存在
		return false
	}
	return s.GetByEmail(email) != nil
}

// isUsernameExists 用户名是否存在
func (s *userService) isUsernameExists(username string) bool {
	return s.GetByUsername(username) != nil
}

// UpdateAvatar 更新头像
func (s *userService) UpdateAvatar(userId int64, avatar string) error {
	return s.UpdateColumn(userId, "avatar", avatar)
}

// UpdateNickname 更新昵称
func (s *userService) UpdateNickname(userId int64, nickname string) error {
	return s.UpdateColumn(userId, "nickname", nickname)
}

// UpdateDescription 更新简介
func (s *userService) UpdateDescription(userId int64, description string) error {
	return s.UpdateColumn(userId, "description", description)
}

// UpdateGender 修改性别
func (s *userService) UpdateGender(userId int64, gender string) error {
	if strs.IsBlank(gender) {
		return s.UpdateColumn(userId, "gender", "")
	} else {
		if gender != string(constants.GenderMale) && gender != string(constants.GenderFemale) {
			return errors.New("invalidate gender value")
		}
		return s.UpdateColumn(userId, "gender", gender)
	}
}

// UpdateBirthday 修改生日
func (s *userService) UpdateBirthday(userId int64, birthdayStr string) error {
	if strs.IsBlank(birthdayStr) {
		return s.UpdateColumn(userId, "birthday", "")
	} else {
		birthday, err := dates.Parse(birthdayStr, dates.FmtDate)
		if err != nil {
			return err
		}
		return s.UpdateColumn(userId, "birthday", birthday)
	}
}

// UpdateBackgroundImage 修改背景图
func (s *userService) UpdateBackgroundImage(userId int64, backgroundImage string) error {
	return s.UpdateColumn(userId, "background_image", backgroundImage)
}

// SetUsername 设置用户名
func (s *userService) SetUsername(userId int64, username string) error {
	username = strings.TrimSpace(username)
	if err := validate.IsUsername(username); err != nil {
		return err
	}

	user := s.Get(userId)
	if len(user.Username.String) > 0 {
		return errors.New("你已设置了用户名，无法重复设置。")
	}
	if s.isUsernameExists(username) {
		return errors.New("用户名：" + username + " 已被占用")
	}
	return s.UpdateColumn(userId, "username", username)
}

// SetEmail 设置密码
func (s *userService) SetEmail(userId int64, email string) error {
	email = strings.TrimSpace(email)
	if err := validate.IsEmail(email); err != nil {
		return err
	}
	user := s.Get(userId)
	if user == nil {
		return errors.New("用户不存在")
	}
	if user.Email.String == email {
		// 用户邮箱没做变更
		return nil
	}
	if s.isEmailExists(email) {
		return errors.New("邮箱：" + email + " 已被占用")
	}
	return s.Updates(userId, map[string]interface{}{
		"email":          email,
		"email_verified": false,
	})
}

// SetPassword 设置密码
func (s *userService) SetPassword(userId int64, password, rePassword string) error {
	if err := validate.IsPassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)
	if len(user.Password) > 0 {
		return errors.New("你已设置了密码，如需修改请前往修改页面。")
	}
	password = passwd.EncodePassword(password)
	return s.UpdateColumn(userId, "password", password)
}

// UpdatePassword 修改密码
func (s *userService) UpdatePassword(userId int64, oldPassword, password, rePassword string) error {
	if err := validate.IsPassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)

	if len(user.Password) == 0 {
		return errors.New("你没设置密码，请先设置密码")
	}

	if !passwd.ValidatePassword(user.Password, oldPassword) {
		return errors.New("旧密码验证失败")
	}

	return s.UpdateColumn(userId, "password", passwd.EncodePassword(password))
}

// IncrTopicCount topic_count + 1
func (s *userService) IncrTopicCount(tx *gorm.DB, userId int64) error {
	if err := repositories.UserRepository.UpdateColumn(sqls.DB(), userId, "topic_count", gorm.Expr("topic_count + 1")); err != nil {
		logrus.Error(err)
		return err
	}
	cache.UserCache.Invalidate(userId)
	return nil
}

// IncrCommentCount comment_count + 1
func (s *userService) IncrCommentCount(userId int64) int {
	t := repositories.UserRepository.Get(sqls.DB(), userId)
	if t == nil {
		return 0
	}
	commentCount := t.CommentCount + 1
	if err := repositories.UserRepository.UpdateColumn(sqls.DB(), userId, "comment_count", commentCount); err != nil {
		logrus.Error(err)
	} else {
		cache.UserCache.Invalidate(userId)
	}
	return commentCount
}

// SyncUserCount 同步用户计数
func (s *userService) SyncUserCount() {
	s.Scan(func(users []model.User) {
		for _, user := range users {
			topicCount := repositories.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).Eq("status", constants.StatusOk))
			commentCount := repositories.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).Eq("status", constants.StatusOk))
			_ = repositories.UserRepository.UpdateColumn(sqls.DB(), user.Id, "topic_count", topicCount)
			_ = repositories.UserRepository.UpdateColumn(sqls.DB(), user.Id, "comment_count", commentCount)
			cache.UserCache.Invalidate(user.Id)
		}
	})
}

// SendEmailVerifyEmail 发送邮箱验证邮件
func (s *userService) SendEmailVerifyEmail(userId int64) error {
	user := s.Get(userId)
	if user == nil {
		return errors.New("用户不存在")
	}
	if user.EmailVerified {
		return errors.New("用户邮箱已验证")
	}
	if err := validate.IsEmail(user.Email.String); err != nil {
		return err
	}
	// 如果设置了邮箱白名单
	if emailWhitelist := SysConfigService.GetEmailWhitelist(); len(emailWhitelist) > 0 {
		isInWhitelist := false
		for _, whitelist := range emailWhitelist {
			if strings.Contains(strings.ToLower(user.Email.String), strings.ToLower(whitelist)) {
				isInWhitelist = true
				break
			}
		}
		if !isInWhitelist {
			// 直接返回，也不抛出异常了，就是不发邮件
			logrus.Error("不支持使用该邮箱进行验证.", user.Email.String)
			return errors.New("不支持该类型邮箱")
		}
	}
	var (
		token     = strs.UUID()
		url       = bbsurls.AbsUrl("/user/email/verify?token=" + token)
		link      = &model.ActionLink{Title: "点击这里验证邮箱>>", Url: url}
		siteTitle = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
		subject   = "邮箱验证 - " + siteTitle
		title     = "邮箱验证 - " + siteTitle
		content   = "该邮件用于验证你在 " + siteTitle + " 中设置邮箱的正确性，请在" + strconv.Itoa(emailVerifyExpireHour) + "小时内完成验证。验证链接：" + url
	)
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.EmailCodeRepository.Create(tx, &model.EmailCode{
			Model:      model.Model{},
			UserId:     userId,
			Email:      user.Email.String,
			Code:       "",
			Token:      token,
			Title:      title,
			Content:    content,
			Used:       false,
			CreateTime: dates.NowTimestamp(),
		}); err != nil {
			return nil
		}
		if err := email.SendTemplateEmail(nil, user.Email.String, subject, title, content, "", link); err != nil {
			return err
		}
		return nil
	})
}

// VerifyEmail 验证邮箱
func (s *userService) VerifyEmail(token string) (string, error) {
	emailCode := EmailCodeService.FindOne(sqls.NewCnd().Eq("token", token))
	if emailCode == nil {
		return "", errs.BadRequest
	}

	if emailCode.Used {
		return "", errs.EmailVerified
	}

	user := s.Get(emailCode.UserId)
	if user == nil || emailCode.Email != user.Email.String {
		return "", errs.EmailTimeout
	}
	if dates.FromTimestamp(emailCode.CreateTime).Add(time.Hour * time.Duration(emailVerifyExpireHour)).Before(time.Now()) {
		return "", errs.EmailTimeout
	}
	err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.UserRepository.UpdateColumn(tx, emailCode.UserId, "email_verified", true); err != nil {
			return err
		}
		cache.UserCache.Invalidate(emailCode.UserId)
		return repositories.EmailCodeRepository.UpdateColumn(tx, emailCode.Id, "used", true)
	})
	if err != nil {
		return "", err
	}

	userReferee := repositories.UserRefereeRepository.FindOne(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).Eq("status", 0))
	if userReferee != nil {
		userReferee.Status = 1
		err := repositories.UserRefereeRepository.Update(sqls.DB(), userReferee)
		if err != nil {
			logrus.Errorf("更新推荐人状态失败 用户ID : [%d] - 推荐人Code : [%d]", userReferee.UserId, userReferee.RefereeCode)
		}
		users := repositories.UserRepository.FindOne(sqls.DB(), sqls.NewCnd().Eq("referee_code", userReferee.RefereeCode))
		users.RefereeCount = users.RefereeCount + 1
		err = repositories.UserRepository.Update(sqls.DB(), users)
		if err != nil {
			logrus.Errorf("更新推荐数量失败 推荐人ID : [%d] - 用户ID : [%d]", users.Id, user.Id)
		}
	}

	return emailCode.Email, nil
}

// CheckPostStatus 用于在发表内容时检查用户状态
func (s *userService) CheckPostStatus(user *model.User) error {
	if user == nil {
		return errs.NotLogin
	}
	if user.Status != constants.StatusOk {
		return errs.UserDisabled
	}
	if user.IsForbidden() {
		return errs.ForbiddenError
	}
	observeSeconds := SysConfigService.GetInt(constants.SysConfigUserObserveSeconds, 0)
	if user.InObservationPeriod(observeSeconds) {
		return web.NewError(errs.InObservationPeriod.Code, "账号尚在观察期，观察期时长："+strconv.Itoa(observeSeconds)+"秒，请稍后再试")
	}
	return nil
}

// IncrScoreForPostTopic 发帖获积分
func (s *userService) IncrScoreForPostTopic(topic *model.Topic) {
	config := SysConfigService.GetConfig()
	if config.ScoreConfig.PostTopicScore <= 0 {
		logrus.Info("请配置发帖积分")
		return
	}
	err := s.addScore(topic.UserId, config.ScoreConfig.PostTopicScore, constants.EntityTopic,
		strconv.FormatInt(topic.Id, 10), "发表话题")
	if err != nil {
		logrus.Error(err)
	}
}

// IncrScoreForPostComment 跟帖获积分
func (s *userService) IncrScoreForPostComment(comment *model.Comment) {
	// 非话题跟帖，跳过
	if comment.EntityType != constants.EntityTopic {
		return
	}
	config := SysConfigService.GetConfig()
	if config.ScoreConfig.PostCommentScore <= 0 {
		logrus.Info("请配置跟帖积分")
		return
	}
	err := s.addScore(comment.UserId, config.ScoreConfig.PostCommentScore, constants.EntityComment,
		strconv.FormatInt(comment.Id, 10), "发表跟帖")
	if err != nil {
		logrus.Error(err)
	}
}

// IncrScore 增加分数
func (s *userService) IncrScore(userId int64, score int64, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, score, sourceType, sourceId, description)
}

// DecrScore 减少分数
func (s *userService) DecrScore(userId int64, score int64, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, -score, sourceType, sourceId, description)
}

// addScore 加分数，也可以加负数
func (s *userService) addScore(userId int64, score int64, sourceType, sourceId, description string) error {
	if score == 0 {
		return errors.New("分数不能为0")
	}
	user := s.Get(userId)
	if user == nil {
		return errors.New("用户不存在")
	}
	if err := s.Updates(userId, map[string]interface{}{
		"score":       gorm.Expr("score + ?", score),
		"update_time": dates.NowTimestamp(),
	}); err != nil {
		return err
	}

	scoreType := constants.ScoreTypeIncr
	if score < 0 {
		scoreType = constants.ScoreTypeDecr
	}
	err := UserScoreLogService.Create(&model.UserScoreLog{
		UserId:      userId,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
		Type:        scoreType,
		Score:       score,
		CreateTime:  dates.NowTimestamp(),
	})
	if err == nil {
		cache.UserCache.Invalidate(userId)
	}
	return err
}

// IncrScore 充值
func (s *userService) PayScore(userId int64, score int64) error {
	if score <= 0 {
		logrus.Errorf("积分数量非法 : [%d] - [%d]", userId, score)
		return errs.ErrScore
	}

	err := s.IncrScore(userId, score, "pay", string(userId), "充值")
	if err != nil {
		logrus.Errorf("充值失败 : [%d] - [%d]", userId, score)
		return errs.ErrPayScore
	}

	event.Send(event.ScorePayEvent{
		ToUserId: userId,
		Score:    score,
	})

	return err
}

// SendEmailEmail 发送找回密码邮箱验证邮件
func (s *userService) SendEmail(e string) error {

	if err := validate.IsEmail(e); err != nil {
		return err
	}

	user := s.FindOne(sqls.NewCnd().Eq("email", e))
	if user == nil {
		return errors.New("用户不存在")
	}

	emailCode := EmailCodeService.FindOne(sqls.NewCnd().Eq("email", e).Where("code <> ?", "").Where("Used = ?", false))
	if emailCode != nil {
		return errors.New("验证码已发送")
	}

	var (
		code      = fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
		siteTitle = cache.SysConfigCache.GetValue(constants.SysConfigSiteTitle)
		subject   = "密码找回 - " + siteTitle
		title     = "密码找回 - " + siteTitle
		content   = "该邮件验证码用于找回你在 " + siteTitle + " 的密码，验证码有效期" + strconv.Itoa(emailVerifyExpireHour) + "小时, 验证码为：" + code
	)
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := repositories.EmailCodeRepository.Create(tx, &model.EmailCode{
			Model:      model.Model{},
			UserId:     user.Id,
			Email:      user.Email.String,
			Code:       code,
			Token:      code,
			Title:      title,
			Content:    content,
			Used:       false,
			CreateTime: dates.NowTimestamp(),
		}); err != nil {
			return nil
		}
		if err := email.SendTemplateEmail(nil, user.Email.String, subject, title, content, "", nil); err != nil {
			return err
		}
		return nil
	})
}

func (s *userService) Forgotpwd(email, email_code, password, rePassword string) error {
	user := s.FindOne(sqls.NewCnd().Eq("email", email))
	if user == nil {
		return errors.New("用户不存在")
	}
	emailCode := EmailCodeService.FindOne(sqls.NewCnd().Eq("email", email).Where("code = ?", email_code).Where("Used = ?", false))
	if emailCode == nil {
		return errors.New("验证码已使用")
	}
	if err := validate.IsPassword(password, rePassword); err != nil {
		return err
	}

	password = passwd.EncodePassword(password)
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		err := repositories.UserRepository.UpdateColumn(tx, user.Id, "password", password)
		if err != nil {
			logrus.Errorf("密码设置失败: [%s]", err.Error())
			return err
		}
		err = repositories.EmailCodeRepository.UpdateColumn(tx, emailCode.Id, "used", true)
		if err != nil {
			logrus.Errorf("code 更新失败: [%s]", err.Error())
			return err
		}
		return nil
	})

}

func (s *userService) CronUserPayScore() {
	logrus.Infof("Vip 积分每日充值任务开启")
	conf := SysConfigService.Get(21).Value
	congMap := make(map[string]int64)
	err := json.Unmarshal([]byte(conf), &congMap)
	if err != nil {
		logrus.Errorf("CronUserPayScore Unmarshal Err : [%s]", err.Error())
	}
	vip1 := s.Find(sqls.NewCnd().Eq("vip", 1))
	go s.userPayScore(vip1, congMap["vip1"])
	vip2 := s.Find(sqls.NewCnd().Eq("vip", 2))
	go s.userPayScore(vip2, congMap["vip2"])
	vip3 := s.Find(sqls.NewCnd().Eq("vip", 3))
	go s.userPayScore(vip3, congMap["vip3"])
	vip4 := s.Find(sqls.NewCnd().Eq("vip", 4))
	go s.userPayScore(vip4, congMap["vip4"])
	vip5 := s.Find(sqls.NewCnd().Eq("vip", 5))
	go s.userPayScore(vip5, congMap["vip5"])
	vip6 := s.Find(sqls.NewCnd().Eq("vip", 6))
	go s.userPayScore(vip6, congMap["vip6"])
}

func (s *userService) userPayScore(user []model.User, score int64) {
	for _, v := range user {
		_ = s.PayScore(v.Id, score)
	}
}
