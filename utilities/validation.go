package utilities

import (
	"fmt"
	"regexp"
	"time"
	"github.com/micro/go-micro/debug/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TIME_FORMAT    = "2006-01-02T15:04:05"
	MAX_TAG_LENGTH = 3
)

// ValidateString if null or empty
func ValidateString(value string) bool {
	if len(value) == 0 {
		return false
	}
	return true
}

// ValidateEmail pop
func ValidateEmail(value string) bool {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegexp.MatchString(value)
}

// DatabaseTimestamp give primitive database timestamp from a string time
func DatabaseTimestamp(dateTimeString string) (time.Time, error) {
	formattedTime, err := time.Parse(TIME_FORMAT, dateTimeString)
	if err != nil {
		log.Error(err)
		return time.Now(), err
	}
	return formattedTime, nil
}

// ValidateTag validates tag format
func ValidateTag(value string) error {
	if len(value) < MAX_TAG_LENGTH {
		return fmt.Errorf("length of the tag should be greater than %d", MAX_TAG_LENGTH)
	}
	tagRegexp := regexp.MustCompile("^[a-zA-Z0-9-_]+$")
	if !tagRegexp.MatchString(value) {
		return fmt.Errorf("tag can have only have alphanumeric characters, \"-\" and \"_\"")
	}
	return nil
}

// CurrentDateTime gives current primitive timestamp
func CurrentDateTime() primitive.Timestamp {
	return primitive.Timestamp{T: uint32(time.Now().Unix())}
}
