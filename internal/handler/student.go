package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/1995parham-teaching/students/internal/model"
	"github.com/1995parham-teaching/students/internal/request"
	"github.com/1995parham-teaching/students/internal/store/student"
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

	st := model.Student{
		Name:    req.Name,
		ID:      fmt.Sprintf("%08d", rand.Int63n(1_000_000_000)),
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
