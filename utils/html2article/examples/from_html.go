package main

import (
	"io/ioutil"
	"net/http"

	"github.com/mlogclub/mlog/utils/html2article"
)

func main() {
	urlStr := "https://www.leiphone.com/news/201602/DsiQtR6c1jCu7iwA.html"
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)

	ext, err := html2article.NewFromHtml(string(bs))
	if err != nil {
		panic(err)
	}
	article, err := ext.ToArticle()
	if err != nil {
		panic(err)
	}
	println("article title is =>", article.Title)
	println("article publishtime is =>", article.Publishtime)
	println("article content is =>", article.Content)
}
