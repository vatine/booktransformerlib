// Copyright 2015- Ingvar Mattsson <imagineaclevernamehere@gmail.com>
// A Go library for transforming text from one format to another
package booktransformlib

import (
	"io"
)

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

// A Backend is a "destination type" (text file, HTML, ePUB, ...) and
// is actuated by a stream of Blobs from a Frontend. When the whole
// work has been emitted, the Close method is called, ceasing the
// output.
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
type WorkData interface {
	Emit(Backend)
}

type Author struct {
	Author string
}

type Chapter struct {
	Name string
}

type Formatting struct {
	Start      bool
	Formatting Format
}

type Punctuation struct {
	Punctuation string
}

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
