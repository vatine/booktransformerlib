package booktransformlib

import (
	"bytes"
	"testing"
	"io"
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
		close(rv)
		
	} ()

	return rv
}

func TestHtmlBackend(t *testing.T) {
	buf := bufferCloser{b: bytes.NewBufferString("")}
	htBackend := NewHtmlBackend(buf)
	

	Convert(fakeParser{}, htBackend, buf)
}
