package main

import (
    "embed"
    "html/template"
    "io"

    "github.com/labstack/echo/v4"
)

//go:embed web/*.html
var tplFS embed.FS

type TemplateRenderer struct{ t *template.Template }

func newRenderer() *TemplateRenderer {
    t := template.Must(template.ParseFS(tplFS, "web/*.html"))
    return &TemplateRenderer{t: t}
}

func (r *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return r.t.ExecuteTemplate(w, name, data)
}
