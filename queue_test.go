package auditrail_test

import (
	"context"
	"testing"
	"time"

	"github.com/botchris/auditrail"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Run("GIVEN a queue instance with a test logger", func(t *testing.T) {
		dst := &auditrail.MemoryLogger{}
		queue := auditrail.NewQueue(dst, auditrail.WithThroughput(10))

		t.Run("WHEN logging 10000 entries in parallel", func(t *testing.T) {
			n := 10000
			syncer := make(chan struct{})

			for i := 0; i < n; i++ {
				go func() {
					<-syncer

					entry := auditrail.NewEntry(
						gofakeit.Username(),
						gofakeit.VerbAction(),
						gofakeit.AppName(),
					)

					require.NoError(t, queue.Log(ctx, entry))
				}()
			}

			close(syncer)

			t.Run("THEN destination logger eventually receives all the entries", func(t *testing.T) {
				require.Eventually(t, func() bool {
					return dst.Size() == n
				}, 10*time.Second, 200*time.Millisecond)
			})
		})
	})
}
