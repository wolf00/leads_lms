package handler

import (
	"context"
	"fmt"

	"github.com/wolf00/leads_lms/client"
	leads "github.com/wolf00/leads_lms/proto/leads"
	"github.com/wolf00/leads_lms/publisher"
	"github.com/wolf00/leads_lms/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LeadsRequestHandler type def
type LeadsRequestHandler struct {
	services.LeadService
}

// NewLead is a single request handler called via client.NewLead or the generated client code
func (e *LeadsRequestHandler) NewLead(ctx context.Context, req *leads.NewLeadRequest, rsp *leads.NewLeadResponse) error {
	// newLeadPublisher, ok := client.NewLeadPublisherFromContext(ctx)
	// if !ok {
	// 	return fmt.Errorf("failed to create publisher")
	// }
	// leadStatusID, err := e.NewLeadStatus(ctx, req)
	// if err != nil {
	// 	return err
	// }
	// go publisher.PublishNewLead("go.micro.service.leads.NewLead", newLeadPublisher, leadStatusID.InsertedID.(primitive.ObjectID).Hex())
	err := e.LeadService.CreateNewLead(ctx, req)
	if err != nil {
		return err
	}

	rsp.Message = "lead creation in progress"
	rsp.Status = true
	return nil
}

func (e *LeadsRequestHandler) NewLeadQueued(ctx context.Context, req *leads.NewLeadRequest, rsp *leads.NewLeadResponse) error {
	newLeadPublisher, ok := client.NewLeadPublisherFromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to create publisher")
	}
	leadStatusID, err := e.NewLeadStatus(ctx, req)
	if err != nil {
		return err
	}
	go publisher.PublishNewLead("go.micro.service.leads.NewLead", newLeadPublisher, leadStatusID.InsertedID.(primitive.ObjectID).Hex())

	rsp.Message = "lead creation in progress"
	rsp.Status = true
	return nil
}

func (e *LeadsRequestHandler) NewLeadByLeadStatusID(ctx context.Context, req *leads.NewLeadEvent, rsp *leads.NewLeadResponse) error {
	leadStatus, err := e.LeadService.LeadStatusFromID(ctx, req.LeadStatusId)
	if err != nil {
		return err
	}

	err = e.LeadService.NewLead(ctx, leadStatus)
	if err != nil {
		return err
	}

	rsp.Message = "lead creation in progress"
	rsp.Status = true
	return nil
}
