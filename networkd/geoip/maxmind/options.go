package maxmind

import (
	"io"

	"github.com/oschwald/geoip2-golang"
)

// Option is a functional option for configuring the Maxmind GeoIP resolver.
type Option func(m *maxmindGeoIPResolver) error

// WithASNDatabase loads the ASN database from the given path.
func WithASNDatabase(reader io.Reader) Option {
	return func(m *maxmindGeoIPResolver) error {
		b, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		db, err := geoip2.FromBytes(b)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &asnResolver{db: db})

		return nil
	}
}

// WithCountryDatabase loads the country database from the given path.
func WithCountryDatabase(reader io.Reader) Option {
	return func(m *maxmindGeoIPResolver) error {
		b, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		db, err := geoip2.FromBytes(b)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &countryResolver{db: db})

		return nil
	}
}

// WithCityDatabase loads the city database from the given path.
func WithCityDatabase(reader io.Reader) Option {
	return func(m *maxmindGeoIPResolver) error {
		b, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		db, err := geoip2.FromBytes(b)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &cityResolver{db: db})

		return nil
	}
}

// WithISPDatabase loads the ISP database from the given path.
func WithISPDatabase(reader io.Reader) Option {
	return func(m *maxmindGeoIPResolver) error {
		b, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		db, err := geoip2.FromBytes(b)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &ispResolver{db: db})

		return nil
	}
}

// WithDomainDatabase loads the Domain database from the given path.
func WithDomainDatabase(reader io.Reader) Option {
	return func(m *maxmindGeoIPResolver) error {
		b, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		db, err := geoip2.FromBytes(b)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &domainResolver{db: db})

		return nil
	}
}

// WithConnectionTypeDatabase loads the Connection-Type database from the given
// path.
func WithConnectionTypeDatabase(reader io.Reader) Option {
	return func(m *maxmindGeoIPResolver) error {
		b, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		db, err := geoip2.FromBytes(b)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &connectionTypeResolver{db: db})

		return nil
	}
}

// WithEnterpriseDatabase loads the Enterprise database from the given path.
func WithEnterpriseDatabase(reader io.Reader) Option {
	return func(m *maxmindGeoIPResolver) error {
		b, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		db, err := geoip2.FromBytes(b)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &enterpriseResolver{db: db})

		return nil
	}
}
