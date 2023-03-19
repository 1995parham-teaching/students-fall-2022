package student

import (
	"errors"

	"github.com/1995parham-teaching/students/internal/model"
)

var (
	ErrStudentAlreadyExists = errors.New("student already exists")
	ErrStudentNotFound      = errors.New("student does not exist")
)

type Student interface {
	GetAll() ([]model.Student, error)
	Create(model.Student) error
	Get(string) (model.Student, error)
	Register(string, string) error
}
