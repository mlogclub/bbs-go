package render

import (
	"bbs-go/common/uploader"
	"bbs-go/model/constants"
	"github.com/mlogclub/simple/json"
	"html"
	"strconv"
	"strings"

	"github.com/mlogclub/simple/markdown"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/PuerkitoBio/goquery"
	"github.com/mlogclub/simple"

	"bbs-go/cache"
	"bbs-go/common"
	"bbs-go/common/avatar"
	"bbs-go/common/urls"
	"bbs-go/config"
	"bbs-go/model"
	"bbs-go/services"
)

func BuildUserDefaultIfNull(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.Id = id
		user.Username = simple.SqlNullString(strconv.FormatInt(id, 10))
		user.Avatar = avatar.DefaultAvatar
		user.CreateTime = simple.NowTimestamp()
	}
	return BuildUser(user)
}

func BuildUserById(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	return BuildUser(user)
}

func BuildUser(user *model.User) *model.UserInfo {
	if user == nil {
		return nil
	}
	a := user.Avatar
	if len(a) == 0 {
		a = avatar.DefaultAvatar
	}
	roles := strings.Split(user.Roles, ",")
	ret := &model.UserInfo{
		Id:                   user.Id,
		Username:             user.Username.String,
		Nickname:             user.Nickname,
		Avatar:               a,
		SmallAvatar:          HandleOssImageStyleAvatar(a),
		BackgroundImage:      user.BackgroundImage,
		SmallBackgroundImage: HandleOssImageStyleSmall(user.BackgroundImage),
		Email:                user.Email.String,
		EmailVerified:        user.EmailVerified,
		Type:                 user.Type,
		Roles:                roles,
		HomePage:             user.HomePage,
		Description:          user.Description,
		TopicCount:           user.TopicCount,
		CommentCount:         user.CommentCount,
		PasswordSet:          len(user.Password) > 0,
		Forbidden:            user.IsForbidden(),
		Status:               user.Status,
		CreateTime:           user.CreateTime,
	}
	if len(ret.Description) == 0 {
		ret.Description = "这家伙很懒，什么都没留下"
	}
	if user.Status == constants.StatusDeleted {
		ret.Username = "blacklist"
		ret.Nickname = "黑名单用户"
		ret.Avatar = avatar.DefaultAvatar
		ret.Email = ""
		ret.HomePage = ""
		ret.Description = ""
		ret.Forbidden = true
	} else {
		ret.Score = cache.UserCache.GetScore(user.Id)
	}
	return ret
}

func BuildUsers(users []model.User) []model.UserInfo {
	if len(users) == 0 {
		return nil
	}
	var responses []model.UserInfo
	for _, user := range users {
		item := BuildUser(&user)
		if item != nil {
			responses = append(responses, *item)
		}
	}
	return responses
}

func BuildArticle(article *model.Article) *model.ArticleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleResponse{}
	rsp.ArticleId = article.Id
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.SourceUrl = article.SourceUrl
	rsp.ViewCount = article.ViewCount
	rsp.CreateTime = article.CreateTime
	rsp.Status = article.Status

	rsp.User = BuildUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.Id)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = BuildTags(tags)

	if article.ContentType == constants.ContentTypeMarkdown {
		content, _ := markdown.New(markdown.SummaryLen(0)).Run(article.Content)
		rsp.Content = BuildHtmlContent(content)
	} else if article.ContentType == constants.ContentTypeHtml {
		rsp.Content = BuildHtmlContent(article.Content)
	}

	return rsp
}

func BuildSimpleArticle(article *model.Article) *model.ArticleSimpleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleSimpleResponse{}
	rsp.ArticleId = article.Id
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.SourceUrl = article.SourceUrl
	rsp.ViewCount = article.ViewCount
	rsp.CreateTime = article.CreateTime
	rsp.Status = article.Status

	rsp.User = BuildUserDefaultIfNull(article.UserId)

	tagIds := cache.ArticleTagCache.Get(article.Id)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = BuildTags(tags)

	if article.ContentType == constants.ContentTypeMarkdown {
		if len(rsp.Summary) == 0 {
			_, rsp.Summary = markdown.New().Run(article.Content)
		}
	} else if article.ContentType == constants.ContentTypeHtml {
		if len(rsp.Summary) == 0 {
			rsp.Summary = simple.GetSummary(simple.GetHtmlText(article.Content), 256)
		}
	}

	return rsp
}

