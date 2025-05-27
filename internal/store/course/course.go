package course

import (
	"context"
	"errors"

	"github.com/1995parham-teaching/students/internal/model"
)

var (
	ErrCourseAlreadyExists = errors.New("course already exists")
	ErrCourseNotFound      = errors.New("course does not exist")
)

type Course interface {
	GetAll(ctx context.Context) ([]model.Course, error)
	Create(ctx context.Context, course model.Course) error
	Get(ctx context.Context, id string) (model.Course, error)
}
