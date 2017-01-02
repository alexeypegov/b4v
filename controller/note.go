package controller

import (
	"fmt"
	"net/http"

	"github.com/alexeypegov/b4v/model"
)

// NoteHandler handles note requests
func NoteHandler(w http.ResponseWriter, r *http.Request, ctx *Context) (int, error) {
	id := r.URL.Query().Get(":id")
	if len(id) == 0 {
		return http.StatusNotFound, fmt.Errorf("Not found")
	}

	note, err := model.GetNote(id, ctx.DB)
	if err != nil {
		fmt.Fprintf(w, "Not found %v\n", note.Title)
		return http.StatusNotFound, err
	}

	fmt.Fprintf(w, "Hello %v\n", note.Title)
	return http.StatusOK, nil
}
