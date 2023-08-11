package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/apsystole/log"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhwyndhamhp/todo-mx/models"
	"github.com/muhwyndhamhp/todo-mx/public"
	"github.com/muhwyndhamhp/todo-mx/utils/resp"
	"golang.org/x/time/rate"
)

func main() {

	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))

	tmpl := template.Must(template.ParseGlob("public/views/*.html"))
	template.Must(tmpl.ParseGlob("public/components/*.html"))

	t := public.NewTemplate(tmpl)
	e.Renderer = t

	e.GET("/hello", Hello)
	e.GET("/", Index)
	e.POST("/add-todo", AddTodos)
	e.Static("/dist", "dist")
	e.Static("/assets", "public/assets")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", "4040")))
}

func AddTodos(c echo.Context) error {
	time.Sleep(1 * time.Second)
	td := models.Todo{
		Title: c.FormValue("title"),
		Body: pgtype.Text{
			String: c.FormValue("body"),
			Valid:  true,
		},
	}
	return c.Render(http.StatusOK, "todo_item", td)
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

func Index(c echo.Context) error {
	todos := map[string][]models.Todo{
		"Todos": {
			{
				Title: "Clean the house!",
				Body: pgtype.Text{
					String: "You need to wipe the floor, mop it, and clean the tables",
					Valid:  true,
				},
			},
			{
				Title: "Kiss Ma Wife!",
				Body: pgtype.Text{
					String: "kiss your wife for happinex x100",
					Valid:  true,
				},
			},
			{
				Title: "The Thing...",
				Body: pgtype.Text{
					String: "the thing is, IDK what to write again in here",
					Valid:  true,
				},
			},
		},
	}
	return c.Render(http.StatusOK, "index", todos)
}

func httpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if code != http.StatusInternalServerError {
		_ = c.JSON(code, err)
	} else {
		log.Error(err)
		_ = resp.HTTPServerError(c)
	}
}
