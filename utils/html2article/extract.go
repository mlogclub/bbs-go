package html2article

import (
	"errors"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strings"

	"github.com/TruthHun/gotil/util"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/net/html/charset"
)

type extractor struct {
	data   map[*Info]*html.Node
	urlStr string
	doc    *html.Node

	maxDensity       float64
	sn               float64
	swn              float64
	title            string
	accurateTitle    string
	titleDistanceMin int
	titleMatchLen    int

	option *Option
}

type Option struct {
	RemoveNoise   bool // remove noise node
	AccurateTitle bool // find the accurate title node
	UserAgent     string
}

func NewFromHtml(htmlStr string) (ext *extractor, err error) {
	return NewFromReader(strings.NewReader(htmlStr))
}

func NewFromReader(reader io.Reader) (ext *extractor, err error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return
	}
	return NewFromNode(doc)
}

func NewFromNode(doc *html.Node) (ext *extractor, err error) {
	ext = &extractor{data: make(map[*Info]*html.Node), doc: doc, option: DEFAULT_OPTION}
	return
}

func NewFromUrl(urlStr string) (ext *extractor, err error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return
	}

	req.Header = make(http.Header)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	// @see : https://github.com/golang/net/blob/master/html/charset/charset.go
	reader, err := charset.NewReader(resp.Body, strings.ToLower(resp.Header.Get("Content-Type")))
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	ext, err = NewFromHtml(string(bs))
	if err != nil {
		return
	}
	ext.urlStr = urlStr
	return
}

func NewFromUrlByHttplib(urlStr string, headers ...map[string]string) (ext *extractor, err error) {
	req := util.BuildRequest("get", urlStr, "", "", "", true, false, headers...)
	str, err := req.String()
	if err != nil {
		return
	}
	ext, err = NewFromHtml(str)
	if err != nil {
		return
	}
	ext.urlStr = urlStr
	return
}

var (
	ERROR_NOTFOUND = errors.New("Content not found")
	DEFAULT_OPTION = &Option{
		RemoveNoise: true,
	}
)

func (ec *extractor) SetOption(option *Option) {
	ec.option = option
}

func (ec *extractor) ToArticle() (article *Article, err error) {
	body := find(ec.doc, isTag(atom.Body))
	if body == nil {
		body = ec.doc
	}

	titleNode := find(ec.doc, isTag(atom.Title))
	if titleNode != nil {
		ec.title = getText(titleNode)
		ec.titleDistanceMin = countChar(ec.title)
		ec.titleMatchLen = ec.titleDistanceMin
	}
	ec.getSn(body)
	ec.getInfo(body)
	node, err := ec.getBestMatch()
	if err != nil {
		return
	}
	if node == nil {
		err = ERROR_NOTFOUND
		return
	}
	ec.tailNode(node)
	article = &Article{}
	article.Publishtime = getPublishTime(node)
	if ec.option.RemoveNoise {
		ec.denoise(node)
	}
	// Get the Content
	article.contentNode = node
	article.Content = getText(node)
	article.Html, err = getHtml(node)
	if err != nil {
		return
	}
	article.Images = getImages(node)
	// find title
	article.Title = ec.title
	if ec.option.AccurateTitle && ec.accurateTitle != "" {
		article.Title = ec.accurateTitle
	}
	article.Images = getImages(node)
	return
}

func (ec *extractor) tailNode(node *html.Node) {
	var densities []float64
	num := 0
	// 第一遍遍历子节点，计算出子节点个数以及保存其密度
	sum := 0.0
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		num++
		d := ec.getInfo(c).Density
		densities = append(densities, d)
		sum += d
	}
	// 文本密度均值
	avg := sum / float64(num)
	articleIndex := make([]int, 0, 100)
	for j := 0; j < len(densities); j++ {
		if densities[j] > avg {
			articleIndex = append(articleIndex, j)
		}
	}
	// 计算节点到文章簇距离的函数，就是计算机最小的索引差
	//
	fn := func(i int) float64 {
		min := math.MaxFloat64
		for _, index := range articleIndex {
			m := math.Abs(float64(index - i))
			if m < min {
				min = m
			}
		}
		return min
	}
	// 得出哪些节点需要remove掉
	rmIndex := make(map[int]bool)
	mins := math.MaxFloat64
	for j := 0; j < len(densities); j++ {
		s := math.Pow(1/(fn(j)+0.1), 2) * ((densities[j] + 1) / avg) // 计算节点的分值，分值对到文章簇的距离以及自身的密度敏感
		if fn(j) < 5 && s < mins {
			mins = s
		}
		if s < mins {
			rmIndex[j] = true
		}
	}
	num = 0
	rmNode := make([]*html.Node, 0, 100)
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if rmIndex[num] == true {
			rmNode = append(rmNode, c)
		}
		num++
	}

	for _, n := range rmNode {
		node.RemoveChild(n)
	}
}

func (ec *extractor) getSn(body *html.Node) {
	txt := getText(body)
	ec.swn = float64(countStopWords(txt))
	ec.sn = float64(countSn(txt))
}

