package backends

import (
	"bytes"
	"fmt"
	"testing"
	"strings"

	"github.com/vatine/booktransformerlib/pkg/frontends/test"
	"github.com/vatine/booktransformerlib/pkg/transformers"
)


type bufferCloser struct {
	b *bytes.Buffer
}

func (b bufferCloser) Write(p []byte) (int, error) {
	return b.b.Write(p)
}

func (b bufferCloser) Close() error {
	_ = b
	return nil
}

func (b bufferCloser) Read(p []byte) (int, error) {
	return b.b.Read(p)
}


func renderStringDiff(expected, seen string) string {
	eLines := strings.Split(expected, "\n")
	sLines := strings.Split(seen, "\n")

	w := bytes.NewBufferString("")
	fmt.Fprintf(w, "len(expected) == %d\nlen(seen) == %d\nexpected lines %d, seen lines %d\n", len(expected), len(seen), len(eLines), len(sLines))

	soff := 0

	for eix, el := range eLines {
		if el == sLines[eix + soff] {
			fmt.Fprintf(w, "  ^%s$\n", el)
		} else {
			fmt.Fprintf(w, "- ^%s$\n+ ^%s$\n", el, sLines[eix + soff])
		}
	}

	return w.String()
}

func TestHtmlBackend(t *testing.T) {
	buf := bufferCloser{b: bytes.NewBufferString("")}
	htBackend := NewHtmlBackend(buf)
	

	transformers.Convert(test.FakeParser(), htBackend, buf)

	expected := `<html>
<head><title>Testing for testers</title></head>
<body><h1>Testing for testers</h1>
<center>Test Testson Testery</center>
<center>Minple elpniM</center>
<h2>1 - This is the chapter that starts. </h2>

<p> This is a word.</p>

<p> blerp<a name="notereturn-c001-n001"><a href="#note-c001-n001"><sup>1</sup></a></a></p>

<p><a name="note-c001-n001"><sup>1</sup></a>This is a test footnote.<small><a href="notereturn-c001-n001">back</a></small></p>


</body>
</html>
`
	seen := buf.b.String()

	if seen != expected {
		t.Errorf("%s", renderStringDiff(expected, seen))
	}
}
