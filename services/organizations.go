package services

import (
	"context"

	"github.com/wolf00/golang_lms/leads/client"
	"github.com/wolf00/golang_lms/leads/constants"
	"github.com/wolf00/golang_lms/leads/db/models"
	organizationProto "github.com/wolf00/golang_lms/user/proto/organization"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/micro/go-micro/errors"
)

func organizationByID(ctx context.Context, organizationID string) (models.Organizations, error) {
	var organization models.Organizations
	organizationContext, ok := client.OrganizationFromContext(ctx)
	if !ok {
		return organization, errors.InternalServerError(constants.UserService, "user client not found")
	}

	orgResponse, err := organizationContext.GetByID(ctx, &organizationProto.OrganizationByIdRequest{
		OrganizationId: organizationID,
	})
	if err != nil {
		return organization, err
	}
	organization.AllowedSources = orgResponse.GetAllowedSources()
	organization.Contact = orgResponse.GetContact()
	organization.Email = orgResponse.GetEmail()
	organization.LogoLink = orgResponse.GetLogoLink()
	organization.Name = orgResponse.GetName()
	organization.SourceTag = orgResponse.GetSourceTag()
	organization.ID, err = primitive.ObjectIDFromHex(orgResponse.GetId())
	if err != nil {
		return organization, err
	}
	return organization, nil
}
