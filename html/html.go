package html

import (
	"html/template"

	"github.com/russross/blackfriday"
)

func MarkdownToTemplateHTML(md []byte) template.HTML {
	flags := blackfriday.HTML_USE_XHTML
	renderer := blackfriday.HtmlRenderer(flags, "", "")

	extensions := blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_AUTOLINK

	html := blackfriday.Markdown([]byte(md), renderer, extensions)
	return template.HTML(string(html))
}

var MainTemplate = template.Must(template.ParseFiles("html/main.html"))

type MainArg struct {
	Body template.HTML
}
