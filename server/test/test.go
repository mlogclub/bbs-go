package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mlogclub/simple"

	"bbs-go/common/config"
	"bbs-go/common/uploader"
)

func main() {
	// key := generateImageKey([]byte("xx"))
	// dir := filepath.Dir(key)
	// fmt.Println(dir)

	config.InitConfig("./bbs-go.yaml")
	// url, err := uploader.NewLocal().CopyImage("https://ss0.bdstatic.com/94oJfD_bAAcT8t7mm9GUKT-xh_/timg?image&quality=100&size=b4000_4000&sec=1585374026&di=d93583edb884226cf7ddb4e4d9354c01&src=http://00.minipic.eastday.com/20170525/20170525151413_87ad5f4c0d7a5798892bdf21d6323cf6_6.jpeg")
	url, err := uploader.NewAliyun().CopyImage("https://ss0.bdstatic.com/94oJfD_bAAcT8t7mm9GUKT-xh_/timg?image&quality=100&size=b4000_4000&sec=1585374026&di=d93583edb884226cf7ddb4e4d9354c01&src=http://00.minipic.eastday.com/20170525/20170525151413_87ad5f4c0d7a5798892bdf21d6323cf6_6.jpeg")
	fmt.Println(url)
	fmt.Println(err)
}

func generateImageKey(data []byte) string {
	md5 := simple.MD5Bytes(data)
	return filepath.Join("images", simple.TimeFormat(time.Now(), "2006/01/02/"), md5+".jpg")
}
