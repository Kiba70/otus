package cpu

import (
	"sync"
	"testing"
	"time"

	"github.com/Kiba70/otus/internal/storage"
	"github.com/stretchr/testify/require"
)

//nolint:gofumpt
var textStat = []string{`cpu  242534 0 395510 44332628 107824 0 49626 0 0 0
cpu0 13860 0 35046 2751637 9042 0 33421 0 0 0
cpu1 11191 0 14954 2778885 14017 0 7652 0 0 0`,
	`cpu  254731 0 418608 47517457 115665 0 51849 0 0 0
cpu0 14521 0 37189 2949327 9585 0 34979 0 0 0
cpu1 11534 0 15701 2978368 15400 0 7893 0 0 0
cpu2 19133 0 33167 2953370 10461 0 3878 0 0 0`,
	`cpu  255177 0 419399 47660468 115855 0 51925 0 0 0
cpu0 14542 0 37290 2958187 9601 0 35035 0 0 0
cpu1 11543 0 15713 2987384 15402 0 7901 0 0 0`,
	`cpu  255600 0 420201 47803112 116062 0 52027 0 0 0
cpu0 14555 0 37359 2967067 9623 0 35109 0 0 0
cpu1 11560 0 15756 2996321 15415 0 7917 0 0 0
cpu2 19200 0 33292 2971065 10594 0 3886 0 0 0`,
	`cpu  256379 0 421423 47914326 118500 0 52201 0 0 0
cpu0 14598 0 37472 2973942 9800 0 35210 0 0 0
cpu1 11589 0 15849 3003312 15531 0 7943 0 0 0
cpu2 19266 0 33390 2977919 10790 0 3905 0 0 0`,
}

var mux sync.Mutex

func TestCPU(t *testing.T) {
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	// Блок должен выполняться как единое целое
	// Тест хотя выполняется как unit test, по факту является интеграционным по всему модулю CPU
	// (парсинг, вычисление, хранение и выбоорка средних значений)
	mux.Lock()
	defer mux.Unlock()

	dataMon = storage.New[cpuStatInternal](100)
	chToParser = make(chan string, chSize)
	defer close(chToParser)
	chToCalculator = make(chan cpuStatInternal, chSize) // Закрывается в парсере

	go parser()
	go calculator()

	// Готовим данные
	// Производится обработка в parser, вычисление в calculator и занесение в storage
	for _, s := range textStat {
		chToParser <- s
	}

	t.Run("Производим получение усреднённых данных для передачи клиенту", func(t *testing.T) {
		time.Sleep(time.Millisecond * 10)

		require.Equal(t, 4, len(dataMon.Get(4)))
		g, err := GetAvg(4)
		require.Nil(t, err)
		require.Equal(t, CPUStat{
			User:   0.41,
			System: 0.71,
			Idle:   98.85,
		}, g)
	})
}
