package store

import (
	"github.com/Sirupsen/logrus"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/onuryilmaz/body-measurement-api/pkg/commons"
	"time"
)

// StormStoreProvider provides data storage with Bolt Database wrapped by Storm layer
type StormStoreProvider struct {
	fileName string
	db       *storm.DB
}

// NewStormStoreProvider initializes and returns new StormStoreProvider as Provider
func NewStormStoreProvider(options commons.Options) Provider {
	sp := &StormStoreProvider{}
	sp.fileName = options.DatabaseFileName
	return sp
}

// Start starts StormStoreProvider and initializes data
func (sp *StormStoreProvider) Start() error {

	logrus.Info("Starting Storm Store Provider..")
	var err error
	sp.db, err = storm.Open(sp.fileName)
	if err != nil {
		logrus.Error("Error creating Storm Data Store: ", err)
		return err
	}
	return nil
}

// Stop stops StormStoreProvider by closing database
func (sp *StormStoreProvider) Stop() error {

	logrus.Warn("Stopping Storm Store Provider..")
	err := sp.db.Close()
	if err != nil {
		logrus.Error("Error during stopping store provider:", err)
		return err
	}
	return nil
}

// Filter filters data from StormStoreProvider for specific user and measurement type with time frame
func (sp *StormStoreProvider) Filter(user string, measurementType string, from time.Time, to time.Time) ([]commons.BodyMeasurement, error) {

	var data []commons.BodyMeasurement
	err := sp.db.Select(q.Eq("UserID", user), q.Eq("Type", measurementType), q.Gte("Timestamp", from), q.Lte("Timestamp", to)).OrderBy("Timestamp").Find(&data)
	if err != nil {
		logrus.Error("Error filtering data:", err)
		return nil, err
	}
	return data, nil
}

// Last gathers the last data from StormStoreProvider for specific user and measurement type
// It is useful for health dashboard types of applications
func (sp *StormStoreProvider) Last(user string, measurementType string) (commons.BodyMeasurement, error) {

	var bm commons.BodyMeasurement
	query := sp.db.Select(q.Eq("UserID", user), q.Eq("Type", measurementType))
	err := query.OrderBy("Timestamp").Reverse().First(&bm)
	if err != nil {
		logrus.Error("Error getting last data:", err)
		return bm, err
	}
	return bm, nil
}

// Put records measurement in StormStoreProvider
func (sp *StormStoreProvider) Put(data commons.BodyMeasurement) error {

	err := sp.db.Save(&data)
	if err != nil {
		logrus.Error("Error saving data:", err)
		return err
	}
	return nil
}
