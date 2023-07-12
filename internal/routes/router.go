package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/hunterwilkins2/todomvc/internal/config"
	"github.com/hunterwilkins2/todomvc/internal/data"
	"github.com/hunterwilkins2/todomvc/internal/templates"
)

type Application struct {
	config    *config.Config
	logger    *log.Logger
	templates *templates.TemplateCache
	Models    *data.Models
	query     data.Query
}

func New(config *config.Config, db *sql.DB, logger *log.Logger) *Application {
	cache, err := templates.New("./ui/html/")
	if err != nil {
		logger.Fatalf("error caching templates: %s", err.Error())
	}

	return &Application{
		config:    config,
		logger:    logger,
		templates: cache,
		Models:    data.NewModels(db),
		query:     data.All,
	}
}

func (app *Application) Routes() http.Handler {
	mux := flow.New()

	mux.HandleFunc("/", app.Home, http.MethodGet)
	mux.HandleFunc("/todo", app.TodosHandler, http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete)
	mux.HandleFunc("/todo/:id", app.TodoHandler, http.MethodGet, http.MethodPatch, http.MethodDelete)
	mux.HandleFunc("/query/:selected", app.UpdateQuery, http.MethodPost)

	if app.config.HotReload {
		mux.HandleFunc("/reload", app.hotReload)
		mux.HandleFunc("/reload-ready", app.reloadReady)
	}

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/...", http.StripPrefix("/static/", fileServer), http.MethodGet)

	return mux
}
