package req

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/params"
	"log/slog"
	"strings"

	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/tidwall/gjson"
)

type CreateTopicReq struct {
	Type          constants.TopicType   `json:"type" form:"type"`
	NodeId        int64                 `json:"nodeId" form:"nodeId"`
	Title         string                `json:"title" form:"title"`
	Content       string                `json:"content" form:"content"`
	ContentType   constants.ContentType `json:"contentType" form:"contentType"`
	HideContent   string                `json:"hideContent" form:"hideContent"`
	Tags          []string              `json:"tags" form:"tags"`
	ImageList     []ImageDTO            `json:"imageList" form:"imageList"`
	Vote          *VoteDTO              `json:"vote" form:"vote"`
	BountyScore   int                   `json:"bountyScore" form:"bountyScore"`     // 悬赏积分（仅问答帖有效，0 表示无悬赏）
	AttachmentIds []string              `json:"attachmentIds" form:"attachmentIds"` // 附件 ID 列表（UUID），发帖时绑定到帖子
	UserAgent     string                `json:"userAgent" form:"userAgent"`
	Ip            string                `json:"ip" form:"ip"`

	CaptchaId       string `json:"captchaId" form:"captchaId"`
	CaptchaCode     string `json:"captchaCode" form:"captchaCode"`
	CaptchaProtocol int    `json:"captchaProtocol" form:"captchaProtocol"`
}

type VoteDTO struct {
	Type      constants.VoteType `json:"type" form:"type"`
	Title     string             `json:"title" form:"title"`
	ExpiredAt int64              `json:"expiredAt" form:"expiredAt"`
	VoteNum   int                `json:"voteNum" form:"voteNum"`
	Options   []VoteOptionDTO    `json:"options" form:"options"`
}

type VoteOptionDTO struct {
	Content string `json:"content" form:"content"`
}

type VoteCastReq struct {
	VoteId    int64   `json:"voteId" form:"voteId"`
	OptionIds []int64 `json:"optionIds" form:"optionIds"`
}

type EditTopicReq struct {
	NodeId        int64    `json:"nodeId" form:"nodeId"`
	Title         string   `json:"title" form:"title"`
	Content       string   `json:"content" form:"content"`
	HideContent   string   `json:"hideContent" form:"hideContent"`
	Tags          []string `json:"tags" form:"tags"`
	AttachmentIds []string `json:"attachmentIds" form:"attachmentIds"` // 附件 ID 列表（UUID），全量替换
}

type CreateArticleReq struct {
	Title       string                `json:"title" form:"title"`
	Summary     string                `json:"summary" form:"summary"`
	Content     string                `json:"content" form:"content"`
	ContentType constants.ContentType `json:"contentType" form:"contentType"`
	Cover       *ImageDTO             `json:"cover" form:"cover"`
	Tags        []string              `json:"tags" form:"tags"`
	SourceUrl   string                `json:"sourceUrl" form:"sourceUrl"`
}

type ImageDTO struct {
	Url string `json:"url" form:"url"`
}

func ParseImageList(imageListStr string) []ImageDTO {
	var imageList []ImageDTO
	if strs.IsNotBlank(imageListStr) {
		ret := gjson.Parse(imageListStr)
		if ret.IsArray() {
			for _, item := range ret.Array() {
				url := item.Get("url").String()
				imageList = append(imageList, ImageDTO{
					Url: url,
				})
			}
		}
	}
	return imageList
}

func ParseImageDTO(str string) (img *ImageDTO) {
	str = strings.TrimSpace(str)
	if strs.IsBlank(str) {
		return nil
	}
	if err := jsons.Parse(str, &img); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
	return img
}

type AdminUserCreateReq struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Nickname string `json:"nickname" form:"nickname"`
	Password string `json:"password" form:"password"`
}

type AdminUserUpdateReq struct {
	Id          int64  `json:"id" form:"id"`
	Username    string `json:"username" form:"username"`
	Email       string `json:"email" form:"email"`
	Nickname    string `json:"nickname" form:"nickname"`
	Avatar      string `json:"avatar" form:"avatar"`
	Gender      string `json:"gender" form:"gender"`
	HomePage    string `json:"homePage" form:"homePage"`
	Description string `json:"description" form:"description"`
	RoleIds     string `json:"roleIds" form:"roleIds"`
	Status      int    `json:"status" form:"status"`
}

type AdminUserForbiddenReq struct {
	UserId int64  `json:"userId" form:"userId"`
	Days   int    `json:"days" form:"days"`
	Reason string `json:"reason" form:"reason"`
}

type PasswordUpdateReq struct {
	OldPassword string `json:"oldPassword" form:"oldPassword"`
	Password    string `json:"password" form:"password"`
	RePassword  string `json:"rePassword" form:"rePassword"`
}

type ArticleTagsReq struct {
	ArticleId int64  `json:"articleId" form:"articleId"`
	Tags      string `json:"tags" form:"tags"`
}

type RolePermissionsReq struct {
	Id            int64  `json:"id" form:"id"`
	PermissionIds string `json:"permissionIds" form:"permissionIds"`
}

