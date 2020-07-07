package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wolf00/golang_lms/leads/constants"
	"github.com/wolf00/golang_lms/leads/db"
	"github.com/wolf00/golang_lms/leads/db/models"

	lead_template_models "github.com/wolf00/golang_lms/lead_template/db/models"

	leads "github.com/wolf00/golang_lms/leads/proto/leads"
	"github.com/wolf00/golang_lms/leads/utilities"

	"github.com/micro/go-micro/v2/util/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewLeadService Handler
type NewLeadService struct {
}

// NewLead create
func (e *NewLeadService) NewLead(ctx context.Context, req *leads.NewLeadRequest, rsp *leads.NewLeadResponse) error {
	validationErrors := e.validateLead(req)
	if len(validationErrors) > 0 {
		e.failure(strings.Join(validationErrors, ", "), rsp)
		return nil
	}
	newLead, err := e.lead(ctx, req)
	if err != nil {
		return err
	}
	existingLead, err := leadByFilter(ctx, bson.M{"contact": newLead.Contact, "email": newLead.Email})
	if err == nil {
		return e.updateExistingLead(ctx, existingLead, newLead, rsp)
	}
	existingLeadWithContact, err := leadByFilter(ctx, bson.M{"contact": newLead.Contact})
	if err == nil {
		newLead.Meta[0].Meta = append(newLead.Meta[0].Meta, models.KeyValue{Key: "contact", Value: newLead.Contact})
		return e.updateExistingLead(ctx, existingLeadWithContact, newLead, rsp)
	}
	existingLeadWithEmail, err := leadByFilter(ctx, bson.M{"email": newLead.Email})
	if err == nil {
		newLead.Meta[0].Meta = append(newLead.Meta[0].Meta, models.KeyValue{Key: "email", Value: newLead.Email})
		return e.updateExistingLead(ctx, existingLeadWithEmail, newLead, rsp)
	}
	// Create a new lead if not existing
	_, err = createLead(ctx, newLead)
	if err != nil {
		log.Error(err)
		return e.failure(err.Error(), rsp)
	}
	rsp.Message = constants.LeadAdded
	rsp.Status = true
	return nil
}

func (e *NewLeadService) updateExistingLead(ctx context.Context, existingLead models.Leads, newLead models.Leads, rsp *leads.NewLeadResponse) error {
	existingLead.LeadSource = append(existingLead.LeadSource, newLead.LeadSource...)
	existingLead.Meta = append(existingLead.Meta, newLead.Meta...)
	existingLead.TemplateValues = append(existingLead.TemplateValues, newLead.TemplateValues...)
	updateData := bson.D{
		{"$set", bson.D{
			{"leadSource", existingLead.LeadSource},
			{"meta", existingLead.Meta},
			{"templateValues", existingLead.TemplateValues},
		}},
	}
	_, err := updateLead(ctx, existingLead.ID, updateData)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("failed to add new lead. details: %s", err.Error())
	}

	rsp.Message = constants.LeadAdded
	rsp.Status = true
	return nil
}

func (e *NewLeadService) validateLead(req *leads.NewLeadRequest) []string {
	validationFailures := []string{}
	if !utilities.ValidateString(req.FirstName) {
		validationFailures = append(validationFailures, constants.EmptyFirstName)
	}
	if !utilities.ValidateString(req.Contact) {
		validationFailures = append(validationFailures, constants.EmptyContact)
	}
	if !utilities.ValidateString(req.Source) {
		validationFailures = append(validationFailures, constants.EmptySource)
	}
	if !utilities.ValidateString(req.CampaignTag) {
		validationFailures = append(validationFailures, constants.EmptyCampaignTag)
	}
	if !utilities.ValidateString(req.Email) {
		validationFailures = append(validationFailures, constants.EmptyEmail)
	} else {
		if !utilities.ValidateEmail(req.Email) {
			validationFailures = append(validationFailures, constants.InvalidEmailFormat)
		}
	}

	return validationFailures
}

func (e *NewLeadService) failure(message string, rsp *leads.NewLeadResponse) error {
	rsp.Message = message
	rsp.Status = false
	return nil
}

func (e *NewLeadService) validateSourceAccess(ctx context.Context, campaignCreator primitive.ObjectID, sourceTag string) error {
	// TO_DO: Validate if a lead creator has access to lead source
	user, err := userByID(ctx, campaignCreator.Hex())
	if err != nil {
		return err
	}
	organization, err := organizationByID(ctx, user.OrganizationID.Hex())
	if err != nil {
		return err
	}
	tagFound := false
	for i := 0; i < len(organization.AllowedSources); i++ {
		if organization.AllowedSources[i] == sourceTag {
			tagFound = true
		}
	}
	if !tagFound {
		return fmt.Errorf("the '%s' source is not allowed to submit the lead", sourceTag)
	}
	return nil
}

