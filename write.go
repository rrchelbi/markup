package markup

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

var writers = sync.Pool{
	New: func() any { return bufio.NewWriterSize(nil, 4<<10) },
}

// Write renders n to w through a pooled buffer. If rendering fails,
// buffered output is discarded and the error is returned, though
// earlier chunks may already have reached w. A nil n writes nothing.
func Write(w io.Writer, n Node) error {
	if n == nil {
		return nil
	}
	bw := writers.Get().(*bufio.Writer)
	bw.Reset(w)
	err := n.Render(bw)
	if err == nil {
		err = bw.Flush()
	}
	bw.Reset(nil)
	writers.Put(bw)
	return err
}

// String renders n and returns the result. On error it returns "".
func String(n Node) (string, error) {
	if n == nil {
		return "", nil
	}
	var sb strings.Builder
	if err := n.Render(&sb); err != nil {
		return "", err
	}
	return sb.String(), nil
}
