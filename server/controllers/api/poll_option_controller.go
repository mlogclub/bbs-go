package api

import (
	"bbs-go/model"
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
	"strconv"
)

type PollOptionController struct {
	Ctx iris.Context
}

type PollOptionWithVoteNumber struct {
	model.PollOption
	VoteNumber int `json:"voteNumber" form:"voteNumber"`
}

// Get poll information by TopicID
func (c *PollOptionController) GetBy(id int64) *simple.JsonResult {
	t := services.PollOptionService.Find(simple.NewSqlCnd().Eq("topic_id", id))

	if t == nil {
		return simple.JsonErrorMsg("Poll info not found, TopicID=" + strconv.FormatInt(id, 10))
	}
	resultWithTotalInfo := make([]PollOptionWithVoteNumber, len(t))

	for index, option := range t {
		cnd := simple.NewSqlCnd("id").Eq("topic_id", id).Eq("poll_option_id", option.OptionId)
		resultWithTotalInfo[index].PollOption = option
		resultWithTotalInfo[index].VoteNumber = services.PollAnswerService.Count(cnd)
	}

	return simple.JsonData(resultWithTotalInfo)
}
