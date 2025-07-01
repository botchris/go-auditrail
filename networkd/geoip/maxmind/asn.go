package maxmind

import (
	"fmt"
	"net"

	"github.com/botchris/go-auditrail/networkd"
	"github.com/oschwald/geoip2-golang"
)

type asnResolver struct {
	db *geoip2.Reader
}

func (a *asnResolver) handle(ip net.IP, geoIP *networkd.GeoIP) {
	record, err := a.db.ASN(ip)
	if err != nil {
		return
	}

	writeIfNotEmpty(&geoIP.AS.Name, record.AutonomousSystemOrganization)

	if geoIP.AS.Number == "" && record.AutonomousSystemNumber != 0 {
		geoIP.AS.Number = fmt.Sprintf("AS%d", record.AutonomousSystemNumber)
	}
}
