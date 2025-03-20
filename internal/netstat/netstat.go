package netstat

import (
	"context"
	"log/slog"
	"os/exec"
	"otus/internal/storage"
	"strings"
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

	var cmdOut, cmdErr strings.Builder
	cmd := exec.CommandContext(ctx, "netstat", "-apeW", "-A", "inet", "--numeric-hosts", "--numeric-ports")
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	if err := cmd.Run(); err != nil {
		return err
	}
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	for _, s := range strings.Split(string(out), "\n") {
		if !(s[:3] == "tcp" || s[:3] == "udp") {
			continue
		}

	}

	return nil
}
