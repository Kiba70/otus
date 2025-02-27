package loadavg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"otus/internal/myerr"
	"otus/internal/process"
	"otus/internal/storage"
	"sync"
	"time"
)

const (
	fileName = "/proc/loadavg"
	chSize   = 5
)

var (
	dataMon    *storage.Storage[AvgStat]
	chToParser = make(chan []byte, chSize)
)

type (
	AvgStat struct {
		One     float32
		Five    float32
		Fifteen float32
	}
)

func Start(ctx context.Context, wgGlobal *sync.WaitGroup) error {
	dataMon = storage.New[AvgStat]()

	slog.Debug("Loalavg Start")

	go parser()

	wgGlobal.Add(1)
	go probber(ctx, wgGlobal)

	return nil
}

func probber(ctx context.Context, wgGlobal *sync.WaitGroup) {
	defer wgGlobal.Done()
	defer close(chToParser)
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
				slog.Error("Load AVG", "error", err)
				process.Stop()
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
	var i int

	data := dataMon.Get(m)
	if data == nil {
		return result, myerr.ErrEmpty
	}

	for i, r = range data {
		result.One += r.One
		result.Five += r.Five
		result.Fifteen += r.Fifteen
	}

	result.One = float32(int(result.One*100)/i) / 100 // Обрезаем 2 знака после запятой
	result.Five = float32(int(result.Five*100)/i) / 100
	result.Fifteen = float32(int(result.Fifteen*100)/i) / 100

	return result, nil
}
