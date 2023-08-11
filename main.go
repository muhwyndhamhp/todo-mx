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
	"github.com/muhwyndhamhp/todo-mx/config"
	"github.com/muhwyndhamhp/todo-mx/models"
	"github.com/muhwyndhamhp/todo-mx/public"
	"github.com/muhwyndhamhp/todo-mx/utils/resp"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

var todos = map[string][]models.Todo{
	"Todos": {
		{
			Model: gorm.Model{
				ID: 1,
			},
			Title: "Clean the house!",
			Body: pgtype.Text{
				String: "You need to wipe the floor, mop it, and clean the tables",
				Valid:  true,
			},
			EncodedBody: "You need to wipe the floor, mop it, and clean the tables",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Title: "Kiss Ma Wife!",
			Body: pgtype.Text{
				String: "kiss your wife for happinex x100",
				Valid:  true,
			},
			EncodedBody: "kiss your wife for happinex x100",
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Title: "The Thing...",
			Body: pgtype.Text{
				String: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tincidunt erat vulputate vehicula gravida. Nullam tincidunt vehicula lorem ac ultricies. Proin elit libero, dignissim sed dolor sed, aliquet euismod velit. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Praesent fringilla luctus consequat. In hac habitasse platea dictumst. Duis efficitur purus ante, sed pretium nibh egestas eget. Curabitur et viverra orci, venenatis rhoncus sapien. Sed varius mattis elit, sit amet sollicitudin turpis vestibulum a. Ut quam leo, lobortis quis maximus quis, blandit eget mi. Duis nisi massa, dictum ut faucibus eu, mollis ac ipsum. Pellentesque tristique id diam et mollis. Curabitur accumsan ipsum nec turpis laoreet, at tincidunt elit euismod. Vivamus ante erat, porttitor id lacus ac, eleifend efficitur mi. Etiam molestie nisl in mollis porttitor. Maecenas ligula eros, placerat vel facilisis eget, varius in dolor. ",
				Valid:  true,
			},
			EncodedBody: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tincidunt erat vulputate vehicula gravida. Nullam tincidunt vehicula lorem ac ultricies. Proin elit libero, dignissim sed dolor sed, aliquet euismod velit. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Praesent fringilla luctus consequat. In hac habitasse platea dictumst. Duis efficitur purus ante, sed pretium nibh egestas eget. Curabitur et viverra orci, venenatis rhoncus sapien. Sed varius mattis elit, sit amet sollicitudin turpis vestibulum a. Ut quam leo, lobortis quis maximus quis, blandit eget mi. Duis nisi massa, dictum ut faucibus eu, mollis ac ipsum. Pellentesque tristique id diam et mollis. Curabitur accumsan ipsum nec turpis laoreet, at tincidunt elit euismod. Vivamus ante erat, porttitor id lacus ac, eleifend efficitur mi. Etiam molestie nisl in mollis porttitor. Maecenas ligula eros, placerat vel facilisis eget, varius in dolor. ",
		},
	},
}

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
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Get("APP_PORT"))))
}

func AddTodos(c echo.Context) error {
	time.Sleep(1 * time.Second)
	td := models.Todo{
		Model: gorm.Model{
			ID: uint(len(todos["Todos"])),
		},
		Title: c.FormValue("title"),
		Body: pgtype.Text{
			String: c.FormValue("body"),
			Valid:  true,
		},
		EncodedBody: template.HTML(c.FormValue("body_encoded")),
	}

	todos["Todos"] = append(todos["Todos"], td)
	fmt.Println(todos)
	return c.Render(http.StatusOK, "todo_item", td)
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

func Index(c echo.Context) error {

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
