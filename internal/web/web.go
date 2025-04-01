package web

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/Kiba70/otus/internal/config"
	"github.com/Kiba70/otus/internal/cpu"
	"github.com/Kiba70/otus/internal/loadavg"
	"github.com/Kiba70/otus/internal/myerr"
	"github.com/Kiba70/otus/internal/netstat"
	"github.com/Kiba70/otus/internal/pb"
	"github.com/Kiba70/otus/internal/process"
	"google.golang.org/grpc"
)

var ctxW context.Context

type server struct {
	pb.UnimplementedMonitoringServer
}

func Start(ctx context.Context, wgGlobal *sync.WaitGroup) error {
	ctxW = ctx

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *config.Port))
	if err != nil {
		return err
	}

	s := grpc.NewServer()

	// loadavg
	pb.RegisterMonitoringServer(s, new(server))

	wgGlobal.Add(1)
	go func() {
		defer wgGlobal.Done()

		if err := s.Serve(lis); err != nil {
			slog.Error("GRPC failed to serve monitoring", "error", err)
			process.Stop()
		}
		slog.DebugContext(ctxW, "Завершение WEB сервера")
	}()

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	return nil
}

func (s *server) LoadAvgGetMon(in *pb.Request, out pb.Monitoring_LoadAvgGetMonServer) error { //nolint:dupl
	if !loadavg.Working.Load() {
		return myerr.ErrNotWork
	}

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
			if errors.Is(err, myerr.ErrEmpty) {
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

func (s *server) CPUGetMon(in *pb.Request, out pb.Monitoring_CPUGetMonServer) error { //nolint:dupl
	if !cpu.Working.Load() {
		return myerr.ErrNotWork
	}

	seconds := in.GetSeconds()

	t := time.NewTicker(time.Duration(in.GetPeriod()) * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctxW.Done(): // Завершаем работу
			// Close
			return myerr.ErrStop
		case <-t.C:
			stat, err := cpu.GetAvg(int(seconds))
			if errors.Is(err, myerr.ErrEmpty) {
				continue
			}
			if err != nil {
				return err
			}
			result := new(pb.CpuReply)
			result.User = stat.User
			result.System = stat.System
			result.Idle = stat.Idle
			if err = out.Send(result); err != nil {
				return nil // Клиент закрыл соединение
			}
		}
	}
}

func (s *server) NetstatGetMon(in *pb.Request, out pb.Monitoring_NetstatGetMonServer) error {
	if !netstat.Working.Load() {
		return myerr.ErrNotWork
	}

	seconds := in.GetSeconds()

	t := time.NewTicker(time.Duration(in.GetPeriod()) * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctxW.Done(): // Завершаем работу
			// Close
			return myerr.ErrStop
		case <-t.C:
			stat, err := netstat.GetSum(int(seconds))
			if errors.Is(err, myerr.ErrEmpty) {
				continue
			}
			if err != nil {
				return err
			}
			result := new(pb.NetstatReply)
			result.Conn = stat.Conn
			result.Socket = make([]*pb.NetstatSocketReply, 0, len(stat.Socket))
			for _, socket := range stat.Socket {
				var pbs pb.NetstatSocketReply
				pbs.Command = socket.Command
				pbs.Pid = socket.Pid
				pbs.Port = socket.Port
				pbs.Protocol = socket.Protocol
				pbs.User = socket.User
				result.Socket = append(result.Socket, &pbs)
			}

			if err = out.Send(result); err != nil {
				return nil // Клиент закрыл соединение
			}
		}
	}
}
