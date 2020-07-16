package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wolf00/leads_lms/constants"
	"github.com/wolf00/leads_lms/db"
	"github.com/wolf00/leads_lms/db/models"

	lead_template_models "github.com/wolf00/lead_template_lms/db/models"

	"github.com/micro/go-micro/v2/util/log"
	leads "github.com/wolf00/leads_lms/proto/leads"
	"github.com/wolf00/leads_lms/utilities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// LeadService Handler
type LeadService struct {
}

func (e *LeadService) leadRequest2LeadStatus(req *leads.NewLeadRequest) models.LeadStatus {
	var leadStatus models.LeadStatus
	leadStatus.LastName = req.LastName
	leadStatus.LeadSource = req.Source
	leadStatus.CampaignTag = req.CampaignTag
	leadStatus.Contact = req.Contact
	leadStatus.Email = req.Email
	leadStatus.FirstName = req.FirstName
	leadStatus.LastName = req.LastName
	leadStatus.Meta = req.Meta
	leadStatus.TemplateValues = req.TemplateValues
	return leadStatus
}

func leadStatusCaptureError(ctx context.Context, leadStatusID primitive.ObjectID, err error) error {
	return leadStatusCapture(ctx, leadStatusID, models.Failed, "failed to add the lead", err)
}

func leadStatusCaptureSuccess(ctx context.Context, leadStatusID primitive.ObjectID) error {
	return leadStatusCapture(ctx, leadStatusID, models.Failed, "lead camptured successfully", nil)
}

func leadStatusCapture(ctx context.Context, leadStatusID primitive.ObjectID, status string, message string, stacktrace error) error {
	var errMessage string
	if stacktrace != nil {
		errMessage = stacktrace.Error()
	}
	updateData := bson.D{
		{"$set", bson.D{
			{"stacktrace", errMessage},
			{"message", message},
			{"status", status},
		}},
	}
	_, dberr := updateLeadStatus(ctx, leadStatusID, updateData)
	if dberr != nil {
		log.Error(dberr)
	}

	return stacktrace
}

func (e *LeadService) NewLeadStatus(ctx context.Context, req *leads.NewLeadRequest) (*mongo.InsertOneResult, error) {
	validationErrors := e.validateLead(req)
	if len(validationErrors) > 0 {
		return nil, fmt.Errorf(strings.Join(validationErrors, ", "))
	}
	leadStatus := e.leadRequest2LeadStatus(req)
	return createLeadStatus(ctx, leadStatus)
}

func (e *LeadService) LeadStatusFromID(ctx context.Context, leadStatusID string) (models.LeadStatus, error) {
	var leadStatus models.LeadStatus
	dbID, err := primitive.ObjectIDFromHex(leadStatusID)
	if err != nil {
		log.Error(err)
		return leadStatus, err
	}
	helper := db.LeadStatus(ctx)

	err = helper.FindOne(ctx, bson.M{"_id": dbID}).Decode(&leadStatus)
	if err != nil {
		log.Error(err)
		return leadStatus, err
	}

	return leadStatus, nil
}

func (e *LeadService) CreateNewLead(ctx context.Context, req *leads.NewLeadRequest) error {
	validationErrors := e.validateLead(req)
	if len(validationErrors) > 0 {
		return fmt.Errorf(strings.Join(validationErrors, ", "))
	}
	leadStatus := e.leadRequest2LeadStatus(req)
	leadStatusID, err := createLeadStatus(ctx, leadStatus)
	if err != nil {
		return err
	}
	leadStatus.ID = leadStatusID.InsertedID.(primitive.ObjectID)
	return e.NewLead(ctx, leadStatus)
}

// NewLead create
func (e *LeadService) NewLead(ctx context.Context, leadStatus models.LeadStatus) error {
	newLead, err := leadFromLeadStatus(ctx, leadStatus)
	leadStatusID := leadStatus.ID
	log.Info("after lead verifed")
	if err != nil {
		return leadStatusCaptureError(ctx, leadStatusID, err)
	}
	existingLead, err := leadByFilter(ctx, bson.M{"contact": newLead.Contact, "email": newLead.Email})
	if err == nil {
		return updateExistingLead(ctx, leadStatusID, existingLead, newLead)
	}
	existingLeadWithContact, err := leadByFilter(ctx, bson.M{"contact": newLead.Contact})
	if err == nil {
		newLead.Meta[0].Meta = append(newLead.Meta[0].Meta, models.KeyValue{Key: "contact", Value: newLead.Contact})
		return updateExistingLead(ctx, leadStatusID, existingLeadWithContact, newLead)
	}
	existingLeadWithEmail, err := leadByFilter(ctx, bson.M{"email": newLead.Email})
	if err == nil {
		newLead.Meta[0].Meta = append(newLead.Meta[0].Meta, models.KeyValue{Key: "email", Value: newLead.Email})
		return updateExistingLead(ctx, leadStatusID, existingLeadWithEmail, newLead)
	}
	log.Info("before create lead")
	// Create a new lead if not existing
	_, err = createLead(ctx, newLead)
	if err != nil {
		log.Error(err)
		leadStatusCaptureError(ctx, leadStatusID, err)
	}
	return leadStatusCaptureSuccess(ctx, leadStatusID)
}

