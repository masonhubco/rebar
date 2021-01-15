package graph

import (
	"context"
	"time"

	"github.com/masonhubco/rebar/samples/graphql/graph/gqlmodels"
)

func (r *mutationResolver) CreateStatus(ctx context.Context, input gqlmodels.NewStatus) (*gqlmodels.Status, error) {
	status := &gqlmodels.Status{
		State:  input.State,
		Redis:  input.Redis,
		Uptime: input.Uptime,
	}
	return status, nil
}

func (r *queryResolver) Status(ctx context.Context) (*gqlmodels.Status, error) {
	status := &gqlmodels.Status{
		State:  "UP",
		Redis:  "Connected",
		Uptime: time.Now().Format(time.RFC3339),
	}
	return status, nil
}
