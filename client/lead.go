package client

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
)

type leadKey struct{}

// NewLeadPublisherFromContext retrieves the client from the Context
func NewLeadPublisherFromContext(ctx context.Context) (micro.Event, bool) {
	e, ok := ctx.Value(leadKey{}).(micro.Event)
	return e, ok
}

// NewLeadPublisherWrapper returns a wrapper for the HeimdallClient
func NewLeadPublisherWrapper(service micro.Service) server.HandlerWrapper {
	publisher := micro.NewPublisher("go.micro.service.leads.NewLead", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, leadKey{}, publisher)
			return fn(ctx, req, rsp)
		}
	}
}
