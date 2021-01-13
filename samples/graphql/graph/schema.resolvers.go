package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/masonhubco/rebar/samples/graphql/graph/generated"
	"github.com/masonhubco/rebar/samples/graphql/graph/model"
	model1 "github.com/masonhubco/rebar/samples/graphql/model"
)

func (r *mutationResolver) CreateStatus(ctx context.Context, input model.NewStatus) (*model1.Status, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *statusResolver) State(ctx context.Context, obj *model1.Status) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Status returns generated.StatusResolver implementation.
func (r *Resolver) Status() generated.StatusResolver { return &statusResolver{r} }

type mutationResolver struct{ *Resolver }
type statusResolver struct{ *Resolver }
