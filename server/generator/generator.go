package main

import (
	"github.com/mlogclub/codegen"

	"bbs-go/model"
)

func main() {
	codegen.Generate("./", "bbs-go", codegen.GetGenerateStruct(&model.CheckIn{}))
}
