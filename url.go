package markup

const unsafeURL = "#unsafe"

// safeURL returns s if its scheme is http, https, mailto or tel, or if
// it has no scheme (relative URL). Otherwise it returns "#unsafe".
//
// It mirrors browser parsing: leading control characters and spaces are
// skipped and tabs and newlines inside the scheme are ignored, so
// "java\tscript:alert(1)" is caught. Any colon before the first '/',
// '?' or '#' counts as a scheme separator, which errs on the side of
// blocking.
func safeURL(s string) string {
	i := 0
	for i < len(s) && s[i] <= ' ' {
		i++
	}
	var buf [16]byte
	n := 0
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\t', '\n', '\r':
			continue
		case '/', '?', '#':
			return s
		case ':':
			if !overflow {
				switch string(buf[:n]) {
				case "http", "https", "mailto", "tel":
					return s
				}
			}
			return unsafeURL
		}
		if n == len(buf) {
			overflow = true
			continue
		}
		if 'A' <= c && c <= 'Z' {
			c |= 0x20
		}
		buf[n] = c
		n++
	}
	return s
}
