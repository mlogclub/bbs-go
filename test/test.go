package main

import (
	"fmt"
	"strings"

	"github.com/mlogclub/simple"
	uuid "github.com/satori/go.uuid"
)

func main() {
	u := uuid.NewV4()
	s := u.String()
	s = strings.ReplaceAll(s, "-", "")
	fmt.Println(s)
	fmt.Println(simple.RuneLen(s))
	// fmt.Println(uuid.NewRandom())
}
