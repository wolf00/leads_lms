package publisher

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"
	leads "github.com/wolf00/leads_lms/proto/leads"
)

// PublishNewLead asda
func PublishNewLead(topic string, p micro.Publisher, leadStatusID string) {
	newLead := &leads.NewLeadEvent{
		LeadStatusId: leadStatusID,
	}

	log.Info("publishing event...")

	if err := p.Publish(context.Background(), newLead); err != nil {
		log.Error("error publishing", err)
	}
}
