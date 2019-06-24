// COPYRIGHT https://github.com/golang/tools/blob/master/cmd/html2article/conv.go
package html2article

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type selector func(*html.Node) bool
type Style string

var (
	timeRegex = []*regexp.Regexp{
		regexp.MustCompile(`([\d]{4})-([\d]{1,2})-([\d]{1,2})\s+([\d]{1,2}:[\d]{1,2})?`),
		regexp.MustCompile(`([\d]{4})-([\d]{1,2})-([\d]{1,2})`),
		regexp.MustCompile(`([\d]{4})\.([\d]{1,2})\.([\d]{1,2})\s+([\d]{1,2}:[\d]{1,2})?`),
		regexp.MustCompile(`([\d]{4})/([\d]{1,2})/([\d]{1,2})\s+([\d]{1,2}:[\d]{1,2})?`),
		regexp.MustCompile(`([\d]{4})\s*年\s*([\d]{1,2})\s*月\s*([\d]{1,2})\s*日\s*([\d]{1,2}:[\d]{1,2})?`),
		regexp.MustCompile(`([\d]{1,2})\s*天前`),
		regexp.MustCompile(`([\d]{1,2})\s*小时前`),
		regexp.MustCompile(`([\d]{1,2})\s*分钟前`),
	}
)

func countStopWords(str string) int {
	var c = 0
	for k, _ := range stopwords {
		c += strings.Count(str, k)
	}
	if c == 0 {
		c = 1
	}
	return c
}

func countSn(str string) int {
	sn := strings.Count(str, ".")
	sn += strings.Count(str, "。")
	if sn == 0 {
		sn = 1
	}
	return sn
}

func countChar(str string) int {
	return len([]rune(Compress(str)))
}

// 如果str是北京时间，返回的时间戳会比当前快8小时
func getTime(str string) int64 {
	fn := func(year int, month int, day int, hour int, minute int) int64 {
		v := fmt.Sprintf("%04d%02d%02d %02d:%02d", year, month, day, hour, minute)
		tm, err := time.Parse("20060102 15:04", v)
		if err == nil {
			return tm.Unix()
		}
		return 0
	}
	for i, t := range timeRegex {
		ts := t.FindStringSubmatch(str)
		switch i {
		case 5:
			if len(ts) == 2 {
				d, _ := strconv.Atoi(ts[1])
				t := time.Now().Add(-time.Hour * time.Duration(24*d))
				return fn(t.Year(), int(t.Month()), t.Day(), 0, 0)
			}
			continue
		case 6:
			if len(ts) == 2 {
				h, _ := strconv.Atoi(ts[1])
				t := time.Now().Add(-time.Hour * time.Duration(h))
				return fn(t.Year(), int(t.Month()), t.Day(), t.Hour(), 0)
			}
			continue

		case 7:
			if len(ts) == 2 {
				h, _ := strconv.Atoi(ts[1])
				t := time.Now().Add(-time.Minute * time.Duration(h))
				return fn(t.Year(), int(t.Month()), t.Day(), t.Hour(), 0)
			}
			continue
		default:
			if len(ts) < 4 {
				continue
			}
			var h string = "00:00"
			if len(ts) > 4 && ts[4] != "" {
				h = ts[4]
			}

			year, _ := strconv.Atoi(ts[1])
			month, _ := strconv.Atoi(ts[2])
			day, _ := strconv.Atoi(ts[3])

			timeAt := strings.Split(h, ":")
			hour, _ := strconv.Atoi(timeAt[0])
			var minute int
			if len(timeAt) > 1 {
				minute, _ = strconv.Atoi(timeAt[1])
			}
			return fn(year, month, day, hour, minute)
		}
	}
	return 0
}

// get Text and transform the charset
func getText(n *html.Node, filter ...selector) string {
	return Compress(strings.TrimSpace(text(n, filter...)))
}

// 这个暂时不用,因为code标签还不好识别
func CompressHtml(str string) string {
	buf := make([]byte, 0, len(str)/2)
	buffer := bytes.NewBuffer(buf)

	flag := false  // 标识当前是否已经有一个空格
	inTag := false // 标识是否在tag里头

	rs := ([]rune)(str)
	for i, r := range rs {
		if r == '<' {
			inTag = true
		} else if (inTag && r == '/') && (i-1 >= 0 && rs[i-1] == '<') {
			inTag = false
		}
		if !inTag {
			if unicode.IsSpace(r) {
				if flag {
					continue
				} else {
					flag = true
					r = ' '
				}
			} else {
				flag = false
			}
		}
		buffer.WriteRune(r)
	}
	return buffer.String()
}

