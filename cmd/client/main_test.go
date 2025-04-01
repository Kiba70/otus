package main

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Kiba70/otus/internal/pb"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestIntegrate(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	t.Run("Запускаем parser & calculator и готовим данные", func(t *testing.T) {
		conn, err := grpc.NewClient("127.0.0.1:8080",
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			slog.Error("Client: can not connect with server", "error", err)
		}
		defer conn.Close()
		client = pb.NewMonitoringClient(conn)

		request := &pb.Request{
			Period:  int32(1),
			Seconds: int32(3),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		require.NoError(t, loadAvg(ctx, request))

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		require.NoError(t, cpu(ctx, request))

		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		require.NoError(t, netstat(ctx, request))
	})
}
