package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/alexeypegov/b4v/controller"
	"github.com/alexeypegov/b4v/model"
	"github.com/alexeypegov/b4v/templates"
	"github.com/alexeypegov/b4v/util"
	"github.com/bmizerany/pat"
	"github.com/golang/glog"
	"github.com/urfave/negroni"
)

// Config contains all of the blog configuration parameters
type Config struct {
	NotesPerPage int               `toml:"notes_per_page"`
	Vars         map[string]string `toml:"vars"`
}

var (
	dataPath   string
	importPath string
	port       int
	config     Config
)

type handler struct {
	*controller.Context
	Handler func(w http.ResponseWriter, r *http.Request, ctx *controller.Context) (int, error)
}

func init() {
	flag.StringVar(&dataPath, "data", "./data", "override blog data path")
	flag.IntVar(&port, "port", 8080, "override server port")
	flag.StringVar(&importPath, "import", "", "import old format json data filename")
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

	db, err := model.OpenDB(filepath.Join(dataPath, "notes.db"))
	if err != nil {
		glog.Fatalln("Unable to initialize database!\n", err)
	}
	defer db.Close()

	if len(importPath) > 0 {
		glog.Infof("Using import file '%s'... ", importPath)
		notesCount, err := model.Populate(importPath, db)
		if err != nil {
			glog.Warning(err)
		} else {
			glog.Infof("ok [%d notes]", notesCount)
		}
	}

	if _, err := toml.DecodeFile("blog.toml", &config); err != nil {
		glog.Fatalln("Unable to read config file!\n", err)
	}

	glog.Info("Rebuilding index... ")
	if err := model.RebuildIndex(config.NotesPerPage, db); err != nil {
		glog.Fatal(err)
	}
	glog.Info("ok")

	ctx := &controller.Context{
		DB:       db,
		Template: templates.New("."),
		Vars:     config.Vars}

	n := negroni.New(negroni.NewRecovery(), util.NewLogMiddleware(0), negroni.NewStatic(http.Dir("public")))

	mux := pat.New()
	mux.Get("/", handler{ctx, controller.IndexHandler})
	mux.Get("/page/:page", handler{ctx, controller.IndexHandler})
	mux.Get("/note/:id", handler{ctx, controller.NoteHandler})
	mux.Get("/rss", handler{ctx, controller.RssHandler})
	n.UseHandler(mux)

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: n,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			glog.Error(err)
		}
	}()

	<-stopChan
	_ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(_ctx)

	glog.Info("Server was shutdown gracefully")
}
