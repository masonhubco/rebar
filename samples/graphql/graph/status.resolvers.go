package graph

import (
	"context"
	"time"

	"github.com/masonhubco/rebar/samples/graphql/graph/gplmodels"
)

func (r *mutationResolver) CreateStatus(ctx context.Context, input gplmodels.NewStatus) (*gplmodels.Status, error) {
	status := &gplmodels.Status{
		State:  input.State,
		Redis:  input.Redis,
		Uptime: input.Uptime,
	}
	return status, nil
}

func (r *queryResolver) Status(ctx context.Context) (*gplmodels.Status, error) {
	status := &gplmodels.Status{
		State:  "UP",
		Redis:  "Connected",
		Uptime: time.Now().Format(time.RFC3339),
	}
	return status, nil
}
