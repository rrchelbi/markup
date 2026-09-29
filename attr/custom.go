// Package attr provides typed constructors for HTML attributes.
//
// Most of this package is generated; see internal/gen.
package attr

import (
	"strings"

	"github.com/rrchelbi/markup"
)

// Class returns a class attribute. Names are joined with spaces, and
// repeated Class attributes on one element are merged.
func Class(names ...string) markup.Attr {
	return markup.NewAttr("class", strings.Join(names, " "))
}

// Data returns a data-* attribute.
func Data(name, value string) markup.Attr { return markup.NewAttr("data-"+name, value) }

// Aria returns an aria-* attribute.
func Aria(name, value string) markup.Attr { return markup.NewAttr("aria-"+name, value) }