func (e *NewLeadService) lead(ctx context.Context, req *leads.NewLeadRequest) (models.Leads, error) {
	newLead := models.Leads{}
	newLead.Contact = req.Contact
	newLead.Email = req.Email
	newLead.FirstName = req.FirstName
	newLead.LastName = req.LastName

	if e.preExistingLead(ctx, req) {
		return newLead, fmt.Errorf("lead already captured")
	}

	// Validate campaign tag
	campaign, err := campaignByTag(ctx, req.CampaignTag)
	if err != nil {
		return newLead, err
	}
	newLead.CampaignID = campaign.ID
	err = e.validateSourceAccess(ctx, campaign.CreatedBy, req.Source)
	if err != nil {
		return newLead, err
	}
	// Template from campaign
	templateID := campaign.TemplateID.Hex()
	leadTemplate, err := leadTemplateByID(ctx, templateID)
	if err != nil {
		log.Error(err)
		return newLead, fmt.Errorf("lead template with the templateId %s is not available", templateID)
	}

	// TO_DO: validate lead template value types
	// Validating lead template keys
	invalidTemplateKeys := []string{}
	leadTemplateValues := models.SourceWithKeyValue{LeadSource: req.Source, Meta: []models.KeyValue{}}
	for mi := 0; mi < len(req.TemplateValues); mi++ {
		templateKey := req.TemplateValues[mi].Key
		if !e.isTemplateKey(leadTemplate, templateKey) {
			invalidTemplateKeys = append(invalidTemplateKeys, templateKey)
		}
		leadTemplateValues.Meta = append(leadTemplateValues.Meta, models.KeyValue{Key: templateKey, Value: req.TemplateValues[mi].Value})
	}
	if len(invalidTemplateKeys) > 0 {
		return newLead, fmt.Errorf("the following fields should not be part of the lead template: %s", strings.Join(invalidTemplateKeys, ","))
	}
	newLead.TemplateValues = append([]models.SourceWithKeyValue{}, leadTemplateValues)
	newLead.TemplateID = campaign.TemplateID

	newLead.LeadSource = append([]string{}, req.Source)
	leadMeta := models.SourceWithKeyValue{LeadSource: req.Source, Meta: []models.KeyValue{}}
	for mi := 0; mi < len(req.Meta); mi++ {
		leadMeta.Meta = append(leadMeta.Meta, models.KeyValue{Key: req.Meta[mi].Key, Value: req.Meta[mi].Value})
	}
	currentTime := time.Now()
	leadMeta.Meta = append(leadMeta.Meta, models.KeyValue{Key: "createdOn", Value: currentTime.Format(utilities.TIME_FORMAT)})
	newLead.Meta = append([]models.SourceWithKeyValue{}, leadMeta)

	newLead.CreatedOn = time.Now()
	return newLead, nil
}

func (e *NewLeadService) isTemplateKey(leadTemplate lead_template_models.LeadTemplate, key string) bool {
	if len(leadTemplate.KeyValueTypes) == 0 {
		return false
	}
	for ti := 0; ti < len(leadTemplate.KeyValueTypes); ti++ {
		if key == leadTemplate.KeyValueTypes[ti].Key {
			return true
		}
	}

	return false
}

func (e *NewLeadService) preExistingLead(ctx context.Context, req *leads.NewLeadRequest) bool {
	filter := bson.M{
		"contact":    req.Contact,
		"email":      req.Email,
		"leadSource": req.Source,
	}
	_, err := leadByFilter(ctx, filter)
	if err != nil {
		return false
	}
	return true
}

func leadByFilter(ctx context.Context, filter interface{}) (models.Leads, error) {
	helper := db.Leads(ctx)

	var existingLead models.Leads

	err := helper.FindOne(ctx, filter).Decode(&existingLead)
	if err != nil {
		log.Error(err)
		return existingLead, err
	}

	return existingLead, nil
}

func createLead(ctx context.Context, newlead models.Leads) (*mongo.InsertOneResult, error) {
	helper := db.Leads(ctx)
	return helper.InsertOne(ctx, newlead)
}

func updateLead(ctx context.Context, leadID primitive.ObjectID, updateFieldAndValues interface{}) (*mongo.UpdateResult, error) {
	helper := db.Leads(ctx)
	return helper.UpdateOne(ctx, bson.M{"_id": leadID}, updateFieldAndValues)
}
