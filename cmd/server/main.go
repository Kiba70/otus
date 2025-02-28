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

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	wgGlobal := &sync.WaitGroup{}

	ctx, cancel := process.Start()
	defer cancel()
	defer process.Stop()

	if err := config.Start(); err != nil {
		slog.Error("Config start error", "erroe", err)
		return
	}

	if err := loadavg.Start(ctx, wgGlobal); err != nil {
		slog.Error("LoadAvg start error", "erroe", err)
		return
	}

	if err := cpu.Start(ctx, wgGlobal); err != nil {
		slog.Error("CPU start error", "error", err)
		return
	}

	if err := web.Start(ctx, wgGlobal); err != nil {
		slog.Error("WEB start error", "error", err)
		return
	}

	wgGlobal.Wait() // Ждём всех
}
