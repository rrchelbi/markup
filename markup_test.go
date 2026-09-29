package markup_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/rrchelbi/markup"
	"github.com/rrchelbi/markup/attr"
	"github.com/rrchelbi/markup/flow"
	"github.com/rrchelbi/markup/html"
)

func TestRender(t *testing.T) {
	items := []string{"a", "b"}
	tests := []struct {
		name string
		node markup.Node
		want string
	}{
		{"text escaped", html.P(markup.Text(`<b>"x" & 'y'</b>`)),
			`<p>&lt;b&gt;&#34;x&#34; &amp; &#39;y&#39;&lt;/b&gt;</p>`},
		{"attr escaped", html.Div(attr.Title(`a"b<c>`)), `<div title="a&#34;b&lt;c&gt;"></div>`},
		{"void", html.Br(), `<br>`},
		{"void attrs", html.Img(attr.Src("/a.png"), attr.Alt("x")), `<img src="/a.png" alt="x">`},
		{"bool attr", html.Input(attr.Disabled()), `<input disabled>`},
		{"class merge", html.Div(attr.Class("a"), attr.Class("b")), `<div class="a b"></div>`},
		{"attr override", html.A(attr.Href("/a"), attr.Href("/b")), `<a href="/b"></a>`},
		{"nil skipped", html.Div(nil, markup.Text("x")), `<div>x</div>`},
		{"group flattens attrs", html.Div(markup.Group{attr.ID("x"), markup.Text("y")}), `<div id="x">y</div>`},
		{"if false attr", html.Div(flow.If(false, attr.ID("x"))), `<div></div>`},
		{"if true attr", html.Div(flow.If(true, attr.ID("x"))), `<div id="x"></div>`},
		{"each", html.Ul(flow.Each(items, func(s string) markup.Node { return html.Li(html.Text(s)) })),
			`<ul><li>a</li><li>b</li></ul>`},
		{"raw", html.Div(markup.Raw("<i>ok</i>")), `<div><i>ok</i></div>`},
		{"data attr", html.Div(attr.Data("id", "7")), `<div data-id="7"></div>`},
		{"script json", html.Script(attr.Type("application/json"), markup.Text(`{"a":"<b>"}`)),
			`<script type="application/json">{"a":"<b>"}</script>`},
		{"unsafe href", html.A(attr.Href("javascript:alert(1)")), `<a href="#unsafe"></a>`},
		{"NewAttr href sanitized", html.A(markup.NewAttr("href", "javascript:x")), `<a href="#unsafe"></a>`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := markup.String(tt.node)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got  %s\nwant %s", got, tt.want)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		node markup.Node
		want error
	}{
		{"tag name", markup.Tag("div onload=x"), markup.ErrInvalidName},
		{"empty tag", markup.Tag(""), markup.ErrInvalidName},
		{"attr name", html.Div(markup.NewAttr(`a" onclick="x`, "v")), markup.ErrInvalidName},
		{"onclick", html.Div(markup.NewAttr("onclick", "x")), markup.ErrUnsafeAttr},
		{"ONCLICK", html.Div(markup.NewAttr("ONCLICK", "x")), markup.ErrUnsafeAttr},
		{"srcdoc", html.Iframe(markup.NewAttr("srcdoc", "<script>")), markup.ErrUnsafeAttr},
		{"hx-on", html.Div(markup.NewAttr("hx-on:click", "x")), markup.ErrUnsafeAttr},
		{"void children", html.Br(markup.Text("x")), markup.ErrVoidChildren},
		{"script close", html.Script(markup.Text("a</script><b>")), markup.ErrUnsafeRawText},
		{"script close case", html.Script(markup.Text("a</SCRIPT >")), markup.ErrUnsafeRawText},
		{"script comment", html.Script(markup.Text("<!--")), markup.ErrUnsafeRawText},
		{"style close", html.Style(markup.Text("</style>")), markup.ErrUnsafeRawText},
		{"script element child", html.Script(html.Div()), markup.ErrUnsafeRawText},
		{"zero attr", html.Div(markup.Attr{}), markup.ErrInvalidName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := markup.String(tt.node); !errors.Is(err, tt.want) {
				t.Fatalf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestUnsafeAttrEscapeHatch(t *testing.T) {
	got, err := markup.String(html.Button(markup.UnsafeAttr("onclick", `go("x")`)))
	if err != nil {
		t.Fatal(err)
	}
	if want := `<button onclick="go(&#34;x&#34;)"></button>`; got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestSafeURL(t *testing.T) {
	const bad = "#unsafe"
	tests := []struct{ in, want string }{
		{"", ""},
		{"/x", "/x"},
		{"x.html", "x.html"},
		{"//host/x", "//host/x"},
		{"https://a.b/c", "https://a.b/c"},
		{"HTTP://a.b", "HTTP://a.b"},
		{"mailto:a@b.c", "mailto:a@b.c"},
		{"tel:+1", "tel:+1"},
		{"/p?a=b:c", "/p?a=b:c"},
		{"p.html#a:b", "p.html#a:b"},
		{"javascript:alert(1)", bad},
		{" \x01JaVaScRiPt:alert(1)", bad},
		{"java\tscript:alert(1)", bad},
		{"java\nscript:alert(1)", bad},
		{"data:text/html,x", bad},
		{"vbscript:x", bad},
		{"averyveryverylongscheme:x", bad},
		{"a:b", bad},
	}
	for _, tt := range tests {
		if got := markup.SafeURL(tt.in); got != tt.want {
			t.Errorf("SafeURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

type failWriter struct{ after int }

func (f *failWriter) Write(p []byte) (int, error) {
	if f.after < len(p) {
		return 0, io.ErrClosedPipe
	}
	f.after -= len(p)
	return len(p), nil
}

func TestWriteError(t *testing.T) {
	page := html.Div(html.P(markup.Text(strings.Repeat("x", 10000))))
	if err := markup.Write(&failWriter{after: 10}, page); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("err = %v, want ErrClosedPipe", err)
	}
	if err := markup.Write(io.Discard, nil); err != nil {
		t.Fatal(err)
	}
}

func FuzzText(f *testing.F) {
	for _, s := range []string{"", "plain", `<script>alert("x")</script>`, "a&b", "\x00", "'\""} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		out, err := markup.String(markup.Text(s))
		if err != nil {
			t.Fatal(err)
		}
		if strings.ContainsAny(out, `<>"'`) {
			t.Fatalf("unescaped output %q for %q", out, s)
		}
		rest := out
		for _, e := range []string{"&amp;", "&lt;", "&gt;", "&#34;", "&#39;"} {
			rest = strings.ReplaceAll(rest, e, "")
		}
		if strings.Contains(rest, "&") {
			t.Fatalf("stray & in %q for %q", out, s)
		}
	})
}

func FuzzURL(f *testing.F) {
	for _, s := range []string{"javascript:x", "/a", "java\tscript:x", " data:x"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		out := markup.SafeURL(s)
		if out == "#unsafe" {
			return
		}
		norm := strings.NewReplacer("\t", "", "\n", "", "\r", "").Replace(strings.TrimLeft(out, "\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f "))
		norm = strings.ToLower(norm)
		for _, p := range []string{"javascript:", "vbscript:", "data:"} {
			if strings.HasPrefix(norm, p) {
				t.Fatalf("SafeURL(%q) = %q passes %s", s, out, p)
			}
		}
	})
}

func page() markup.Node {
	rows := make([]string, 20)
	for i := range rows {
		rows[i] = "row"
	}
	return html.HTML(
		html.Head(html.Title(html.Text("t")), html.Meta(attr.Charset("utf-8"))),
		html.Body(attr.Class("app"),
			html.Nav(html.A(attr.Href("/"), html.Text("home"))),
			html.Ul(flow.Each(rows, func(s string) markup.Node {
				return html.Li(attr.Class("item"), html.A(attr.Href("/i/"+s), html.Text(s)))
			})),
		),
	)
}

func BenchmarkBuildAndRender(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := markup.Write(io.Discard, page()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRenderOnly(b *testing.B) {
	p := page()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := markup.Write(io.Discard, p); err != nil {
			b.Fatal(err)
		}
	}
}
