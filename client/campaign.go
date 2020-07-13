package client

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
	campaign "github.com/wolf00/campaign_lms/proto/campaign"
	"github.com/wolf00/leads_lms/constants"
)

type campaignKey struct{}

// CampaignFromContext retrieves the client from the Context
func CampaignFromContext(ctx context.Context) (campaign.CampaignService, bool) {
	c, ok := ctx.Value(campaignKey{}).(campaign.CampaignService)
	return c, ok
}

// CampaignWrapper returns a wrapper for the HeimdallClient
func CampaignWrapper(service micro.Service) server.HandlerWrapper {
	client := campaign.NewCampaignService(constants.CampaignService, service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, campaignKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
