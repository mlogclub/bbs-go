package render

import (
	"bbs-go/model"
	"bbs-go/model/constants"
	"bbs-go/pkg/common"
	"bbs-go/pkg/html"
	"bbs-go/pkg/markdown"
	"bbs-go/pkg/text"
)

func BuildProject(project *model.Project) *model.ProjectResponse {
	if project == nil {
		return nil
	}
	rsp := &model.ProjectResponse{}
	rsp.ProjectId = project.Id
	rsp.User = BuildUserInfoDefaultIfNull(project.UserId)
	rsp.Name = project.Name
	rsp.Title = project.Title
	rsp.Logo = project.Logo
	rsp.Url = project.Url
	rsp.Url = project.Url
	rsp.DocUrl = project.DocUrl
	rsp.CreateTime = project.CreateTime

	if project.ContentType == constants.ContentTypeHtml {
		rsp.Content = handleHtmlContent(project.Content)
		rsp.Summary = text.GetSummary(html.GetHtmlText(project.Content), constants.SummaryLen)
	} else {
		content := markdown.ToHTML(project.Content)
		summary := html.GetSummary(content, constants.SummaryLen)
		rsp.Content = handleHtmlContent(content)
		rsp.Summary = summary
	}

	return rsp
}

func BuildSimpleProjects(projects []model.Project) []model.ProjectSimpleResponse {
	if len(projects) == 0 {
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
	rsp.User = BuildUserInfoDefaultIfNull(project.UserId)
	rsp.Name = project.Name
	rsp.Title = project.Title
	rsp.Logo = project.Logo
	rsp.Url = project.Url
	rsp.DocUrl = project.DocUrl
	rsp.DownloadUrl = project.DownloadUrl
	rsp.CreateTime = project.CreateTime

	if project.ContentType == constants.ContentTypeHtml {
		rsp.Summary = text.GetSummary(html.GetHtmlText(project.Content), constants.SummaryLen)
	} else {
		rsp.Summary = common.GetMarkdownSummary(project.Content)
	}

	return rsp
}
