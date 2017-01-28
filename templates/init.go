package templates

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/alexeypegov/b4v/model"
)

const (
	ext = ".tpl"
)

// Data for templates
type Data struct {
	Notes []*model.Note
	Note  *model.Note
	Vars  *Vars
}

// Vars contrains general blog info for templates
type Vars struct {
	Title     string
	Copyright string
}

var funcMap = template.FuncMap{
	"html": func(s string) template.HTML {
		return template.HTML(s)
	},
	"timestamp": func(ts time.Time) string {
		return fmt.Sprintf("%02d/%02d/%4d", ts.Day(), ts.Month(), ts.Year())
	},
	"minus": func(a, b int) int {
		return a - b
	},
	"plus": func(a, b int) int {
		return a + b
	},
}

// New initializes templates
func New(path string) *template.Template {
	pattern := filepath.Join(path, fmt.Sprintf("*%s", ext))
	tpl := template.Must(template.New("main").Funcs(funcMap).ParseGlob(pattern))
	return tpl
}
