package main

import (
	"fmt"
	"net/http"

	"github.com/apsystole/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/muhwyndhamhp/todo-mx/config"
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

	e.Static("/dist", "dist")
	e.Static("/assets", "public/assets")
	public.NewTemplateHandler(e)
	public.NewFrontendHandler(e)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Get(config.APP_PORT))))
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
