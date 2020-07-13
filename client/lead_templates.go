package client

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	lead_template "github.com/wolf00/lead_template_lms/proto/lead_template"
	"github.com/wolf00/leads_lms/constants"
)

type leadTemplateKey struct{}

// FromContext retrieves the client from the Context
func LeadTemplateFromContext(ctx context.Context) (lead_template.LeadTemplateService, bool) {
	c, ok := ctx.Value(leadTemplateKey{}).(lead_template.LeadTemplateService)
	return c, ok
}

// Client returns a wrapper for the HeimdallClient
func LeadTemplateWrapper(service micro.Service) server.HandlerWrapper {
	client := lead_template.NewLeadTemplateService(constants.LeadTemplateService, service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, leadTemplateKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
