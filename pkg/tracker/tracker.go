package tracker

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/onuryilmaz/body-measurement-api/pkg/commons"
	"net/http"
	"time"
)

type TrackerGateway struct {
	trackingServer string
}

func NewTrackerGateway(options commons.Options) *TrackerGateway {
	tr := &TrackerGateway{}
	tr.trackingServer = options.TrackingAddress
	return tr
}

func (tr *TrackerGateway) Track(dataConsumer string, dataOwner string, measurementType string, from time.Time, to time.Time) error {

	if dataConsumer != "" && dataOwner != "" && measurementType != "" {
		trackingData := &commons.TrackingData{}
		trackingData.DataConsumerId = dataConsumer
		trackingData.DataOwnerId = dataOwner
		trackingData.Timestamp = time.Now()
		trackingData.Type = measurementType

		jsonValue, err := json.Marshal(trackingData)


		resp, err := http.Post(tr.trackingServer, "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			logrus.Error("Error during tracking request:", err)
			return err
		}

		if resp.StatusCode != http.StatusOK {
			logrus.Error("Error response code of tracking request:", resp.StatusCode)
			return errors.New("tracking request is not successful")
		}

	} else {
		return errors.New("insufficient data to call tracker")
	}

	logrus.Debugf("Tracking request is successful for consumer %s | owner %s | type %s | from %v | to %v", dataConsumer, dataOwner, measurementType, from, to)
	return nil
}
