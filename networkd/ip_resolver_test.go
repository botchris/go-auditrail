package networkd_test

import (
	"os"
	"testing"

	"github.com/botchris/go-auditrail/networkd"
	"github.com/botchris/go-auditrail/networkd/geoip/maxmind"
	"github.com/botchris/go-auditrail/pkg/must"
	"github.com/stretchr/testify/require"
)

func TestNewCachedIPResolver(t *testing.T) {
	inner, err := maxmind.NewMaxmindGeoIPResolver(
		maxmind.WithASNDatabase(must.Read(os.Open("geoip/maxmind/testdata/GeoLite2-ASN.mmdb"))),
		maxmind.WithCityDatabase(must.Read(os.Open("geoip/maxmind/testdata/GeoIP2-City-Test.mmdb"))),
		maxmind.WithCountryDatabase(must.Read(os.Open("geoip/maxmind/testdata/GeoIP2-Country-Test.mmdb"))),
		maxmind.WithISPDatabase(must.Read(os.Open("geoip/maxmind/testdata/GeoIP2-ISP-Test.mmdb"))),
		maxmind.WithDomainDatabase(must.Read(os.Open("geoip/maxmind/testdata/GeoIP2-Domain-Test.mmdb"))),
		maxmind.WithConnectionTypeDatabase(must.Read(os.Open("geoip/maxmind/testdata/GeoIP2-Connection-Type-Test.mmdb"))),
		maxmind.WithEnterpriseDatabase(must.Read(os.Open("geoip/maxmind/testdata/GeoIP2-Enterprise-Test.mmdb"))),
	)
	require.NoError(t, err)

	resolver, err := networkd.NewCachedIPResolver(inner, 128)
	require.NoError(t, err)

	geoip := resolver.Resolve("81.2.69.142")
	require.NotEmpty(t, geoip.Timezone)

	require.Equal(t, 1, resolver.Size())
}
