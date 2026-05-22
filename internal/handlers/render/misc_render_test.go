package render

import (
	"bbs-go/internal/models/resp"
	"reflect"
	"strings"
	"testing"
)

func TestHandleTopicHtmlContentBuildsTocAndHeadingIds(t *testing.T) {
	htmlContent := `
<h1>Page title</h1>
<h2>Intro</h2>
<p>body</p>
<h3>Intro</h3>
<h4>中文 小节</h4>
<h5>Ignored</h5>
<h2>!!!</h2>
<h2>   </h2>
`

	content, toc := handleTopicHtmlContent(htmlContent)

	expected := []resp.TopicTocItem{
		{Id: "topic-heading-intro", Title: "Intro", Level: 2},
		{Id: "topic-heading-intro-2", Title: "Intro", Level: 3},
		{Id: "topic-heading-中文-小节", Title: "中文 小节", Level: 4},
		{Id: "section", Title: "!!!", Level: 2},
	}
	if !reflect.DeepEqual(toc, expected) {
		t.Fatalf("unexpected toc: %#v", toc)
	}

	for _, item := range expected {
		if !strings.Contains(content, `id="`+item.Id+`"`) {
			t.Fatalf("content does not contain id %q: %s", item.Id, content)
		}
	}
	if strings.Contains(content, `id="Page title"`) || strings.Contains(content, `id="Ignored"`) {
		t.Fatalf("content should not assign ids to h1/h5: %s", content)
	}
}

func TestHandleTopicHtmlContentReturnsUpdatedContent(t *testing.T) {
	content, toc := handleTopicHtmlContent(`<h2>Title</h2><p>body</p>`)

	if len(toc) != 1 || toc[0].Id != "topic-heading-title" {
		t.Fatalf("unexpected toc: %#v", toc)
	}
	if !strings.Contains(content, `<h2 id="topic-heading-title">Title</h2>`) {
		t.Fatalf("heading id was not written to content: %s", content)
	}
}
