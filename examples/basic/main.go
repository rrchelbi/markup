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
