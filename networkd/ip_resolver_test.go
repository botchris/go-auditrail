package networkd_test

import (
	"testing"

	"github.com/botchris/go-auditrail/networkd"
	"github.com/botchris/go-auditrail/networkd/geoip/maxmind"
	"github.com/stretchr/testify/require"
)

func TestNewCachedIPResolver(t *testing.T) {
	inner, err := maxmind.NewMaxmindGeoIPResolver(
		maxmind.WithASNDatabase("geoip/maxmind/testdata/GeoLite2-ASN.mmdb"),
		maxmind.WithCityDatabase("geoip/maxmind/testdata/GeoIP2-City-Test.mmdb"),
		maxmind.WithCountryDatabase("geoip/maxmind/testdata/GeoIP2-Country-Test.mmdb"),
		maxmind.WithISPDatabase("geoip/maxmind/testdata/GeoIP2-ISP-Test.mmdb"),
		maxmind.WithDomainDatabase("geoip/maxmind/testdata/GeoIP2-Domain-Test.mmdb"),
		maxmind.WithConnectionTypeDatabase("geoip/maxmind/testdata/GeoIP2-Connection-Type-Test.mmdb"),
		maxmind.WithEnterpriseDatabase("geoip/maxmind/testdata/GeoIP2-Enterprise-Test.mmdb"),
	)
	require.NoError(t, err)

	resolver, err := networkd.NewCachedIPResolver(inner, 128)
	require.NoError(t, err)

	geoip := resolver.Resolve("81.2.69.142")
	require.NotEmpty(t, geoip.Timezone)

	require.Equal(t, 1, resolver.Size())
}
