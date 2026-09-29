package markup

import (
	"fmt"
	"io"
)

type attrKind uint8

const (
	kindText   attrKind = iota // escaped value
	kindBool                   // name only
	kindURL                    // escaped value, scheme checked
	kindUnsafe                 // escaped value, name not screened
)

// Attr is an HTML attribute. Its zero value is invalid and fails to
// render. Attr is also a Node, so it can be passed to Tag alongside
// children.
type Attr struct {
	name, value string
	kind        attrKind
}

// NewAttr returns an attribute whose value is escaped on output.
// Attributes that hold URLs (href, src, action, ...) are scheme-checked
// automatically. Rendering fails with ErrUnsafeAttr for event handlers
// and srcdoc; use UnsafeAttr if you need them.
func NewAttr(name, value string) Attr {
	kind := kindText
	if isURLAttr(name) {
		kind = kindURL
	}
	return Attr{name: name, value: value, kind: kind}
}

// URL returns an attribute holding a URL, for names NewAttr does not
// already treat as URL-valued. Schemes other than http, https, mailto
// and tel render as "#unsafe".
func URL(name, value string) Attr {
	return Attr{name: name, value: value, kind: kindURL}
}

// Bool returns a boolean attribute such as disabled.
func Bool(name string) Attr { return Attr{name: name, kind: kindBool} }

// UnsafeAttr returns an attribute that skips the event-handler and URL
// checks. The name is still validated and the value is still escaped,
// but the browser may execute the value as code.
//
// UNSAFE: never pass untrusted input.
func UnsafeAttr(name, value string) Attr {
	return Attr{name: name, value: value, kind: kindUnsafe}
}

// Name returns the attribute name.
func (a Attr) Name() string { return a.name }

// Value returns the attribute value before escaping.
func (a Attr) Value() string { return a.value }

// Render writes ` name="value"` with a leading space.
func (a Attr) Render(w io.Writer) error {
	if !validAttrName(a.name) {
		return fmt.Errorf("%w: attribute %q", ErrInvalidName, a.name)
	}
	if a.kind != kindUnsafe && blockedAttr(a.name) {
		return fmt.Errorf("%w: %q", ErrUnsafeAttr, a.name)
	}
	if _, err := io.WriteString(w, " "); err != nil {
		return err
	}
	if _, err := io.WriteString(w, a.name); err != nil {
		return err
	}
	if a.kind == kindBool {
		return nil
	}
	v := a.value
	if a.kind == kindURL {
		v = safeURL(v)
	}
	if _, err := io.WriteString(w, `="`); err != nil {
		return err
	}
	if err := writeEscaped(w, v); err != nil {
		return err
	}
	_, err := io.WriteString(w, `"`)
	return err
}
