package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alexedwards/flow"
	"github.com/hunterwilkins2/todomvc/internal/data"
)

var (
	ErrServerError = errors.New("could not complete request")
)

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	err := app.templates.Render(w, "home.page.html", struct {
		HotReload bool
	}{HotReload: app.config.HotReload})
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) RenderTodoComponent(w http.ResponseWriter, err error) {
	if err != nil {
		app.serverError(w, err)
	}

	todoData, err := app.Models.Todo.GetAll(app.query)
	if err != nil {
		app.serverError(w, err)
	}

	err = app.templates.Render(w, "todos.partial.html", todoData)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) UpdateQuery(w http.ResponseWriter, r *http.Request) {
	query := flow.Param(r.Context(), "selected")
	switch query {
	case "all":
		app.query = data.All
	case "active":
		app.query = data.Active
	case "completed":
		app.query = data.Completed
	}

	app.RenderTodoComponent(w, nil)
}

func (app *Application) TodosHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case http.MethodPost:
		err = app.addTodo(r)
	case http.MethodPatch:
		err = app.toggleCompleted(r)
	case http.MethodDelete:
		err = app.deleteCompleted(r)
	}

	app.RenderTodoComponent(w, err)
}

func (app *Application) TodoHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	idStr := flow.Param(r.Context(), "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.serverError(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.getTodoInput(id, w)
		return
	case http.MethodPatch:
		err = app.updateTodo(id, r)
	case http.MethodDelete:
		err = app.deleteTodo(id)
	}

	app.RenderTodoComponent(w, err)
}

func (app *Application) addTodo(r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return ErrServerError
	}

	name := r.PostForm.Get("new-todo")
	return app.Models.Todo.Insert(name)
}

func (app *Application) toggleCompleted(r *http.Request) error {
	todosData, err := app.Models.Todo.GetAll(app.query)
	if err != nil {
		return err
	}

	return app.Models.Todo.UpdateStatus(todosData.NotCompleted != 0)
}

func (app *Application) deleteCompleted(r *http.Request) error {
	return app.Models.Todo.DeleteCompleted()
}

func (app *Application) getTodoInput(id int, w http.ResponseWriter) {
	todo, err := app.Models.Todo.GetTodo(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.templates.Render(w, "todo-edit.partial.html", todo)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) updateTodo(id int, r *http.Request) error {
	todo, err := app.Models.Todo.GetTodo(id)
	if err != nil {
		return err
	}

	err = r.ParseForm()
	if err != nil {
		return err
	}
	editedName := r.PostForm.Get("edit-todo")
	if editedName != "" {
		todo.Name = editedName
	} else {
		todo.Completed = !todo.Completed
	}

	return app.Models.Todo.UpdateTodo(&todo)
}

func (app *Application) deleteTodo(id int) error {
	return app.Models.Todo.Delete(id)
}
