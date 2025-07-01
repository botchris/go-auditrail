package maxmind

import (
	"fmt"
	"net"

	"github.com/botchris/go-auditrail/networkd"
	"github.com/oschwald/geoip2-golang"
)

type enterpriseResolver struct {
	db *geoip2.Reader
}

func (a *enterpriseResolver) handle(ip net.IP, geoIP *networkd.GeoIP) {
	record, err := a.db.City(ip)
	if err != nil {
		return
	}

	writeIfNotEmpty(&geoIP.Continent.Code, record.Continent.Code)
	writeIfNotEmpty(&geoIP.Continent.Name, record.Continent.Names["en"])
	writeIfNotEmpty(&geoIP.Country.Code, record.Country.IsoCode)
	writeIfNotEmpty(&geoIP.Country.Name, record.Country.Names["en"])
	writeIfNotEmpty(&geoIP.City.Code, fmt.Sprintf("%d", record.City.GeoNameID))
	writeIfNotEmpty(&geoIP.City.Name, record.City.Names["en"])
	writeIfNotEmpty(&geoIP.Timezone, record.Location.TimeZone)

	writeIfNotEmpty(&geoIP.Location.Latitude, record.Location.Latitude)
	writeIfNotEmpty(&geoIP.Location.Longitude, record.Location.Longitude)
}
