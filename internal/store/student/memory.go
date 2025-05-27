package student

import (
	"context"

	"github.com/1995parham-teaching/students/internal/model"
)

type inMemoryItem struct {
	Name    string
	Courses []string
}

type InMemory struct {
	students map[string]inMemoryItem
}

func NewInMemory() Student {
	return &InMemory{
		students: make(map[string]inMemoryItem),
	}
}

func (im *InMemory) GetAll(_ context.Context) ([]model.Student, error) {
	students := make([]model.Student, 0, len(im.students))

	for id, i := range im.students {
		students = append(students, model.Student{
			Name:    i.Name,
			ID:      id,
			Courses: nil,
		})
	}

	return students, nil
}

func (im *InMemory) Create(_ context.Context, s model.Student) error {
	if _, ok := im.students[s.ID]; ok {
		return ErrStudentAlreadyExists
	}

	courses := make([]string, 0)

	for _, c := range s.Courses {
		courses = append(courses, c.ID)
	}

	im.students[s.ID] = inMemoryItem{
		Name:    s.Name,
		Courses: courses,
	}

	return nil
}

func (im *InMemory) Register(_ context.Context, _ string, _ string) error {
	return nil
}

func (im *InMemory) Get(_ context.Context, id string) (model.Student, error) {
	s, ok := im.students[id]
	if !ok {
		return model.Student{}, ErrStudentNotFound
	}

	return model.Student{
		Name:    s.Name,
		ID:      id,
		Courses: nil,
	}, nil
}
