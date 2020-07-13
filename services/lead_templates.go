package services

import (
	"context"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/wolf00/lead_template_lms/db/models"
	"github.com/wolf00/leads_lms/client"
	"github.com/wolf00/leads_lms/constants"

	leadTemplateProto "github.com/wolf00/lead_template_lms/proto/lead_template"
)

// LeadTemplateClientService to get
type LeadTemplateClientService struct{}

func leadTemplateByID(ctx context.Context, templateID string) (models.LeadTemplate, error) {
	var leadTemplate models.LeadTemplate

	leadTemplateClient, ok := client.LeadTemplateFromContext(ctx)

	if !ok {
		return leadTemplate, errors.InternalServerError(constants.LeadTemplateService, "lead template client not found")
	}

	leadTemplateResp, err := leadTemplateClient.Get(ctx, &leadTemplateProto.LeadTemplateByIdRequest{
		TemplateId: templateID,
	})

	if err != nil {
		log.Error(err)
		return leadTemplate, err
	}
	leadTemplate.Name = leadTemplateResp.GetName()
	keyValueTypes := []models.TemplateKeyValueTypes{}

	for i := 0; i < len(leadTemplateResp.GetKeyValueTypes()); i++ {
		respKeyValue := leadTemplateResp.GetKeyValueTypes()[i]
		keyValueTypes = append(keyValueTypes, models.TemplateKeyValueTypes{
			Key:       respKeyValue.GetKey(),
			ValueType: respKeyValue.GetValueType(),
		})
	}

	leadTemplate.KeyValueTypes = keyValueTypes

	return leadTemplate, nil
}
