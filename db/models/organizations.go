package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Organizations typ definition
type Organizations struct {
	Name           string
	Email          string
	Contact        string
	LogoLink       string
	ID             primitive.ObjectID
	CreatedOn      primitive.Timestamp
	Active         bool
	AllowedSources []string
	SourceTag      string
}
