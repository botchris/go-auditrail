package maxmind

import (
	"fmt"
	"net"

	"github.com/botchris/go-auditrail/networkd"
	"github.com/oschwald/geoip2-golang"
)

type cityResolver struct {
	db *geoip2.Reader
}

func (a *cityResolver) handle(ip net.IP, geoIP *networkd.GeoIP) {
	record, err := a.db.City(ip)
	if err != nil {
		return
	}

	if geoIP.City.Code == "" && record.City.GeoNameID != 0 {
		geoIP.City.Code = fmt.Sprintf("%d", record.City.GeoNameID)
	}

	if geoIP.City.Name == "" && record.City.Names["en"] != "" {
		geoIP.City.Name = record.City.Names["en"]
	}

	if geoIP.Continent.Code == "" && record.Continent.Code != "" {
		geoIP.Continent.Code = record.Continent.Code
	}

	if geoIP.Continent.Name == "" && record.Country.Names["en"] != "" {
		geoIP.Country.Name = record.Country.Names["en"]
	}

	if geoIP.Location.Latitude == 0 && record.Location.Latitude != 0 {
		geoIP.Location.Latitude = record.Location.Latitude
	}

	if geoIP.Location.Longitude == 0 && record.Location.Longitude != 0 {
		geoIP.Location.Longitude = record.Location.Longitude
	}

	if geoIP.Timezone == "" && record.Location.TimeZone != "" {
		geoIP.Timezone = record.Location.TimeZone
	}

	if geoIP.Subdivision.Code == "" && len(record.Subdivisions) > 0 {
		geoIP.Subdivision.Code = record.Subdivisions[0].IsoCode
		geoIP.Subdivision.Name = record.Subdivisions[0].Names["en"]
	}
}
