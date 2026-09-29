// Package html provides typed constructors for HTML elements.
//
// Most of this package is generated; see internal/gen.
package html

import "github.com/rrchelbi/markup"

// Aliases so callers need not import package markup for common types.
type (
	Node  = markup.Node
	Text  = markup.Text
	Raw   = markup.Raw
	Group = markup.Group
)

// Doctype is the HTML5 doctype declaration.
const Doctype = markup.Raw("<!doctype html>")

// Textf formats like fmt.Sprintf and returns the result as Text.
func Textf(format string, a ...any) Text { return markup.Textf(format, a...) }
