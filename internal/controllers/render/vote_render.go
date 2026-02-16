package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/resp"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
)

func BuildVote(ctx iris.Context, vote *models.Vote) *resp.VoteResponse {
	if vote == nil {
		return nil
	}

	ret := &resp.VoteResponse{
		Id:          vote.Id,
		Type:        vote.Type,
		Title:       vote.Title,
		ExpiredAt:   vote.ExpiredAt,
		VoteNum:     vote.VoteNum,
		OptionCount: vote.OptionCount,
		VoteCount:   vote.VoteCount,
		Expired:     dates.NowTimestamp() > vote.ExpiredAt,
	}

	currentUserId := common.GetCurrentUserID(ctx)

	var selectedMap map[int64]bool
	if currentUserId > 0 {
		if record := services.VoteRecordService.GetBy(currentUserId, vote.Id); record != nil {
			ret.Voted = true
			ret.OptionIds = services.VoteService.ParseOptionIds(record.OptionIds)
			selectedMap = make(map[int64]bool, len(ret.OptionIds))
			for _, optionId := range ret.OptionIds {
				selectedMap[optionId] = true
			}
		}
	}

	options := services.VoteOptionService.FindByVoteId(vote.Id)
	for _, option := range options {
		item := resp.VoteOptionResponse{
			Id:        option.Id,
			Content:   option.Content,
			SortNo:    option.SortNo,
			VoteCount: option.VoteCount,
		}
		if vote.VoteCount > 0 {
			item.Percent = float64(option.VoteCount) / float64(vote.VoteCount) * 100
		}
		if selectedMap != nil {
			item.Voted = selectedMap[option.Id]
		}
		ret.Options = append(ret.Options, item)
	}
	return ret
}
