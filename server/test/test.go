package main

import (
	"fmt"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/common/avatar"
)

func main() {
	// var fores []color.Color
	// for _, hexColor := range hexColors {
	// 	c, _ := colorToRGB(hexColor)
	// 	fores = append(fores, c)
	// }
	//
	// instance, _ := identicon.New(100, color.Transparent, fores...)
	// for i := 0; i < 10; i++ {
	// 	img := instance.Make([]byte(strconv.Itoa(i)))
	// 	fi, err := os.Create("/Users/gaoyoubo/Downloads/avatar/" + strconv.Itoa(i) + ".png")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	_ = png.Encode(fi, img)
	// 	_ = fi.Close()
	// }

	hs := hashset.New()
	var i int64 = 0
	for ; i < 1000; i++ {
		data, _ := avatar.Generate(i)
		md5 := simple.MD5Bytes(data)
		hs.Add(md5)

		fmt.Println(hs.Size())
	}

	// fmt.Println(len(data))
	// if err != nil {
	// 	panic(err)
	// }
	// {
	// 	fi, _ := os.Create("/Users/gaoyoubo/Downloads/avatar/fuck.png")
	// 	fi.Write(data)
	// }
}
