package fly

import (
	"embed"
	"html/template"
	"net/http"
	"os"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

func RenderIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Region": os.Getenv("FLY_REGION"),
	}
	t.ExecuteTemplate(w, "index.html.tmpl", data)
}
