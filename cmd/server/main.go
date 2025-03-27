package main

//nolint:gofumpt,gci
import (
	"log/slog"
	"sync" //nolint:gofumpt

	"otus/internal/config"
	"otus/internal/cpu"
	"otus/internal/loadavg"
	"otus/internal/netstat"
	"otus/internal/process"
	"otus/internal/web" //nolint:gofumpt,nolintlint
)

func main() {
	wgGlobal := &sync.WaitGroup{}

	ctx, cancel := process.Start()
	defer cancel()
	defer process.Stop()

	modules := config.Start()
	if len(modules) == 0 {
		slog.Error("Server: Config start error", "error", "nothing modules for start")
		return
	}

	if _, on := modules["loadavg"]; on {
		if err := loadavg.Start(ctx, wgGlobal); err != nil {
			slog.Error("Server: LoadAvg start error", "error", err)
			return
		}
	}

	if _, on := modules["cpu"]; on {
		if err := cpu.Start(ctx, wgGlobal); err != nil {
			slog.Error("Server: CPU start error", "error", err)
			return
		}
	}

	if _, on := modules["netstat"]; on {
		if err := netstat.Start(ctx, wgGlobal); err != nil {
			slog.Error("Server: Netstat start error", "error", err)
			return
		}
	}

	if err := web.Start(ctx, wgGlobal); err != nil {
		slog.Error("Server: WEB start error", "error", err)
		return
	}

	wgGlobal.Wait() // Ждём всех
}
