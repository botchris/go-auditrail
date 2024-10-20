package maxmind

import (
	"net"

	"github.com/botchris/auditrail/networkd"
	"github.com/oschwald/geoip2-golang"
)

type domainResolver struct {
	db *geoip2.Reader
}

func (a *domainResolver) handle(ip net.IP, geoIP *networkd.GeoIP) {
	record, err := a.db.Domain(ip)
	if err != nil {
		return
	}

	if geoIP.AS.Domain == "" && record.Domain != "" {
		geoIP.AS.Domain = record.Domain
	}
}
