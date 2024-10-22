package auditrail_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/botchris/auditrail"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dst := auditrail.NewMemoryLogger()
	dl := &delayed{
		Logger: dst,
		delay:  time.Millisecond * 1,
	}

	queue := auditrail.NewQueue(dl, auditrail.WithQueueThroughput(3))

	time.Sleep(20 * time.Millisecond) // let's queue settle to wait condition.

	n := 1000
	wg := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			entry := auditrail.NewEntry(
				gofakeit.Username(),
				gofakeit.VerbAction(),
				gofakeit.AppName(),
			)

			require.NoError(t, queue.Log(ctx, entry))
		}()
	}

	wg.Wait()
	checkClose(t, ctx, queue)

	require.EqualValues(t, n, dst.Size())
	require.True(t, queue.IsClosed())
}

func TestQueueDrop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const nm = 10

	cc := &atomic.Int64{}
	eq := auditrail.NewQueue(
		&dropper{err: errors.New("dropped")},
		auditrail.WithQueueThroughput(1),
		auditrail.WithQueueDropHandler(func(*auditrail.Entry, error) { cc.Add(1) }),
	)

	time.Sleep(10 * time.Millisecond) // let's queue settle to wait condition.

	var wg sync.WaitGroup

	for i := 1; i <= nm; i++ {
		wg.Add(1)

		go func(m *auditrail.Entry) {
			defer wg.Done()

			if err := eq.Log(ctx, m); err != nil {
				t.Errorf("error writing message: %v", err)
			}
		}(auditrail.NewEntry(
			gofakeit.Username(),
			gofakeit.VerbAction(),
			gofakeit.AppName(),
		))
	}

	wg.Wait()
	checkClose(t, ctx, eq)
}

type dropper struct {
	auditrail.Logger
	err    error
	closed bool
	mu     sync.Mutex
}

func (d *dropper) Log(context.Context, *auditrail.Entry) error {
	return d.err
}

func (d *dropper) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.closed {
		return nil
	}

	d.closed = true

	return nil
}

type delayed struct {
	auditrail.Logger
	delay time.Duration
}

func (d *delayed) Log(ctx context.Context, e *auditrail.Entry) error {
	time.Sleep(d.delay)

	return d.Logger.Log(ctx, e)
}
