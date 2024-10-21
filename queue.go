package auditrail

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"
)

// ErrQueueClosed is returned when the queue is closed.
var ErrQueueClosed = fmt.Errorf("queue is closed")

// DropHandlerFunc is a function that will be called when a message is dropped
// from the queue.
type DropHandlerFunc func(*Entry, error)

// QueueOption is a function that configures a queue.
type QueueOption func(options *queueOptions)

// WithTimeout controls the maximum amount of time a worker will wait for the target
// logger to process a message. If the timeout is exceeded, the message will be dropped.
// If the timeout is less than or equal to zero, it will be set to 3 seconds.
func WithTimeout(timeout time.Duration) QueueOption {
	return func(opts *queueOptions) {
		if timeout <= 0 {
			timeout = 3 * time.Second
		}

		opts.timeout = timeout
	}
}

// WithDropHandler sets a function that will be called when a message is dropped. The
// function is called with the dropped message and the error that caused the drop.
func WithDropHandler(handler DropHandlerFunc) QueueOption {
	return func(opts *queueOptions) {
		if handler == nil {
			handler = func(*Entry, error) {}
		}

		opts.dropHandling = handler
	}
}

// WithThroughput controls the number of concurrent workers that will process messages
// from the queue. If throughput is less than or equal to zero, it will be set to 1.
func WithThroughput(throughput int) QueueOption {
	return func(opts *queueOptions) {
		if throughput <= 0 {
			throughput = 1
		}

		opts.throughput = throughput
	}
}

var _ Logger = (*Queue)(nil)

type queueEnvelope struct {
	message *Entry
	closed  bool
}

type queueOptions struct {
	timeout      time.Duration
	dropHandling DropHandlerFunc
	throughput   int
}

var defaultQueueOptions = queueOptions{
	timeout:      3 * time.Second,
	dropHandling: func(*Entry, error) {},
	throughput:   1,
}

// Queue is a logger that buffers log entries and processes them asynchronously.
type Queue struct {
	dst     Logger
	opts    queueOptions
	list    *list.List
	cond    *sync.Cond
	mu      sync.Mutex
	once    sync.Once
	closed  chan struct{}
	closing bool
}

// NewQueue builds a new logger queue which provides a buffer for entries to be
// processed asynchronously.
//
// The throughput parameter controls the number of concurrent workers that will
// process messages from the queue. If throughput is less than or equal to zero,
// it will be set to 1.
//
// The timeout parameter controls the maximum amount of time a worker will wait
// for the target logger to process a message. If the timeout is exceeded, the
// message will be dropped.
//

func NewQueue(dst Logger, options ...QueueOption) *Queue {
	opts := defaultQueueOptions
	q := &Queue{
		dst:    dst,
		list:   list.New(),
		closed: make(chan struct{}),
	}

	for _, option := range options {
		option(&opts)
	}

	q.opts = opts
	q.cond = sync.NewCond(&q.mu)

	for i := 0; i < q.opts.throughput; i++ {
		go q.run()
	}

	return q
}

// Log writes the given log entry to the queue for asynchronous processing.
func (lq *Queue) Log(_ context.Context, entry *Entry) error {
	lq.mu.Lock()
	defer lq.mu.Unlock()

	if lq.IsClosed() {
		return fmt.Errorf("%w: queue is closed", ErrQueueClosed)
	}

	lq.list.PushBack(queueEnvelope{message: entry})
	lq.cond.Signal() // signal waiters

	return nil
}

// Close shutdown the logger queue.
func (lq *Queue) Close() error {
	lq.mu.Lock()
	defer lq.mu.Unlock()

	if lq.IsClosed() {
		return nil
	}

	// set closing flag
	lq.closing = true
	lq.cond.Signal() // signal flushes queue
	lq.cond.Wait()   // wait for signal from last flush

	lq.once.Do(func() {
		close(lq.closed)
	})

	return nil
}

// IsClosed returns true if the queue is closed.
func (lq *Queue) IsClosed() bool {
	select {
	case <-lq.closed:
		return true
	default:
		return false
	}
}

// run is the main goroutine to flush messages to the target logger.
func (lq *Queue) run() {
	baseCtx := context.Background()

	for {
		envelope := lq.next()
		if envelope.closed {
			return // queueClosed block means event queue is closed.
		}

		ctx, cancel := context.WithTimeout(baseCtx, lq.opts.timeout)
		if err := lq.dst.Log(ctx, envelope.message); err != nil {
			lq.opts.dropHandling(envelope.message, err)
		}

		cancel()
	}
}

// next encompasses the critical section of the run loop. When the queue is
// empty, it will block on the condition. If new data arrives, it will wake
// and return a block. When closed, queueClosed constant will be returned.
func (lq *Queue) next() queueEnvelope {
	lq.mu.Lock()
	defer lq.mu.Unlock()

	for lq.list.Len() < 1 {
		if lq.closing || lq.IsClosed() {
			lq.cond.Broadcast()

			return queueEnvelope{closed: true}
		}

		lq.cond.Wait()
	}

	front := lq.list.Front()
	block, ok := front.Value.(queueEnvelope)

	if !ok {
		return queueEnvelope{closed: true}
	}

	lq.list.Remove(front)

	return block
}
