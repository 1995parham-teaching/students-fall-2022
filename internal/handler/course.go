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
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
)

const (
	CourseIDLen = 8
	CourseIDMax = 100_000_000
)

type Course struct {
	Store course.Course
}

func (s Course) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req request.CourseCreate

	if err := c.Bind(&req); err != nil {
		log.Println(err)

		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {
		log.Println(err)

		return echo.ErrBadRequest
	}

	idBig, err := rand.Int(rand.Reader, big.NewInt(CourseIDLen))
	if err != nil {
		panic(err)
	}

	cr := model.Course{
		Name: req.Name,
		ID:   fmt.Sprintf("%08d", idBig.Int64()),
	}

	if err := s.Store.Create(ctx, cr); err != nil {
		if errors.Is(err, course.ErrCourseAlreadyExists) {
			return echo.ErrBadRequest
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, cr)
}

func (s Course) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	ss, err := s.Store.GetAll(ctx)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, ss)
}

func (s Course) Get(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	if err := validation.Validate(id, validation.Length(CourseIDLen, CourseIDLen), is.Digit); err != nil {
		return echo.ErrBadRequest
	}

	st, err := s.Store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, course.ErrCourseNotFound) {
			return echo.ErrNotFound
		}

		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, st)
}

func (s Course) Register(g *echo.Group) {
	g.POST("/courses", s.Create)
	g.GET("/courses", s.GetAll)
	g.GET("/courses/:id", s.Get)
}