func updateExistingLead(ctx context.Context, leadStatusID primitive.ObjectID, existingLead models.Leads, newLead models.Leads) error {
	existingLead.LeadSource = append(existingLead.LeadSource, newLead.LeadSource...)
	existingLead.Meta = append(existingLead.Meta, newLead.Meta...)
	existingLead.TemplateValues = append(existingLead.TemplateValues, newLead.TemplateValues...)
	log.Info("before update exising lead")
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
		leadStatusCaptureError(ctx, leadStatusID, err)
		return fmt.Errorf("failed to add new lead. details: %s", err.Error())
	}
	return leadStatusCaptureSuccess(ctx, leadStatusID)
}

func (e *LeadService) validateLead(req *leads.NewLeadRequest) []string {
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

func validateSourceAccess(ctx context.Context, campaignCreator primitive.ObjectID, sourceTag string) error {
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

func leadFromLeadStatus(ctx context.Context, req models.LeadStatus) (models.Leads, error) {
	newLead := models.Leads{}
	newLead.Contact = req.Contact
	newLead.Email = req.Email
	newLead.FirstName = req.FirstName
	newLead.LastName = req.LastName

	if preExistingLead(ctx, req) {
		return newLead, fmt.Errorf("lead already captured")
	}

	// Validate campaign tag
	campaign, err := campaignByTag(ctx, req.CampaignTag)
	if err != nil {
		return newLead, err
	}
	newLead.CampaignID = campaign.ID
	err = validateSourceAccess(ctx, campaign.CreatedBy, req.LeadSource)
	if err != nil {
		return newLead, err
	}
	log.Info("lead verifed")
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
	leadTemplateValues := models.SourceWithKeyValue{LeadSource: req.LeadSource, Meta: []models.KeyValue{}}
	for mi := 0; mi < len(req.TemplateValues); mi++ {
		templateKey := req.TemplateValues[mi].Key
		if !isTemplateKey(leadTemplate, templateKey) {
			invalidTemplateKeys = append(invalidTemplateKeys, templateKey)
		}
		leadTemplateValues.Meta = append(leadTemplateValues.Meta, models.KeyValue{Key: templateKey, Value: req.TemplateValues[mi].Value})
	}
	if len(invalidTemplateKeys) > 0 {
		return newLead, fmt.Errorf("the following fields should not be part of the lead template: %s", strings.Join(invalidTemplateKeys, ","))
	}
	newLead.TemplateValues = append([]models.SourceWithKeyValue{}, leadTemplateValues)
	newLead.TemplateID = campaign.TemplateID

	newLead.LeadSource = append([]string{}, req.LeadSource)
	leadMeta := models.SourceWithKeyValue{LeadSource: req.LeadSource, Meta: []models.KeyValue{}}
	for mi := 0; mi < len(req.Meta); mi++ {
		leadMeta.Meta = append(leadMeta.Meta, models.KeyValue{Key: req.Meta[mi].Key, Value: req.Meta[mi].Value})
	}
	currentTime := time.Now()
	leadMeta.Meta = append(leadMeta.Meta, models.KeyValue{Key: "createdOn", Value: currentTime.Format(utilities.TIME_FORMAT)})
	newLead.Meta = append([]models.SourceWithKeyValue{}, leadMeta)

	newLead.CreatedOn = time.Now()
	return newLead, nil
}

func isTemplateKey(leadTemplate lead_template_models.LeadTemplate, key string) bool {
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

func preExistingLead(ctx context.Context, req models.LeadStatus) bool {
	filter := bson.M{
		"contact":    req.Contact,
		"email":      req.Email,
		"leadSource": req.LeadSource,
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

func createLeadStatus(ctx context.Context, newlead models.LeadStatus) (*mongo.InsertOneResult, error) {
	helper := db.LeadStatus(ctx)
	newlead.CreatedOn = time.Now()
	newlead.Status = models.Processing
	return helper.InsertOne(ctx, newlead)
}

func updateLeadStatus(ctx context.Context, leadStatusID primitive.ObjectID, updateFieldAndValues interface{}) (*mongo.UpdateResult, error) {
	helper := db.LeadStatus(ctx)
	return helper.UpdateOne(ctx, bson.M{"_id": leadStatusID}, updateFieldAndValues)
}

func updateLead(ctx context.Context, leadID primitive.ObjectID, updateFieldAndValues interface{}) (*mongo.UpdateResult, error) {
	helper := db.Leads(ctx)
	return helper.UpdateOne(ctx, bson.M{"_id": leadID}, updateFieldAndValues)
}
