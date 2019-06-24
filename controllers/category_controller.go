package controllers

import (
	"github.com/kataras/iris"

	"github.com/mlogclub/mlog/services"
)

type CategoryController struct {
	Ctx             iris.Context
	CategoryService *services.CategoryService
	TagService      *services.TagService
	ArticleService  *services.ArticleService
}
