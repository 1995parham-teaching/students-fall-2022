package main

import (
	"log"

	"github.com/1995parham-teaching/students/internal/graph"
	"github.com/1995parham-teaching/students/internal/graph/resolver"
	"github.com/1995parham-teaching/students/internal/handler"
	"github.com/1995parham-teaching/students/internal/store/course"
	"github.com/1995parham-teaching/students/internal/store/student"
	gHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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

	// start debug mode.
	db = db.Debug()

	ss := student.NewSQL(db)

	{
		h := handler.Student{
			Store: ss,
		}

		h.Register(app.Group("/v1"))
	}

	sc := course.NewSQL(db)

	{
		h := handler.Course{
			Store: sc,
		}

		h.Register(app.Group("/v1"))
	}

	{
		srv := gHandler.NewDefaultServer(graph.NewExecutableSchema(resolver.New(ss)))
		g := app.Group("/v2")

		g.POST("/query", echo.WrapHandler(srv))
		g.GET("/graphiql", echo.WrapHandler(playground.Handler("students-fall-2022", "/v2/query")))
	}

	if err := app.Start("127.0.0.1:1373"); err != nil {
		log.Fatal(err)
	}
}
