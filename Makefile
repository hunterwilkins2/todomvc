run:
	go run ./cmd/web/main.go

run/live:
	gow -c -e=go,mod,html,js,css run ./cmd/web/main.go -hot-reload

build:
	go build -o bin/todo ./cmd/web/main.go 

create-db:
	sqlite3 db/todo.db ".read db/todo.sql"