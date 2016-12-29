package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alexeypegov/b4v/model"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

func main() {
  db, err := model.OpenDB("./b4v.db")
  if err != nil {
    panic(err)
  }

	defer db.Close()

	if len(os.Args) < 2 {
		fmt.Println("Usage: b4v command")
		return
	}

	switch os.Args[1] {
	case "populate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: b4v populate <backup.json>")
		} else {
			model.Populate(os.Args[2], db)
		}
		break
	case "start":
		mux := httprouter.New()
		mux.GET("/", index)
		mux.GET("/note/:id", printNote)
		mux.ServeFiles("/static/*filepath", http.Dir("public"))
		n := negroni.Classic()
		n.UseHandler(mux)
		n.Run(":3000")
		break
	}
}

func printNote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello %v\n", ps.ByName("id"))
}

func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello world!\n")
}
