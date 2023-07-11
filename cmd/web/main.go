package main

import (
	"flag"

	"github.com/hunterwilkins2/todomvc/internal/config"
	"github.com/hunterwilkins2/todomvc/internal/server"
)

func main() {
	cfg := config.Config{}
	flag.IntVar(&cfg.Port, "port", 4000, "Server port")
	flag.StringVar(&cfg.DSN, "dsn", "./db/todo.db", "Sqlite3 DSN")

	server := server.New(&cfg)
	server.Start()
}
