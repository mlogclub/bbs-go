package render

import (
	"html/template"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/PuerkitoBio/goquery"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/common"
	"github.com/mlogclub/bbs-go/common/avatar"
	"github.com/mlogclub/bbs-go/common/urls"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/bbs-go/services/cache"
)

func BuildUserDefaultIfNull(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.Id = id
		user.Username = simple.SqlNullString(strconv.FormatInt(id, 10))
		user.Avatar = avatar.GetDefaultAvatar(id)
		user.CreateTime = simple.NowTimestamp()
	}
	return BuildUser(user)
}

func BuildUserById(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		return nil
	}
	return BuildUser(user)
}

func BuildUser(user *model.User) *model.UserInfo {
	if user == nil {
		return nil
	}
	a := user.Avatar
	if len(a) == 0 {
		a = avatar.GetDefaultAvatar(user.Id)
	}
	roles := strings.Split(user.Roles, ",")
	return &model.UserInfo{
		Id:          user.Id,
		Username:    user.Username.String,
		Nickname:    user.Nickname,
		Avatar:      a,
		Email:       user.Email.String,
		Type:        user.Type,
		Roles:       roles,
		Description: user.Description,
		PasswordSet: len(user.Password) > 0,
		CreateTime:  user.CreateTime,
	}
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

func BuildCategory(category *model.Category) *model.CategoryResponse {
	if category == nil {
		return nil
	}
	return &model.CategoryResponse{CategoryId: category.Id, CategoryName: category.Name}
}

func BuildArticle(article *model.Article) *model.ArticleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleResponse{}
	rsp.ArticleId = article.Id
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.Share = article.Share
	rsp.SourceUrl = article.SourceUrl
	rsp.CreateTime = article.CreateTime

	rsp.User = BuildUserDefaultIfNull(article.UserId)

	if article.CategoryId > 0 {
		category := cache.CategoryCache.Get(article.CategoryId)
		rsp.Category = BuildCategory(category)
	}

	tagIds := cache.ArticleTagCache.Get(article.Id)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = BuildTags(tags)

	if article.ContentType == model.ContentTypeMarkdown {
		mr := simple.NewMd(simple.MdWithTOC()).Run(article.Content)
		rsp.Content = template.HTML(BuildHtmlContent(mr.ContentHtml))
		rsp.Toc = template.HTML(mr.TocHtml)
		if len(rsp.Summary) == 0 {
			rsp.Summary = mr.SummaryText
		}
	} else {
		rsp.Content = template.HTML(BuildHtmlContent(article.Content))
		if len(rsp.Summary) == 0 {
			rsp.Summary = simple.GetSummary(article.Content, 256)
		}
	}

	return rsp
}

func BuildArticles(articles []model.Article) []model.ArticleResponse {
	if articles == nil || len(articles) == 0 {
		return nil
	}
	var responses []model.ArticleResponse
	for _, article := range articles {
		responses = append(responses, *BuildArticle(&article))
	}
	return responses
}

