package main

import (
	"github.com/wolf00/golang_lms/leads/handler"
	leads "github.com/wolf00/golang_lms/leads/proto/leads"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/wolf00/golang_lms/leads/client"
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
	)

	// Register Handler
	leads.RegisterLeadsHandler(service.Server(), new(handler.LeadsRequestHandler))

	// Register Struct as Subscriber
	// micro.RegisterSubscriber("go.micro.service.leads", service.Server(), new(subscriber.Leads))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
