package subscriber

import (
	"context"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	// "github.com/micro/go-plugins/registry/consul/v2"
	leads "github.com/wolf00/leads_lms/proto/leads"
)

type NewLeadHandler struct{}

func (e *NewLeadHandler) HandleNewLead(ctx context.Context, newLead *leads.NewLeadEvent) error {
	log.Info("Function Received message: ", newLead.LeadStatusId)
	// Consul registry
	// registry := consul.NewRegistry()

	service := micro.NewService(
	// micro.Registry(registry),
	)
	service.Init()

	leadService := leads.NewLeadsService("go.micro.service.leads", service.Client())
	_, err := leadService.NewLeadByLeadStatusID(context.Background(), newLead)

	return err
}
