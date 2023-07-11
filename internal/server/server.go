package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hunterwilkins2/todomvc/internal/config"
	"github.com/hunterwilkins2/todomvc/internal/routes"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	logger *log.Logger
	config *config.Config
	server *http.Server
	db     *sql.DB
}

func New(config *config.Config) *Server {
	logger := log.New(os.Stdout, "[TODO MVC] ", log.Ldate|log.Ltime)

	db, err := openDB(config.DSN)
	if err != nil {
		logger.Fatalf("error opening database: %s", err.Error())
	}

	app := routes.New(config, db, logger)
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      app.Routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	return &Server{
		logger: logger,
		config: config,
		server: &server,
		db:     db,
	}
}

func (s *Server) Start() {
	defer s.db.Close()
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		s.logger.Printf("shutting down server. received signal %s", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- s.server.Shutdown(ctx)
	}()

	s.logger.Printf("starting server on %s\n", s.server.Addr)
	err := s.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatalf("an uncaught error occurred: %s", err.Error())
	}

	s.logger.Println("server stopped")
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
