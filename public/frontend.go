package public

import (
	"html/template"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/todo-mx/db"
	"github.com/muhwyndhamhp/todo-mx/models"
)

type FrontendHandler struct {
}

func NewFrontendHandler(e *echo.Echo) {
	handler := &FrontendHandler{}

	e.GET("/hello", handler.Hello)
	e.GET("/", handler.Index)
	e.POST("/add-todo", handler.AddTodos)
}

func (m *FrontendHandler) AddTodos(c echo.Context) error {
	time.Sleep(1 * time.Second)
	td := models.Todo{
		Title: c.FormValue("title"),
		Body: pgtype.Text{
			String: c.FormValue("body"),
			Valid:  true,
		},
		EncodedBody: template.HTML(c.FormValue("body_encoded")),
	}

	if err := m.SaveTodo(&td); err != nil {
		return err
	}
	return c.Render(http.StatusOK, "todo_item", td)
}

func (*FrontendHandler) Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

func (m *FrontendHandler) Index(c echo.Context) error {
	todos, err := m.GetTodos()
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "index", todos)
}

func (*FrontendHandler) SaveTodo(value *models.Todo) error {
	err := db.GetDB().Save(value).Error

	return err
}

func (*FrontendHandler) GetTodos() ([]models.Todo, error) {
	var res []models.Todo
	err := db.GetDB().Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
