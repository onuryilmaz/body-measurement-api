package server

import (
	"fmt"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/onuryilmaz/body-measurement-api/pkg/commons"
	"github.com/onuryilmaz/body-measurement-api/pkg/store"
	"github.com/phayes/freeport"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/onuryilmaz/body-measurement-api/pkg/tracker"
)

func TestRESTServer(t *testing.T) {

	options := commons.Options{}
	options.ServerPort = fmt.Sprintf("%v", freeport.GetPort())
	dataProvider := &store.InMemoryDataProvider{}
	trackerGateway := tracker.NewTrackerGateway(options)
	RESTServer := NewREST(options, dataProvider, trackerGateway)
	RESTServer.Start()

	Convey("Start and check RESTServer", t, func() {
		So(RESTServer, ShouldNotBeNil)

		// Wait for server is up
		time.Sleep(time.Second)

		Convey("Record a measurement with GET", func() {
			So(len(dataProvider.DB), ShouldEqual, 0)
			response, err := http.Get("http://localhost:" + options.ServerPort + "/api/record/testUser/testType/1.1")
			So(err, ShouldBeNil)
			So(response.StatusCode, ShouldEqual, 200)
			So(len(dataProvider.DB), ShouldEqual, 1)
			So(dataProvider.DB[0].Type, ShouldEqual, "testType")
			So(dataProvider.DB[0].UserID, ShouldEqual, "testUser")
			So(dataProvider.DB[0].Value, ShouldEqual, 1.1)
		})

		Convey("Record a measurement with POST", func() {
			bm := &commons.BodyMeasurement{}
			bm.Value = 1.2
			bm.UserID = "testUser"
			bm.Type = "testType"
			bm.Timestamp = time.Now().AddDate(-1, 0, 0)
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(bm)
			res, err := http.Post("http://localhost:" + options.ServerPort + "/api/save", "application/json; charset=utf-8", b)
			So(err, ShouldBeNil)
			So(res.StatusCode, ShouldEqual, 200)

		})
	})

}
