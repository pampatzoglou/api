package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/pampatzoglou/api/graph/generated"
	"github.com/pampatzoglou/api/graph/model"
)

// Shops is the resolver for the shops field.
func (r *queryResolver) Shops(ctx context.Context) ([]*model.Shop, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
