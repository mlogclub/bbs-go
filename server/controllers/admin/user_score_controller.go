package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/controllers/render"
	"bbs-go/services"
)

type UserScoreController struct {
	Ctx iris.Context
}

func (c *UserScoreController) GetBy(id int64) *simple.JsonResult {
	t := services.UserScoreService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (c *UserScoreController) AnyList() *simple.JsonResult {
	list, paging := services.UserScoreService.FindPageByParams(simple.NewQueryParams(c.Ctx).EqByReq("user_id").PageByReq().Desc("id"))
	var results []map[string]interface{}
	for _, userScore := range list {
		user := render.BuildUserDefaultIfNull(userScore.UserId)
		item := simple.NewRspBuilder(userScore).Put("user", user).Build()
		results = append(results, item)
	}
	return simple.JsonData(&simple.PageResult{Results: results, Page: paging})
}
