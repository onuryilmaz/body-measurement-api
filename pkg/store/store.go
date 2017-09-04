// Package store provides functionality for storing and querying data
package store

import (
	"time"

	"github.com/onuryilmaz/body-measurement-api/pkg/commons"
)

// Provider interface defines required actions for data store layer
type Provider interface {
	Start() error
	Stop() error
	Filter(user string, measurementType string, from time.Time, to time.Time) ([]commons.BodyMeasurement, error)
	Last(user string, measurementType string) (commons.BodyMeasurement, error)
	Put(commons.BodyMeasurement) error
}
