package resolver

import (
	"github.com/1995parham-teaching/students/internal/graph"
	"github.com/1995parham-teaching/students/internal/store/student"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Store student.Student
}

func NewResolver(store student.Student) *Resolver {
	return &Resolver{
		Store: store,
	}
}

func New(store student.Student) graph.Config {
	// nolint: exhaustruct
	c := graph.Config{
		Schema:     nil,
		Resolvers:  NewResolver(store),
		Directives: graph.DirectiveRoot{},
	}

	return c
}
