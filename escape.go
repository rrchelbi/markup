package markup

import "io"

// writeEscaped writes s with HTML metacharacters escaped. It suits both
// text nodes and double-quoted attribute values. Runs of safe bytes are
// written in one call, so clean strings cost a single write.
func writeEscaped(w io.Writer, s string) error {
	last := 0
	for i := 0; i < len(s); i++ {
		var rep string
		switch s[i] {
		case '&':
			rep = "&amp;"
		case '<':
			rep = "&lt;"
		case '>':
			rep = "&gt;"
		case '"':
			rep = "&#34;"
		case '\'':
			rep = "&#39;"
		case 0:
			rep = "\uFFFD"
		default:
			continue
		}
		if last < i {
			if _, err := io.WriteString(w, s[last:i]); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, rep); err != nil {
			return err
		}
		last = i + 1
	}
	if last < len(s) {
		_, err := io.WriteString(w, s[last:])
		return err
	}
	return nil
}
