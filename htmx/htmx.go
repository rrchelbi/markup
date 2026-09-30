// Package htmx provides attributes for htmx (https://htmx.org).
//
// hx-on:* attributes are deliberately absent: they run inline code.
package htmx

import "github.com/rrchelbi/markup"

func Get(url string) markup.Attr    { return markup.URL("hx-get", url) }
func Post(url string) markup.Attr   { return markup.URL("hx-post", url) }
func Put(url string) markup.Attr    { return markup.URL("hx-put", url) }
func Patch(url string) markup.Attr  { return markup.URL("hx-patch", url) }
func Delete(url string) markup.Attr { return markup.URL("hx-delete", url) }

func Target(selector string) markup.Attr    { return markup.NewAttr("hx-target", selector) }
func Swap(strategy string) markup.Attr      { return markup.NewAttr("hx-swap", strategy) }
func Trigger(spec string) markup.Attr       { return markup.NewAttr("hx-trigger", spec) }
func Select(selector string) markup.Attr    { return markup.NewAttr("hx-select", selector) }
func Indicator(selector string) markup.Attr { return markup.NewAttr("hx-indicator", selector) }
func Confirm(message string) markup.Attr    { return markup.NewAttr("hx-confirm", message) }
func PushURL(v string) markup.Attr          { return markup.NewAttr("hx-push-url", v) }
func Boost() markup.Attr                    { return markup.NewAttr("hx-boost", "true") }
