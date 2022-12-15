package course

import (
	"errors"

	"github.com/1995parham-teaching/students/internal/model"
)

var (
	ErrCourseAlreadyExists = errors.New("course already exists")
	ErrCourseNotFound      = errors.New("course does not exist")
)

type Course interface {
	GetAll() ([]model.Course, error)
	Create(model.Course) error
	Get(string) (model.Course, error)
}
