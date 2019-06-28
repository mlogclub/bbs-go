# html2md
convert html to markdown

在GitHub上搜了下HTML转markdown的工具，并按照star从高到低逐个试了下，感觉不怎么符合自己的预期，索性自己写一个。


## HTML标签
并不是对所有的标签都做处理，比如`pre`、`blockquote`和`code`标签，这些没转成` ``` ` 或 `>` 或者是` ` `，因为markdown解析器解析不会有问题。

### 处理的标签
- h1~h6
- del
- b
- strong
- i
- em
- dfn
- var
- cite
- br
- span
- div
- figure
- p
- article
- nav
- footer
- header
- section
- table、thead、tbody、tr、th、td
- ul、ol、li
- hr

### 不作处理的标签
- pre
- blockquote
- code


## 转化效果
为了测试转化效果，我分别从github.io、gitbook、readthedoc三个站点随便提取了个正文的HTML内容，然后并将其转成markdown，看了下效果，比较符合自己的预期。
当然，效果并不可能是100%的。

## 使用方法

### go语言使用
1. 安装
`go get github.com/TruthHun/html2md`
1. 引入和调用
```go
mdStr:=html2md.Convert(htmlStr)
```

### 其他
已经编译打包了win、mac和linux的64位的可执行文件，在`bin`目录下

windows使用：
```
html2md.exe input.html output.md
```

mac/linux使用：
```
html2md input.html output.md
```

其它语言，直接使用cmd调用二进制可执行文件对文档进行处理即可

## 支持我
如果您使用了当前包或程序，遇到问题，向我反馈就是对我最好的支持；如果项目帮到了您，给当前项目一个star，也是对我莫大的认可与支持。
