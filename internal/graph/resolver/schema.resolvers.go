package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.50

import (
	"context"
	"fmt"

	"github.com/1995parham-teaching/students/internal/graph"
	"github.com/1995parham-teaching/students/internal/graph/model"
)

// StudentsByName is the resolver for the studentsByName field.
func (r *queryResolver) StudentsByName(ctx context.Context, name *string) ([]*model.Student, error) {
	panic(fmt.Errorf("not implemented: StudentsByName - studentsByName"))
}

// StudentByID is the resolver for the studentByID field.
func (r *queryResolver) StudentByID(ctx context.Context, id string) (*model.Student, error) {
	panic(fmt.Errorf("not implemented: StudentByID - studentByID"))
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
