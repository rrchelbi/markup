// Package flow provides conditionals and loops for building markup trees.
package flow

import "github.com/rrchelbi/markup"

// If returns n when cond is true and an empty Group otherwise. n is
// evaluated by the caller either way; use When to defer construction.
func If(cond bool, n markup.Node) markup.Node {
	if cond {
		return n
	}
	return markup.Group(nil)
}

// IfElse returns a when cond is true and b otherwise.
func IfElse(cond bool, a, b markup.Node) markup.Node {
	if cond {
		return a
	}
	return b
}

// When calls f and returns its result only when cond is true.
func When(cond bool, f func() markup.Node) markup.Node {
	if cond {
		return f()
	}
	return markup.Group(nil)
}

// Each returns the nodes produced by f for every element of s.
func Each[S ~[]E, E any](s S, f func(E) markup.Node) markup.Node {
	g := make(markup.Group, 0, len(s))
	for _, e := range s {
		g = append(g, f(e))
	}
	return g
}
