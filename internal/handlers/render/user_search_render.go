package render

import (
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/search"
)

func BuildSearchUsers(docs []search.UserDocument) []resp.SearchUserResponse {
	var items []resp.SearchUserResponse
	for _, doc := range docs {
		items = append(items, BuildSearchUser(doc))
	}
	return items
}

func BuildSearchUser(doc search.UserDocument) resp.SearchUserResponse {
	user := BuildUserInfoDefaultIfNull(doc.Id)
	if doc.Nickname != "" {
		user.Nickname = doc.Nickname
	}
	if doc.Description != "" {
		user.Description = doc.Description
	}
	return resp.SearchUserResponse{
		User:        user,
		Nickname:    doc.Nickname,
		Username:    doc.Username,
		Description: doc.Description,
		CreateTime:  doc.CreateTime,
	}
}
