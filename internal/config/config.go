package config

import (
	"flag"
	"log/slog"
	"os"
	"strings"
)

const (
	all_modules = "cpu,loadavg,netstat"
)

var (
	Port    = flag.Int("port", 8080, "The server port")
	EnvName = flag.String("mod", "OTUS_MOD_START", "ENV name for start modules")
	Debug   = flag.Bool("d", false, "Debug")
)

func Start() map[string]struct{} {
	flag.Parse()

	if *Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	modules_env := os.Getenv(*EnvName)

	if modules_env == "" {
		modules_env = all_modules // Переменная окружения не установлена - запускаем всё
	}

	modules := make(map[string]struct{})

	for _, s := range strings.Split(modules_env, ",") {
		if strings.Contains(all_modules, s) { // Проверяем на корректность
			modules[s] = struct{}{}
		}
	}

	return modules
}
