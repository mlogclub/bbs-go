package main

import (
	"fmt"

	"github.com/mlogclub/bbs-go/common"
)

func main() {
	fmt.Println(common.ApplyImageStyle("http://static.mlog.club/images/default-avatar/default.png", "avatar"))
	fmt.Println(common.ApplyImageStyle("https://file.mlog.club/images/2019/08/07/ff14571273b239543a1ce6b4de48d85d.jpg", "avatar"))
	fmt.Println(common.ApplyImageStyle("https://file.mlog.club/images/2019/08/07/ff14571273b239543a1ce6b4de48d85d.jpg!avatar", "avatar"))
}
