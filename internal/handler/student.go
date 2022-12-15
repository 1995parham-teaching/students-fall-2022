package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/1995parham-teaching/students/internal/model"
	"github.com/1995parham-teaching/students/internal/request"
	"github.com/1995parham-teaching/students/internal/store/student"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
)

type Student struct {
	Store student.Student
}

func (s Student) Create(c echo.Context) error {
	var req request.StudentCreate

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {
		return echo.ErrBadRequest
	}

	st := model.Student{
		Name:    req.Name,
		ID:      fmt.Sprintf("%08d", rand.Int63n(100_000_000)),
		Courses: nil,
	}

	if err := s.Store.Create(st); err != nil {
		if errors.Is(err, student.ErrStudentAlreadyExists) {
			return echo.ErrBadRequest
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, st)
}

func (s Student) GetAll(c echo.Context) error {
	ss, err := s.Store.GetAll()
	if err != nil {
		return echo.ErrInternalServerError
	}

	h := c.Request().Header.Get("STUDENTS-FALL-2022")
	fmt.Println(h)

	c.Response().Header().Add("STUDENTS-FALL-2022", "123")

	return c.JSON(http.StatusOK, ss)
}

func (s Student) Get(c echo.Context) error {
	id := c.Param("id")

	if err := validation.Validate(id, validation.Length(8, 8), is.Digit); err != nil {
		return echo.ErrBadRequest
	}

	st, err := s.Store.Get(id)
	if err != nil {
		if errors.Is(err, student.ErrStudentNotFound) {
			return echo.ErrNotFound
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, st)
}

func (s Student) Register(g *echo.Group) {
	g.POST("/students", s.Create)
	g.GET("/students", s.GetAll)
	g.GET("/students/:id", s.Get)
}
