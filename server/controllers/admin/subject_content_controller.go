package admin

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/services"
)

type SubjectContentController struct {
	Ctx iris.Context
}

func (this *SubjectContentController) GetBy(id int64) *simple.JsonResult {
	t := services.SubjectContentService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *SubjectContentController) AnyList() *simple.JsonResult {
	list, paging := services.SubjectContentService.FindPageByParams(simple.NewQueryParams(this.Ctx).EqByReq("subject_id").
		EqByReq("entity_type").EqByReq("entity_id").EqByReq("deleted").PageByReq().Desc("id"))

	var itemList []map[string]interface{}
	for _, v := range list {
		subject := services.SubjectService.Get(v.SubjectId)
		itemList = append(itemList, simple.NewRspBuilder(v).Put("subject", subject).Build())
	}
	return simple.JsonData(&simple.PageResult{Results: itemList, Page: paging})
}
