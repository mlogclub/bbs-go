package main

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/codegen"
)

func main() {
	codegen.Generate(
		"./",
		"bbs-go",
		1,
		codegen.GetGenerateStruct(&models.Migration{}),
		codegen.GetGenerateStruct(&models.TaskConfig{}),
		codegen.GetGenerateStruct(&models.UserTaskEvent{}),
		codegen.GetGenerateStruct(&models.UserTaskLog{}),
		codegen.GetGenerateStruct(&models.Badge{}),
		codegen.GetGenerateStruct(&models.UserBadge{}),
		codegen.GetGenerateStruct(&models.LevelConfig{}),
		codegen.GetGenerateStruct(&models.UserExpLog{}),
		codegen.GetGenerateStruct(&models.Vote{}),
		codegen.GetGenerateStruct(&models.VoteOption{}),
		codegen.GetGenerateStruct(&models.VoteRecord{}),
	)
}
