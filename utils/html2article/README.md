## 基于文本密度的html2article实现[golang] 

## Install
	go get -u -v github.com/sundy-li/html2article


## Performance
 - Accuracy: `>= 98% `
 - Qps: 2w/s , 0.06ms/op ```
         go test -bench=.
	      BenchmarkExtract-4   	   20000	     66341 ns/op
	    ```
	      
 - 说明(对比其他开源实现,可能是目前最快的html2article实现,我们测试的数据集约3kw来自于微信公众号,各大类中文科技媒体历史文章,目前能达到98%以上准确率)
 - 除了必要dom解析以及时间解析, 为了高效率实现, 避免了过多的正则匹配


## Examples
参考examples
[from_url.go][1]

	
	package main

	import (
		"github.com/sundy-li/html2article"
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
		println("article publishtime is =>", article.Publishtime) //using UTC timezone
		println("article content is =>", article.Content)

		//parse the article to be readability
		article.Readable(urlStr)
		println("read=>", article.ReadContent)
	}

## Options

```
	ext.SetOption(&html2article.Option{
		AccurateTitle: true,  //Get the accurate title instead of from title tag
		RemoveNoise: false,  //Remove the noise node such as some footer
	})
```



## Algorithm
- [参考论文][2]
- [Java实现][3]


[1]: https://github.com/sundy-li/html2article/blob/master/examples/from_url.go
[2]: http://www.doc88.com/p-7714009813182.html
[3]: https://github.com/CrawlScript/WebCollector
 
