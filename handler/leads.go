package handler

import (
	"context"

	leads "github.com/wolf00/golang_lms/leads/proto/leads"
	"github.com/wolf00/golang_lms/leads/services"
)

// LeadsRequestHandler type def
type LeadsRequestHandler struct {
	services.NewLeadService
}

// NewLead is a single request handler called via client.NewLead or the generated client code
func (e *LeadsRequestHandler) NewLead(ctx context.Context, req *leads.NewLeadRequest, rsp *leads.NewLeadResponse) error {
	return e.NewLeadService.NewLead(ctx, req, rsp)
}
