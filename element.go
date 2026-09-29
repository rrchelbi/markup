package markup

import (
	"fmt"
	"io"
	"strings"
)

// Element is an HTML element. Elements are immutable once built and
// safe for concurrent rendering.
type Element struct {
	name     string
	attrs    []Attr
	children []Node
}

// Tag returns an element with the given tag name. Items may be Attrs,
// Groups (flattened) or any other Node; nil items are ignored. A later
// attribute replaces an earlier one of the same name, except class,
// which is joined with a space.
//
// The name is validated when the element is rendered.
func Tag(name string, items ...Node) *Element {
	e := &Element{name: name}
	if len(items) > 0 {
		e.children = make([]Node, 0, len(items))
		e.add(items)
	}
	return e
}

// Name returns the tag name.
func (e *Element) Name() string { return e.name }

func (e *Element) add(items []Node) {
	for _, it := range items {
		switch v := it.(type) {
		case nil:
		case Attr:
			e.addAttr(v)
		case Group:
			e.add(v)
		case *Element:
			if v != nil {
				e.children = append(e.children, v)
			}
		default:
			e.children = append(e.children, it)
		}
	}
}

func (e *Element) addAttr(a Attr) {
	for i := range e.attrs {
		old := &e.attrs[i]
		if !strings.EqualFold(old.name, a.name) {
			continue
		}
		if a.name == "class" && old.kind == kindText && a.kind == kindText && old.value != "" && a.value != "" {
			old.value += " " + a.value
		} else {
			*old = a
		}
		return
	}
	e.attrs = append(e.attrs, a)
}

// Render writes the element and its subtree.
func (e *Element) Render(w io.Writer) error {
	if !validTagName(e.name) {
		return fmt.Errorf("%w: tag %q", ErrInvalidName, e.name)
	}
	void := isVoid(e.name)
	if void && len(e.children) > 0 {
		return fmt.Errorf("%w: <%s>", ErrVoidChildren, e.name)
	}
	if _, err := io.WriteString(w, "<"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, e.name); err != nil {
		return err
	}
	for i := range e.attrs {
		if err := e.attrs[i].Render(w); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}
	if void {
		return nil
	}
	if isRawText(e.name) {
		if err := renderRawText(w, e.name, e.children); err != nil {
			return err
		}
	} else {
		for _, c := range e.children {
			if err := c.Render(w); err != nil {
				return err
			}
		}
	}
	if _, err := io.WriteString(w, "</"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, e.name); err != nil {
		return err
	}
	_, err := io.WriteString(w, ">")
	return err
}

// renderRawText writes the children of a script or style element. Only
// Text and Raw children are allowed. Text is written unescaped, but only
// if it cannot terminate the element or open an HTML comment.
func renderRawText(w io.Writer, tag string, children []Node) error {
	for _, c := range children {
		var s string
		switch v := c.(type) {
		case Text:
			s = string(v)
			if closesRawText(s, tag) {
				return fmt.Errorf("%w: <%s>", ErrUnsafeRawText, tag)
			}
		case Raw:
			s = string(v)
		default:
			return fmt.Errorf("%w: <%s> accepts only Text and Raw", ErrUnsafeRawText, tag)
		}
		if _, err := io.WriteString(w, s); err != nil {
			return err
		}
	}
	return nil
}

// closesRawText reports whether s contains "</tag" (any case) or "<!--".
func closesRawText(s, tag string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] != '<' {
			continue
		}
		rest := s[i+1:]
		if strings.HasPrefix(rest, "!--") {
			return true
		}
		if len(rest) > len(tag) && rest[0] == '/' && strings.EqualFold(rest[1:1+len(tag)], tag) {
			return true
		}
	}
	return false
}