type TopicAcceptAnswerReq struct {
	Id        int64 `json:"id" form:"id"`
	CommentId int64 `json:"commentId" form:"commentId"`
}

type LoginSignupReq struct {
	CaptchaId       string `json:"captchaId" form:"captchaId"`
	CaptchaCode     string `json:"captchaCode" form:"captchaCode"`
	CaptchaProtocol int    `json:"captchaProtocol" form:"captchaProtocol"`
	Email           string `json:"email" form:"email"`
	Username        string `json:"username" form:"username"`
	Password        string `json:"password" form:"password"`
	RePassword      string `json:"rePassword" form:"rePassword"`
	Nickname        string `json:"nickname" form:"nickname"`
	Redirect        string `json:"redirect" form:"redirect"`
}

type LoginSigninReq struct {
	CaptchaId       string `json:"captchaId" form:"captchaId"`
	CaptchaCode     string `json:"captchaCode" form:"captchaCode"`
	CaptchaProtocol int    `json:"captchaProtocol" form:"captchaProtocol"`
	Username        string `json:"username" form:"username"`
	Password        string `json:"password" form:"password"`
	Redirect        string `json:"redirect" form:"redirect"`
}

type LoginResetEmailReq struct {
	CaptchaId       string `json:"captchaId" form:"captchaId"`
	CaptchaCode     string `json:"captchaCode" form:"captchaCode"`
	CaptchaProtocol int    `json:"captchaProtocol" form:"captchaProtocol"`
	Email           string `json:"email" form:"email"`
}

type LoginResetPasswordReq struct {
	Token      string `json:"token" form:"token"`
	Password   string `json:"password" form:"password"`
	RePassword string `json:"rePassword" form:"rePassword"`
}

type LoginSmsCodeReq struct {
	Phone       string `json:"phone" form:"phone"`
	CaptchaId   string `json:"captchaId" form:"captchaId"`
	CaptchaCode string `json:"captchaCode" form:"captchaCode"`
}

type LoginSmsReq struct {
	SmsId    string `json:"smsId" form:"smsId"`
	SmsCode  string `json:"smsCode" form:"smsCode"`
	Redirect string `json:"redirect" form:"redirect"`
}

type OAuthConfigReq struct {
	Redirect string `form:"redirect"`
	Bind     bool   `form:"bind"`
}

type OAuthCodeStateReq struct {
	Code  string `json:"code" form:"code"`
	State string `json:"state" form:"state"`
}

type UserUpdateReq struct {
	Nickname    string `json:"nickname" form:"nickname"`
	HomePage    string `json:"homePage" form:"homePage"`
	Description string `json:"description" form:"description"`
	Gender      string `json:"gender" form:"gender"`
}

type UserForbiddenReq struct {
	UserId string `json:"userId" form:"userId"`
	Days   int    `json:"days" form:"days"`
	Reason string `json:"reason" form:"reason"`
}

func (r UserForbiddenReq) DecodedUserId() int64 {
	return idcodec.Decode(r.UserId)
}

type ArticleReq struct {
	Title   string `json:"title" form:"title"`
	Summary string `json:"summary" form:"summary"`
	Content string `json:"content" form:"content"`
	Cover   string `json:"cover" form:"cover"`
	Tags    string `json:"tags" form:"tags"`
}

func (r ArticleReq) ParsedTags() []string {
	return SplitCommaStrings(r.Tags)
}

type CreateCommentReq struct {
	EntityType string `json:"entityType" form:"entityType"`
	EntityId   string `json:"entityId" form:"entityId"`
	Content    string `json:"content" form:"content"`
	ImageList  string `json:"imageList" form:"imageList"`
	QuoteId    int64  `json:"quoteId" form:"quoteId"`
	UserAgent  string `json:"userAgent" form:"userAgent"`
	Ip         string `json:"ip" form:"ip"`
}

func (r CreateCommentReq) DecodedEntityId() int64 {
	return idcodec.Decode(r.EntityId)
}

func (r CreateCommentReq) ParsedImageList() []ImageDTO {
	return ParseImageList(r.ImageList)
}

type EntityActionReq struct {
	EntityType string `json:"entityType" form:"entityType"`
	EntityId   string `json:"entityId" form:"entityId"`
}

func (r EntityActionReq) DecodedEntityId() int64 {
	return idcodec.Decode(r.EntityId)
}

type LikedIdsReq struct {
	EntityType string `json:"entityType" form:"entityType"`
	EntityIds  string `json:"entityIds" form:"entityIds"`
}

func (r LikedIdsReq) ParsedEntityIds() []int64 {
	return SplitCommaInt64s(r.EntityIds)
}

type UserReportReq struct {
	DataId   int64  `json:"dataId" form:"dataId"`
	DataType string `json:"dataType" form:"dataType"`
	Reason   string `json:"reason" form:"reason"`
}

type PatchDownloadScoreReq struct {
	Id            string `json:"id" form:"id"`
	DownloadScore int    `json:"downloadScore" form:"downloadScore"`
}

func SplitCommaStrings(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	ret := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			ret = append(ret, part)
		}
	}
	return ret
}

func SplitCommaInt64s(value string) []int64 {
	return params.StrSplitToInt64Arr(value)
}
