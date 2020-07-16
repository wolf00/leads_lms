package models

import (
	"time"

	leads "github.com/wolf00/leads_lms/proto/leads"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Lead status status
const (
	Processing = "Processing"
	Success    = "Success"
	Failed     = "Failed"
)

// LeadStatus type definition
type LeadStatus struct {
	CampaignTag    string                           `bson:"campaignTag"`
	Contact        string                           `bson:"contact"`
	CreatedOn      time.Time                        `bson:"createdOn,omitempty"`
	Email          string                           `bson:"email"`
	FirstName      string                           `bson:"firstName"`
	ID             primitive.ObjectID               `bson:"_id,omitempty"`
	LastName       string                           `bson:"lastname"`
	LeadSource     string                           `bson:"leadSource,omitempty"`
	Meta           []*leads.NewLeadRequest_KeyValue `bson:"meta,omitempty"`
	TemplateValues []*leads.NewLeadRequest_KeyValue `bson:"templateValues,omitempty"`
	Status         string                           `bson:"status,omitempty"`
	Message        string                           `bson:"message,omitempty"`
	Stacktrace     string                           `bson:"stacktrace,omitempty"`
}
