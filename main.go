package main

import (
	"github.com/wolf00/leads_lms/handler"
	leads "github.com/wolf00/leads_lms/proto/leads"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	// "github.com/micro/go-micro/v2/server"
	"github.com/wolf00/leads_lms/client"
	"github.com/wolf00/leads_lms/subscriber"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.leads"),
		micro.Version("0.1"),
	)

	// Initialise service
	service.Init(
		micro.WrapHandler(client.CampaignWrapper(service)),
		micro.WrapHandler(client.LeadTemplateWrapper(service)),
		micro.WrapHandler(client.OrganizationWrapper(service)),
		micro.WrapHandler(client.SourceWrapper(service)),
		micro.WrapHandler(client.UserWrapper(service)),
		micro.WrapHandler(client.NewLeadPublisherWrapper(service)),
	)

	// Register Handler
	leads.RegisterLeadsHandler(service.Server(), new(handler.LeadsRequestHandler))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.leads.NewLead", service.Server(), new(subscriber.NewLeadHandler))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
