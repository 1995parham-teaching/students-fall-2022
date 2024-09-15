package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.50

import (
	"context"
	"fmt"
	rand "math/rand/v2"

	"github.com/1995parham-teaching/students/internal/graph"
	"github.com/1995parham-teaching/students/internal/graph/model"
	imodel "github.com/1995parham-teaching/students/internal/model"
	"github.com/1995parham-teaching/students/internal/request"
)

// CreateStudent is the resolver for the createStudent field.
func (r *mutationResolver) CreateStudent(ctx context.Context, name string) (*model.Student, error) {
	req := request.StudentCreate{
		Name: name,
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	st := imodel.Student{
		Name:    req.Name,
		ID:      fmt.Sprintf("%08d", rand.Int64()%StudentIDMax),
		Courses: nil,
	}

	if err := r.Store.Create(st); err != nil {
	}

	return &model.Student{
		ID:      st.ID,
		Name:    st.Name,
		Courses: nil,
	}, nil
}

// StudentsByName is the resolver for the studentsByName field.
func (r *queryResolver) StudentsByName(ctx context.Context, name *string) ([]*model.Student, error) {
	return nil, nil
}

// StudentByID is the resolver for the studentByID field.
func (r *queryResolver) StudentByID(ctx context.Context, id string) (*model.Student, error) {
	s, err := r.Store.Get(id)
	if err != nil {
		return nil, err
	}

	courses := make([]*model.Course, 0)
	for _, c := range s.Courses {
		courses = append(courses, &model.Course{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	return &model.Student{
		ID:      s.ID,
		Name:    s.Name,
		Courses: courses,
	}, nil
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }