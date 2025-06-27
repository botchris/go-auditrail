package networkd_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/botchris/go-auditrail/networkd"
	"github.com/botchris/go-auditrail/networkd/geoip/maxmind"
	"github.com/stretchr/testify/require"
)

func TestNewCachedIPResolver(t *testing.T) {
	inner, err := maxmind.NewMaxmindGeoIPResolver(
		maxmind.WithASNDatabase(Must(os.Open("geoip/maxmind/testdata/GeoLite2-ASN.mmdb"))),
		maxmind.WithCityDatabase(Must(os.Open("geoip/maxmind/testdata/GeoIP2-City-Test.mmdb"))),
		maxmind.WithCountryDatabase(Must(os.Open("geoip/maxmind/testdata/GeoIP2-Country-Test.mmdb"))),
		maxmind.WithISPDatabase(Must(os.Open("geoip/maxmind/testdata/GeoIP2-ISP-Test.mmdb"))),
		maxmind.WithDomainDatabase(Must(os.Open("geoip/maxmind/testdata/GeoIP2-Domain-Test.mmdb"))),
		maxmind.WithConnectionTypeDatabase(Must(os.Open("geoip/maxmind/testdata/GeoIP2-Connection-Type-Test.mmdb"))),
		maxmind.WithEnterpriseDatabase(Must(os.Open("geoip/maxmind/testdata/GeoIP2-Enterprise-Test.mmdb"))),
	)
	require.NoError(t, err)

	resolver, err := networkd.NewCachedIPResolver(inner, 128)
	require.NoError(t, err)

	geoip := resolver.Resolve("81.2.69.142")
	require.NotEmpty(t, geoip.Timezone)

	require.Equal(t, 1, resolver.Size())
}

// Must is a helper function that asserts that an error is nil and returns the reader.
func Must(reader io.ReadCloser, err error) io.Reader {
	if err != nil {
		log.Fatalf("%v: unexpected error", err)
	}

	defer reader.Close()

	buffer := &bytes.Buffer{}
	_, err = io.Copy(buffer, reader)
	if err != nil {
		log.Fatalf("%v: failed to read from reader", err)
	}

	return buffer
}
