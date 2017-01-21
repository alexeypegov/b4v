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
  page, _ := strconv.Atoi(pageString)

  notes, err := model.GetNotes(page, ctx.DB)
  if err != nil {
    return http.StatusInternalServerError, err
  }

  if err := ctx.Template.ExecuteTemplate(w, "index.tpl", &templates.Data{Notes: notes, Vars: ctx.Vars}); err != nil {
    return http.StatusInternalServerError, err
  }

  return http.StatusOK, nil
}
