module github.com/mlogclub/bbs-go

go 1.13

replace github.com/go-check/check => github.com/go-check/check v0.0.0-20180628173108-788fd7840127

require (
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/aliyun/aliyun-oss-go-sdk v2.0.3+incompatible
	github.com/antchfx/htmlquery v1.1.0 // indirect
	github.com/antchfx/xmlquery v1.1.0 // indirect
	github.com/antchfx/xpath v1.1.0 // indirect
	github.com/emirpasic/gods v1.12.0
	github.com/go-resty/resty/v2 v2.1.0
	github.com/goburrow/cache v0.1.0
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gocolly/colly v1.2.0
	github.com/gorilla/feeds v1.1.1
	github.com/ikeikeikeike/go-sitemap-generator/v2 v2.0.2
	github.com/iris-contrib/middleware/cors v0.0.0-20191028172159-41f72a73786a
	github.com/issue9/identicon v1.0.1
	github.com/jinzhu/gorm v1.9.11
	github.com/jordan-wright/email v0.0.0-20190819015918-041e0cec78b0
	github.com/kataras/iris/v12 v12.0.1
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/mattn/go-runewidth v0.0.6 // indirect
	github.com/mattn/godown v0.0.0-20180312012330-2e9e17e0ea51
	github.com/mlogclub/simple v1.0.46 // currently it is v1.0.42 but based on the latest PR it should be updated to 1.0.43 (that's why it is presented as this here) in order to be compatible with this project.
	github.com/robfig/cron v1.2.0
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/sundy-li/html2article v0.0.0-20180131134645-09ac198090c2
	github.com/temoto/robotstxt v1.1.1 // indirect
	github.com/tidwall/gjson v1.3.4
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/resty.v1 v1.12.0
	gopkg.in/yaml.v2 v2.2.5
)
