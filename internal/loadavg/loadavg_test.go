//go:build integration

package loadavg

import (
	"log/slog"
	"otus/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegrate(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	t.Run("Запускаем parser и готовим данные", func(t *testing.T) {
		dataMon = storage.New[AvgStat](100)
		go parser()
		for range 100 {
			chToParser <- []byte("0.16 0.21 0.21 1/575 139321")
		}
	})

	t.Run("Производим получение усреднённых данных для передачи клиенту", func(t *testing.T) {
		time.Sleep(time.Second)

		require.Equal(t, 100, len(dataMon.Get(100)))
		g, err := GetAvg(100)
		require.Nil(t, err)
		require.Equal(t, AvgStat{
			One:     0.16,
			Five:    0.21,
			Fifteen: 0.21,
		}, g)

		// close(chToParser) Специально не закрываем канал
	})
}
