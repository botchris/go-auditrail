package maxmind

import "github.com/oschwald/geoip2-golang"

// Option is a functional option for configuring the Maxmind GeoIP resolver.
type Option func(m *maxmindGeoIPResolver) error

// WithASNDatabase loads the ASN database from the given path.
func WithASNDatabase(path string) Option {
	return func(m *maxmindGeoIPResolver) error {
		db, err := geoip2.Open(path)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &asnResolver{db: db})

		return nil
	}
}

// WithCountryDatabase loads the country database from the given path.
func WithCountryDatabase(path string) Option {
	return func(m *maxmindGeoIPResolver) error {
		db, err := geoip2.Open(path)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &countryResolver{db: db})

		return nil
	}
}

// WithCityDatabase loads the city database from the given path.
func WithCityDatabase(path string) Option {
	return func(m *maxmindGeoIPResolver) error {
		db, err := geoip2.Open(path)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &cityResolver{db: db})

		return nil
	}
}

// WithISPDatabase loads the ISP database from the given path.
func WithISPDatabase(path string) Option {
	return func(m *maxmindGeoIPResolver) error {
		db, err := geoip2.Open(path)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &ispResolver{db: db})

		return nil
	}
}

// WithDomainDatabase loads the Domain database from the given path.
func WithDomainDatabase(path string) Option {
	return func(m *maxmindGeoIPResolver) error {
		db, err := geoip2.Open(path)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &domainResolver{db: db})

		return nil
	}
}

// WithConnectionTypeDatabase loads the Connection-Type database from the given
// path.
func WithConnectionTypeDatabase(path string) Option {
	return func(m *maxmindGeoIPResolver) error {
		db, err := geoip2.Open(path)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &connectionTypeResolver{db: db})

		return nil
	}
}

// WithEnterpriseDatabase loads the Enterprise database from the given path.
func WithEnterpriseDatabase(path string) Option {
	return func(m *maxmindGeoIPResolver) error {
		db, err := geoip2.Open(path)
		if err != nil {
			return err
		}

		m.handlers = append(m.handlers, &enterpriseResolver{db: db})

		return nil
	}
}
