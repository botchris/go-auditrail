package auditrail

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch"
)

type elasticLogger struct {
	index  string
	client *elasticsearch.Client
}

// NewElasticLogger creates a new ElasticSearch logger.
func NewElasticLogger(index string, client *elasticsearch.Client) Logger {
	return &elasticLogger{
		index:  index,
		client: client,
	}
}

func (e *elasticLogger) Log(ctx context.Context, entry *Entry) error {
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
