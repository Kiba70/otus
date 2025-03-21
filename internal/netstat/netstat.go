package netstat

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"otus/internal/storage"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	dataMon    *storage.Storage[Netstat]
	Working    atomic.Bool
	chToParser = make(chan []byte, 10)
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

	// Запускаем parser
	go parser()

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
	ctx, cancel := context.WithTimeout(ctxGlobal, 300*time.Millisecond)
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

	chToParser <- out

	return nil
}

func parser() {

	for out := range chToParser {
		for _, s := range strings.Split(string(out), "\n") {
			switch s[:3] {
			case "tcp":
				parseLineTCP(s)
			case "udp":
				parseLineUDP(s)
			}
		}
	}

}

func parseLineTCP(line string) {
	// r := regexp.MustCompile(`tcp.+:(\d+)\s.*:[*|\d+]\s+([[:graph:]]+)\s+([[:alpha:]]+)\s+\d+\s+(\d+)/?([[:graph:]]*)\s*`)
	r := regexp.MustCompile(`tcp.+:(\d+)\s.*:[*|\d+]\s+([[:graph:]]+)\s+([[:alpha:]]+)\s+\d+\s+([[:graph:]]*).*`)
	ss := r.FindStringSubmatch(line)
	fmt.Println("TCP:", len(ss), ss)
	for i, s := range ss {
		fmt.Println("TCP:", i, s)
	}
}

func parseLineUDP(line string) {
	s := strings.Split(line, ":")
	fmt.Println("UDP:", s)
}
