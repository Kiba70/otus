package config

import (
	"flag"
	"log/slog"
	"os"
	"strings"
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

	modulesEnv := os.Getenv(*EnvName)

	if modulesEnv == "" {
		modulesEnv = allModules // Переменная окружения не установлена - запускаем всё
	}

	modules := make(map[string]struct{})

	for s := range strings.SplitSeq(modulesEnv, ",") {
		if strings.Contains(allModules, s) { // Проверяем на корректность
			modules[s] = struct{}{}
		}
	}

	return modules
}
