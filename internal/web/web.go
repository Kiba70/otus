package web

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"otus/internal/config"
	"otus/internal/loadavg"
	"otus/internal/myerr"
	"otus/internal/pb"
	"otus/internal/process"

	"google.golang.org/grpc"
)

var (
	ctxW context.Context
)

type server struct {
	pb.UnimplementedLoadAvgServer
}

func Start(ctx context.Context, wgGlobal *sync.WaitGroup) error {
	ctxW = context.WithoutCancel(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *config.Port))
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	// loadavg
	pb.RegisterLoadAvgServer(s, new(server))

	wgGlobal.Add(1)
	go func() {
		defer wgGlobal.Done()

		if err := s.Serve(lis); err != nil {
			slog.Error("GRPC failed to serve loadavg", "error", err)
			process.Stop()
		}
	}()

	// Gracefull shutdown
	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	return nil
}

func (s *server) LoadAvgGetMon(in *pb.LoadAvgRequest, out pb.LoadAvg_LoadAvgGetMonServer) error {
	seconds := in.GetSeconds()

	t := time.NewTicker(time.Duration(in.GetPeriod()) * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctxW.Done(): // Завершаем работу
			// Close
			return myerr.ErrStop
		case <-t.C:
			stat, err := loadavg.GetAvg(int(seconds))
			if err == myerr.ErrEmpty {
				continue
			}
			if err != nil {
				return err
			}
			result := new(pb.LoadAvgReply)
			result.One = stat.One
			result.Five = stat.Five
			result.Fifteen = stat.Fifteen
			if err = out.Send(result); err != nil {
				return nil // Клиент закрыл соединение
			}
		}
	}
}
