package store

import (
	"time"

	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/onuryilmaz/body-measurement-api/pkg/commons"
)

// StormStoreTrackingProvider provides data storage with Bolt Database wrapped by Storm layer
type StormStoreTrackingProvider struct {
	fileName string
	db       *storm.DB
}

// NewStormStoreProvider initializes and returns new StormStoreTrackingProvider as Provider
func NewStormStoreTrackingProvider(options commons.Options) TrackingProvider {
	sp := &StormStoreTrackingProvider{}
	sp.fileName = options.DatabaseFileName
	return sp
}

// Start starts StormStoreTrackingProvider and initializes data
func (sp *StormStoreTrackingProvider) Start() error {

	logrus.Info("Starting Storm Store Provider..")
	var err error
	sp.db, err = storm.Open(sp.fileName)
	if err != nil {
		logrus.Error("Error creating Storm Data Store: ", err)
		return err
	}
	return nil
}

// Stop stops StormStoreTrackingProvider by closing database
func (sp *StormStoreTrackingProvider) Stop() error {

	logrus.Warn("Stopping Storm Store Provider..")
	err := sp.db.Close()
	if err != nil {
		logrus.Error("Error during stopping store provider:", err)
		return err
	}
	return nil
}

// Filter filters data from StormStoreTrackingProvider for specific user and measurement type with time frame
func (sp *StormStoreTrackingProvider) Filter(dataConsumer string, dataOwner string, measurementType string, from time.Time, to time.Time) ([]commons.TrackingData, error) {

	var data []commons.TrackingData

	if dataOwner == "" {

		return nil, errors.New("empty data owner for filtering")
	}

	query := []q.Matcher{q.Eq("DataOwnerId", dataOwner), q.Gte("Timestamp", from), q.Lte("Timestamp", to)}
	if dataConsumer != "all" {
		query = append(query, q.Eq("DataConsumerId", dataConsumer))
	}

	if measurementType != "all" {
		query = append(query, q.Eq("Type", measurementType))
	}

	err := sp.db.Select(query...).OrderBy("Timestamp").Find(&data)

	if err != nil && err.Error() == "not found" { //TODO also be retuned no data as HTTP method
		return data, nil
	}
	if err != nil {
		logrus.Error("Error filtering data:", err)
		return nil, err
	}
	return data, nil
}

// Put records measurement in StormStoreProvider
func (sp *StormStoreTrackingProvider) Put(data commons.TrackingData) error {

	err := sp.db.Save(&data)
	if err != nil {
		logrus.Error("Error saving data:", err)
		return err
	}
	return nil
}
