package auditrail

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// RetryerOption is a function that configures a retryer.
type RetryerOption func(options *retryer)

// RetryStrategy defines a strategy for retrying trail writes.
//
// All methods should be goroutine safe.
type RetryStrategy interface {
	// Proceed is called before every message send. If proceed returns a
	// positive, non-zero integer, the retryer will back off by the provided
	// duration.
	//
	// A message is provided, by may be ignored.
	Proceed(*Entry) time.Duration

	// Failure reports a failure to the strategy. If this method returns true,
	// the message should be dropped.
	Failure(*Entry, error) bool

	// Success should be called when a message is sent successfully.
	Success(*Entry)
}

type retryer struct {
	dst          Logger
	strategy     RetryStrategy
	dropHandling DropHandlerFunc
	closed       bool
	closedChan   chan struct{}
	mu           sync.RWMutex
}

// NewRetryer creates a new retryer that will retry failed log writes using the
// provided strategy.
func NewRetryer(dst Logger, opts ...RetryerOption) Logger {
	r := &retryer{
		dst:          dst,
		strategy:     NewExponentialBackoff(DefaultExponentialBackoffConfig),
		dropHandling: func(entry *Entry, err error) {},
		closedChan:   make(chan struct{}),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *retryer) Log(ctx context.Context, entry *Entry) error {
retry:
	r.mu.RLock()

	if r.closed {
		r.mu.RUnlock()

		return fmt.Errorf("%w: retriyer could not log the given entry", ErrTrailClosed)
	}

	r.mu.RUnlock()

	if backoff := r.strategy.Proceed(entry); backoff > 0 {
		select {
		case <-time.After(backoff):
			// TODO: This branch holds up the next try. Before, we
			// would simply break to the "retry" label and then possibly wait
			// again. However, this requires all retry strategies to have a
			// large probability of probing the sync for success, rather than
			// just backing off and sending the request.
		case <-r.Closed():
			return ErrTrailClosed
		}
	}

	if err := r.dst.Log(ctx, entry); err != nil {
		if errors.Is(err, ErrTrailClosed) {
			// terminal!
			return err
		}

		if r.strategy.Failure(entry, err) {
			r.dropHandling(entry, err)

			return nil
		}

		goto retry
	}

	r.strategy.Success(entry)

	return nil
}

func (r *retryer) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}

	if err := r.dst.Close(); err != nil {
		return fmt.Errorf("%w: retrying sink could not close underlying sink", err)
	}

	r.closed = true

	close(r.closedChan)

	return nil
}

func (r *retryer) Closed() <-chan struct{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.closedChan
}

func (r *retryer) IsClosed() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.closed
}

// WithRetryStrategy configures the retry strategy for the retryer. If strategy is
// nil, a default exponential backoff strategy is used.
func WithRetryStrategy(strategy RetryStrategy) RetryerOption {
	return func(options *retryer) {
		if strategy == nil {
			strategy = NewExponentialBackoff(DefaultExponentialBackoffConfig)
		}

		options.strategy = strategy
	}
}

// WithRetryDropHandler configures the drop handler for the retryer. If handler is
// nil, a no-op handler is used.
func WithRetryDropHandler(handler DropHandlerFunc) RetryerOption {
	return func(options *retryer) {
		if handler == nil {
			handler = func(*Entry, error) {}
		}

		options.dropHandling = handler
	}
}

// NewBreakerStrategy returns a breaker that will backoff after the threshold has been
// tripped. A Breaker is thread safe and may be shared by many goroutines.
func NewBreakerStrategy(threshold int, backoff time.Duration) RetryStrategy {
	return &breakerStrategy{
		threshold: threshold,
		backoff:   backoff,
	}
}

// Breaker implements a circuit breaker retry strategy.
//
// The current implementation never drops messages.
type breakerStrategy struct {
	threshold int
	recent    int
	last      time.Time
	backoff   time.Duration // time after which we retry after failure.
	mu        sync.Mutex
}

// Proceed checks the failures against the threshold.
func (b *breakerStrategy) Proceed(*Entry) time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.recent < b.threshold {
		return 0
	}

	return time.Until(b.last.Add(b.backoff))
}

// Success resets the breaker.
func (b *breakerStrategy) Success(*Entry) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.recent = 0
	b.last = time.Time{}
}

// Failure records the failure and latest failure time.
func (b *breakerStrategy) Failure(*Entry, error) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.recent++
	b.last = time.Now().UTC()

	return false // never drop messages.
}

// ExponentialBackoffConfig configures backoff parameters.
//
// Note that these parameters operate on the upper bound for choosing a random
// value. For example, at Base=1s, a random value in [0,1s) will be chosen for
// the backoff value.
type ExponentialBackoffConfig struct {
	// Base is the minimum bound for backing off after failure.
	Base time.Duration

	// Factor sets the amount of time by which the backoff grows with each
	// failure.
	Factor time.Duration

	// Max is the absolute maximum bound for a single backoff.
	Max time.Duration
}

// DefaultExponentialBackoffConfig provides a default configuration for
// exponential backoff.
var DefaultExponentialBackoffConfig = ExponentialBackoffConfig{
	Base:   time.Second,
	Factor: time.Second,
	Max:    20 * time.Second,
}

// NewExponentialBackoff returns an exponential backoff strategy with the
// desired config. If config is nil, the default is returned.
func NewExponentialBackoff(config ExponentialBackoffConfig) RetryStrategy {
	return &exponentialBackoffStrategy{
		config: config,
	}
}

// exponentialBackoffStrategy implements random backoff with exponentially increasing
// bounds as the number consecutive failures increase.
type exponentialBackoffStrategy struct {
	failures uint64 // consecutive failure counter (needs to be 64-bit aligned)
	config   ExponentialBackoffConfig
}

// Proceed returns the next randomly bound exponential backoff time.
func (b *exponentialBackoffStrategy) Proceed(*Entry) time.Duration {
	return b.backoff(atomic.LoadUint64(&b.failures))
}

// Success resets the failures counter.
func (b *exponentialBackoffStrategy) Success(*Entry) {
	atomic.StoreUint64(&b.failures, 0)
}

// Failure increments the failure counter.
func (b *exponentialBackoffStrategy) Failure(*Entry, error) bool {
	atomic.AddUint64(&b.failures, 1)

	return false
}

// backoff calculates the amount of time to wait based on the number of
// consecutive failures.
func (b *exponentialBackoffStrategy) backoff(failures uint64) time.Duration {
	if failures <= 0 {
		// proceed normally when there are no failures.
		return 0
	}

	factor := b.config.Factor
	if factor <= 0 {
		factor = DefaultExponentialBackoffConfig.Factor
	}

	backoff := b.config.Base + factor*time.Duration(1<<(failures-1))

	mx := b.config.Max
	if mx <= 0 {
		mx = DefaultExponentialBackoffConfig.Max
	}

	if backoff > mx || backoff < 0 {
		backoff = mx
	}

	// Choose a uniformly distributed value from [0, backoff).
	return time.Duration(rand.Int63n(int64(backoff)))
}
