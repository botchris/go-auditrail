package maxmind

import (
	"net"

	"github.com/botchris/go-auditrail/networkd"
)

type handler interface {
	handle(ip net.IP, geoIP *networkd.GeoIP)
}

func writeIfNotEmpty[T comparable](target *T, value T) bool {
	var zero T

	if target != nil && *target == zero && value != zero {
		*target = value

		return true
	}

	return false
}