func BuildSimpleArticle(article *model.Article) *model.ArticleSimpleResponse {
	if article == nil {
		return nil
	}

	rsp := &model.ArticleSimpleResponse{}
	rsp.ArticleId = article.Id
	rsp.Title = article.Title
	rsp.Summary = article.Summary
	rsp.Share = article.Share
	rsp.SourceUrl = article.SourceUrl
	rsp.CreateTime = article.CreateTime

	rsp.User = BuildUserDefaultIfNull(article.UserId)

	if article.CategoryId > 0 {
		category := cache.CategoryCache.Get(article.CategoryId)
		rsp.Category = BuildCategory(category)
	}

	tagIds := cache.ArticleTagCache.Get(article.Id)
	tags := cache.TagCache.GetList(tagIds)
	rsp.Tags = BuildTags(tags)

	if article.ContentType == model.ContentTypeMarkdown {
		if len(rsp.Summary) == 0 {
			mr := simple.NewMd(simple.MdWithTOC()).Run(article.Content)
			rsp.Summary = mr.SummaryText
		}
	} else {
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

	tags := services.TopicService.GetTopicTags(topic.Id)
	rsp.Tags = BuildTags(tags)

	mr := simple.NewMd(simple.MdWithTOC()).Run(topic.Content)
	rsp.Content = template.HTML(BuildHtmlContent(mr.ContentHtml))
	rsp.Toc = template.HTML(mr.TocHtml)

	return rsp
}

func BuildTopics(topics []model.Topic) []model.TopicResponse {
	if topics == nil || len(topics) == 0 {
		return nil
	}
	var responses []model.TopicResponse
	for _, topic := range topics {
		responses = append(responses, *BuildTopic(&topic))
	}
	return responses
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

	if project.ContentType == model.ContentTypeHtml {
		rsp.Content = template.HTML(BuildHtmlContent(project.Content))
		rsp.Summary = simple.GetSummary(simple.GetHtmlText(project.Content), 256)
	} else {
		mr := simple.NewMd().Run(project.Content)
		rsp.Content = template.HTML(BuildHtmlContent(mr.ContentHtml))
		rsp.Summary = mr.SummaryText
	}

	return rsp
}

func BuildSimpleProjects(projects []model.Project) [] model.ProjectSimpleResponse {
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

	if project.ContentType == model.ContentTypeHtml {
		rsp.Summary = simple.GetSummary(simple.GetHtmlText(project.Content), 256)
	} else {
		rsp.Summary = common.GetMarkdownSummary(project.Content)
	}

	return rsp
}

func BuildComment(comment model.Comment) *model.CommentResponse {
	return _buildComment(&comment, true)
}

func _buildComment(comment *model.Comment, buildQuote bool) *model.CommentResponse {
	if comment == nil {
		return nil
	}
	markdownResult := simple.NewMd().Run(comment.Content)
	content := template.HTML(markdownResult.ContentHtml)

	ret := &model.CommentResponse{
		CommentId:  comment.Id,
		User:       BuildUserDefaultIfNull(comment.UserId),
		EntityType: comment.EntityType,
		EntityId:   comment.EntityId,
		Content:    content,
		QuoteId:    comment.QuoteId,
		Status:     comment.Status,
		CreateTime: comment.CreateTime,
	}

	if buildQuote && comment.QuoteId > 0 {
		quote := _buildComment(services.CommentService.Get(comment.QuoteId), false)
		if quote != nil {
			ret.Quote = quote
			ret.QuoteContent = template.HTML(quote.User.Nickname+"：") + quote.Content
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

	if favorite.EntityType == model.EntityTypeArticle {
		article := services.ArticleService.Get(favorite.EntityId)
		if article == nil || article.Status != model.ArticleStatusPublished {
			rsp.Deleted = true
		} else {
			rsp.Url = urls.ArticleUrl(article.Id)
			rsp.User = BuildUserById(article.UserId)
			rsp.Title = article.Title
			if article.ContentType == model.ContentTypeMarkdown {
				rsp.Content = common.GetMarkdownSummary(article.Content)
			} else {
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(article.Content))
				if err == nil {
					text := doc.Text()
					rsp.Content = simple.GetSummary(text, 256)
				}
			}
		}
	} else {
		topic := services.TopicService.Get(favorite.EntityId)
		if topic == nil || topic.Status != model.TopicStatusOk {
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
	if message.Type == model.MsgTypeComment {
		entityType := gjson.Get(message.ExtraData, "entityType")
		entityId := gjson.Get(message.ExtraData, "entityId")
		if entityType.String() == model.EntityTypeArticle {
			detailUrl = urls.ArticleUrl(entityId.Int())
		} else if entityType.String() == model.EntityTypeTopic {
			detailUrl = urls.TopicUrl(entityId.Int())
		}
	}
	from := BuildUserDefaultIfNull(message.FromId)
	if message.FromId <= 0 {
		from.Nickname = "系统通知"
		from.Avatar = avatar.DefaultAvatars[0]
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

		// // 标记站外链接，搜索引擎爬虫不传递权重值
		// if !strings.Contains(href, "mlog.club") {
		// 	selection.SetAttr("rel", "external nofollow")
		// }

		// 内部跳转
		if len(href) > 0 && !urls.IsInternalUrl(href) {
			newHref := simple.ParseUrl(urls.AbsUrl("/redirect")).AddQuery("url", href).BuildStr()
			selection.SetAttr("href", newHref)
			selection.SetAttr("target", "_blank")
		}

		// 如果是锚链接
		if strings.Index(href, "#") == 0 {
			selection.ReplaceWithHtml(selection.Text())
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
		if strings.Contains(src, "qpic.cn") {
			newSrc := simple.ParseUrl("/api/img/proxy").AddQuery("url", src).BuildStr()
			selection.SetAttr("src", newSrc)
		}
	})

	html, err := doc.Html()
	if err != nil {
		return htmlContent
	}
	return html
}
