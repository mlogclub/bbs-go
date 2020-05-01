pangu.go
========

[![](https://img.shields.io/travis/vinta/pangu.svg?style=flat-square)](https://travis-ci.org/vinta/pangu)
[![](https://img.shields.io/codecov/c/github/vinta/pangu/master.svg?style=flat-square)](https://codecov.io/github/vinta/pangu)
[![](https://img.shields.io/badge/made%20with-%e2%9d%a4-ff69b4.svg?style=flat-square)](http://vinta.ws)

Paranoid text spacing for good readability, to automatically insert whitespace between CJK (Chinese, Japanese, Korean) and half-width characters (alphabetical letters, numerical digits and symbols).

* Go version: [pangu.go](https://github.com/vinta/pangu)
* Java version: [pangu.java](https://github.com/vinta/pangu.java)
* JavaScript version: [pangu.js](https://github.com/vinta/paranoid-auto-spacing)
* Node.js version: [pangu.node](https://github.com/huei90/pangu.node)
* Python version: [pangu.py](https://github.com/vinta/pangu.py)
* Ruby version: [pangu.rb](https://github.com/dlackty/pangu.rb)

## Installation

To install the package, `pangu`, for using in your Go programs:

``` bash
$ go get -u github.com/vinta/pangu
```

To install the command-line tool, `pangu-axe`:

``` bash
$ go get -u github.com/vinta/pangu/pangu-axe
```

## Usage

### Package

``` go
package main

import (
    "fmt"
    "github.com/vinta/pangu"
)

func main() {
    s := pangu.SpacingText("新八的構造成分有95%是眼鏡、3%是水、2%是垃圾")
    fmt.Println(s)
    // Output:
    // 新八的構造成分有 95% 是眼鏡、3% 是水、2% 是垃圾
}
```

### Command-line Interface

``` bash
$ pangu-axe text "所以,請問Jackey的鼻子有幾個?3.14個"
所以, 請問 Jackey 的鼻子有幾個? 3.14 個

$ pangu-axe file 銀河便車指南.txt
$ pangu-axe file 宇宙盡頭的餐廳.txt -o 宇宙盡頭的餐廳（好讀版）.txt
$ pangu-axe file 生命、宇宙及萬事萬物.txt 再見，謝謝你的魚.txt 基本無害.txt
```

## Documentation

* `pangu` on [GoDoc](http://godoc.org/github.com/vinta/pangu)
* `pangu-axe` on [GoDoc](http://godoc.org/github.com/vinta/pangu/pangu-axe)

Have a question? Ask it on the [GitHub issues](https://github.com/vinta/pangu/issues)!
