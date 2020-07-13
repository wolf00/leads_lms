package services

import (
	"context"
	"fmt"

	"github.com/wolf00/leads_lms/client"
	"github.com/wolf00/leads_lms/constants"
	"github.com/wolf00/user_lms/db/models"
	sourceProto "github.com/wolf00/user_lms/proto/source"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/util/log"
)

func sourceByTag(ctx context.Context, tag string) (models.Sources, error) {
	var source models.Sources
	sourceClient, ok := client.SourceFromContext(ctx)
	if ok {
		return source, errors.InternalServerError(constants.UserService, "user client not found")
	}

	sourceResponse, err := sourceClient.GetBySourceTag(ctx, &sourceProto.SourceRequest{
		SourceTag: tag,
	})

	if err != nil {
		log.Error(err)
		return source, fmt.Errorf("source with the %s tag is not available", tag)
	}
	source.SourceTag = sourceResponse.GetSourceTag()
	source.SystemSource = sourceResponse.SystemSource
	return source, nil
}
