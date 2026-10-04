# markup

HTML as Go values. Safe by default, no codegen at build time, no reflection.

```go
package main

import (
	"log"
	"net/http"

	"github.com/rrchelbi/markup"
	"github.com/rrchelbi/markup/attr"
	"github.com/rrchelbi/markup/flow"
	. "github.com/rrchelbi/markup/html"
	"github.com/rrchelbi/markup/htmx"
	"github.com/rrchelbi/markup/mhttp"
)

func card(title string, tags []string, featured bool) markup.Node {
	return Div(attr.Class("card"),
		flow.If(featured, attr.Class("featured")),
		H2(Text(title)),
		Ul(flow.Each(tags, func(t string) Node { return Li(Text(t)) })),
		Button(htmx.Post("/like"), htmx.Swap("outerHTML"), Text("Like")),
	)
}

func index(r *http.Request) markup.Node {
	name := r.URL.Query().Get("name") // untrusted: escaped automatically
	return Group{
		Doctype,
		HTML(attr.Lang("en"),
			Head(Meta(attr.Charset("utf-8")), Title(Text("markup"))),
			Body(
				H1(Textf("Hello, %s", name)),
				card("First", []string{"go", "html"}, true),
			),
		),
	}
}

func main() {
	http.Handle("/", mhttp.HandlerFunc(index))
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
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
