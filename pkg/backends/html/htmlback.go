package backends

import (
	"fmt"
	"io"
)

type HtmlBackend struct {
	w          io.WriteCloser
	chapterNo  int
	footnoteNo int
	authors    []string
	title      string
	inPara     bool
	footnotes  map[string]string
	WordCount  int
}

func NewHtmlBackend(w io.WriteCloser) *HtmlBackend {
	rv := HtmlBackend{w: w}
	rv.chapterNo = 0
	rv.footnotes = make(map[string]string)
	rv.inPara = false
	rv.authors = nil
	rv.WordCount = 0

	return &rv
}

func (b *HtmlBackend) AddAuthor(author string) {
	b.authors = append(b.authors, author)
}

func (b *HtmlBackend) Close() {
	if b.inPara {
		b.EndParagraph()
	}
	b.writeFootnotes()
	fmt.Fprintf(b.w, "\n</body>\n</html>\n")
	b.w.Close()
}

func (b HtmlBackend) noteAnchor(chapter, count int) string {
	return fmt.Sprintf("c%03d-n%03d", chapter, count)
}

func (b *HtmlBackend) EmitFootnote(note string) {
	b.footnoteNo += 1
	anchor := b.noteAnchor(b.chapterNo, b.footnoteNo)
	b.footnotes[anchor] = note
	fmt.Fprintf(b.w, "<a name=\"notereturn-%s\"><a href=\"#note-%s\"><sup>%d</sup></a></a>", anchor, anchor, b.footnoteNo)
}

func (b *HtmlBackend) EmitPunctuation(p string) {
	fmt.Fprintf(b.w, "%s", p)
}

func (b *HtmlBackend) EmitWord(w string) {
	b.WordCount += 1
	fmt.Fprintf(b.w, " %s", w)
}

func (b *HtmlBackend) EndBlockQuote() {
	fmt.Fprintf(b.w, "</blockquote>")
}

func (b *HtmlBackend) EndBold() {
	fmt.Fprintf(b.w, "</b>")
}

func (b *HtmlBackend) EndItalic() {
	fmt.Fprintf(b.w, "</i>")
}

func (b *HtmlBackend) EndParagraph() {
	if b.inPara {
		b.inPara = false
		fmt.Fprintf(b.w, "</p>\n\n")
	}
}

func (b *HtmlBackend) NewLine() {
	fmt.Fprintf(b.w, "<br />\n")
}

func (b *HtmlBackend) writeFootnotes() {
	for c := 1; c <= b.footnoteNo; c++ {
		suffix := b.noteAnchor(b.chapterNo, c)
		fmt.Fprintf(b.w, "<p><a name=\"note-%s\"><sup>%d</sup></a>%s<small><a href=\"notereturn-%s\">back</a></small></p>\n\n", suffix, c, b.footnotes[suffix], suffix)
	}
	b.footnoteNo = 0
}

func (b *HtmlBackend) NewChapter(title string) {
	if (b.chapterNo == 0) {
		fmt.Fprintf(b.w, "<html>\n<head><title>%s</title></head>\n<body><h1>%s</h1>\n", b.title, b.title)
		for _, a := range b.authors {
			fmt.Fprintf(b.w, "<center>%s</center>\n", a)
		}
	}
	b.writeFootnotes()
	b.chapterNo += 1
	if b.inPara {
		b.inPara = false
		fmt.Fprintf(b.w, "</p>\n\n")
	}
	fmt.Fprintf(b.w, "<h2>%d - %s </h2>\n\n", b.chapterNo, title)
}

func (b *HtmlBackend) SetTitle(title string) {
	b.title = title
}

func (b *HtmlBackend) StartBlockQuote() {
	fmt.Fprintf(b.w, "<blockquote>")
}

func (b *HtmlBackend) StartBold() {
	fmt.Fprintf(b.w, "<b>")
}

func (b *HtmlBackend) StartItalic() {
	fmt.Fprintf(b.w, "<i>")
}

func (b *HtmlBackend) StartParagraph() {
	if b.inPara {
		b.EndParagraph()
	}
	b.inPara = true
	fmt.Fprintf(b.w, "<p>")
}
