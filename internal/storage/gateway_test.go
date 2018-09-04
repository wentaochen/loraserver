package storage

import (
	"testing"
	"time"

	"github.com/brocaar/loraserver/internal/common"
	"github.com/brocaar/loraserver/internal/test"
	"github.com/brocaar/lorawan"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGateway(t *testing.T) {
	conf := test.GetConfig()

	Convey("Given a clean database", t, func() {
		db, err := common.OpenDatabase(conf.PostgresDSN)
		So(err, ShouldBeNil)
		test.MustResetDB(db)

		Convey("When creating a gateway", func() {
			gw := Gateway{
				GatewayID: lorawan.EUI64{29, 238, 8, 208, 182, 145, 209, 73},
				Location: GPSPoint{
					Latitude:  1.23456789,
					Longitude: 4.56789012,
				},
			}
			So(CreateGateway(db, &gw), ShouldBeNil)

			// some precicion will get lost when writing to the db
			// truncate it to ms precision for comparison
			gw.CreatedAt = gw.CreatedAt.UTC().Truncate(time.Millisecond)
			gw.UpdatedAt = gw.UpdatedAt.UTC().Truncate(time.Millisecond)

			Convey("Then it can be retrieved", func() {
				gw2, err := GetGateway(db, gw.GatewayID)
				So(err, ShouldBeNil)

				gw2.CreatedAt = gw2.CreatedAt.UTC().Truncate(time.Millisecond)
				gw2.UpdatedAt = gw2.UpdatedAt.UTC().Truncate(time.Millisecond)

				So(gw2, ShouldResemble, gw)

				gws, err := GetGatewaysForIDs(db, []lorawan.EUI64{gw.GatewayID})
				So(err, ShouldBeNil)
				gw3, ok := gws[gw.GatewayID]
				So(ok, ShouldBeTrue)
				So(gw3.GatewayID, ShouldResemble, gw.GatewayID)
			})

			Convey("Then it can be updated", func() {
				now := time.Now().UTC().Truncate(time.Millisecond)

				gw.FirstSeenAt = &now
				gw.LastSeenAt = &now
				gw.Location.Latitude = 1.23456780
				gw.Location.Longitude = 5.56789012
				gw.Altitude = 100.5

				fpgaID := lorawan.EUI64{1, 2, 3, 4, 5, 6, 7, 8}
				ftKey := lorawan.AES128Key{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

				gw.FPGAID = &fpgaID
				gw.FineTimestampAESKey = &ftKey

				So(UpdateGateway(db, &gw), ShouldBeNil)

				gw2, err := GetGateway(db, gw.GatewayID)
				So(err, ShouldBeNil)

				So(gw2.GatewayID, ShouldEqual, gw.GatewayID)
				So(gw2.CreatedAt.UTC().Truncate(time.Millisecond), ShouldResemble, gw.CreatedAt.UTC())
				So(gw2.UpdatedAt.UTC().Truncate(time.Millisecond), ShouldResemble, gw.UpdatedAt.UTC().Truncate(time.Millisecond))
				So(gw2.FirstSeenAt.UTC().Truncate(time.Millisecond), ShouldResemble, gw.FirstSeenAt.UTC().Truncate(time.Millisecond))
				So(gw2.LastSeenAt.UTC().Truncate(time.Millisecond), ShouldResemble, gw.LastSeenAt.UTC().Truncate(time.Millisecond))
				So(gw2.Location, ShouldResemble, gw.Location)
				So(gw2.Altitude, ShouldResemble, gw.Altitude)
			})

			Convey("Then it can be deleted", func() {
				So(DeleteGateway(db, gw.GatewayID), ShouldBeNil)
				_, err := GetGateway(db, gw.GatewayID)
				So(err, ShouldResemble, ErrDoesNotExist)
			})
		})
	})
}