func BuildSimpleArticles(articles []model.Article) []model.ArticleSimpleResponse {
	if articles == nil || len(articles) == 0 {
		return nil
	}
	var responses []model.ArticleSimpleResponse
	for _, article := range articles {
		responses = append(responses, *BuildSimpleArticle(&article))
	}
	return responses
}

func BuildNode(node *model.TopicNode) *model.NodeResponse {
	if node == nil {
		return nil
	}
	return &model.NodeResponse{
		NodeId:      node.Id,
		Name:        node.Name,
		Description: node.Description,
	}
}

func BuildNodes(nodes []model.TopicNode) []model.NodeResponse {
	if len(nodes) == 0 {
		return nil
	}
	var ret []model.NodeResponse
	for _, node := range nodes {
		ret = append(ret, *BuildNode(&node))
	}
	return ret
}

func BuildTopic(topic *model.Topic) *model.TopicResponse {
	if topic == nil {
		return nil
	}

	rsp := &model.TopicResponse{}

	rsp.TopicId = topic.Id
	rsp.Title = topic.Title
	rsp.User = BuildUserDefaultIfNull(topic.UserId)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount

	if topic.NodeId > 0 {
		node := services.TopicNodeService.Get(topic.NodeId)
		rsp.Node = BuildNode(node)
	}

	tags := services.TopicService.GetTopicTags(topic.Id)
	rsp.Tags = BuildTags(tags)

	content, _ := markdown.New(markdown.SummaryLen(0)).Run(topic.Content)
	rsp.Content = BuildHtmlContent(content)

	return rsp
}

func BuildSimpleTopic(topic *model.Topic) *model.TopicSimpleResponse {
	if topic == nil {
		return nil
	}

	rsp := &model.TopicSimpleResponse{}

	rsp.TopicId = topic.Id
	rsp.Title = topic.Title
	rsp.User = BuildUserDefaultIfNull(topic.UserId)
	rsp.LastCommentTime = topic.LastCommentTime
	rsp.CreateTime = topic.CreateTime
	rsp.ViewCount = topic.ViewCount
	rsp.CommentCount = topic.CommentCount
	rsp.LikeCount = topic.LikeCount

	if topic.NodeId > 0 {
		node := services.TopicNodeService.Get(topic.NodeId)
		rsp.Node = BuildNode(node)
	}

	tags := services.TopicService.GetTopicTags(topic.Id)
	rsp.Tags = BuildTags(tags)
	return rsp
}

func BuildSimpleTopics(topics []model.Topic) []model.TopicSimpleResponse {
	if topics == nil || len(topics) == 0 {
		return nil
	}
	var responses []model.TopicSimpleResponse
	for _, topic := range topics {
		responses = append(responses, *BuildSimpleTopic(&topic))
	}
	return responses
}

func BuildTweet(tweet *model.Tweet) *model.TweetResponse {
	if tweet == nil {
		return nil
	}

	rsp := &model.TweetResponse{
		TweetId:      tweet.Id,
		User:         BuildUserDefaultIfNull(tweet.UserId),
		Content:      tweet.Content,
		CommentCount: tweet.CommentCount,
		LikeCount:    tweet.LikeCount,
		Status:       tweet.Status,
		CreateTime:   tweet.CreateTime,
	}
	if simple.IsNotBlank(tweet.ImageList) {
		var images []string
		if err := json.Parse(tweet.ImageList, &images); err == nil {
			if len(images) > 0 {
				var imageList []model.ImageInfo
				for _, image := range images {
					imageList = append(imageList, model.ImageInfo{
						Url:     HandleOssImageStyleDetail(image),
						Preview: HandleOssImageStylePreview(image),
					})
				}
				rsp.ImageList = imageList
			}
		} else {
			logrus.Error(err)
		}
	}
	return rsp
}

func BuildTweets(tweets []model.Tweet) []model.TweetResponse {
	var ret []model.TweetResponse
	for _, tweet := range tweets {
		ret = append(ret, *BuildTweet(&tweet))
	}
	return ret
}

func BuildProject(project *model.Project) *model.ProjectResponse {
	if project == nil {
		return nil
	}
	rsp := &model.ProjectResponse{}
	rsp.ProjectId = project.Id
	rsp.User = BuildUserDefaultIfNull(project.UserId)
	rsp.Name = project.Name
	rsp.Title = project.Title
	rsp.Logo = project.Logo
	rsp.Url = project.Url
	rsp.Url = project.Url
	rsp.DocUrl = project.DocUrl
	rsp.CreateTime = project.CreateTime

	if project.ContentType == constants.ContentTypeHtml {
		rsp.Content = BuildHtmlContent(project.Content)
		rsp.Summary = simple.GetSummary(simple.GetHtmlText(project.Content), 256)
	} else {
		content, summary := markdown.New().Run(project.Content)
		rsp.Content = BuildHtmlContent(content)
		rsp.Summary = summary
	}

	return rsp
}

