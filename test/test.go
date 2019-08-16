package main

import (
	"fmt"

	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/services/collect"
)

func main() {
	p := collect.CollectProject("https://studygolang.com/p/gomybatis")
	fmt.Println(simple.FormatJson(p))
}
