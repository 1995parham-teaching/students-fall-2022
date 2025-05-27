package student

import (
	"context"
	"errors"

	"github.com/1995parham-teaching/students/internal/model"
)

var (
	ErrStudentAlreadyExists = errors.New("student already exists")
	ErrStudentNotFound      = errors.New("student does not exist")
)

type Student interface {
	GetAll(ctx context.Context) ([]model.Student, error)
	Create(ctx context.Context, student model.Student) error
	Get(ctx context.Context, id string) (model.Student, error)
	Register(ctx context.Context, sid string, cid string) error
}
