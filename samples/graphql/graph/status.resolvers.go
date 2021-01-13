package graph

import (
	"context"
	"time"

	"github.com/masonhubco/rebar/samples/graphql/graph/gplmodels"
	"github.com/masonhubco/rebar/samples/graphql/model"
)

func (r *mutationResolver) CreateStatus(ctx context.Context, input gplmodels.NewStatus) (*model.Status, error) {
	status := &model.Status{
		Status: input.State,
		Redis:  input.Redis,
		Uptime: input.Uptime,
	}
	return status, nil
}

func (r *queryResolver) Status(ctx context.Context) (*model.Status, error) {
	status := &model.Status{
		Status: "UP",
		Redis:  "Connected",
		Uptime: time.Now().Format(time.RFC3339),
	}
	return status, nil
}
