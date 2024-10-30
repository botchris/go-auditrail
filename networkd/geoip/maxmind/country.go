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

	if geoIP.Continent.Code == "" && record.Continent.Code != "" {
		geoIP.Continent.Code = record.Continent.Code
	}

	if geoIP.Continent.Name == "" && record.Continent.Names["en"] != "" {
		geoIP.Continent.Name = record.Continent.Names["en"]
	}

	if geoIP.Country.Code == "" && record.Country.IsoCode != "" {
		geoIP.Country.Code = record.Country.IsoCode
	}

	if geoIP.Country.Name == "" && record.Country.Names["en"] != "" {
		geoIP.Country.Name = record.Country.Names["en"]
	}
}
