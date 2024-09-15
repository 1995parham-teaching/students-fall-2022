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
	Create(student model.Student) error
	Get(id string) (model.Student, error)
	Register(sid string, cid string) error
}
