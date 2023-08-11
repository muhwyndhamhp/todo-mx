package public

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/todo-mx/db"
	"github.com/muhwyndhamhp/todo-mx/models"
	"github.com/muhwyndhamhp/todo-mx/utils/scopes"
)

type FrontendHandler struct {
}

func NewFrontendHandler(e *echo.Echo) {
	handler := &FrontendHandler{}

	e.GET("/hello", handler.Hello)
	e.GET("/", handler.Index)
	e.GET("/todos", handler.GetTodos)
	e.POST("/add-todo", handler.AddTodo)
}

func (m *FrontendHandler) GetTodos(c echo.Context) error {
	time.Sleep(1 * time.Second)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))

	todos, err := m.GetTodosFromDB(page, pageSize)
	if err != nil {
		return err
	}
	todos[len(todos)-1].Meta["IsLastItem"] = true
	todos[len(todos)-1].Meta["Page"] = page + 1

	return c.Render(http.StatusOK, "todo_list", todos)
}

func (m *FrontendHandler) AddTodo(c echo.Context) error {
	time.Sleep(1 * time.Second)
	td := models.Todo{
		Title: c.FormValue("title"),
		Body: pgtype.Text{
			String: c.FormValue("body"),
			Valid:  true,
		},
		EncodedBody: template.HTML(c.FormValue("body_encoded")),
	}

	if err := m.SaveTodoFromDB(&td); err != nil {
		return err
	}
	return c.Render(http.StatusOK, "todo_item", td)
}

func (*FrontendHandler) Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

func (m *FrontendHandler) Index(c echo.Context) error {
	todos, err := m.GetTodosFromDB(1, 5)
	if err != nil {
		return err
	}
	todos[len(todos)-1].Meta["IsLastItem"] = true
	todos[len(todos)-1].Meta["Page"] = 2

	temp := map[string]interface{}{
		"Todos": todos,
		"NewTodo": models.Todo{
			Title:       "",
			Body:        pgtype.Text{},
			EncodedBody: "",
			Meta:        models.BuildTodoMeta(),
		},
	}
	return c.Render(http.StatusOK, "index", temp)
}

func (*FrontendHandler) SaveTodoFromDB(value *models.Todo) error {
	err := db.GetDB().Save(value).Error

	return err
}

func (*FrontendHandler) GetTodosFromDB(page, pageSize int) ([]models.Todo, error) {
	var res []models.Todo
	err := db.GetDB().
		Scopes(scopes.Paginate(page, pageSize)).
		Order("updated_at desc").
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