func (ec *extractor) getInfo(node *html.Node) (info *Info) {
	info = NewInfo()

	// remove unused element
	switch node.DataAtom {
	case atom.Script, atom.Object, atom.Style, atom.Iframe, atom.Ins:
		travesRemove(node)
		return
	case atom.Br:
		brRemove(node)
		return
	}
	if display := attr(node, "style"); len(display) > 0 && strings.Contains(display, "display: none") {
		travesRemove(node)
		return
	}
	if node.Type == html.TextNode {
		if node.Parent != nil {
			info.Data = getText(node)
			info.TextCount = countChar(info.Data)
			info.LeafList = append(info.LeafList, info.TextCount)
			info.TextCount = countChar(info.Data)
			info.Density = float64(info.TextCount) / 0.1

			if isTitleNode(node.Parent) {
				ec.filterTitle(node.Parent)
			}
		}
		return
	} else if node.Type == html.ElementNode {
		if isTag(atom.Style)(node) || isTag(atom.Script)(node) {
			return
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			cInfo := ec.getInfo(c)
			info.TextCount += cInfo.TextCount
			info.LinkTextCount += cInfo.LinkTextCount
			info.TagCount += cInfo.TagCount
			info.LinkTagCount += cInfo.LinkTagCount
			info.LeafList = append(info.LeafList, cInfo.LeafList...)
			info.Data += cInfo.Data
			info.Pcount += cInfo.Pcount
			info.ImageCount += cInfo.ImageCount
			info.InputCount += cInfo.InputCount
			info.DensitySum += cInfo.Density

			// cls := attr(node, "class")
			// if cls == "main fl blkContainer" {
			// 	println("adding", cInfo.Data, cInfo.Density)
			// }
		}

		info.TagCount++
		switch node.DataAtom {
		case atom.A:
			info.LinkTagCount++
			info.LinkTextCount += countChar(info.Data)
			if node.Parent != nil && isTag(atom.P)(node.Parent) {
				info.LinkTextCount /= 3
			}
		case atom.P:
			info.Pcount++
		case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
			ec.filterTitle(node)
		case atom.Img, atom.Image:
			info.ImageCount++
		case atom.Input, atom.Textarea, atom.Button:
			info.InputCount++
		}
		var pureLen = info.TextCount - info.LinkTextCount
		var size = info.TagCount - info.LinkTagCount

		if pureLen != 0 && size != 0 {
			info.Density = float64(pureLen) / float64(size)
		}
		if isContentNode(node) {
			ec.addNode(node, info)
		}
		return
	} else if node.Type == html.CommentNode {
		travesRemove(node)
	}
	return
}

// 正文去掉title 编辑距离太近的节点,设置title
func (ec *extractor) filterTitle(n *html.Node) {
	if ec.titleDistanceMin == 0 {
		return
	}
	if isHNode(n) {
		txt := getText(n)
		size := diffString(txt, ec.title)
		size = size - size/3
		// println("count2=>", size, txt, "title=>", ec.title)
		if ec.option.AccurateTitle && size < ec.titleDistanceMin {
			travesRemove(n)
			ec.accurateTitle = txt
			ec.titleDistanceMin = size
			return
		}
	}
	txt := getText(n, func(s *html.Node) bool { return s.Type == html.TextNode })
	maxValue := ec.titleMatchLen / 3
	count := countChar(txt)
	if count >= maxValue && count <= maxValue*3+2 {
		size := diffString(txt, ec.title)
		// println("count=>", count, maxValue, "size=>", size, txt, "title=>", ec.title)
		if n.Parent != nil {
			if isHNode(n.Parent) {
				size /= 2
			}
		}
		if size < maxValue*2 {
			travesRemove(n)
			if ec.option.AccurateTitle && size < ec.titleDistanceMin {
				ec.accurateTitle = txt
				ec.titleDistanceMin = size
			}
		}
	}
}

// 正文去噪
// 去噪即删掉正文中文文本方差小于density * 0.3的非文本节点
// 只清洗前后三个节点
func (ec *extractor) denoise(node *html.Node) {
	avgm := ec.maxDensity * 0.3
	var i = -1
	for n := node.FirstChild; n != nil && i < 3; n = n.NextSibling {
		if n.Type == html.TextNode {
			continue
		}
		i++
		if isNoisingNode(n) {
			info := ec.getInfo(n)
			if info.Density <= avgm {
				travesRemove(n)
				continue
			}
		}
	}

	i = -1
	for n := node.LastChild; n != nil && i < 3; n = n.PrevSibling {
		if n.Type == html.TextNode {
			continue
		}
		i++
		if isNoisingNode(n) {
			info := ec.getInfo(n)
			info.avg = info.getAvg()
			if info.avg < avgm {
				travesRemove(n)
				continue
			}
		}
	}
}

func (ec *extractor) addNode(node *html.Node, info *Info) {
	info.node = node
	info.CalScore(ec.sn, ec.swn)
	ec.data[info] = node
}

func (ec *extractor) getBestMatch() (node *html.Node, err error) {
	if len(ec.data) < 1 {
		err = ERROR_NOTFOUND
		return
	}
	var maxScore float64 = -100
	for kinfo, v := range ec.data {
		if kinfo.score >= maxScore {
			maxScore = kinfo.score
			node = v
		}
		if kinfo.Density > ec.maxDensity {
			ec.maxDensity = kinfo.Density
		}
	}
	if node == nil {
		err = ERROR_NOTFOUND
	}
	return
}

func getPublishTime(node *html.Node) (ts int64) {
	pnode := node.Parent
	sel := func(n *html.Node) bool {
		return n != node
	}
	for i := 0; i < 6 && pnode != nil; i++ {
		h := getText(pnode, sel)
		ts = getTime(h)
		if ts > 0 {
			break
		}
		pnode = pnode.Parent
	}
	if ts == 0 {
		h := getText(node)
		ts = getTime(h)
	}
	return
}
