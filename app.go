package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexeypegov/b4v/controller"
	"github.com/alexeypegov/b4v/model"
	"github.com/bmizerany/pat"
	"github.com/urfave/negroni"
)

const (
	port   = 8080
	dbFile = "./b4v.db"
)

type handler struct {
	*controller.Context
	H func(w http.ResponseWriter, r *http.Request, ctx *controller.Context) (int, error)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := h.H(w, r, h.Context)
	if err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			break
		default:
			log.Printf("HTTP %d: %q", status, err)
			break
		}
	}
}

func main() {
	db, err := model.OpenDB(dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if len(os.Args) == 2 {
		fmt.Print(fmt.Sprintf("Using import file '%s'... ", os.Args[1]))
		notesCount, err := model.Populate(os.Args[1], db)
		if err != nil {
			fmt.Println(fmt.Sprintf("fail (%s)", err.Error()))
			return
		}

		fmt.Println(fmt.Sprintf("ok [%d notes]", notesCount))
	}

	fmt.Print("Rebuilding index... ")
	if err := model.RebuildIndex(db); err != nil {
		fmt.Println(fmt.Sprintf("fail (%s)", err.Error()))
		return
	}

	fmt.Println("ok")

	context := &controller.Context{DB: db}
	mux := pat.New()
	mux.Get("/note/:id", handler{context, controller.NoteHandler})
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(fmt.Sprintf(":%d", port))
}
