package main

//nolint:gofumpt,gci,nolintlint
import (
	"log/slog"
	"sync" //nolint:gofumpt

	"github.com/Kiba70/otus/internal/config"
	"github.com/Kiba70/otus/internal/cpu"
	"github.com/Kiba70/otus/internal/loadavg"
	"github.com/Kiba70/otus/internal/netstat"
	"github.com/Kiba70/otus/internal/process"
	"github.com/Kiba70/otus/internal/web" //nolint:gofumpt,nolintlint
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
