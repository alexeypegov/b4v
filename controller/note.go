package controller

import (
	"fmt"
	"net/http"

	"github.com/alexeypegov/b4v/model"
	"github.com/alexeypegov/b4v/templates"
)

// NoteHandler handles note requests
func NoteHandler(w http.ResponseWriter, r *http.Request, ctx *Context) (int, error) {
	id := r.URL.Query().Get(":id")
	if len(id) == 0 {
		return http.StatusNotFound, fmt.Errorf("Not found")
	}

	note, err := model.GetNote(id, ctx.DB)
	if err != nil {
		return http.StatusNotFound, err
	}

	if err := ctx.Template.ExecuteTemplate(w, templates.Template, &templates.Data{Note: note, Vars: ctx.Vars}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
