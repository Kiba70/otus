package netstat

import (
	"context"
	"log/slog"
	"otus/internal/storage"
	"sync"
	"sync/atomic"
	"time"
)

var (
	dataMon *storage.Storage[Netstat]
	Working atomic.Bool
)

type (
	Netstat struct {
		Socket []Socket
		Conn   map[string]int32
	}

	Socket struct {
		Command  string
		Pid      int32
		User     string
		Protocol string
		Port     int32
	}
)

func Start(ctx context.Context, wgGlobal *sync.WaitGroup) error {
	dataMon = storage.New[Netstat]()

	slog.Debug("CPU Start")

	wgGlobal.Add(1)
	go probber(ctx, wgGlobal)

	return nil
}

func probber(ctx context.Context, wgGlobal *sync.WaitGroup) {
	defer wgGlobal.Done()

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
			if err := getData(ctx); err != nil {
				slog.Error("Netstat", "error read data from netstat", err)
				// process.Stop() // Останавливаем работу всего сервера или только сбор данного параметра? Если всего сервера - снять комментарий
				return
			}
		}
	}
}

func getData(ctxGlobal context.Context) error {
	ctx, cancel := context.WithTimeout(ctxGlobal, 500*time.Millisecond)
	defer cancel()

	var cmdOut, errOut strings.Builder
	cmd := exec.CommandContext(ctx,"netstat", "-lntup")
	cmd.Stdout = &cmdOut
	cmd.Error = &errOut
	if err := cmd.Run(); err != nil {
		return err
	}
	out := cmd.Output()

	return nil
}