func BuildSimpleProjects(projects []model.Project) []model.ProjectSimpleResponse {
	if projects == nil || len(projects) == 0 {
		return nil
	}
	var responses []model.ProjectSimpleResponse
	for _, project := range projects {
		responses = append(responses, *BuildSimpleProject(&project))
	}
	return responses
}

func BuildSimpleProject(project *model.Project) *model.ProjectSimpleResponse {
	if project == nil {
		return nil
	}
	rsp := &model.ProjectSimpleResponse{}
	rsp.ProjectId = project.Id
	rsp.User = BuildUserDefaultIfNull(project.UserId)
	rsp.Name = project.Name
	rsp.Title = project.Title
	rsp.Logo = project.Logo
	rsp.Url = project.Url
	rsp.DocUrl = project.DocUrl
	rsp.DownloadUrl = project.DownloadUrl
	rsp.CreateTime = project.CreateTime

	if project.ContentType == constants.ContentTypeHtml {
		rsp.Summary = simple.GetSummary(simple.GetHtmlText(project.Content), 256)
	} else {
		rsp.Summary = common.GetMarkdownSummary(project.Content)
	}

	return rsp
}

func BuildComments(comments []model.Comment) []model.CommentResponse {
	var ret []model.CommentResponse
	for _, comment := range comments {
		ret = append(ret, *BuildComment(comment))
	}
	return ret
}

func BuildComment(comment model.Comment) *model.CommentResponse {
	return _buildComment(&comment, true)
}

func _buildComment(comment *model.Comment, buildQuote bool) *model.CommentResponse {
	if comment == nil {
		return nil
	}

	ret := &model.CommentResponse{
		CommentId:  comment.Id,
		User:       BuildUserDefaultIfNull(comment.UserId),
		EntityType: comment.EntityType,
		EntityId:   comment.EntityId,
		QuoteId:    comment.QuoteId,
		Status:     comment.Status,
		CreateTime: comment.CreateTime,
	}

	if comment.ContentType == constants.ContentTypeMarkdown {
		content, _ := markdown.New().Run(comment.Content)
		ret.Content = BuildHtmlContent(content)
	} else if comment.ContentType == constants.ContentTypeHtml {
		ret.Content = BuildHtmlContent(comment.Content)
	} else {
		ret.Content = html.EscapeString(comment.Content)
	}

	if buildQuote && comment.QuoteId > 0 {
		quote := _buildComment(services.CommentService.Get(comment.QuoteId), false)
		if quote != nil {
			ret.Quote = quote
			ret.QuoteContent = quote.User.Nickname + "：" + quote.Content
		}
	}
	return ret
}

func BuildTag(tag *model.Tag) *model.TagResponse {
	if tag == nil {
		return nil
	}
	return &model.TagResponse{TagId: tag.Id, TagName: tag.Name}
}

func BuildTags(tags []model.Tag) *[]model.TagResponse {
	if len(tags) == 0 {
		return nil
	}
	var responses []model.TagResponse
	for _, tag := range tags {
		responses = append(responses, *BuildTag(&tag))
	}
	return &responses
}

func BuildFavorite(favorite *model.Favorite) *model.FavoriteResponse {
	rsp := &model.FavoriteResponse{}
	rsp.FavoriteId = favorite.Id
	rsp.EntityType = favorite.EntityType
	rsp.CreateTime = favorite.CreateTime

	if favorite.EntityType == constants.EntityArticle {
		article := services.ArticleService.Get(favorite.EntityId)
		if article == nil || article.Status != constants.StatusOk {
			rsp.Deleted = true
		} else {
			rsp.Url = urls.ArticleUrl(article.Id)
			rsp.User = BuildUserById(article.UserId)
			rsp.Title = article.Title
			if article.ContentType == constants.ContentTypeMarkdown {
				rsp.Content = common.GetMarkdownSummary(article.Content)
			} else if article.ContentType == constants.ContentTypeHtml {
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content))
				if err == nil {
					text := doc.Text()
					rsp.Content = simple.GetSummary(text, 256)
				}
			}
		}
	} else {
		topic := services.TopicService.Get(favorite.EntityId)
		if topic == nil || topic.Status != constants.StatusOk {
			rsp.Deleted = true
		} else {
			rsp.Url = urls.TopicUrl(topic.Id)
			rsp.User = BuildUserById(topic.UserId)
			rsp.Title = topic.Title
			rsp.Content = common.GetMarkdownSummary(topic.Content)
		}
	}
	return rsp
}

