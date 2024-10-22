package auditrail

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"
)

// DropHandlerFunc is a function that will be called when a message is dropped
// from the queue.
type DropHandlerFunc func(*Entry, error)

// QueueOption is a function that configures a queue.
type QueueOption func(options *queueOptions)

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

type queue struct {
	dst     Logger
	opts    queueOptions
	list    *list.List
	cond    *sync.Cond
	once    sync.Once
	closed  chan struct{}
	closing bool
	mu      sync.Mutex
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
func NewQueue(dst Logger, options ...QueueOption) Closable {
	opts := defaultQueueOptions
	q := &queue{
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
func (q *queue) Log(_ context.Context, entry *Entry) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsClosed() {
		return fmt.Errorf("%w: queue is closed", ErrTrailClosed)
	}

	q.list.PushBack(queueEnvelope{message: entry})
	q.cond.Signal() // signal waiters

	return nil
}

// Close shutdown the logger queue.
func (q *queue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.IsClosed() {
		return nil
	}

	// set closing flag
	q.closing = true
	q.cond.Signal() // signal flushes queue
	q.cond.Wait()   // wait for signal from last flush

	q.once.Do(func() {
		close(q.closed)
	})

	return nil
}

func (q *queue) Closed() <-chan struct{} {
	return q.closed
}

// IsClosed returns true if the queue is closed.
func (q *queue) IsClosed() bool {
	select {
	case <-q.closed:
		return true
	default:
		return false
	}
}

// run is the main goroutine to flush messages to the target logger.
func (q *queue) run() {
	baseCtx := context.Background()

	for {
		envelope := q.next()
		if envelope.closed {
			return // queueClosed block means event queue is closed.
		}

		ctx, cancel := context.WithTimeout(baseCtx, q.opts.timeout)
		if err := q.dst.Log(ctx, envelope.message); err != nil {
			q.opts.dropHandling(envelope.message, err)
		}

		cancel()
	}
}

// next encompasses the critical section of the run loop. When the queue is
// empty, it will block on the condition. If new data arrives, it will wake
// and return a block. When closed, queueClosed constant will be returned.
func (q *queue) next() queueEnvelope {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.list.Len() < 1 {
		if q.closing || q.IsClosed() {
			q.cond.Broadcast()

			return queueEnvelope{closed: true}
		}

		q.cond.Wait()
	}

	front := q.list.Front()
	block, ok := front.Value.(queueEnvelope)

	if !ok {
		return queueEnvelope{closed: true}
	}

	q.list.Remove(front)

	return block
}

// WithQueueTimeout controls the maximum amount of time a worker will wait for the target
// logger to process a message. If the timeout is exceeded, the message will be dropped.
// If the timeout is less than or equal to zero, it will be set to 3 seconds.
func WithQueueTimeout(timeout time.Duration) QueueOption {
	return func(opts *queueOptions) {
		if timeout <= 0 {
			timeout = 3 * time.Second
		}

		opts.timeout = timeout
	}
}

// WithQueueDropHandler sets a function that will be called when a message is dropped. The
// function is called with the dropped message and the error that caused the drop.
func WithQueueDropHandler(handler DropHandlerFunc) QueueOption {
	return func(opts *queueOptions) {
		if handler == nil {
			handler = func(*Entry, error) {}
		}

		opts.dropHandling = handler
	}
}

// WithQueueThroughput controls the number of concurrent workers that will process messages
// from the queue. If throughput is less than or equal to zero, it will be set to 1.
func WithQueueThroughput(throughput int) QueueOption {
	return func(opts *queueOptions) {
		if throughput <= 0 {
			throughput = 1
		}

		opts.throughput = throughput
	}
}
