package public

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/muhwyndhamhp/todo-mx/db"
	"github.com/muhwyndhamhp/todo-mx/models"
	"github.com/muhwyndhamhp/todo-mx/utils/constants"
	"github.com/muhwyndhamhp/todo-mx/utils/errs"
	"github.com/muhwyndhamhp/todo-mx/utils/markd"
	"github.com/muhwyndhamhp/todo-mx/utils/scopes"
	"gorm.io/gorm"
)

type FrontendHandler struct {
}

const (
	INDEX_PATH    = "/"
	GET_TODO_PATH = "/todos"
	ADD_TODO_PATH = "/todos"

	EDIT_TODO_PATH_PREFIX = "/todos/"
	EDIT_TODO_PATH_SUFFIX = "/edit"

	UPDATE_TODO_PATH_PREFIX = "/todos/"
)

func NewFrontendHandler(e *echo.Echo) {
	handler := &FrontendHandler{}

	e.GET("/hello", handler.Hello)

	e.GET(INDEX_PATH, handler.Index)
	e.GET(GET_TODO_PATH, handler.GetTodos)

	e.POST(ADD_TODO_PATH, handler.AddTodo)
	e.GET(EDIT_TODO_PATH_PREFIX+":id"+EDIT_TODO_PATH_SUFFIX, handler.EditTodo)
	e.PUT(UPDATE_TODO_PATH_PREFIX+":id", handler.UpdateTodo)
}

func (*FrontendHandler) Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

func (m *FrontendHandler) Index(c echo.Context) error {
	todos, err := m.GetTodosFromDB(1, 5)
	if err != nil {
		return err
	}
	AppendLastItemMetadata(1, todos)

	temp := map[string]interface{}{
		"Todos": todos,
		"NewTodo": models.Todo{
			Title:       "",
			Body:        pgtype.Text{},
			EncodedBody: "",
			Meta:        models.BuildTodoMeta(ADD_TODO_PATH, &models.Todo{}),
		},
		"AddPath": ADD_TODO_PATH,
	}
	return c.Render(http.StatusOK, "index", temp)
}

func (m *FrontendHandler) GetTodos(c echo.Context) error {
	time.Sleep(1 * time.Second)
	page, _ := strconv.Atoi(c.QueryParam(constants.PAGE))
	pageSize, _ := strconv.Atoi(c.QueryParam(constants.PAGE_SIZE))

	todos, err := m.GetTodosFromDB(page, pageSize)
	if err != nil {
		return err
	}
	AppendLastItemMetadata(page, todos)

	return c.Render(http.StatusOK, "todo_list", todos)
}

func (m *FrontendHandler) UpdateTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}
	td := models.Todo{
		Model: gorm.Model{
			ID: uint(id),
		},
		Title: c.FormValue("title"),
		Body: pgtype.Text{
			String: c.FormValue("body"),
			Valid:  true,
		},
	}

	md, err := markd.ParseMD(c.FormValue("body"))
	if err != nil {
		return err
	}

	td.EncodedBody = template.HTML(md)

	if err := m.UpdateTodoDB(&td); err != nil {
		return err
	}
	return c.Render(http.StatusOK, "todo_item", td)
}

func (m *FrontendHandler) AddTodo(c echo.Context) error {
	time.Sleep(1 * time.Second)
	td := models.Todo{
		Title: c.FormValue("title"),
		Body: pgtype.Text{
			String: c.FormValue("body"),
			Valid:  true,
		},
	}

	md, err := markd.ParseMD(c.FormValue("body"))
	if err != nil {
		return err
	}

	td.EncodedBody = template.HTML(md)

	if err := m.SaveTodoFromDB(&td); err != nil {
		return err
	}
	return c.Render(http.StatusOK, "todo_item", td)
}

func (m *FrontendHandler) EditTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id <= 0 {
		return c.JSON(http.StatusBadRequest, nil)
	}

	todo, err := m.GetTodoByID(uint(id))
	if err != nil {
		return errs.Wrap(err)
	}

	todo.Meta = models.BuildTodoMeta(fmt.Sprintf("%s%d", UPDATE_TODO_PATH_PREFIX, todo.ID), todo)

	return c.Render(http.StatusOK, "todo_edit", todo)
}

func AppendLastItemMetadata(lastPage int, todos []models.Todo) {
	if len(todos) <= 1 {
		return
	}
	for i := range todos {
		todos[i].Meta["EditPath"] = fmt.Sprintf(
			"%s%d%s",
			EDIT_TODO_PATH_PREFIX,
			todos[i].ID,
			EDIT_TODO_PATH_SUFFIX,
		)
	}
	todos[len(todos)-1].Meta["IsLastItem"] = true
	todos[len(todos)-1].Meta["Page"] = lastPage + 1
}

func (*FrontendHandler) GetTodoByID(id uint) (*models.Todo, error) {
	var res models.Todo
	if err := db.GetDB().First(&res, id).Error; err != nil {
		return nil, errs.Wrap(err)
	}

	return &res, nil
}

func (*FrontendHandler) SaveTodoFromDB(value *models.Todo) error {
	err := db.GetDB().Save(value).Error

	return err
}

func (*FrontendHandler) UpdateTodoDB(value *models.Todo) error {
	err := db.GetDB().Save(value).Error

	return err
}

func (*FrontendHandler) GetTodosFromDB(page, pageSize int) ([]models.Todo, error) {
	var res []models.Todo
	err := db.GetDB().
		Scopes(scopes.Paginate(page, pageSize)).
		Order("created_at desc").
		Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
