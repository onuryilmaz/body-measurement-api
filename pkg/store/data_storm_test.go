package store

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/onuryilmaz/body-measurement-api/pkg/commons"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStormProvider(t *testing.T) {

	options := commons.Options{}
	fileName := fmt.Sprintf("%v.db", time.Now().Unix())
	options.DatabaseFileName = fileName
	defer os.Remove(fileName)

	datastore := NewStormStoreProvider(options)

	Convey("Create Storm Store Provider", t, func() {

		Convey("Start and check provider", func() {

			So(datastore, ShouldNotBeNil)
			err := datastore.Start()
			So(err, ShouldBeNil)
		})

		Convey("Store data", func() {

			bm := commons.BodyMeasurement{}
			bm.Timestamp = time.Now()
			bm.Value = 1.1
			bm.UserID = "1001"
			bm.Type = "test"
			err := datastore.Put(bm)
			So(err, ShouldBeNil)

			bm = commons.BodyMeasurement{}
			bm.Timestamp = time.Now().AddDate(-1, 0, 0) // Data for last year
			bm.Value = 1.2
			bm.UserID = "1001"
			bm.Type = "test"
			err = datastore.Put(bm)
			So(err, ShouldBeNil)
		})

		Convey("Check filtered data", func() {

			// Filter for the last year's data
			bm, err := datastore.Filter("1001", "test", time.Now().AddDate(-2, 0, 0), time.Now().AddDate(0, -1, 0))
			So(err, ShouldBeNil)
			So(len(bm), ShouldEqual, 1)
			So(bm[0].Type, ShouldEqual, "test")
			So(bm[0].Value, ShouldEqual, 1.2)
			So(bm[0].UserID, ShouldEqual, "1001")
		})

		Convey("Close data store", func() {

			err := datastore.Stop()
			So(err, ShouldBeNil)
		})

	})

}
