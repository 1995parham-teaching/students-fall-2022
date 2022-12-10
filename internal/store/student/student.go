package student

import "github.com/1995parham-teaching/students/internal/model"

type Student interface {
	GetAll() ([]model.Student, error)
	Create(model.Student) error
	Get(id string) (model.Student, error)
}
