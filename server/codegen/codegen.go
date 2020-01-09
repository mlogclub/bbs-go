package main

import (
	"bbs-go/simple"

	"bbs-go/model"
)

func main() {
	simple.Generate("./", "bbs-go", simple.GetGenerateStruct(&model.TopicNode{}))
}
