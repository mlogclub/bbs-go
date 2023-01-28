package order

import (
	"fmt"
	"math/rand"
	"time"
)

func GetOrder() string {
	t := time.Now()
	date := t.Format("2006010215")
	orderNum := fmt.Sprintf("%s%06v", date, rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	return orderNum
}
