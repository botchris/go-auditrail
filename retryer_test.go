package auditrail_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/botchris/auditrail"
	"github.com/brianvoe/gofakeit/v6"
)

func TestNewRetryerSinkBreaker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testRetryerStrategy(t, ctx, auditrail.NewBreakerStrategy(3, 10*time.Millisecond))
}

func TestRetryerExponentialBackoff(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testRetryerStrategy(t, ctx, auditrail.NewExponentialBackoff(auditrail.ExponentialBackoffConfig{
		Base:   time.Millisecond,
		Factor: time.Millisecond,
		Max:    time.Millisecond * 5,
	}))
}

func TestRetryerExponentialBackoffV2(t *testing.T) {
	config := auditrail.DefaultExponentialBackoffConfig
	strategy := auditrail.NewExponentialBackoff(config)
	backoff := strategy.Proceed(nil)

	if backoff != 0 {
		t.Errorf("untouched backoff should be zero-wait: %v != 0", backoff)
	}

	expected := config.Base + config.Factor

	for i := 1; i <= 10; i++ {
		if strategy.Failure(nil, nil) {
			t.Errorf("no facilities for dropping messages in ExponentialBackoffStrategy")
		}

		for j := 0; j < 1000; j++ {
			// sample this several thousand times.
			if bo := strategy.Proceed(nil); bo > expected {
				t.Fatalf("expected must be bounded by %v after %v failures: %v", expected, i, bo)
			}
		}

		expected = config.Base + config.Factor*time.Duration(1<<uint64(i))
		if expected > config.Max {
			expected = config.Max
		}
	}

	strategy.Success(nil) // recovery!

	backoff = strategy.Proceed(nil)
	if backoff != 0 {
		t.Errorf("should have recovered: %v != 0", backoff)
	}
}

func testRetryerStrategy(t *testing.T, ctx context.Context, strategy auditrail.RetryStrategy) {
	const nm = 100

	tl := &auditrail.MemoryLogger{}

	// Make a sync that fails most of the time, ensuring that all the messages
	// make it through.
	flaky := &flakyLogger{
		rate:   1.0, // start out always failing.
		Logger: tl,
	}

	s := auditrail.NewRetryer(flaky, auditrail.WithRetryStrategy(strategy))

	var wg sync.WaitGroup

	for i := 1; i <= nm; i++ {
		entry := auditrail.NewEntry(
			gofakeit.Username(),
			gofakeit.VerbAction(),
			gofakeit.AppName(),
		)

		// Above 50, set the failure rate lower
		if i > 50 {
			flaky.mu.Lock()
			flaky.rate = 0.9
			flaky.mu.Unlock()
		}

		wg.Add(1)

		go func(e *auditrail.Entry) {
			defer wg.Done()

			if err := s.Log(ctx, e); err != nil {
				t.Errorf("error writing message: %v", err)
			}
		}(entry)
	}

	wg.Wait()
	checkClose(t, ctx, s)
}

type flakyLogger struct {
	auditrail.Logger
	rate float64
	mu   sync.Mutex
}

func (fs *flakyLogger) Log(ctx context.Context, e *auditrail.Entry) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if rand.Float64() < fs.rate {
		return fmt.Errorf("error logging entry: %v", e)
	}

	return fs.Logger.Log(ctx, e)
}

func checkClose(t *testing.T, ctx context.Context, trail auditrail.Closable) {
	if err := trail.Close(); err != nil {
		t.Fatalf("unexpected error closing: %v", err)
	}

	// second close should not crash and should not return error.
	if err := trail.Close(); err != nil {
		t.Fatalf("unexpected error on double close: %v", err)
	}

	var fail *auditrail.Entry

	// Write after closed should be an error
	if err := trail.Log(ctx, fail); err == nil {
		t.Fatalf("write after closed did not have an error")
	} else if !errors.Is(err, auditrail.ErrTrailClosed) {
		t.Fatalf("error should be ErrSinkClosed")
	}
}
