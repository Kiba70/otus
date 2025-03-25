package main

import (
	"log/slog"
	"otus/internal/config"
	"otus/internal/cpu"
	"otus/internal/loadavg"
	"otus/internal/netstat"
	"otus/internal/process"
	"otus/internal/web"
	"sync"
)

// const (
// 	logLevel = slog.LevelDebug
// )

func main() {
	wgGlobal := &sync.WaitGroup{}

	ctx, cancel := process.Start()
	defer cancel()
	defer process.Stop()

	if err := config.Start(); err != nil {
		slog.Error("Server: Config start error", "erroe", err)
		return
	}

	// if *config.Debug {
	// 	slog.SetLogLoggerLevel(logLevel)
	// }

	if err := loadavg.Start(ctx, wgGlobal); err != nil {
		slog.Error("Server: LoadAvg start error", "erroe", err)
		return
	}

	if err := cpu.Start(ctx, wgGlobal); err != nil {
		slog.Error("Server: CPU start error", "error", err)
		return
	}

	if err := netstat.Start(ctx, wgGlobal); err != nil {
		slog.Error("Server: CPU start error", "error", err)
		return
	}

	if err := web.Start(ctx, wgGlobal); err != nil {
		slog.Error("Server: WEB start error", "error", err)
		return
	}

	wgGlobal.Wait() // Ждём всех
}
