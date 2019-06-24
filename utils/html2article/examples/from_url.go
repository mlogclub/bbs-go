package main

import (
	"github.com/mlogclub/mlog/utils/html2article"
)

func main() {
	urlStr := "https://www.leiphone.com/news/201602/DsiQtR6c1jCu7iwA.html"
	ext, err := html2article.NewFromUrl(urlStr)
	if err != nil {
		panic(err)
	}
	article, err := ext.ToArticle()
	if err != nil {
		panic(err)
	}
	println("article title is =>", article.Title)
	println("article publishtime is =>", article.Publishtime) // using UTC timezone
	println("article content is =>", article.Content)

	// parse the article to be readability
	article.Readable(urlStr)
	println("read=>", article.ReadContent)
}
