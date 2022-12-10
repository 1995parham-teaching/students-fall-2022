package main

import (
	"github.com/1995parham-teaching/students/internal/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	app.GET("/hello", handler.Hello)

	app.Start("127.0.0.1:1373")
}
