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

	if len(os.Args) < 2 {
		fmt.Println("Usage: b4v command (commands are: populate, start)")
		return
	}

	switch os.Args[1] {
	case "populate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: b4v populate <backup.json>")
		} else {
			if err := model.Populate(os.Args[2], db); err != nil {
				panic(err)
			}
		}
		break
	case "start":
		context := &controller.Context{DB: db}
		mux := pat.New()
		mux.Get("/note/:id", handler{context, controller.NoteHandler})
		n := negroni.Classic()
		n.UseHandler(mux)
		n.Run(":3000")
		break
	}
}
