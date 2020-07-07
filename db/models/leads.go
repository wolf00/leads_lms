package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Lead types
const (
	HOT  = iota
	WARM = iota
	COLD = iota
)

// Verification mediums
const (
	SMS   = iota
	EMAIL = iota
	CALL  = iota
)

// Leads type definition
type Leads struct {
	CampaignID         primitive.ObjectID   `bson:"campaignId"`
	Contact            string               `bson:"contact"`
	CreatedOn          time.Time            `bson:"createdOn,omitempty"`
	Email              string               `bson:"email"`
	FirstName          string               `bson:"firstName"`
	ID                 primitive.ObjectID   `bson:"_id,omitempty"`
	LastName           string               `bson:"lastname"`
	LeadScore          int32                `bson:"leadScore,omitempty"`
	LeadSource         []string             `bson:"leadSource,omitempty"`
	LeadType           string               `bson:"leadType,omitempty"`
	Meta               []SourceWithKeyValue `bson:"meta,omitempty"`
	VerificationMedium string               `bson:"verificationMedium,omitempty"`
	TemplateID         primitive.ObjectID   `bson:"templateId"`
	TemplateValues     []SourceWithKeyValue `bson:"templateValues,omitempty"`
}

// SourceWithKeyValue type def
type SourceWithKeyValue struct {
	LeadSource string     `bson:"leadSource,omitempty"`
	Meta       []KeyValue `bson:"meta,omitempty"`
}

// KeyValue type def
type KeyValue struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}
