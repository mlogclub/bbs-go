package godown

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestGodown(t *testing.T) {
	m, err := filepath.Glob("testdata/*.html")
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(m)
	for _, file := range m {
		f, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		if err = Convert(&buf, f, nil); err != nil {
			t.Fatal(err)
		}

		b, err := ioutil.ReadFile(file[:len(file)-4] + "md")
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != buf.String() {
			t.Errorf("(%s):\nwant:\n%s}}}\ngot:\n%s}}}\n", file, string(b), buf.String())
		}
		f.Close()
	}
}

type errReader int

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	var e errReader
	err := Convert(&buf, e, nil)
	if err == nil {
		t.Fatal("should be an error")
	}
}

func TestGuessLang(t *testing.T) {
	var buf bytes.Buffer
	err := Convert(&buf, strings.NewReader(`
<pre>
def do_something():
  pass
</pre>
	`), &Option{
		GuessLang: func(s string) (string, error) { return "python", nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "```python\ndef do_something():\n  pass\n```\n\n\n"
	if buf.String() != want {
		t.Errorf("\nwant:\n%s}}}\ngot:\n%s}}}\n", want, buf.String())
	}
}

func TestGuessLangBq(t *testing.T) {
	var buf bytes.Buffer
	err := Convert(&buf, strings.NewReader(`
<blockquote class="code">
<b>def</b> do_something():
  <i>pass</i>
</blockquote>
	`), &Option{
		GuessLang: func(s string) (string, error) { return "python", nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "```python\ndef do_something():\n  pass\n```\n\n\n"
	if buf.String() != want {
		t.Errorf("\nwant:\n%s}}}\ngot:\n%s}}}\n", want, buf.String())
	}
}

func TestScript(t *testing.T) {
	var buf bytes.Buffer
	err := Convert(&buf, strings.NewReader(`
<p>here is script</p>

<script type="text/javascript" src="https://code.jqeury.com/jquery-latest.js"></script>

<script type="text/javascript"><!--
alert(1)
--></script>
	`), &Option{
		Script: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `here is script

<script type="text/javascript" src="https://code.jqeury.com/jquery-latest.js"></script>

<script type="text/javascript"><!--
alert(1)
--></script>


`
	if buf.String() != want {
		t.Errorf("\nwant:\n%s}}}\ngot:\n%s}}}\n", want, buf.String())
	}
}

func TestStyle(t *testing.T) {
	var buf bytes.Buffer
	err := Convert(&buf, strings.NewReader(`
<p>here is style</p>

<style><!--
body {
	background-color: red;
}
--></style>
	`), &Option{
		Style: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `here is style

<style><!--
body {
	background-color: red;
}
--></style>


`
	if buf.String() != want {
		t.Errorf("\nwant:\n%s}}}\ngot:\n%s}}}\n", want, buf.String())
	}
}
