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
)

// RESTTracking provides functionality for HTTP REST API Server
type RESTTracking struct {
	router    *httprouter.Router
	server    *http.Server
	port      string
	datastore store.TrackingProvider
}

// NewRESTTracking creates a REST Tracking API server instance with the provided options and datastore layer
func NewRESTTracking(options commons.Options, datastore store.TrackingProvider) *RESTTracking {
	rest := &RESTTracking{}
	rest.datastore = datastore
	rest.port = options.ServerPort
	rest.router = httprouter.New()
	return rest
}

// Start starts REST API server and connects handlers to the router on port
func (r *RESTTracking) Start() {

	logrus.Info("Starting REST server...")
	logrus.Infof("REST server connecting to port %v", r.port)

	r.router.GET("/api/filter/:owner/:consumer/:type/:from/:to", r.filterHandler)
	r.router.POST("/api/record", r.recordHandler)

	r.server = &http.Server{Addr: ":" + r.port, Handler: r.router}
	go r.server.ListenAndServe()
}

// Stop stops REST API gracefully
func (r *RESTTracking) Stop() {
	logrus.Warn("Stopping REST server..")
	r.server.Shutdown(context.TODO())
}

func (r *RESTTracking) filterHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

	owner := p.ByName("owner")
	measurementType := p.ByName("type")
	consumer := p.ByName("consumer")
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

	logrus.Debugf("Filter REST handler for user %s | type %s | from %v to %v", owner, measurementType, fromTime, toTime)

	data, err := r.datastore.Filter(consumer, owner, measurementType, fromTime, toTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (r *RESTTracking) recordHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

	data := &commons.TrackingData{}
	 if req.Method == "POST" {
		err := json.NewDecoder(req.Body).Decode(data)
		if err != nil {
			logrus.Error("Error decoding received data:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		 err = r.datastore.Put(*data)
		 if err != nil {
			 logrus.Error("Error recording data:", err)
			 w.WriteHeader(http.StatusInternalServerError)
			 return
		 }
		 w.WriteHeader(http.StatusOK)
		 return
	} else {
		 logrus.Error("Not supported request type: ", req.Method)
		 w.WriteHeader(http.StatusBadRequest)
		 return
	 }



}
