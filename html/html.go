package html

import (
	"html/template"

	"github.com/russross/blackfriday"
)

func MarkdownToTemplateHTML(md []byte) template.HTML {
	html := blackfriday.MarkdownCommon([]byte(md))
	return template.HTML(string(html))
}

var MainTemplate = template.Must(template.ParseFiles("html/main.html"))

type MainArg struct {
	Body template.HTML
}
