package auditrail_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/botchris/go-auditrail"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestKinesisLogger(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("GIVEN a kinesis logger", func(t *testing.T) {
		api := &mockKinesisAPI{}
		streamName := gofakeit.UUID()
		logger, err := auditrail.NewKinesisLogger(api, streamName)

		require.NoError(t, err)

		t.Run("WHEN logging 100 entries", func(t *testing.T) {
			n := 100

			for i := 0; i < n; i++ {
				entry := auditrail.NewEntry(
					gofakeit.UUID(),
					gofakeit.Emoji(),
					gofakeit.Emoji(),
				).
					WithAuthMethod(gofakeit.UUID()).
					WithOccurredAt(time.Now())

				require.NoError(t, logger.Log(ctx, entry))
			}

			t.Run("THEN entries should be sent to the kinesis stream in JSON format", func(t *testing.T) {
				require.Len(t, api.putCalls, n)

				for i := range api.putCalls {
					require.Equal(t, streamName, *api.putCalls[i].StreamName)
					require.NotEmpty(t, api.putCalls[i].PartitionKey)
					require.NotEmpty(t, api.putCalls[i].Data)
					require.True(t, json.Valid(api.putCalls[i].Data))
				}
			})
		})
	})
}

type mockKinesisAPI struct {
	putCalls []*kinesis.PutRecordInput
}

func (m *mockKinesisAPI) PutRecord(_ context.Context, params *kinesis.PutRecordInput, _ ...func(*kinesis.Options)) (*kinesis.PutRecordOutput, error) {
	m.putCalls = append(m.putCalls, params)

	return &kinesis.PutRecordOutput{}, nil
}
