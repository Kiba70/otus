// Статистика по сетевым соединениям:
//
// 1. слушающие TCP & UDP сокеты: command, pid, user, protocol, port;
// 2. количество TCP соединений, находящихся в разных состояниях (ESTAB, FIN_WAIT, SYN_RCV и пр.).
//
// Раелизован алгоритм:
// 1. По данному пункту: Все открытые порты, которые встречались за период М
// 2. По данному пункту: По каждому статусу соединения сумма значений за период М (не среднее)

package netstat

import (
	"context"
	"errors"
	"log/slog"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Kiba70/otus/internal/myerr"
	"github.com/Kiba70/otus/internal/storage"
)

const (
	chSize = 10
)

var (
	dataMon    *storage.Storage[Netstat]
	Working    atomic.Bool
	chToParser chan []byte
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
				// process.Stop()
				// Останавливаем работу всего сервера или только сбор данного параметра?
				// Если всего сервера - снять комментарий
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

	out, err := exec.CommandContext(ctx, netstatCommand, netstatARGS...).CombinedOutput()
	if err != nil {
		return err
	}

	chToParser <- out

	return nil
}

//nolint:gocognit
func parser() {
	for out := range chToParser {
		sockets := make([]Socket, 0, strings.Count(string(out), lineDelim)+1)
		conn := make(map[string]int32)

		for s := range strings.SplitSeq(string(out), lineDelim) {
			if len(s) < 3 {
				continue
			}

			switch strings.Trim(s[:5], " \t") {
			case "tcp", "TCP":
				socket, status, err := parseLineTCP(s)
				if err != nil {
					if !errors.Is(err, errNetV6) {
						slog.Error("Netstat", "error", err)
					}
					continue // Ошибка
				}
				if status[:6] == "LISTEN" { // Слушаем порт
					socket.Protocol = "tcp"
					sockets = append(sockets, socket)
					continue
				}
				if status != "" {
					conn[status]++ // Счётчик
				}
			case "udp", "UDP":
				socket, err := parseLineUDP(s)
				if err != nil {
					if !errors.Is(err, errNetV6) {
						slog.Error("Netstat", "error", err)
					}
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

func GetSum(m int) (Netstat, error) {
	var result Netstat

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
