# markup

HTML as Go values. Safe by default, no codegen at build time, no reflection.

```go
import (
	"markup/attr"
	"markup/flow"
	. "markup/html"
)

Div(attr.Class("card"),
	H1(Text("Hello")),
	Ul(flow.Each(items, func(s string) Node { return Li(Text(s)) })),
)
```

## Packages

| Package | Purpose |
|---|---|
| `markup` | Core: `Node`, `Tag`, `Attr`, `Text`, `Raw`, `Group`, `Write`, `String` |
| `markup/html` | Element constructors (generated) |
| `markup/attr` | Attribute constructors (generated) |
| `markup/flow` | `If`, `IfElse`, `When`, `Each` |
| `markup/htmx` | htmx attributes |
| `markup/mhttp` | `http.Handler` that renders to a buffer first |

## Safety

- Text and attribute values are escaped.
- Tag and attribute names are validated; `on*` and `srcdoc` are rejected.
- URL attributes allow only relative URLs and `http`, `https`, `mailto`, `tel`.
- `<script>`/`<style>` text is rejected if it could close the element.
- Escape hatches: `markup.Raw`, `markup.UnsafeAttr`. Never feed them user input.

Not covered: CSS in `style` attributes and expressions in framework
attributes (Alpine `x-data`, etc.) are escaped but not interpreted.

## Development

```bash
go generate ./...   # regenerate html/ and attr/
go test -race ./...
```
