package baiduai

import (
	"fmt"
	"testing"

	"github.com/mlogclub/simple"
)

func TestAi(t *testing.T) {
	title := "iphone手机出现“白苹果”原因及解决办法，用苹果手机的可以看下"
	content := `如果下面的方法还是没有解决你的问题建议来我们门店看下成都市锦江区红星路三段99号银石广场24层01室。在通电的情况下掉进清水，这种情况一不需要拆机处理。尽快断电。用力甩干，但别把机器甩掉，主意要把屏幕内的水甩出来。如果屏幕残留有水滴，干后会有痕迹。^H3 放在台灯，射灯等轻微热源下让水分慢慢散去。`
	tags := GetTags(title, content)
	fmt.Println(simple.FormatJson(tags))

	// categories := GetCategories(title, content)
	// fmt.Println(simple.FormatJson(categories))
	//
	// summary, err := GetNewsSummary(title, content, 100)
	// fmt.Println(summary, err)

	ret, err := AnalyzeHtml(title, content)
	fmt.Println(simple.FormatJson(ret))
	fmt.Println(err)
}
