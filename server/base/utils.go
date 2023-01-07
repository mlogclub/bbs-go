package base

import (
	"github.com/mlogclub/simple/common/digests"
	"math/rand"
	"time"
)

func GetScore(score int64) int64 {
	return score / 2
}

func GetLevel(score int64) int64 {
	return score / 1000
}

func Get16MD5Encode(str string) string {
	return digests.MD5Bytes([]byte(str))[8:24]
}

func RandScore(max int64) int64 {
	var timeStamp = time.Now().Unix()
	r := rand.New(rand.NewSource(timeStamp))
	score := r.Int63n(max)
	if score < 5 {
		return 5
	}
	return score
}
