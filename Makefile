run:
	go run ./cmd/web/main.go

build:
	go build -o bin/todo ./cmd/web/main.go 

create-db:
	sqlite3 db/todo.db ".read db/todo.sql"