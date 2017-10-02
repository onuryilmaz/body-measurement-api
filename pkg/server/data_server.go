// Package server provides functionality for handling data input and queries
package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/onuryilmaz/body-measurement-api/pkg/commons"
	"github.com/onuryilmaz/body-measurement-api/pkg/store"
	"github.com/onuryilmaz/body-measurement-api/pkg/tracker"
)

// REST provides functionality for HTTP REST API Server
type REST struct {
	router          *httprouter.Router
	server          *http.Server
	port            string
	datastore       store.DataProvider
	trackerGateway *tracker.TrackerGateway
}

// NewREST creates a REST API server instance with the provided options and datastore layer
func NewREST(options commons.Options, datastore store.DataProvider, trackerGateway *tracker.TrackerGateway) *REST {
	rest := &REST{}
	rest.datastore = datastore
	rest.port = options.ServerPort
	rest.trackerGateway = trackerGateway
	rest.router = httprouter.New()
	return rest
}

// Start starts REST API server and connects handlers to the router on port
func (r *REST) Start() {

	logrus.Info("Starting REST server...")
	logrus.Infof("REST server connecting to port %v", r.port)

	r.router.GET("/api/filter/:user/:type/:from/:to", r.filterHandler)
	r.router.GET("/api/record/:user/:type/:value", r.recordHandler)
	r.router.POST("/api/save", r.recordHandler)

	r.router.GET("/api/access/:consumer/:user/:type/:from/:to", r.consumerFilterHandler)

	r.server = &http.Server{Addr: ":" + r.port, Handler: r.router}
	go r.server.ListenAndServe()
}

// Stop stops REST API gracefully
func (r *REST) Stop() {
	logrus.Warn("Stopping REST server..")
	r.server.Shutdown(context.TODO())
}

func (r *REST) consumerFilterHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

	consumer := p.ByName("consumer")
	user := p.ByName("user")
	measurementType := p.ByName("type")
	from := p.ByName("from")
	to := p.ByName("to")

	fromInt, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		logrus.Error("Error during from parsing:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fromTime := time.Unix(fromInt, 0)

	toInt, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		logrus.Error("Error during to parsing:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	toTime := time.Unix(toInt, 0)

	logrus.Debugf("Filter REST handler for user %s | type %s | from %v to %v", user, measurementType, fromTime, toTime)

	r.trackerGateway.Track(consumer,user,measurementType, fromTime, toTime)

	data, err := r.datastore.Filter(user, measurementType, fromTime, toTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (r *REST) filterHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

	user := p.ByName("user")
	measurementType := p.ByName("type")
	from := p.ByName("from")
	to := p.ByName("to")

	fromInt, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		logrus.Error("Error during from parsing:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fromTime := time.Unix(fromInt, 0)

	toInt, err := strconv.ParseInt(to, 10, 64)
	if err != nil {
		logrus.Error("Error during to parsing:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	toTime := time.Unix(toInt, 0)

	logrus.Debugf("Filter REST handler for user %s | type %s | from %v to %v", user, measurementType, fromTime, toTime)

	data, err := r.datastore.Filter(user, measurementType, fromTime, toTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (r *REST) recordHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

	data := &commons.BodyMeasurement{}
	if req.Method == "GET" {
		user := p.ByName("user")
		measurementType := p.ByName("type")
		value := p.ByName("value")

		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			logrus.Error("Error parsing value:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		data = &commons.BodyMeasurement{UserID: user, Value: valueFloat, Type: measurementType, Timestamp: time.Now()}

		logrus.Debugf("Record handler for user %s | type %s | type %v", user, measurementType, valueFloat)
	} else if req.Method == "POST" {

		err := json.NewDecoder(req.Body).Decode(data)
		if err != nil {
			logrus.Error("Error decoding received data:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err := r.datastore.Put(*data)
	if err != nil {
		logrus.Error("Error recording data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return

}
