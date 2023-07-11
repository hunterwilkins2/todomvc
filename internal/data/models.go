package data

import "database/sql"

type Models struct {
	Todo TodoModelInterface
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Todo: TodoModel{db},
	}
}
