package main

import (
	"bbs-go/common/config"
	"bbs-go/common/uploader"
	"fmt"
	"testing"
)

func TestMinioClient(t *testing.T){
	config.InitConfig("../bbs-go.yaml")

	// url, err := uploader.NewLocal().CopyImage("https://ss0.bdstatic.com/94oJfD_bAAcT8t7mm9GUKT-xh_/timg?image&quality=100&size=b4000_4000&sec=1585374026&di=d93583edb884226cf7ddb4e4d9354c01&src=http://00.minipic.eastday.com/20170525/20170525151413_87ad5f4c0d7a5798892bdf21d6323cf6_6.jpeg")
	url, err := uploader.CopyImage("https://ss0.bdstatic.com/94oJfD_bAAcT8t7mm9GUKT-xh_/timg?image&quality=100&size=b4000_4000&sec=1585374026&di=d93583edb884226cf7ddb4e4d9354c01&src=http://00.minipic.eastday.com/20170525/20170525151413_87ad5f4c0d7a5798892bdf21d6323cf6_6.jpeg")
	fmt.Println(url)
	fmt.Println(err)
}

