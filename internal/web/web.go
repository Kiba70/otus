package web

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"otus/internal/config"
	"otus/internal/pb"
	"otus/internal/process"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedLoadAvgServer
}

func Start(ctx context.Context, wgGlobal *sync.WaitGroup) error {
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

func (s *server) LoadAvgGetMon(in *pb.LoadAvgRequest, cln pb.LoadAvg_LoadAvgGetMonServer) error {
	_ = in.GetPeriod()
	_ = in.GetSeconds()

	return nil
}