func BuildFavorites(favorites []model.Favorite) []model.FavoriteResponse {
	if favorites == nil || len(favorites) == 0 {
		return nil
	}
	var responses []model.FavoriteResponse
	for _, favorite := range favorites {
		responses = append(responses, *BuildFavorite(&favorite))
	}
	return responses
}

func BuildMessage(message *model.Message) *model.MessageResponse {
	if message == nil {
		return nil
	}

	detailUrl := ""
	if message.Type == constants.MsgTypeComment {
		entityType := gjson.Get(message.ExtraData, "entityType")
		entityId := gjson.Get(message.ExtraData, "entityId")
		if entityType.String() == constants.EntityArticle {
			detailUrl = urls.ArticleUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTopic {
			detailUrl = urls.TopicUrl(entityId.Int())
		} else if entityType.String() == constants.EntityTweet {
			detailUrl = urls.TweetUrl(entityId.Int())
		}
	}
	from := BuildUserDefaultIfNull(message.FromId)
	if message.FromId <= 0 {
		from.Nickname = "系统通知"
		from.Avatar = avatar.DefaultAvatar
	}
	return &model.MessageResponse{
		MessageId:    message.Id,
		From:         from,
		UserId:       message.UserId,
		Content:      message.Content,
		QuoteContent: message.QuoteContent,
		Type:         message.Type,
		DetailUrl:    detailUrl,
		ExtraData:    message.ExtraData,
		Status:       message.Status,
		CreateTime:   message.CreateTime,
	}
}

func BuildMessages(messages []model.Message) []model.MessageResponse {
	if len(messages) == 0 {
		return nil
	}
	var responses []model.MessageResponse
	for _, message := range messages {
		responses = append(responses, *BuildMessage(&message))
	}
	return responses
}

func BuildHtmlContent(htmlContent string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return htmlContent
	}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")

		if simple.IsBlank(href) {
			return
		}

		// 不是内部链接
		if !urls.IsInternalUrl(href) {
			selection.SetAttr("target", "_blank")
			selection.SetAttr("rel", "external nofollow") // 标记站外链接，搜索引擎爬虫不传递权重值

			_config := services.SysConfigService.GetConfig()
			if _config.UrlRedirect { // 开启非内部链接跳转
				newHref := simple.ParseUrl(urls.AbsUrl("/redirect")).AddQuery("url", href).BuildStr()
				selection.SetAttr("href", newHref)
			}
		}

		// 如果a标签没有title，那么设置title
		title := selection.AttrOr("title", "")
		if len(title) == 0 {
			selection.SetAttr("title", selection.Text())
		}
	})

	// 处理图片
	doc.Find("img").Each(func(i int, selection *goquery.Selection) {
		src := selection.AttrOr("src", "")
		// 处理第三方图片
		if strings.Contains(src, "qpic.cn") {
			src = simple.ParseUrl("/api/img/proxy").AddQuery("url", src).BuildStr()
			// selection.SetAttr("src", src)
		}

		// 处理图片样式
		src = HandleOssImageStyleDetail(src)

		// 处理lazyload
		selection.SetAttr("data-src", src)
		selection.RemoveAttr("src")
	})

	if htmlStr, err := doc.Find("body").Html(); err == nil {
		return htmlStr
	}
	return htmlContent
}

func HandleOssImageStyleAvatar(url string) string {
	if !uploader.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StyleAvatar)
}

func HandleOssImageStyleDetail(url string) string {
	if !uploader.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StyleDetail)
}

func HandleOssImageStyleSmall(url string) string {
	if !uploader.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StyleSmall)
}

func HandleOssImageStylePreview(url string) string {
	if !uploader.IsEnabledOss() {
		return url
	}
	return HandleOssImageStyle(url, config.Instance.Uploader.AliyunOss.StylePreview)
}

func HandleOssImageStyle(url, style string) string {
	if simple.IsBlank(style) || simple.IsBlank(url) {
		return url
	}
	if !uploader.IsOssImageUrl(url) {
		return url
	}
	sep := config.Instance.Uploader.AliyunOss.StyleSplitter
	if simple.IsBlank(sep) {
		return url
	}
	return strings.Join([]string{url, style}, sep)
}
