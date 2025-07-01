package maxmind

import (
	"net"

	"github.com/botchris/go-auditrail/networkd"
	"github.com/oschwald/geoip2-golang"
)

type ispResolver struct {
	db *geoip2.Reader
}

func (a *ispResolver) handle(ip net.IP, geoIP *networkd.GeoIP) {
	record, err := a.db.ISP(ip)
	if err != nil {
		return
	}

	if writeIfNotEmpty(&geoIP.AS.Name, record.ISP) {
		geoIP.AS.Type = "isp"
	}
}
