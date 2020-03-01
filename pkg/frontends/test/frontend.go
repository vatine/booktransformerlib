package test

import (
	"io"

	"github.com/vatine/booktransformerlib/pkg/transformers"
)

// This is useful as a data source during backend testing, it provides us with a source of easy, basic, test data.
type fakeParser struct{}

func FakeParser() transformers.Frontend {
	return fakeParser{}
}

func (f fakeParser) Parse(r io.Reader) <-chan transformers.WorkData {
	_ = r

	rv := make(chan transformers.WorkData)

	go func () {
		rv <- transformers.Author{Author: "Test Testson Testery"}
		rv <- transformers.Title{Title: "Testing for testers"}
		rv <- transformers.Author{Author: "Minple elpniM"}
		rv <- transformers.Chapter{Name: "This is the chapter that starts."}
		rv <- transformers.Formatting{Start: true, Formatting: transformers.Paragraph}
		for _, w := range []string{"This", "is", "a", "word"} {
			rv <- transformers.Word{Word: w}
		}
		rv <- transformers.Punctuation{Punctuation: "."}
		rv <- transformers.Formatting{Start: false, Formatting: transformers.Paragraph}
		rv <- transformers.Formatting{Start: true, Formatting: transformers.Paragraph}
		rv <- transformers.Word{Word: "blerp"}
		rv <- transformers.Footnote{Text: "This is a test footnote."}
		rv <- transformers.Formatting{Start: false, Formatting: transformers.Paragraph}
		close(rv)
		
	} ()

	return rv
}
