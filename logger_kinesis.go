package auditrail

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

// KinesisAPI captures the kinesis client part that we need.
type KinesisAPI interface {
	PutRecord(ctx context.Context, params *kinesis.PutRecordInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error)
}

type kinesisLogger struct {
	client       KinesisAPI
	streamName   string
	closed       bool
	closeChannel chan struct{}
	mu           sync.RWMutex
}

// NewKinesisLogger builds a new logger that writes log entries to a Kinesis
// stream as JSON objects separated by newlines.
func NewKinesisLogger(client KinesisAPI, streamName string) (Logger, error) {
	return &kinesisLogger{
		client:       client,
		streamName:   streamName,
		closeChannel: make(chan struct{}),
	}, nil
}

func (l *kinesisLogger) Log(ctx context.Context, entry *Entry) error {
	l.mu.RLock()
	closed := l.closed
	l.mu.RUnlock()

	if closed {
		return ErrTrailClosed
	}

	log, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = l.client.PutRecord(ctx, &kinesis.PutRecordInput{
		Data:         append(log, '\n'),
		PartitionKey: aws.String(entry.GetModule()),
		StreamName:   &l.streamName,
	})

	return err
}

func (l *kinesisLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return nil
	}

	l.closed = true

	close(l.closeChannel)

	return nil
}

func (l *kinesisLogger) Closed() <-chan struct{} {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.closeChannel
}

func (l *kinesisLogger) IsClosed() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.closed
}
