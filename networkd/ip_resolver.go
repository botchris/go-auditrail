package networkd

import (
	"github.com/golang/groupcache/singleflight"
	lru "github.com/hashicorp/golang-lru/v2"
)

var _ IPResolver = (*CachedIPResolver)(nil)

// IPResolver is capable of resolving geoip information from a given IP address.
//
// Popular implementations may include Maxmind GeoIP and IP2Location.
//
// NOTE: enriching [GeoIP] information is an expensive operation and should be
// used judiciously. It is recommended to only enrich the IP address and then
// use post-processing to enrich the log entries with GeoIP information. For
// example: https://www.elastic.co/guide/en/elasticsearch/reference/current/geoip-processor.html
type IPResolver interface {
	Resolve(ip string) GeoIP
}

// CachedIPResolver is an IPResolver that caches the results of the inner
// IPResolver so that subsequent calls to Resolve for the same IP address are
// faster.
type CachedIPResolver struct {
	inner IPResolver
	cache *lru.Cache[string, GeoIP]
	sf    *singleflight.Group
}

// NewCachedIPResolver returns a new [IPResolver] that caches the results of the
// inner [IPResolver] so that subsequent calls to Resolve are faster.
//
// The size parameter specifies the maximum number of entries in the cache.
func NewCachedIPResolver(inner IPResolver, size int) (*CachedIPResolver, error) {
	cache, err := lru.New[string, GeoIP](size)
	if err != nil {
		return nil, err
	}

	return &CachedIPResolver{
		inner: inner,
		cache: cache,
		sf:    &singleflight.Group{},
	}, nil
}

// Resolve resolves the geoip information for the given IP address with cache
// support.
func (c *CachedIPResolver) Resolve(ip string) GeoIP {
	if gd, ok := c.cache.Get(ip); ok {
		return gd
	}

	value, err := c.sf.Do(ip, func() (interface{}, error) {
		d := c.inner.Resolve(ip)
		c.cache.Add(ip, d)

		return d, nil
	})

	if err != nil {
		return GeoIP{}
	}

	return value.(GeoIP)
}

// Remove invalidates the cache for the given IP address.
func (c *CachedIPResolver) Remove(ip string) {
	c.cache.Remove(ip)
}

// Purge purges the cache.
func (c *CachedIPResolver) Purge() {
	c.cache.Purge()
}

// Size returns the number of entries in the cache.
func (c *CachedIPResolver) Size() int {
	return c.cache.Len()
}

// Contains returns true if the cache contains the given IP address.
func (c *CachedIPResolver) Contains(ip string) bool {
	return c.cache.Contains(ip)
}
