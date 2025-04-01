package cpu

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Kiba70/otus/internal/myerr"
	"github.com/Kiba70/otus/internal/storage"
)

const (
	fileName = "/proc/stat"
	chSize   = 5
)

var (
	dataMon        *storage.Storage[cpuStatInternal]
	chToParser     chan string
	chToCalculator chan cpuStatInternal
	Working        atomic.Bool
)

type (
	CPUStat struct { //nolint:revive
		User   float32
		System float32
		Idle   float32
	}
	cpuStatInternal struct {
		User   int
		System int
		Idle   int
	}
)

func Start(ctx context.Context, wgGlobal *sync.WaitGroup) error {
	slog.Info("Start CPU collector")

	dataMon = storage.New[cpuStatInternal]()
	chToParser = make(chan string, chSize)
	chToCalculator = make(chan cpuStatInternal, chSize)

	go parser()
	go calculator()

	wgGlobal.Add(1)
	go probber(ctx, wgGlobal)

	return nil
}

func probber(ctx context.Context, wgGlobal *sync.WaitGroup) {
	defer wgGlobal.Done()
	defer slog.Info("CPU collector stopped")
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
				slog.Error("CPU", "error read data from "+fileName, err)
				// process.Stop()
				// Останавливаем работу всего сервера или только сбор данного параметра?
				// Если всего сервера - снять комментарий
				return
			}
		}
	}
}

func getData() error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	// Читаем только 1 строку
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return scanner.Err()
	}

	chToParser <- scanner.Text()

	return nil
}

func parser() {
	defer close(chToCalculator)

	var i int

	for data := range chToParser {
		if data[:4] != "cpu " {
			slog.Error("CPU", "incorrect data in parser", data[:4])
			continue
		}

		var (
			s    cpuStatInternal
			nice int
		)

		i++

		_, err := fmt.Sscanf(data, "cpu %d %d %d %d", &s.User, &nice, &s.System, &s.Idle)
		if err != nil {
			slog.Error("CPU", "sscanf error", err)
		}

		chToCalculator <- s // Дальше, на вычисление
	}
}

func calculator() {
	var (
		prev               cpuStatInternal
		user, system, idle int
	)

	for data := range chToCalculator {
		if prev.User == 0 && prev.System == 0 && prev.Idle == 0 {
			// Первая итерация - пропускаем, готовим данные для вычисления
			prev = data
			continue
		}

		user = data.User - prev.User
		system = data.System - prev.System
		idle = data.Idle - prev.Idle

		prev = data // Для последующей итерации

		data.User = user * 100 * 100 / (user + system + idle)
		data.System = system * 100 * 100 / (user + system + idle)
		data.Idle = idle * 100 * 100 / (user + system + idle)

		dataMon.Add(data) // В хранилище
	}
}

func GetAvg(m int) (CPUStat, error) {
	var (
		result                CPUStat
		r                     cpuStatInternal
		i, user, system, idle int
	)

	data := dataMon.Get(m)
	if data == nil {
		return result, myerr.ErrEmpty
	}

	for i, r = range data {
		user += r.User
		system += r.System
		idle += r.Idle
	}

	i++
	result.User = float32(user/i) / 100
	result.System = float32(system/i) / 100
	result.Idle = float32(idle/i) / 100

	return result, nil
}
