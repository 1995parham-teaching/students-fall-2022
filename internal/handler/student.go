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
	var req request.StudentCreate

	if err := c.Bind(&req); err != nil {
		log.Println(err)

		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {
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

	if err := s.Store.Create(st); err != nil {
		if errors.Is(err, student.ErrStudentAlreadyExists) {
			return echo.ErrBadRequest
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, st)
}

func (s Student) GetAll(c echo.Context) error {
	ss, err := s.Store.GetAll()
	if err != nil {
		return echo.ErrInternalServerError
	}

	h := c.Request().Header.Get("STUDENTS-FALL-2022")
	log.Printf("STUDENTS-FALL-2022: %s\n", h)

	c.Response().Header().Add("STUDENTS-FALL-2022", "123")

	return c.JSON(http.StatusOK, ss)
}

func (s Student) Get(c echo.Context) error {
	id := c.Param("id")

	if err := validation.Validate(id, validation.Length(StudentIDLen, StudentIDLen), is.Digit); err != nil {
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

func (s Student) Fill(c echo.Context) error {
	sid := c.Param("sid")
	cid := c.Param("cid")

	if err := validation.Validate(sid, validation.Length(StudentIDLen, StudentIDLen), is.Digit); err != nil {
		return echo.ErrBadRequest
	}

	if err := validation.Validate(cid, validation.Length(CourseIDLen, CourseIDLen), is.Digit); err != nil {
		return echo.ErrBadRequest
	}

	if err := s.Store.Register(sid, cid); err != nil {
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
