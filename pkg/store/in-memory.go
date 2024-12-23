package store

import (
	"errors"
	"time"

	"github.com/onuryilmaz/body-measurement-api/pkg/commons"
)

// InMemoryDataProvider provides in-memory data storage for testing and prototyping purposes only
// It keeps all the data in-memory without any persistency
type InMemoryDataProvider struct {
	DB []commons.BodyMeasurement
}

// Start starts InMemoryDataProvider
func (tdp *InMemoryDataProvider) Start() error {
	return nil
}

// Stop stops InMemoryDataProvider
func (tdp *InMemoryDataProvider) Stop() error {
	return nil
}

// Put naively appends BodyMeasurement to InMemoryDataProvider
func (tdp *InMemoryDataProvider) Put(bm commons.BodyMeasurement) error {
	tdp.DB = append(tdp.DB, bm)
	return nil
}


// Filter returns the data instances falling into the time frame
func (tdp *InMemoryDataProvider) Filter(userId string, measurementType string, from time.Time, to time.Time) ([]commons.BodyMeasurement, error) {
	bmFiltered := make([]commons.BodyMeasurement, 0)

	for _, bm := range tdp.DB {
		if bm.UserID == userId && bm.Type == measurementType && bm.Timestamp.After(from) && bm.Timestamp.Before(to) {
			bmFiltered = append(bmFiltered, bm)
		}
	}
	if len(bmFiltered) < 1 {
		return bmFiltered, errors.New("Not found any filtered data for the customer!")
	}
	return bmFiltered, nil

}
