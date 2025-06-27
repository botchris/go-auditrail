package maxmind_test

import (
	"os"
	"testing"

	"github.com/botchris/go-auditrail/networkd/geoip/maxmind"
	"github.com/botchris/go-auditrail/pkg/must"
	"github.com/stretchr/testify/require"
)

func TestMaxmind(t *testing.T) {
	resolver, err := maxmind.NewMaxmindGeoIPResolver(
		maxmind.WithASNDatabase(must.Read(os.Open("testdata/GeoLite2-ASN.mmdb"))),
		maxmind.WithCityDatabase(must.Read(os.Open("testdata/GeoIP2-City-Test.mmdb"))),
		maxmind.WithCountryDatabase(must.Read(os.Open("testdata/GeoIP2-Country-Test.mmdb"))),
		maxmind.WithISPDatabase(must.Read(os.Open("testdata/GeoIP2-ISP-Test.mmdb"))),
		maxmind.WithDomainDatabase(must.Read(os.Open("testdata/GeoIP2-Domain-Test.mmdb"))),
		maxmind.WithConnectionTypeDatabase(must.Read(os.Open("testdata/GeoIP2-Connection-Type-Test.mmdb"))),
		maxmind.WithEnterpriseDatabase(must.Read(os.Open("testdata/GeoIP2-Enterprise-Test.mmdb"))),
	)

	require.NoError(t, err)

	geoIP := resolver.Resolve("81.2.69.142")

	require.NotEmpty(t, geoIP.Continent.Code)
	require.NotEmpty(t, geoIP.Continent.Name)

	require.NotEmpty(t, geoIP.Country.Code)
	require.NotEmpty(t, geoIP.Country.Name)

	require.NotEmpty(t, geoIP.City.Code)
	require.NotEmpty(t, geoIP.City.Name)

	require.NotEmpty(t, geoIP.AS.Domain)
	require.NotEmpty(t, geoIP.AS.Name)
	require.NotEmpty(t, geoIP.AS.Number)

	require.NotEmpty(t, geoIP.Timezone)
}

func BenchmarkMaxmindResolve(b *testing.B) {
	resolver, err := maxmind.NewMaxmindGeoIPResolver(
		maxmind.WithASNDatabase(must.Read(os.Open("testdata/GeoLite2-ASN.mmdb"))),
		maxmind.WithCityDatabase(must.Read(os.Open("testdata/GeoIP2-City-Test.mmdb"))),
		maxmind.WithCountryDatabase(must.Read(os.Open("testdata/GeoIP2-Country-Test.mmdb"))),
		maxmind.WithISPDatabase(must.Read(os.Open("testdata/GeoIP2-ISP-Test.mmdb"))),
		maxmind.WithDomainDatabase(must.Read(os.Open("testdata/GeoIP2-Domain-Test.mmdb"))),
		maxmind.WithConnectionTypeDatabase(must.Read(os.Open("testdata/GeoIP2-Connection-Type-Test.mmdb"))),
		maxmind.WithEnterpriseDatabase(must.Read(os.Open("testdata/GeoIP2-Enterprise-Test.mmdb"))),
	)

	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = resolver.Resolve("81.2.69.142")
	}
}
