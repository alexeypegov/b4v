package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/alexeypegov/b4v/controller"
	"github.com/alexeypegov/b4v/model"
	"github.com/alexeypegov/b4v/templates"
	"github.com/bmizerany/pat"
	"github.com/golang/glog"
	"github.com/urfave/negroni"
	"github.com/tylerb/graceful"
)

// Config contains all of the blog configuration parameters
type Config struct {
	Port         int
	Database     string
	Templates    string
	NotesPerPage int               `toml:"notes_per_page"`
	Vars         map[string]string `toml:"vars"`
}

var (
	importPath string
	configPath string

	config Config
)

type handler struct {
	*controller.Context
	Handler func(w http.ResponseWriter, r *http.Request, ctx *controller.Context) (int, error)
}

func init() {
	flag.StringVar(&configPath, "config", "./blog.toml", "specify blog configuration path")
	flag.StringVar(&importPath, "import", "", "old format json data filename")
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
	defer glog.Flush()

	flag.Parse()

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		glog.Fatalln("Unable to read config file!\n", err)
	}

	db, err := model.OpenDB(config.Database)
	if err != nil {
		glog.Fatalln("Unable to initialize database!\n", err)
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
	if err := model.RebuildIndex(config.NotesPerPage, db); err != nil {
		glog.Fatal(err)
	}
	glog.Info("ok")

	context := &controller.Context{
		DB:       db,
		Template: templates.New(config.Templates),
		Vars:     config.Vars}

	n := negroni.New(negroni.NewRecovery(), NewLogMiddleware(0), negroni.NewStatic(http.Dir("public")))

	mux := pat.New()
	mux.Get("/", handler{context, controller.IndexHandler})
	mux.Get("/page/:page", handler{context, controller.IndexHandler})
	mux.Get("/note/:id", handler{context, controller.NoteHandler})
	n.UseHandler(mux)

	server := &graceful.Server{
		Timeout: 5 * time.Second,
		Server: &http.Server{
			Addr: fmt.Sprintf(":%d", config.Port),
			Handler: n},
	}
	
	server.ListenAndServe()
}
