package subscriber

import (
	"context"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	leads "github.com/wolf00/leads_lms/proto/leads"
)

type NewLeadHandler struct{}

func (e *NewLeadHandler) HandleNewLead(ctx context.Context, newLead *leads.NewLeadEvent) error {
	log.Info("Function Received message: ", newLead.LeadStatusId)
	service := micro.NewService()
	service.Init()

	leadService := leads.NewLeadsService("go.micro.service.leads", service.Client())
	_, err := leadService.NewLeadByLeadStatusID(context.Background(), newLead)

	return err
}
