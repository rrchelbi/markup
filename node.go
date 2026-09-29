package markup

import (
	"fmt"
	"io"
)

// A Node can render itself as HTML. Implementations must escape any
// untrusted data they write and must not retain w.
type Node interface {
	Render(w io.Writer) error
}

// NodeFunc adapts a function to a Node.
type NodeFunc func(w io.Writer) error

// Render calls f(w).
func (f NodeFunc) Render(w io.Writer) error { return f(w) }

// Text is a text node. It is HTML-escaped when rendered.
type Text string

// Render writes the escaped text.
func (t Text) Render(w io.Writer) error { return writeEscaped(w, string(t)) }

// Textf formats like fmt.Sprintf and returns the result as Text.
func Textf(format string, a ...any) Text { return Text(fmt.Sprintf(format, a...)) }

// Raw is markup written verbatim, with no escaping.
//
// UNSAFE: never construct a Raw from untrusted input.
type Raw string

// Render writes the string unchanged.
func (r Raw) Render(w io.Writer) error {
	_, err := io.WriteString(w, string(r))
	return err
}

// Group is a sequence of nodes with no wrapper element. When passed to
// Tag, a Group is flattened, so attributes inside it apply to the
// parent element.
type Group []Node

// Render renders each non-nil node in order.
func (g Group) Render(w io.Writer) error {
	for _, n := range g {
		if n == nil {
			continue
		}
		if err := n.Render(w); err != nil {
			return err
		}
	}
	return nil
}
