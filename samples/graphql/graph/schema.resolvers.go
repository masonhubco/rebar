package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	generated1 "github.com/masonhubco/rebar/samples/graphql/graph/generated"
	model1 "github.com/masonhubco/rebar/samples/graphql/graph/model"
	"github.com/masonhubco/rebar/samples/graphql/model"
)

func (r *mutationResolver) CreateStatus(ctx context.Context, input model1.NewStatus) (*model.Status, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *statusResolver) State(ctx context.Context, obj *model.Status) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated1.MutationResolver implementation.
func (r *Resolver) Mutation() generated1.MutationResolver { return &mutationResolver{r} }

// Status returns generated1.StatusResolver implementation.
func (r *Resolver) Status() generated1.StatusResolver { return &statusResolver{r} }

type mutationResolver struct{ *Resolver }
type statusResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
