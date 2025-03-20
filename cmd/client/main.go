package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"otus/internal/pb"
	"otus/internal/process"
	"strings"

	// "golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	logLevel = slog.LevelDebug
)

var (
	Port, Seconds, Period *int
	Server, MonParametr   *string
	// conn                  *grpc.ClientConn
	client pb.MonitoringClient
)

func init() {
	Server = flag.String("server", "127.0.0.1", "Адрес сервера")
	Port = flag.Int("port", 8080, "Порт сервера")
	Seconds = flag.Int("seconds", 1, "Количество секунд, за которые производится усреднение")
	Period = flag.Int("period", 1, "Период получения информации")
	MonParametr = flag.String("param", "", "Параметр мониторинга (одно значение): (load/cpu)")
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
		slog.Debug("Client: Param loadAvg")
		if err := cpu(ctx, request); err != nil {
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
		if err == io.EOF {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
				fmt.Println("ExCTX")
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

	stream, err := client.CpuGetMon(ctx, request)
	if err != nil {
		return err
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
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
