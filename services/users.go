package services

import (
	"context"

	"github.com/wolf00/golang_lms/leads/client"
	"github.com/wolf00/golang_lms/leads/constants"
	"github.com/wolf00/golang_lms/user/db/models"
	userProto "github.com/wolf00/golang_lms/user/proto/user"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/micro/go-micro/errors"
)

func userByID(ctx context.Context, userID string) (models.Users, error) {
	var user models.Users
	userContext, ok := client.UserFromContext(ctx)
	if !ok {
		return user, errors.InternalServerError(constants.UserService, "user client not found")
	}

	userResponse, err := userContext.GetByID(ctx, &userProto.UserByIdRequest{
		UserId: userID,
	})
	if err != nil {
		return user, err
	}
	user.Contact = userResponse.GetContact()
	user.Email = userResponse.GetEmail()
	user.FirstName = userResponse.GetFirstName()
	user.LastName = userResponse.GetLastName()
	user.OrganizationID, err = primitive.ObjectIDFromHex(userResponse.GetOrganizationId())
	if err != nil {
		return user, err
	}
	user.ID, err = primitive.ObjectIDFromHex(userResponse.GetId())
	if err != nil {
		return user, err
	}
	return user, nil
}
