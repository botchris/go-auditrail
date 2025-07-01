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

//nolint:gocyclo
func (a *cityResolver) handle(ip net.IP, geoIP *networkd.GeoIP) {
	record, err := a.db.City(ip)
	if err != nil {
		return
	}

	writeIfNotEmpty(&geoIP.City.Code, fmt.Sprintf("%d", record.City.GeoNameID))
	writeIfNotEmpty(&geoIP.City.Name, record.City.Names["en"])
	writeIfNotEmpty(&geoIP.Continent.Code, record.Continent.Code)
	writeIfNotEmpty(&geoIP.Continent.Name, record.Continent.Names["en"])
	writeIfNotEmpty(&geoIP.Country.Code, record.Country.IsoCode)
	writeIfNotEmpty(&geoIP.Country.Name, record.Country.Names["en"])
	writeIfNotEmpty(&geoIP.Timezone, record.Location.TimeZone)
	writeIfNotEmpty(&geoIP.Location.Latitude, record.Location.Latitude)
	writeIfNotEmpty(&geoIP.Location.Longitude, record.Location.Longitude)

	// Subdivision (only if present)
	if len(record.Subdivisions) > 0 {
		writeIfNotEmpty(&geoIP.Subdivision.Code, record.Subdivisions[0].IsoCode)
		writeIfNotEmpty(&geoIP.Subdivision.Name, record.Subdivisions[0].Names["en"])
	}
}
