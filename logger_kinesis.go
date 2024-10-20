package auditrail

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

// KinesisAPI captures the kinesis client part that we need.
type KinesisAPI interface {
	PutRecord(ctx context.Context, params *kinesis.PutRecordInput, optFns ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error)
}

type kinesisLogger struct {
	client     KinesisAPI
	streamName string
}

// NewKinesisLogger builds a new logger that writes log entries to a Kinesis
// stream as JSON objects separated by newlines.
func NewKinesisLogger(client KinesisAPI, streamName string) (Logger, error) {
	return &kinesisLogger{
		client:     client,
		streamName: streamName,
	}, nil
}

func (l *kinesisLogger) Log(ctx context.Context, entry *Entry) error {
	log, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = l.client.PutRecord(ctx, &kinesis.PutRecordInput{
		Data:         append(log, '\n'),
		PartitionKey: aws.String(entry.Module),
		StreamName:   &l.streamName,
	})

	return err
}
