package maxmind

import (
	"net"

	"github.com/botchris/go-auditrail/networkd"
	"github.com/oschwald/geoip2-golang"
)

type countryResolver struct {
	db *geoip2.Reader
}

func (a *countryResolver) handle(ip net.IP, geoIP *networkd.GeoIP) {
	record, err := a.db.Country(ip)
	if err != nil {
		return
	}

	writeIfNotEmpty(&geoIP.Continent.Code, record.Continent.Code)
	writeIfNotEmpty(&geoIP.Continent.Name, record.Continent.Names["en"])
	writeIfNotEmpty(&geoIP.Country.Code, record.Country.IsoCode)
	writeIfNotEmpty(&geoIP.Country.Name, record.Country.Names["en"])
}
