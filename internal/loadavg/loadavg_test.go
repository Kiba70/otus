package loadavg

import (
	"otus/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("Test parser prepare", func(t *testing.T) {
		dataMon = storage.New[AvgStat]()
		go parser()
		for range 100 {
			chToParser <- []byte("0.16 0.21 0.21 1/575 139321")
		}
	})

	t.Run("Test parser getting", func(t *testing.T) {
		time.Sleep(time.Second)

		require.Equal(t, 100, len(dataMon.Get(100)))
		g, err := GetAvg(100)
		require.Nil(t, err)
		require.Equal(t, AvgStat{
			one:     0.16,
			five:    0.21,
			fifteen: 0.21,
		}, g)

		close(chToParser)
	})
}
