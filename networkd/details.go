package networkd

// Details represents the details of a network connection.
type Details struct {
	Client Client `json:"client"`
}

// Client capture details of the client connecting to the system.
type Client struct {
	IP    string `json:"ip"`
	GeoIP *GeoIP `json:"geoip,omitempty"`
}

// GeoIP capture details of the location of the connection.
type GeoIP struct {
	AS          AS          `json:"as"`
	Continent   Continent   `json:"continent"`
	Country     Country     `json:"country"`
	City        City        `json:"city"`
	Location    Location    `json:"location"`
	Subdivision Subdivision `json:"subdivision"`
	Timezone    string      `json:"timezone"`
}

// AS autonomous system details.
type AS struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Number string `json:"number"`
	Route  string `json:"route"`
	Type   string `json:"type"`
}

// Continent represents the continent from which the connection was originated.
type Continent struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Country represents the country from which the connection was originated.
type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// City represents the city from which the connection was originated.
type City struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Location represents GPS coordinates of the connection.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Subdivision represents a subdivision of a country.
type Subdivision struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
