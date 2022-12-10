package student

import "github.com/1995parham-teaching/students/internal/model"

type InMemory struct {
	students map[string]model.Student
}

func (im InMemory) GetAll() ([]model.Student, error) {
	students := make([]model.Student, 0, len(im.students))

	for _, s := range im.students {
		students = append(students, s)
	}

	return students, nil
}
func (im InMemory) Create(model.Student) error           {}
func (im InMemory) Get(id string) (model.Student, error) {}
