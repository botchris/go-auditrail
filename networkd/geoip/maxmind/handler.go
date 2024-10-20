package maxmind

import (
	"net"

	"github.com/botchris/auditrail/networkd"
)

type handler interface {
	handle(ip net.IP, geoIP *networkd.GeoIP)
}