// 压缩字符串
// 将多个空格字符压缩为一个空格
func Compress(str string) string {
	buf := make([]byte, 0, len(str)/2)
	buffer := bytes.NewBuffer(buf)

	flag := false // 标识当前是否已经有一个空格
	for _, r := range str {
		if unicode.IsSpace(r) {
			if flag {
				continue
			} else {
				flag = true
				r = ' '
			}
		} else {
			flag = false
		}
		buffer.WriteRune(r)
	}
	return buffer.String()
}

func text(n *html.Node, filter ...selector) string {
	if isTag(atom.Style)(n) || isTag(atom.Script)(n) || isTag(atom.Image)(n) || isTag(atom.Img)(n) || isTag(atom.Textarea)(n) || isTag(atom.Input)(n) || isTag(atom.Noscript)(n) {
		return ""
	}
	var buf bytes.Buffer

	switch n.Type {
	case html.TextNode:
		buf.WriteString(n.Data)
	case html.ElementNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			ok := true
			for _, f := range filter {
				if !f(c) {
					ok = false
					break
				}
			}
			if ok {
				fmt.Fprint(&buf, text(c, filter...))
			}
		}
	}
	return buf.String()
}

func getHtml(n *html.Node) (str string, err error) {
	var buf bytes.Buffer
	err = html.Render(&buf, n)
	str = buf.String()
	return
}

func getImages(node *html.Node) []string {
	res := []string{}
	mp := make(map[string]bool)
	walk(node, func(n *html.Node) bool {
		if isTag(atom.Img)(n) {
			if width, err := strconv.Atoi(attr(n, "width")); err == nil {
				if width != 0 && width < 30 {
					return false
				}
			}

			if height, err := strconv.Atoi(attr(n, "height")); err == nil {
				if height != 0 && height < 30 {
					return false
				}
			}

			// 不抓取默认不展示图片
			if display := attr(n.Parent, "style"); len(display) > 0 && strings.Contains(display, "display: none") {
				return false
			}

			if display := attr(n, "style"); len(display) > 0 && strings.Contains(display, "display: none") {
				return false
			}

			atts := []string{"data-original", "data-echo", "data-src"}
			src := ""
			for _, a := range atts {
				src = attr(n, a)
				if len(src) > 0 {
					setAttr(n, "src", src)
					removeAttr(n, a)
					break
				}
			}
			if len(src) == 0 {
				src = attr(n, "src")
			}
			excludeStrs := []string{
				"w16_h16.png",
				"logo.png",
				"icon.png",
			}

			if len(src) > 0 {
				for _, exc := range excludeStrs {
					if strings.Contains(src, exc) {
						return false
					}
				}
			}

			if _, ok := mp[src]; !ok && len(src) > 0 {
				mp[src] = true
				res = append(res, src)
			}
			return false
		} else {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				res = append(res, getImages(c)...)
			}
			return false
		}
	})
	return res
}

func isContentNode(n *html.Node) bool {
	return n.DataAtom == atom.Div || n.DataAtom == atom.Section || n.DataAtom == atom.Table || n.DataAtom == atom.Tr || n.DataAtom == atom.Td || n.DataAtom == atom.Tbody || n.DataAtom == atom.Article
}

func isNoisingNode(n *html.Node) bool {
	switch n.DataAtom {
	case atom.Div, atom.Script, atom.Ul, atom.Tr, atom.Section, atom.Footer:
		return true
	}
	return false
}

func isTitleNode(n *html.Node) bool {
	switch n.DataAtom {
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		return true
	case atom.A, atom.Span, atom.Td, atom.Tr, atom.Div, atom.Header:
		return true
	}

	return false
}

func isHNode(n *html.Node) bool {
	switch n.DataAtom {
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		return true
	}
	return false
}

func isTag(a atom.Atom) selector {
	return func(n *html.Node) bool {
		return n.DataAtom == a
	}
}

// func hasContent(str string) selector {
// 	return func(n *html.Node) bool {
// 		return n.Data
// 	}
// }

func alwaysTrue() selector {
	return func(n *html.Node) bool {
		return true
	}
}

func hasAttr(key, val string) selector {
	return func(n *html.Node) bool {
		for _, a := range n.Attr {
			if a.Key == key && a.Val == val {
				return true
			}
		}
		return false
	}
}

