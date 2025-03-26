////go:build integration

package loadavg

import (
	"otus/internal/storage"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var mux sync.Mutex

func TestLoadAVG(t *testing.T) {
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	// Блок должен выполняться как единое целое
	// Тест хотя выполняется как unit test, по факту является интеграционным по всему модулю LoadAVG
	// (парсинг, хранение и выбоорка средних значений)
	mux.Lock()
	defer mux.Unlock()

	dataMon = storage.New[AvgStat](100)
	chToParser = make(chan []byte, chSize)
	defer close(chToParser)

	go parser()
	t.Run("Готовим данные", func(t *testing.T) {

		for range 100 {
			chToParser <- []byte("0.16 0.21 0.21 1/575 139321")
		}
	})

	t.Run("Производим получение усреднённых данных для передачи клиенту", func(t *testing.T) {
		time.Sleep(time.Millisecond * 10)

		require.Equal(t, 100, len(dataMon.Get(100)))
		g, err := GetAvg(100)
		require.Nil(t, err)
		require.Equal(t, AvgStat{
			One:     0.16,
			Five:    0.21,
			Fifteen: 0.21,
		}, g)
	})
}
