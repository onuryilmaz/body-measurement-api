// Package store provides functionality for storing and querying data
package store

import (
	"time"

	"github.com/onuryilmaz/body-measurement-api/pkg/commons"
)

// DataProvider interface defines required actions for data store layer
type DataProvider interface {
	Start() error
	Stop() error
	Filter(user string, measurementType string, from time.Time, to time.Time) ([]commons.BodyMeasurement, error)
	Put(commons.BodyMeasurement) error
}

// TrackingProvider interface defines required actions for data store layer of tracking events
type TrackingProvider interface {
	Start() error
	Stop() error
	Filter(dataConsumer string, dataOwner string, measurementType string, from time.Time, to time.Time) ([]commons.TrackingData, error)
	Put(data commons.TrackingData) error
}
