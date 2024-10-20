package maxmind

import (
	"net"

	"github.com/botchris/auditrail/networkd"
	"github.com/oschwald/geoip2-golang"
)

type connectionTypeResolver struct {
	db *geoip2.Reader
}

func (a *connectionTypeResolver) handle(ip net.IP, geoIP *networkd.GeoIP) {
	record, err := a.db.ConnectionType(ip)
	if err != nil {
		return
	}

	if geoIP.AS.Type == "" && record.ConnectionType != "" {
		geoIP.AS.Type = record.ConnectionType
	}
}
