package client

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	"github.com/wolf00/golang_lms/leads/constants"
	organization "github.com/wolf00/golang_lms/user/proto/organization"
)

type organizationKey struct{}

// OrganizationFromContext retrieves the client from the Context
func OrganizationFromContext(ctx context.Context) (organization.OrganizationService, bool) {
	c, ok := ctx.Value(organizationKey{}).(organization.OrganizationService)
	return c, ok
}

// OrganizationWrapper returns a wrapper for the HeimdallClient
func OrganizationWrapper(service micro.Service) server.HandlerWrapper {
	client := organization.NewOrganizationService(constants.UserService, service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, organizationKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
