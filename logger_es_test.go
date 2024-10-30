package auditrail_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/botchris/go-auditrail"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/elastic/go-elasticsearch"
	"github.com/stretchr/testify/require"
)

func TestElasticLogger(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("GIVEN a elastic logger", func(t *testing.T) {
		jsonCalls := make([]string, 0)
		cfg := elasticsearch.Config{
			Transport: &fakeTransport{
				RoundTripFn: func(r *http.Request) (*http.Response, error) {
					buf := new(strings.Builder)

					if _, err := io.Copy(buf, r.Body); err != nil {
						return nil, err
					}

					jsonCalls = append(jsonCalls, buf.String())

					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader("MOCK")),
					}, nil
				},
			},
		}
		client, err := elasticsearch.NewClient(cfg)
		require.NoError(t, err)

		indexName := gofakeit.UUID()
		logger := auditrail.NewElasticLogger(indexName, client)

		t.Run("WHEN logging 100 entries", func(t *testing.T) {
			n := 100

			for i := 0; i < n; i++ {
				entry := auditrail.NewEntry(
					gofakeit.Username(),
					gofakeit.VerbAction(),
					gofakeit.IPv4Address(),
				)

				require.NoError(t, logger.Log(ctx, entry))
			}

			t.Run("THEN entries are successfully logged AND can be unmarshalled back", func(t *testing.T) {
				require.Len(t, jsonCalls, n)

				for i := range jsonCalls {
					require.NotEmpty(t, jsonCalls[i])
					require.True(t, json.Valid([]byte(jsonCalls[i])))

					var entry auditrail.Entry

					require.NoError(t, json.Unmarshal([]byte(jsonCalls[i]), &entry))
				}
			})
		})
	})
}

type fakeTransport struct {
	RoundTripFn func(r *http.Request) (*http.Response, error)
}

func (t *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripFn(req)
}
