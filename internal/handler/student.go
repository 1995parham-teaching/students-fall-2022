package handler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/1995parham-teaching/students/internal/model"
	"github.com/1995parham-teaching/students/internal/request"
	"github.com/1995parham-teaching/students/internal/store/course"
	"github.com/1995parham-teaching/students/internal/store/student"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
)

const (
	StudentIDLen = 8
	StudentIDMax = 100_000_000
)

type Student struct {
	Store student.Student
}

func (s Student) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req request.StudentCreate

	err := c.Bind(&req)
	if err != nil {
		log.Println(err)

		return echo.ErrBadRequest
	}

	err = req.Validate()
	if err != nil {
		log.Println(err)

		return echo.ErrBadRequest
	}

	idBig, err := rand.Int(rand.Reader, big.NewInt(StudentIDMax))
	if err != nil {
		panic(err)
	}

	st := model.Student{
		Name:    req.Name,
		ID:      fmt.Sprintf("%08d", idBig.Int64()),
		Courses: nil,
	}

	err = s.Store.Create(ctx, st)
	if err != nil {
		if errors.Is(err, student.ErrStudentAlreadyExists) {
			return echo.ErrBadRequest
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, st)
}

func (s Student) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	ss, err := s.Store.GetAll(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}

	h := c.Request().Header.Get("Students-Fall-2022")
	log.Printf("Students-Fall-2022: %s\n", h)

	c.Response().Header().Add("Students-Fall-2022", "123")

	return c.JSON(http.StatusOK, ss)
}

func (s Student) Get(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	err := validation.Validate(id, validation.Length(StudentIDLen, StudentIDLen), is.Digit)
	if err != nil {
		return echo.ErrBadRequest
	}

	st, err := s.Store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, student.ErrStudentNotFound) {
			return echo.ErrNotFound
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, st)
}

func (s Student) Fill(c echo.Context) error {
	ctx := c.Request().Context()

	sid := c.Param("sid")
	cid := c.Param("cid")

	err := validation.Validate(sid, validation.Length(StudentIDLen, StudentIDLen), is.Digit)
	if err != nil {
		return echo.ErrBadRequest
	}

	err = validation.Validate(cid, validation.Length(CourseIDLen, CourseIDLen), is.Digit)
	if err != nil {
		return echo.ErrBadRequest
	}

	err = s.Store.Register(ctx, sid, cid)
	if err != nil {
		if errors.Is(err, student.ErrStudentNotFound) {
			return echo.ErrNotFound
		}

		if errors.Is(err, course.ErrCourseNotFound) {
			return echo.ErrNotFound
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, nil)
}

func (s Student) Register(g *echo.Group) {
	g.POST("/students", s.Create)
	g.GET("/students", s.GetAll)
	g.GET("/students/:id", s.Get)
	g.GET("/students/:sid/register/:cid", s.Fill)
}
