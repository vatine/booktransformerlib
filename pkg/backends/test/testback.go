package test

import (
	"io"
	"testing"

	"github.com/vatine/booktransformerlib/pkg/transformers"
)

type TestBackend struct {
	Expected  []transformers.WorkData
	Authors   []string
	WordCount int
	title     string
	t         *testing.T
}

func SetExpectations(tb *TestBackend, t *testing.T, title string, expected []transformers.WorkData) {
	tb.t = t
	tb.Expected = expected
	tb.title = title
}

func SetAuthors(tb *TestBackend, authors []string) {
	tb.Authors = authors
}

func NewTestBackend(w io.WriteCloser) *TestBackend {
	rv := TestBackend{}
	rv.WordCount = 0

	return &rv
}

func (b *TestBackend) AddAuthor(author string) {
	posIx := -1
	for ix, a := range b.Authors {
		if author == a {
			posIx = ix
			break
		}
	}

	if posIx == -1 {
		b.t.Errorf("Expected author %s to be present.", author)
	} else {
		b.Authors = append(b.Authors[:posIx], b.Authors[posIx+1:]...)
	}
}

func (b *TestBackend) Close() {
	for _, a := range b.Authors {
		b.t.Errorf("Author %s was never encountered.", a)
	}
}

func (b *TestBackend) EmitFootnote(note string) {
	expected, ok := b.Expected[0].(transformers.Footnote)
	b.Expected = b.Expected[1:]
	if !ok {
		b.t.Errorf("Saw a footnote, expected %s", expected)
	} else {
		if expected.Text != note {
			b.t.Errorf("Expected footnote '%s', saw footnote '%s'", expected.Text, note)
		}
	}
}

func (b *TestBackend) EmitPunctuation(p string) {
	expected, ok := b.Expected[0].(transformers.Punctuation)

	if !ok {
		b.t.Errorf("Expected %s, got Punctuation.", expected)
	} else {
		if expected.Punctuation != p {
			b.t.Errorf("Expected punctuation '%s', saw '%s'", expected.Punctuation, p)
		}
	}
}

func (b *TestBackend) EmitWord(w string) {
	b.WordCount += 1
}

func (b *TestBackend) EndBlockQuote() {
	expected, ok := b.Expected[0].(transformers.Formatting)

	if !ok {
		b.t.Errorf("Expected %s, got EndBlockQuote.", expected)
	} else {
		if !expected.Start && (expected.Formatting == transformers.BlockQuote) {
			return
		} else {
			b.t.Errorf("Expected %s, got EndBlockQuote.", expected)
		}
	}
}

func (b *TestBackend) EndBold() {
	expected, ok := b.Expected[0].(transformers.Formatting)

	if !ok {
		b.t.Errorf("Expected %s, got EndBold.", expected)
	} else {
		if !expected.Start && (expected.Formatting == transformers.Bold) {
			return
		} else {
			b.t.Errorf("Expected %s, got EndBold.", expected)
		}
	}
}

func (b *TestBackend) EndItalic() {
	expected, ok := b.Expected[0].(transformers.Formatting)

	if !ok {
		b.t.Errorf("Expected %s, got EndItalic.", expected)
	} else {
		if !expected.Start && (expected.Formatting == transformers.Italic) {
			return
		} else {
			b.t.Errorf("Expected %s, got EndItalic.", expected)
		}
	}
}

func (b *TestBackend) EndParagraph() {
	expected, ok := b.Expected[0].(transformers.Formatting)

	if !ok {
		b.t.Errorf("Expected %s, got EndParagraph.", expected)
	} else {
		if !expected.Start && (expected.Formatting == transformers.Paragraph) {
			return
		} else {
			b.t.Errorf("Expected %s, got EndParagraph.", expected)
		}
	}
}

func (b *TestBackend) NewLine() {
	expected, ok := b.Expected[0].(transformers.Formatting)

	if !ok {
		b.t.Errorf("Expected %s, got End.", expected)
	} else {
		if (expected.Formatting == transformers.NewLine) {
			return
		} else {
			b.t.Errorf("Expected %s, got NewLine.", expected)
		}
	}
}

func (b *TestBackend) NewChapter(title string) {
	// This should probably have some logic...
}

func (b *TestBackend) SetTitle(title string) {
	if b.title != title {
		b.t.Errorf("Expected title '%s', saw '%s'", b.title, title)
	}
}

func (b *TestBackend) StartBlockQuote() {
	expected, ok := b.Expected[0].(transformers.Formatting)

	if !ok {
		b.t.Errorf("Expected %s, got StartBlockQuote.", expected)
	} else {
		if expected.Start && (expected.Formatting == transformers.BlockQuote) {
			return
		} else {
			b.t.Errorf("Expected %s, got StartBlockQuote.", expected)
		}
	}
}

func (b *TestBackend) StartBold() {
	expected, ok := b.Expected[0].(transformers.Formatting)

	if !ok {
		b.t.Errorf("Expected %s, got StartBold.", expected)
	} else {
		if expected.Start && (expected.Formatting == transformers.Bold) {
			return
		} else {
			b.t.Errorf("Expected %s, got StartBold.", expected)
		}
	}
}

func (b *TestBackend) StartItalic() {
	expected, ok := b.Expected[0].(transformers.Formatting)

	if !ok {
		b.t.Errorf("Expected %s, got StartItalic.", expected)
	} else {
		if expected.Start && (expected.Formatting == transformers.Italic) {
			return
		} else {
			b.t.Errorf("Expected %s, got StartItalic.", expected)
		}
	}
}

func (b *TestBackend) StartParagraph() {
	expected, ok := b.Expected[0].(transformers.Formatting)

	if !ok {
		b.t.Errorf("Expected %s, got StartParagraph.", expected)
	} else {
		if expected.Start && (expected.Formatting == transformers.Paragraph) {
			return
		} else {
			b.t.Errorf("Expected %s, got StartParagraph.", expected)
		}
	}
}
