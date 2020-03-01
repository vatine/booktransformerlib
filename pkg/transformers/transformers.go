// A library for transforming text from one format to another.

// The general principle of operation is that a front-end generates a
// stream of WorkData items, these are then fed to a backend, where
// they're emitted in a fashion suitable for whatever format the
// backend encodes. This library does not create a structured
// intermediate representation of the work, instead operating on
// things like "I see an author", "here is a word", "this is a new
// paragraph", "oh, there's nothing left"...
package transformers

import (
	"io"
)

// Type to allow us to specify tetx formatting
type Format int

const (
	Italic Format = iota
	Bold
	Underline
	Strikethrough
	NewLine
	BlockQuote
	Paragraph
)

// This is useful as a data source during backend testing, it provides us with a source of easy, basic, test data.
type fakeParser struct{}

func FakeParser() Frontend {
	return fakeParser{}
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
		rv <- Formatting{Start: true, Formatting: Paragraph}
		rv <- Word{Word: "blerp"}
		rv <- Footnote{Text: "This is a test footnote."}
		rv <- Formatting{Start: false, Formatting: Paragraph}
		close(rv)
		
	} ()

	return rv
}


// A Backend is a "destination type" (text file, HTML, ePUB, ...) and
// is actuated by a stream of Blobs from a Frontend. When the whole
// work has been emitted, the Close method is called, which does
// assorted cleanup and closes the output stream.
type Backend interface {
 	AddAuthor(string)
	Close()
	EmitFootnote(string)
	EmitPunctuation(string)
	EmitWord(string)
	EndBlockQuote()
	EndBold()
	EndItalic()
	EndParagraph()
	NewChapter(string)
	NewLine()
	SetTitle(string)
	StartBlockQuote()
	StartBold()
	StartItalic()
	StartParagraph()
}

// A Frontend is a converter of source work data. The work is passed
// in via an io.Reader to Parse, which should return a channel through
// which WorkData items are passed. These are then emitted through a
// Backend.
type Frontend interface {
	Parse(io.Reader) <-chan WorkData
}

// Generic wrapper for anything that is expected to come from a work.
// The Emit function is responsible for passing the relevant data from
// the frontend to the backend.
type WorkData interface {
	Emit(Backend)
}

// An author. We'd normally expect at least one per work.
type Author struct {
	Author string
}

// This signifies the start of a new chapter. The name is not compulsory,
// but it is customary.
type Chapter struct {
	Name string
}

// A footnote.
type Footnote struct {
	Text string
}

// This is somewhat of a wrapper type, indicating the start (or end)
// of a specific type of formatting.
type Formatting struct {
	Start      bool
	Formatting Format
}

// Punctuation marks.
type Punctuation struct {
	Punctuation string
}

// The title of a work.
type Title struct {
	Title string
}

type Word struct {
	Word string
}

func (a Author) Emit(b Backend) {
	b.AddAuthor(a.Author)
}

func (c Chapter) Emit(b Backend) {
	b.NewChapter(c.Name)
}

func (f Footnote) Emit(b Backend) {
	b.EmitFootnote(f.Text)
}

func (f Formatting) Emit(b Backend) {
	switch f.Formatting {
	case NewLine:
		b.NewLine()
	case BlockQuote:
		if f.Start {
			b.StartBlockQuote()
		} else {
			b.EndBlockQuote()
		}
	case Italic:
		if f.Start {
			b.StartItalic()
		} else {
			b.EndItalic()
		}
	case Bold:
		if f.Start {
			b.StartBold()
		} else {
			b.EndBold()
		}
	case Paragraph:
		if f.Start {
			b.StartParagraph()
		} else {
			b.EndParagraph()
		}
	}
}

func (p Punctuation) Emit(b Backend) {
	b.EmitPunctuation(p.Punctuation)
}

func (t Title) Emit(b Backend) {
	b.SetTitle(t.Title)
}

func (w Word) Emit(b Backend) {
	b.EmitWord(w.Word)
}

// Convert is the driver that couples a frontend to a backend,
func Convert(f Frontend, b Backend, r io.Reader) {
	c := f.Parse(r)

	for wd := range c {
		wd.Emit(b)
	}

	b.Close()
}
