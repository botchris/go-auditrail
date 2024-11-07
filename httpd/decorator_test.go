package httpd_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/botchris/go-auditrail"
	"github.com/botchris/go-auditrail/httpd"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestDecorator(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx = httpd.AddToContext(ctx, httpd.Details{
		Method:     http.MethodPost,
		StatusCode: "200 OK",
		UserAgent:  gofakeit.UserAgent(),
		URL: httpd.URL{
			Host: gofakeit.DomainName(),
			Path: gofakeit.URL(),
		},
	})

	logger := auditrail.NewMemoryLogger()
	d := httpd.Decorator(logger)

	for i := 0; i < 100; i++ {
		event := auditrail.NewEntry(gofakeit.Username(), gofakeit.VerbAction(), gofakeit.AppName())
		require.NoError(t, d.Log(ctx, event))
	}

	require.Equal(t, 100, logger.Size())

	for _, log := range logger.Trail() {
		details := log.GetDetails()

		require.NotEmpty(t, details["http"])
		require.IsType(t, httpd.Details{}, details["http"])

		s, err := json.Marshal(log)

		require.NoError(t, err)
		require.NotEmpty(t, string(s))
		require.NotEqual(t, "{}", string(s))
	}
}