func attr(node *html.Node, key string) (value string) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func findAll(node *html.Node, fn selector) (nodes []*html.Node) {
	walk(node, func(n *html.Node) bool {
		if fn(n) {
			nodes = append(nodes, n)
		}
		return true
	})
	return
}

func find(n *html.Node, fn selector) *html.Node {
	var result *html.Node
	walk(n, func(n *html.Node) bool {
		if result != nil {
			return false
		}
		if fn(n) {
			result = n
			return false
		}
		return true
	})
	return result
}

func walk(n *html.Node, fn selector) {
	if fn(n) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c, fn)
		}
	}
}

func walkRemove(n *html.Node, fn selector, del selector) {
	if del(n) {
		travesRemove(n)
		return
	}
	if fn(n) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walkRemove(c, fn, del)
		}
	}
}

// remove node n when using ` c := node.FirstChild; c != nil; c = c.NextSibling` traves
func travesRemove(n *html.Node) {
	next := n.NextSibling
	if n.Parent != nil {
		n.Parent.RemoveChild(n)
	}
	n.NextSibling = next
}

// remove br node
func brRemove(n *html.Node) {
	next := n.NextSibling
	if n.Parent != nil {
		if n.Parent.NextSibling != nil && n.Parent.Type == html.TextNode && n.Parent.NextSibling.Type == html.TextNode {
			n.Parent.Data += fmt.Sprintf("<%s>", n.Data)
			n.Parent.Data += n.Parent.NextSibling.Data
			travesRemove(n.Parent.NextSibling)
		}
	}
	n.NextSibling = next
}

func removeAttr(n *html.Node, attrName string) {
	for i, a := range n.Attr {
		if a.Key == attrName {
			n.Attr[i], n.Attr[len(n.Attr)-1], n.Attr =
				n.Attr[len(n.Attr)-1], html.Attribute{}, n.Attr[:len(n.Attr)-1]
			return
		}
	}
}

func setAttr(n *html.Node, attrName, value string) {
	for i, a := range n.Attr {
		if a.Key == attrName {
			n.Attr[i].Val = value
			return
		}
	}
	if len(n.Attr) == 0 {
		n.Attr = make([]html.Attribute, 1)
	}
	n.Attr = append(n.Attr, html.Attribute{
		Key: attrName,
		Val: value,
	})
}

func diffString(a, b string) int {
	aa := []rune(a)
	bb := []rune(b)
	return diffRune(aa, bb)
}

func diffRune(a, b []rune) int {
	var tmp []rune
	if len(a) < len(b) {
		tmp = a
		a = b
		b = tmp
	}
	// a should be bigger
	mp := make(map[rune]bool, len(a))
	mpb := make(map[rune]bool, len(b))
	for _, r := range a {
		mp[r] = true
	}
	for _, r := range b {
		mpb[r] = true
	}
	for k, _ := range mpb {
		if _, ok := mp[k]; ok {
			delete(mp, k)
			delete(mpb, k)
		}
	}
	return len(mpb) + len(mp)
}

func distanceExit(a, b string, maxValue int) (dis int, ok bool) {
	aa := []rune(a)
	bb := []rune(b)
	return distanceRuneExit(aa, bb, len(aa), len(bb), maxValue)
}

func distanceRuneExit(a, b []rune, s1 int, s2 int, maxValue int) (dis int, ok bool) {
	if maxValue < 0 {
		return
	}
	if s1 == 0 || s2 == 0 {
		dis = max(s1, s2)
		if dis > maxValue {
			return
		}
		ok = true
		return
	}
	if a[0] == b[0] {
		return distanceRuneExit(a[1:], b[1:], s1-1, s2-1, maxValue)
	}
	var v = (1 << 31) - 1

	if vv, ok := distanceRuneExit(a[1:], b[1:], s1-1, s2-1, maxValue-1); ok && vv+1 < v {
		v = vv + 1
	}
	if vv, ok := distanceRuneExit(a, b[1:], s1, s2-1, maxValue-1); ok && vv+1 < v {
		v = vv + 1
	}
	if vv, ok := distanceRuneExit(a[1:], b, s1-1, s2, maxValue-1); ok && vv+1 < v {
		v = vv + 1
	}
	if v == (1<<31)-1 {
		return
	}
	dis = v
	ok = true
	return
}
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
