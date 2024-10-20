package auditrail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/botchris/auditrail"
	"github.com/botchris/auditrail/httpd"
	"github.com/botchris/auditrail/networkd"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestFilepathLogger(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("GIVEN a file audit logger for an empty file", func(t *testing.T) {
		path := t.TempDir() + "/audit.log"
		logger, fErr := auditrail.NewFilePathLogger(path)

		require.NoError(t, fErr)

		fi, err := os.Stat(path)
		require.NoError(t, err)
		require.Zero(t, fi.Size())

		t.Run("WHEN logging an entry", func(t *testing.T) {
			userID := gofakeit.UUID()
			entry := auditrail.NewEntry(
				userID,
				"order_create",
				"ordering_service",
			).
				AppendDetails("http", httpd.Details{
					Method:     "POST",
					StatusCode: "201",
					UserAgent:  "curl/7.68.0",
					URL: httpd.URL{
						Host: "api.example.com",
						Path: "/orders",
					},
				})

			lErr := logger.Log(ctx, entry)
			require.NoError(t, lErr)

			t.Run("THEN the file should contain the entry", func(t *testing.T) {
				fi, err = os.Stat(path)
				require.NoError(t, err)
				require.NotZero(t, fi.Size())

				b, err := os.ReadFile(path)
				require.NoError(t, err)

				require.Contains(t, string(b), userID)
				require.Contains(t, string(b), "POST")
				require.Contains(t, string(b), "curl/7.68.0")
				require.Contains(t, string(b), "api.example.com")
			})
		})
	})
}

func TestFDLogger(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("GIVEN a FD pointing to stdout logger WHEN logging THEN entry is written", func(t *testing.T) {
		logger, fErr := auditrail.NewFileLogger(os.Stdout)
		require.NoError(t, fErr)

		userID := gofakeit.UUID()
		entry := auditrail.NewEntry(
			userID,
			"order_create",
			"ordering_service",
		).
			AppendDetails("http", httpd.Details{
				Method:     "POST",
				StatusCode: "201",
				UserAgent:  "curl/7.68.0",
				URL: httpd.URL{
					Host: "api.example.com",
					Path: "/orders",
				},
			})

		lErr := logger.Log(ctx, entry)
		require.NoError(t, lErr)
	})

	t.Run("GIVEN a FD pointing to stderr logger WHEN logging THEN entry is written", func(t *testing.T) {
		logger, fErr := auditrail.NewFileLogger(os.Stderr)
		require.NoError(t, fErr)

		userID := gofakeit.UUID()
		entry := auditrail.NewEntry(
			userID,
			"order_create",
			"ordering_service",
		).
			AppendDetails("http", httpd.Details{
				Method:     "POST",
				StatusCode: "201",
				UserAgent:  "curl/7.68.0",
				URL: httpd.URL{
					Host: "api.example.com",
					Path: "/orders",
				},
			})

		lErr := logger.Log(ctx, entry)
		require.NoError(t, lErr)
	})
}

func TestDecorators(t *testing.T) {
	mainCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("GIVEN a file audit logger for an empty file with HTTP, Network and User decorators AND a context with http details", func(t *testing.T) {
		httpDetails := httpd.Details{
			Method:     "GET",
			StatusCode: "201",
			UserAgent:  "curl/7.68.0",
			URL: httpd.URL{
				Host: "api.example.com",
				Path: "/orders",
			},
		}

		networkDetails := networkd.Details{
			Client: networkd.Client{
				IP: gofakeit.IPv4Address(),
			},
		}

		ctx := httpd.AddToContext(mainCtx, httpDetails)
		ctx = networkd.AddToContext(ctx, networkDetails)

		path := t.TempDir() + "/audit.log"
		logger, fErr := auditrail.NewFilePathLogger(path)

		require.NoError(t, fErr)

		logger = httpd.Decorator(
			networkd.Decorator(logger, nil),
		)

		fi, err := os.Stat(path)
		require.NoError(t, err)
		require.Zero(t, fi.Size())

		t.Run("WHEN logging an entry", func(t *testing.T) {
			userID := gofakeit.UUID()
			entry := auditrail.NewEntry(
				userID,
				"order_create",
				"ordering_service",
			)

			lErr := logger.Log(ctx, entry)
			require.NoError(t, lErr)

			t.Run("THEN the file should contain the entry and http details", func(t *testing.T) {
				fi, err = os.Stat(path)
				require.NoError(t, err)
				require.NotZero(t, fi.Size())

				b, err := os.ReadFile(path)
				require.NoError(t, err)

				require.Contains(t, string(b), userID)
				require.Contains(t, string(b), "GET")
				require.Contains(t, string(b), "curl/7.68.0")
				require.Contains(t, string(b), "api.example.com")
			})
		})
	})
}
