package services

import (
	"context"
	"fmt"

	"github.com/wolf00/golang_lms/campaign/db/models"
	campaignProto "github.com/wolf00/golang_lms/campaign/proto/campaign"
	"github.com/wolf00/golang_lms/leads/client"
	"github.com/wolf00/golang_lms/leads/constants"
	"github.com/wolf00/golang_lms/leads/utilities"

	"github.com/micro/go-micro/debug/log"
	"github.com/micro/go-micro/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CampaignIDByTag retrieves campaign id by campaign tag
func campaignByTag(ctx context.Context, tag string) (models.Campaign, error) {
	var campaign models.Campaign

	campaignContext, ok := client.CampaignFromContext(ctx)
	if !ok {
		return campaign, errors.InternalServerError(constants.CampaignService, "campaign client not found")
	}
	campaignResponse, err := campaignContext.GetByTag(ctx, &campaignProto.CampaignByTagRequest{
		CampaignTag: tag,
	})
	if err != nil {
		log.Error(err)
		return campaign, fmt.Errorf("campaign with the %s tag is not available", tag)
	}
	return campaignResponse2Campaign(*campaignResponse)
}

func campaignResponse2Campaign(campaignResponse campaignProto.CampaignResponse) (models.Campaign, error) {
	var campaign models.Campaign
	var err error
	campaign.CampaignName = campaignResponse.GetCampaignName()
	campaign.CampaignTag = campaignResponse.GetCampaignTag()
	campaign.CreatedBy, err = primitive.ObjectIDFromHex(campaignResponse.GetCreatedBy())
	if err != nil {
		log.Error(err)
		return campaign, err
	}
	campaign.Description = campaignResponse.GetDescription()
	campaign.EndDateTime, err = utilities.DatabaseTimestamp(campaignResponse.GetEndDateTime())
	if err != nil {
		log.Error(err)
		return campaign, err
	}
	campaign.Purpose = campaignResponse.GetPurpose()
	campaign.StartDateTime, err = utilities.DatabaseTimestamp(campaignResponse.GetStartDateTime())
	if err != nil {
		log.Error(err)
		return campaign, err
	}
	campaign.TemplateID, err = primitive.ObjectIDFromHex(campaignResponse.GetTemplateId())
	if err != nil {
		log.Error(err)
		return campaign, err
	}
	campaign.ID, err = primitive.ObjectIDFromHex(campaignResponse.GetId())
	if err != nil {
		log.Error(err)
		return campaign, err
	}
	return campaign, nil
}
