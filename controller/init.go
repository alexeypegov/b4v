package controller

import (
	"html/template"

	"github.com/alexeypegov/b4v/model"
	"github.com/alexeypegov/b4v/templates"
)

// Context holds handler context parameters
type Context struct {
	DB       *model.DB
	Template *template.Template
	Vars     *templates.Vars
}
