package maxmind

import (
	"net"

	"github.com/botchris/go-auditrail/networkd"
)

type maxmindGeoIPResolver struct {
	handlers []handler
}

// NewMaxmindGeoIPResolver returns a new IPResolver that uses the Maxmind
// GeoIP databases.
//
// see: https://github.com/P3TERX/GeoLite.mmdb
func NewMaxmindGeoIPResolver(opts ...Option) (networkd.IPResolver, error) {
	r := &maxmindGeoIPResolver{}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (m *maxmindGeoIPResolver) Resolve(ip string) networkd.GeoIP {
	geoIP := &networkd.GeoIP{}
	netIP := net.ParseIP(ip)

	for _, h := range m.handlers {
		h.handle(netIP, geoIP)
	}

	return *geoIP
}
