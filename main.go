package main

import (
	"github.com/1995parham-teaching/students/internal/handler"
	"github.com/1995parham-teaching/students/internal/store/student"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	s := student.NewInMemory()

	h := handler.Student{
		Store: s,
	}

	app.POST("/student", h.Create)

	app.Start("127.0.0.1:1373")
}
