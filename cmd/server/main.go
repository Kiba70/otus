package main

import (
	"log/slog"
	"otus/internal/config"
	"otus/internal/cpu"
	"otus/internal/loadavg"
	"otus/internal/process"
	"otus/internal/web"
	"sync"
)

const (
	logLevel = slog.LevelDebug
)

func main() {
	slog.SetLogLoggerLevel(logLevel)
	wgGlobal := &sync.WaitGroup{}

	ctx, cancel := process.Start()
	defer cancel()
	defer process.Stop()

	if err := config.Start(); err != nil {
		slog.Error("Server: Config start error", "erroe", err)
		return
	}

	if err := loadavg.Start(ctx, wgGlobal); err != nil {
		slog.Error("Server: LoadAvg start error", "erroe", err)
		return
	}

	if err := cpu.Start(ctx, wgGlobal); err != nil {
		slog.Error("Server: CPU start error", "error", err)
		return
	}

	if err := web.Start(ctx, wgGlobal); err != nil {
		slog.Error("Server: WEB start error", "error", err)
		return
	}

	wgGlobal.Wait() // Ждём всех
}
