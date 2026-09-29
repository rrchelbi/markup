// Package markup builds HTML as Go values and renders it safely.
//
// A document is a tree of Nodes. Elements are created with Tag (or the
// typed helpers in package html), and accept attributes and children in
// one variadic list:
//
//	html.Div(attr.Class("card"),
//		html.H1(markup.Text("Hello")),
//	)
//
// # Safety model
//
// Everything is safe by default:
//
//   - Text and attribute values are HTML-escaped.
//   - Tag and attribute names are validated at render time.
//   - Event-handler attributes (on*) and srcdoc are rejected.
//   - URL attributes (href, src, action, ...) with a scheme other than
//     http, https, mailto or tel are replaced with "#unsafe".
//   - Text inside <script> and <style> is written verbatim but rejected
//     if it could close the element early.
//   - Void elements (br, img, ...) cannot have children.
//
// The escape hatches are explicit and greppable: Raw and UnsafeAttr.
// Never pass untrusted input to them.
//
// Render errors are returned, never panicked. When rendering to a
// network connection, render to a buffer first (see package mhttp) so a
// failure does not leave a half-written page.
package markup

//go:generate go run ./internal/gen
