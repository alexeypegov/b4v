package controller

import (
	"net/http"
  "strconv"

	"github.com/alexeypegov/b4v/model"
  "github.com/alexeypegov/b4v/templates"
)

// IndexHandler handles index page
func IndexHandler(w http.ResponseWriter, r *http.Request, ctx *Context) (int, error) {
  pageString := r.URL.Query().Get(":page")
  page, err := strconv.Atoi(pageString)
	if err != nil {
		page = 1
	}

  notes, err := model.GetNotes(page, ctx.DB)
  if err != nil {
    return http.StatusInternalServerError, err
  }

	total, err := model.GetPagesCount(ctx.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	paging := &templates.Paging{Current: page, Total: total}
	data := &templates.Data{Notes: notes, Vars: ctx.Vars, Paging: paging}
  if err := ctx.Template.ExecuteTemplate(w, "index.tpl", data); err != nil {
    return http.StatusInternalServerError, err
  }

  return http.StatusOK, nil
}
