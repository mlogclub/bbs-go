package main

import (
	"github.com/mlogclub/simple/codegen"

	"bbs-go/model"
)

func main() {
	codegen.Generate("./", "bbs-go", codegen.GetGenerateStruct(&model.CheckIn{}))
}
