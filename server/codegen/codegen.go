package main

import (
	"github.com/mlogclub/simple"

	"bbs-go/model"
)

func main() {
	simple.Generate("./", "bbs-go", simple.GetGenerateStruct(&model.UserScore{}))
	simple.Generate("./", "bbs-go", simple.GetGenerateStruct(&model.UserScoreLog{}))
}
