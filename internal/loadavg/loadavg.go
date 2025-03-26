package loadavg

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"otus/internal/myerr"
	"otus/internal/storage"
	"sync"
	"sync/atomic"
	"time"
)

const (
	fileName = "/proc/loadavg"
	chSize   = 5
)

var (
	dataMon    *storage.Storage[AvgStat]
	chToParser chan []byte
	Working    atomic.Bool
)

type (
	AvgStat struct {
		One     float32
		Five    float32
		Fifteen float32
	}
)

func Start(ctx context.Context, wgGlobal *sync.WaitGroup) error {
	slog.Info("Start Load AVG collector")

	dataMon = storage.New[AvgStat]()
	chToParser = make(chan []byte, chSize)

	go parser()

	wgGlobal.Add(1)
	go probber(ctx, wgGlobal)

	return nil
}

func probber(ctx context.Context, wgGlobal *sync.WaitGroup) {
	defer wgGlobal.Done()
	defer slog.Info("Load AVG collector stopped")
	defer close(chToParser)

	// Признак работы сборщика данных
	Working.Store(true)
	defer Working.Store(false)

	// Используем time.Ticker для точного периода в 1 секунду
	// Исключаем накапливающуюся ошибку которая возникает при использовании time.After в цикле
	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := getData(); err != nil {
				slog.Error("Load AVG", "error read data from "+fileName, err)
				// process.Stop() // Останавливаем работу всего сервера или только данного параметра? Если всего сервера - снять комментарий
				return
			}
		}
	}
}

func getData() error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	chToParser <- data

	return nil
}

func parser() {
	for data := range chToParser {
		var s AvgStat

		_, err := fmt.Sscanf(string(data), "%f %f %f", &s.One, &s.Five, &s.Fifteen)
		if err != nil {
			slog.Error("Load AVG", "sscanf error", err)
		}
		dataMon.Add(s)
	}
}

func GetAvg(m int) (AvgStat, error) {
	var result, r AvgStat
	var i, one, five, fifteen int

	data := dataMon.Get(m)
	if data == nil {
		return result, myerr.ErrEmpty
	}

	// Сразу переводим в int для исключени ошибки плавающей точки
	for i, r = range data {
		one += int(math.Round(float64(r.One) * 100))
		five += int(math.Round(float64(r.Five) * 100))
		fifteen += int(math.Round(float64(r.Fifteen) * 100))
	}

	i++
	result.One = float32(one/i) / 100 // Обрезаем 2 знака после запятой
	result.Five = float32(five/i) / 100
	result.Fifteen = float32(fifteen/i) / 100

	return result, nil
}
