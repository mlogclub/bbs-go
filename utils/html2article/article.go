package html2article

import (
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Article struct {
	// Basic
	Html        string `json:"content_html"`
	Content     string `json:"content"`
	Title       string `json:"title"`
	Publishtime int64  `json:"publish_time"`

	// Others
	Images      []string `json:"images"`
	ReadContent string   `json:"read_content"`
	contentNode *html.Node
}

func (a *Article) Readable(urlStr string) {
	a.ParseReadContent()
	a.ParseImage(urlStr)
}

// ParseReadContent parse the ReadContent to be readability
func (a *Article) ParseReadContent() {
	a.cleanAttrs(a.contentNode, "class", "id", "style", "width", "height", "onclick", "onmouseover", "border")
	a.clean(a.contentNode, func(n *html.Node) bool {
		return n.Type == html.CommentNode || n.DataAtom == atom.Script || n.DataAtom == atom.Object
	})
	a.ReadContent, _ = getHtml(a.contentNode)
	// a.ReadContent = CompressHtml(a.ReadContent)
}

// ParseImage parse the image src to the absolute path
func (a *Article) ParseImage(urlStr string) {
	_url, err := url.Parse(urlStr)
	if err != nil {
		return
	}
	mp := make(map[string]string)
	for i, _ := range a.Images {
		if strings.Index(a.Images[i], "http") != 0 {
			var newImg string
			if strings.Index(a.Images[i], "//") == 0 {
				newImg = _url.Scheme + ":" + a.Images[i]
			} else if strings.Index(a.Images[i], "/") == 0 {
				newImg = _url.Scheme + "://" + _url.Host + a.Images[i]
			} else {
				newImg = _url.Scheme + "://" + _url.Host + path.Join(path.Dir(_url.RequestURI()), a.Images[i])
			}
			mp[a.Images[i]] = newImg
			a.Images[i] = newImg
		}
	}
	for k, v := range mp {
		a.Html = strings.Replace(a.Html, k, v, -1)
		a.ReadContent = strings.Replace(a.ReadContent, k, v, -1)
	}
}

func (a *Article) Paragraphs() []string {
	paras := []string{}
	walk(a.contentNode, func(n *html.Node) bool {
		if isTag(atom.P)(n) {
			text := Compress(text(n))
			if text != "" {
				paras = append(paras, text)
			}
			return false
		} else {
			return true
		}
	})
	return paras
}

func (a *Article) clean(sel *html.Node, toClean selector) {
	for c := sel.FirstChild; c != nil; c = c.NextSibling {
		if toClean(c) {
			pre := c.PrevSibling
			sel.RemoveChild(c)
			c = pre
		} else {
			a.clean(c, toClean)
		}
		if c == nil {
			c = sel.FirstChild
			if c == nil {
				break
			}
		}
	}
}

func (a *Article) cleanAttrs(sel *html.Node, attrs ...string) {
	for _, attr := range attrs {
		removeAttr(sel, attr)
	}
	for c := sel.FirstChild; c != nil; c = c.NextSibling {
		a.cleanAttrs(c, attrs...)
	}
}

func (a *Article) GetContentNode() *html.Node {
	return a.contentNode
}
