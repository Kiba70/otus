package config

import (
	"flag"
	"log/slog"
)

var (
	Port  = flag.Int("port", 8080, "The server port")
	Debug = flag.Bool("d", false, "Debug")
)

func Start() error {
	flag.Parse()

	if *Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	return nil
}
