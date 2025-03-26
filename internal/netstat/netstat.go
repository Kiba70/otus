package netstat

import (
	"context"
	"errors"
	"log/slog"
	"os/exec"
	"otus/internal/myerr"
	"otus/internal/storage"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	chSize = 10
)

var (
	dataMon    *storage.Storage[Netstat]
	Working    atomic.Bool
	chToParser chan []byte
	regTCP     = regexp.MustCompile(`tcp.+:(\d+)\s.*:(\*|\d+)\s+([[:graph:]]+)\s+([\w-_@\.]+)\s+\d+\s+([[:graph:]]*).*`)
	regUDP     = regexp.MustCompile(`udp.+:(\d+)\s.*:(\*|\d+)\s+([\w-_@\.]+)\s+\d+\s+([[:graph:]]*).*`)
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
	chToParser = make(chan []byte, chSize)

	slog.Info("Start Netstat collector")

	wgGlobal.Add(1)
	go probber(ctx, wgGlobal)

	return nil
}

func probber(ctx context.Context, wgGlobal *sync.WaitGroup) {
	defer wgGlobal.Done()
	defer slog.Info("Netstat collector stopped")
	defer close(chToParser)

	errorNetstatCount := 100 // Счётчик ошибочных запусков Netstat

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
				if errorNetstatCount <= 0 { // Иногда netstat завершается с ошибкой - игнорируем несколько ошибок
					return
				}
				errorNetstatCount--
			}
		}
	}
}

func getData(ctxGlobal context.Context) error {
	ctx, cancel := context.WithTimeout(ctxGlobal, 300*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, "netstat", "-apeW", "-A", "inet", "--numeric-hosts", "--numeric-ports")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	chToParser <- out

	return nil
}

func parser() {

	for out := range chToParser {
		sockets := make([]Socket, 0, strings.Count(string(out), "\n")+1)
		conn := make(map[string]int32)

		for _, s := range strings.Split(string(out), "\n") {
			if len(s) < 3 {
				continue
			}

			switch s[:3] {
			case "tcp":
				socket, status, err := parseLineTCP(s)
				if err != nil {
					slog.Error("Netstat", "error", err)
					continue // Ошибка
				}
				if status == "LISTEN" { // Слушаем порт
					socket.Protocol = "tcp"
					sockets = append(sockets, socket)
					continue
				}
				if status != "" {
					conn[status]++ // Счётчик
				}
			case "udp":
				socket, err := parseLineUDP(s)
				if err != nil {
					slog.Error("Netstat", "error", err)
					continue // Ошибка
				}

				// У UDP нет статуса соединения - отображаются только порты, которые открыты для входящих пакетов

				socket.Protocol = "udp"
				sockets = append(sockets, socket)
			}
		}

		var netstat Netstat

		netstat.Socket = sockets
		netstat.Conn = conn

		dataMon.Add(netstat)
	}

}

func parseLineTCP(line string) (Socket, string, error) {
	var sock Socket

	slog.Debug("Netstat", "line", line)

	splitLine := regTCP.FindStringSubmatch(line)

	if len(splitLine) != 6 { // столько должно быть распарсенных элементов
		return sock, "", errors.New("error in parsing line TCP")
	}
	slog.Debug("Netstat", "port", splitLine[1], "status", splitLine[3], "user", splitLine[4], "pid", splitLine[5])

	sock.User = splitLine[4]

	if i32, err := strconv.Atoi(splitLine[1]); err == nil {
		sock.Port = int32(i32)
	}

	if splitLine[5] != "-" {
		s2 := strings.Split(splitLine[5], "/")
		if len(s2) != 2 { // Не в формате
			return sock, splitLine[3], nil // Значения по умолчанию
		}
		if i32, err := strconv.Atoi(s2[0]); err == nil {
			sock.Pid = int32(i32)
		}
		sock.Command = s2[1]
	}

	return sock, splitLine[3], nil
}

func parseLineUDP(line string) (Socket, error) {
	var sock Socket

	slog.Debug("Netstat", "line", line)

	splitLine := regUDP.FindStringSubmatch(line)

	if len(splitLine) != 5 { // столько должно быть распарсенных элементов
		return sock, errors.New("error in parsing line UDP")
	}
	slog.Debug("Netstat", "port", splitLine[1], "user", splitLine[3], "pid", splitLine[4])

	sock.User = splitLine[3]

	if i32, err := strconv.Atoi(splitLine[1]); err == nil {
		sock.Port = int32(i32)
	}

	if splitLine[4] != "-" {
		s2 := strings.Split(splitLine[4], "/")
		if len(s2) != 2 { // Не в формате
			return sock, nil // Значения по умолчанию
		}
		if i32, err := strconv.Atoi(s2[0]); err == nil {
			sock.Pid = int32(i32)
		}
		sock.Command = s2[1]
	}

	return sock, nil
}

func GetSum(m int) (Netstat, error) {
	var (
		result Netstat
	)

	data := dataMon.Get(m)
	if data == nil {
		return result, myerr.ErrEmpty
	}

	result.Socket = make([]Socket, 0)
	result.Conn = make(map[string]int32)

	for _, netstat := range data {
		result.addSockets(netstat.Socket)
		for conn, count := range netstat.Conn {
			result.Conn[conn] += count
		}
	}

	return result, nil
}

func (n *Netstat) addSockets(newSockets []Socket) {
	for _, socket := range newSockets {
		have := false
		for _, s := range n.Socket {
			if equalSocket(s, socket) {
				have = true
			}
		}
		if !have {
			n.Socket = append(n.Socket, socket)
		}
	}
}

func equalSocket(s1, s2 Socket) bool {
	return s1 == s2
}
