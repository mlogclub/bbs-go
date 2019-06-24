package main

import (
	"github.com/mlogclub/simple/codegen"

	"github.com/mlogclub/mlog/model"
)

func main() {
	codegen.Generate("./", "github.com/mlogclub/mlog", codegen.GetGenerateStruct(&model.OauthClient{}))
	codegen.Generate("./", "github.com/mlogclub/mlog", codegen.GetGenerateStruct(&model.OauthToken{}))
}
