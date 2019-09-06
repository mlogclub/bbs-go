package main

import (
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

func main() {
	simple.Generate("./", "github.com/mlogclub/bbs-go", simple.GetGenerateStruct(&model.CollectRule{}))
	simple.Generate("./", "github.com/mlogclub/bbs-go", simple.GetGenerateStruct(&model.CollectLink{}))
	simple.Generate("./", "github.com/mlogclub/bbs-go", simple.GetGenerateStruct(&model.CollectArticle{}))
}
