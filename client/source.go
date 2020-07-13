package client

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	"github.com/wolf00/leads_lms/constants"
	source "github.com/wolf00/user_lms/proto/source"
)

type sourceKey struct{}

// SourceFromContext retrieves the client from the Context
func SourceFromContext(ctx context.Context) (source.SourceService, bool) {
	c, ok := ctx.Value(sourceKey{}).(source.SourceService)
	return c, ok
}

// SourceWrapper returns a wrapper for the HeimdallClient
func SourceWrapper(service micro.Service) server.HandlerWrapper {
	client := source.NewSourceService(constants.UserService, service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, sourceKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
