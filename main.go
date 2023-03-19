package main

import (
	"log"

	"github.com/1995parham-teaching/students/internal/handler"
	"github.com/1995parham-teaching/students/internal/model"
	"github.com/1995parham-teaching/students/internal/store/course"
	"github.com/1995parham-teaching/students/internal/store/student"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := echo.New()

	db, err := gorm.Open(sqlite.Open("students.db"), new(gorm.Config))
	if err != nil {
		log.Fatal(err)
	}

	s := student.NewSQL(db)
	course.NewSQL(db)

	h := handler.Student{
		Store: s,
	}

	h.Register(app.Group("/v1"))

	if err := app.Start("127.0.0.1:1373"); err != nil {
		log.Fatal(err)
	}
}
