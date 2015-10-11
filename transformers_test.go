package booktransformlib

import (
	"bytes"
	"fmt"
	"testing"
	"io"
	"strings"
)

type fakeParser struct {}

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


func (f fakeParser) Parse(r io.Reader) <-chan WorkData {
	_ = r

	rv := make(chan WorkData)

	go func () {
		rv <- Author{Author: "Test Testson Testery"}
		rv <- Title{Title: "Testing for testers"}
		rv <- Author{Author: "Minple elpniM"}
		rv <- Chapter{Name: "This is the chapter that starts."}
		rv <- Formatting{Start: true, Formatting: Paragraph}
		for _, w := range []string{"This", "is", "a", "word"} {
			rv <- Word{Word: w}
		}
		rv <- Punctuation{Punctuation: "."}
		rv <- Formatting{Start: false, Formatting: Paragraph}
		close(rv)
		
	} ()

	return rv
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
	

	Convert(fakeParser{}, htBackend, buf)

	expected := `<html>
<head><title>Testing for testers</title></head>
<body><h1>Testing for testers</h1>
<center>Test Testson Testery</center>
<center>Minple elpniM</center>
<h2>1 - This is the chapter that starts. </h2>

<p> This is a word.</p>


</body>
</html>
`
	seen := buf.b.String()

	if seen != expected {
		t.Errorf("%s", renderStringDiff(expected, seen))
	}
}
