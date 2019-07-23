package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

var h = `
<nav>
	<ul>
		<li>
			<a href="#toc_0">Fuck</a>
			<ul>
				<li><a href="#toc_0">为什么是 Go 语言</a></li>
	
				<li><a href="#toc_1">Go 语言简介</a></li>
	
				<li><a href="#toc_2">数据类型</a></li>
	
				<li><a href="#toc_3">基本语法</a>
					<ul>
						<li><a href="#toc_4">HelloWorld</a></li>
	
						<li><a href="#toc_5">变量</a>
							<ul>
								<li><a href="#toc_6">变量声明</a></li>
	
								<li><a href="#toc_7">类型推断</a></li>
							</ul>
						</li>
	
						<li><a href="#toc_8">函数</a></li>
	
						<li><a href="#toc_9">结构体</a></li>
	
						<li><a href="#toc_10">指针类型和值类型</a></li>
	
						<li><a href="#toc_11">并发编程</a></li>
					</ul>
				</li>
	
				<li><a href="#toc_12">Java 程序员觉得不好用的地方</a></li>
	
				<li><a href="#toc_13">参考</a></li>
			</ul>
		</li>
	</ul>
</nav>
`

func main() {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(h))
	// doc.Find("nav > ul > li").Each(func(i int, selection *goquery.Selection) {
	// 	fmt.Println(i)
	// })
	fmt.Println(toc(doc))
}

func toc(doc *goquery.Document) string {
	top := doc.Find("nav > ul > li")
	if top.Size() != 1 {
		return ""
	}
	topA := doc.Find("nav > ul > li > a")
	if topA.Size() != 0 {
		return ""
	}
	tocHtml, _ := top.Html()
	return tocHtml
}
