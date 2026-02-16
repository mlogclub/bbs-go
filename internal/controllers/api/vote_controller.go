package api

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

type VoteController struct {
	Ctx iris.Context
}

func (c *VoteController) GetBy(voteId int64) *web.JsonResult {
	vote := services.VoteService.Get(voteId)
	if vote == nil {
		return web.JsonErrorMsg("vote not found")
	}
	return web.JsonData(render.BuildVote(c.Ctx, vote))
}

func (c *VoteController) PostCast() *web.JsonResult {
	user := common.GetCurrentUser(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}

	var form req.VoteCastForm
	if err := c.Ctx.ReadJSON(&form); err != nil {
		return web.JsonError(err)
	}

	err := services.VoteService.Cast(user.Id, form)
	if err != nil {
		return web.JsonError(err)
	}

	vote := services.VoteService.Get(form.VoteId)
	return web.JsonData(render.BuildVote(c.Ctx, vote))
}
