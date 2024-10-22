package auditrail

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch"
)

type elasticLogger struct {
	index        string
	client       *elasticsearch.Client
	closed       bool
	closeChannel chan struct{}
	mu           sync.RWMutex
}

// NewElasticLogger creates a new ElasticSearch logger.
func NewElasticLogger(index string, client *elasticsearch.Client) Logger {
	return &elasticLogger{
		index:        index,
		client:       client,
		closeChannel: make(chan struct{}),
	}
}

func (e *elasticLogger) Log(ctx context.Context, entry *Entry) error {
	e.mu.RLock()
	closed := e.closed
	e.mu.RUnlock()

	if closed {
		return ErrTrailClosed
	}

	body, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = e.client.Index(
		e.index,
		strings.NewReader(string(body)),
		e.client.Index.WithContext(ctx),
		e.client.Index.WithDocumentID(entry.IdempotencyID),
	)

	return err
}

func (e *elasticLogger) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil
	}

	e.closed = true

	close(e.closeChannel)

	return nil
}

func (e *elasticLogger) Closed() <-chan struct{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.closeChannel
}

func (e *elasticLogger) IsClosed() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.closed
}
