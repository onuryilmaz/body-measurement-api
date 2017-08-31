// Package store provides functionality for storing and querying data
package store

import (
	"github.com/onuryilmaz/body-measurement-api/pkg/commons"
	"time"
)

// Provider interface defines required actions for data store layer
type Provider interface {
	Start() error
	Stop() error
	Filter(string, string, time.Time, time.Time) ([]commons.BodyMeasurement, error)
	Last(string, string) (commons.BodyMeasurement, error)
	Put(commons.BodyMeasurement) error
}
