package html2md

import (
	"io/ioutil"
	"testing"
)

func TestConvert(t *testing.T) {
	b, _ := ioutil.ReadFile("example/code.html")
	md := Convert(string(b))
	ioutil.WriteFile("example/code.md", []byte(md), 0777)

	b, _ = ioutil.ReadFile("example/hello.html")
	md = Convert(string(b))
	ioutil.WriteFile("example/hello.md", []byte(md), 0777)

	b, _ = ioutil.ReadFile("example/gitbook.html")
	md = Convert(string(b))
	ioutil.WriteFile("example/gitbook.md", []byte(md), 0777)

	b, _ = ioutil.ReadFile("example/github.io.html")
	md = Convert(string(b))
	ioutil.WriteFile("example/github.io.md", []byte(md), 0777)

	b, _ = ioutil.ReadFile("example/readthedoc.html")
	md = Convert(string(b))
	ioutil.WriteFile("example/readthedoc.md", []byte(md), 0777)

	b, _ = ioutil.ReadFile("example/gitee.html")
	md = Convert(string(b))
	ioutil.WriteFile("example/gitee.md", []byte(md), 0777)

	b, _ = ioutil.ReadFile("example/gitbook1.html")
	md = Convert(string(b))
	ioutil.WriteFile("example/gitbook1.md", []byte(md), 0777)

	b, _ = ioutil.ReadFile("example/django.html")
	md = Convert(string(b))
	ioutil.WriteFile("example/django.md", []byte(md), 0777)

	b, _ = ioutil.ReadFile("example/blockquote.html")
	md = Convert(string(b))
	ioutil.WriteFile("example/blockquote.md", []byte(md), 0777)
}
