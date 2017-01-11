package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/alexeypegov/b4v/controller"
	"github.com/alexeypegov/b4v/model"
	"github.com/alexeypegov/b4v/templates"
	"github.com/bmizerany/pat"
	"github.com/urfave/negroni"
	"github.com/golang/glog"
)

var (
	port          int
	importPath    string
	databasePath  string
	templatesPath string
)

type handler struct {
	*controller.Context
	Handler func(w http.ResponseWriter, r *http.Request, ctx *controller.Context) (int, error)
}

func init() {
	flag.IntVar(&port, "port", 8080, "server port")
	flag.StringVar(&importPath, "import", "", "old format json data filename")
	flag.StringVar(&databasePath, "db", "./b4v.db", "database path")
	flag.StringVar(&templatesPath, "templates", "./templates", "path to template files")
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := h.Handler(w, r, h.Context)
	if err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			break
		default:
			glog.Infof("HTTP %d: %q", status, err)
			break
		}
	}
}

func main() {
	flag.Parse()

	db, err := model.OpenDB(databasePath)
	if err != nil {
		glog.Fatal("Unable to initialize db", err)
	}
	defer db.Close()

	if len(importPath) > 0 {
		glog.Infof("Using import file '%s'... ", importPath)
		notesCount, err := model.Populate(importPath, db)
		if err != nil {
			glog.Fatal(err)
		}

		glog.Infof("ok [%d notes]", notesCount)
	}

	glog.Info("Rebuilding index... ")
	if err := model.RebuildIndex(db); err != nil {
		glog.Fatal(err)
	}
	glog.Info("ok")

	context := &controller.Context{DB: db, Template: templates.New("./templates")}
	mux := pat.New()
	mux.Get("/", handler{context, controller.IndexHandler})
	mux.Get("/note/:id", handler{context, controller.NoteHandler})
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(fmt.Sprintf(":%d", port))
}
