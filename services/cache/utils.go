package cache

import (
	"github.com/goburrow/cache"
	"github.com/sirupsen/logrus"
	"strconv"
)

func Key2Int64(key cache.Key) int64 {
	// str := key.(string)
	// return ToInt64(str)
	return key.(int64)
}

func ToInt64(str string) int64 {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		logrus.Error(err)
		return 0
	}
	return val
}
