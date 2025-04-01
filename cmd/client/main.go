package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/Kiba70/otus/internal/pb"
	"github.com/Kiba70/otus/internal/process"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	logLevel = slog.LevelDebug
)

var (
	Port        *int
	Seconds     *int
	Period      *int
	Server      *string
	MonParametr *string
	client      pb.MonitoringClient
)

func init() {
	Server = flag.String("server", "127.0.0.1", "Адрес сервера")
	Port = flag.Int("port", 8080, "Порт сервера")
	Seconds = flag.Int("seconds", 1, "Количество секунд, за которые производится усреднение")
	Period = flag.Int("period", 1, "Период получения информации")
	MonParametr = flag.String("param", "", "Параметр мониторинга (одно значение): (load/cpu/netstat)")
}

func main() {
	slog.SetLogLoggerLevel(logLevel)

	ctx, cancel := process.Start()
	defer cancel()

	flag.Parse()

	// Готовим структуру для GRPC соединения. Само соединение не организуется
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", *Server, *Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Client: can not connect with server", "error", err)
	}
	defer conn.Close()
	client = pb.NewMonitoringClient(conn)

	//nolint:gosec
	request := &pb.Request{
		Period:  int32(*Period),
		Seconds: int32(*Seconds),
	}

	switch *MonParametr {
	case "load":
		slog.Debug("Client: Param loadAvg")
		if err := loadAvg(ctx, request); err != nil {
			slog.Error("Client: LoadAvg", "error", err)
			return
		}
	case "cpu":
		slog.Debug("Client: Param CPU")
		if err := cpu(ctx, request); err != nil {
			slog.Error("Client: CPU", "error", err)
			return
		}
	case "netstat":
		slog.Debug("Client: Param Netstat")
		if err := netstat(ctx, request); err != nil {
			slog.Error("Client: CPU", "error", err)
			return
		}
	default:
		flag.Usage()
		return
	}

	slog.InfoContext(ctx, "Client: End of program")
}

func loadAvg(ctx context.Context, request *pb.Request) error {
	slog.Debug("Client: Streaming started")

	stream, err := client.LoadAvgGetMon(ctx, request)
	if err != nil {
		return err
	}

	for {
		message, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
				slog.Debug("ExCTX")
				return nil // Нормальное завершение
			}
			fmt.Println("ERROR:", err)
			return err
		}
		log.Println("LoadAVG:", message)
	}
	slog.Debug("Client: Streaming finished")

	return nil
}

func cpu(ctx context.Context, request *pb.Request) error {
	slog.Debug("Client: Streaming started")

	stream, err := client.CPUGetMon(ctx, request)
	if err != nil {
		return err
	}

	for {
		message, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
				return nil // Нормальное завершение
			}
			return err
		}
		log.Println("CPU:", message)
	}
	slog.Debug("Client: Streaming finished")

	return nil
}

func netstat(ctx context.Context, request *pb.Request) error {
	slog.Debug("Client: Streaming started")
	defer slog.Debug("Client: Streaming finished")

	stream, err := client.NetstatGetMon(ctx, request)
	if err != nil {
		return err
	}

	for {
		message, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
				return nil // Нормальное завершение
			}
			return err
		}
		// log.Println("Netstat:", message)

		fmt.Println("Current connections:")
		i := len(message.Conn)
		for statusConn, num := range message.Conn {
			i--
			fmt.Printf("Type: %s = %d", statusConn, num)
			if i != 0 {
				fmt.Print(" | ")
			}
		}
		fmt.Printf("\nLISTENS:\n%4s | %5s | %25s | %6s | %20s\n", "Prot", "Port", "User", "PID", "Command")
		for _, socket := range message.Socket {
			fmt.Printf("%4s | %5d | %25s | %6d | %20s\n", socket.GetProtocol(),
				socket.GetPort(), socket.GetUser(), socket.GetPid(), socket.GetCommand())
		}
		fmt.Printf("-------------------------------------\n\n")
	}

	return nil
}
