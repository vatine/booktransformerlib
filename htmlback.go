package booktransformerlib

import (
	"fmt"
	"io"
)

type HtmlBackend struct {
	w io.writer
	chapterNo  int
	footnoteNo int
	authors    []string
	inPara     bool
	footnotes  map[string]string
}


func (b *HtmlBackend) AddAuthor(author string) {
	b.authors = append(b.authors, author)
}

func (b *HtmlBackend) Close() {
	b.w.Close()
}

func (b HtmlBackend) noteAnchor(chapter, count int) string {
	return fmt.Sprintf("c%03d-n%03d", chapter, count)
}

func (b *HtmlBackend) EmitFootnote(note string) {
	b.footnoteNo += 1
	anchor := b.noteAnchor(chapter, b.footnoteNo)
	b.footnotes[anchor] = note
	fmt.Fprintf(b.w, "<a name=\"notereturn-%s\"><a href=\"note-%s\"><sup>%d</sup></a></a>", anchor, anchor, b.footnoteNo)
}
