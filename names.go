package markup

import "strings"

// validTagName reports whether s is [A-Za-z][A-Za-z0-9-]*.
func validTagName(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case 'a' <= c && c <= 'z', 'A' <= c && c <= 'Z':
		case i > 0 && ('0' <= c && c <= '9' || c == '-'):
		default:
			return false
		}
	}
	return true
}

// validAttrName reports whether s is [A-Za-z_:][A-Za-z0-9_:.-]*.
func validAttrName(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case 'a' <= c && c <= 'z', 'A' <= c && c <= 'Z', c == '_', c == ':':
		case i > 0 && ('0' <= c && c <= '9' || c == '-' || c == '.'):
		default:
			return false
		}
	}
	return true
}

// blockedAttr reports whether name can carry executable code: inline
// event handlers (onclick, hx-on:click, x-on:click) and srcdoc.
func blockedAttr(name string) bool {
	if len(name) >= 2 && name[0]|0x20 == 'o' && name[1]|0x20 == 'n' {
		return true
	}
	lower := strings.ToLower(name)
	return lower == "srcdoc" ||
		strings.Contains(lower, "-on:") ||
		strings.Contains(lower, "-on-")
}

// isURLAttr reports whether name holds a URL.
func isURLAttr(name string) bool {
	switch strings.ToLower(name) {
	case "href", "src", "action", "formaction", "cite", "poster", "data",
		"manifest", "longdesc", "codebase", "icon", "background", "xlink:href":
		return true
	}
	return false
}

// isVoid reports whether tag is an HTML void element.
func isVoid(tag string) bool {
	switch tag {
	case "area", "base", "br", "col", "embed", "hr", "img", "input",
		"link", "meta", "source", "track", "wbr":
		return true
	}
	return false
}

// isRawText reports whether tag holds raw text that entity escaping
// would corrupt.
func isRawText(tag string) bool { return tag == "script" || tag == "style" }
