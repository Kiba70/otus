package main

import (
	"log/slog"
	"otus/internal/loadavg"
	"otus/internal/process"
	"sync"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	wgGlobal := &sync.WaitGroup{}

	ctx, cancel := process.Start()
	defer cancel()

	loadavg.Start(ctx, wgGlobal)

	wgGlobal.Wait() // Ждём всех
}
