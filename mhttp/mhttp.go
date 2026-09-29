// Package mhttp serves markup trees over net/http.
package mhttp

import (
	"bytes"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	"github.com/rrchelbi/markup"
)

const maxPooled = 64 << 10

var bufs = sync.Pool{New: func() any { return new(bytes.Buffer) }}

// HandlerFunc renders the node returned by f as text/html.
//
// The page is rendered to a buffer first. If rendering fails the client
// gets a plain 500 and the error is logged, never a half-written page.
type HandlerFunc func(r *http.Request) markup.Node

// ServeHTTP implements http.Handler.
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := bufs.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		if buf.Cap() <= maxPooled {
			bufs.Put(buf)
		}
	}()

	if n := f(r); n != nil {
		if err := n.Render(buf); err != nil {
			slog.ErrorContext(r.Context(), "markup: render failed", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Content-Length", strconv.Itoa(buf.Len()))
	w.Write(buf.Bytes())
}
