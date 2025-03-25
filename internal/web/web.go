package web

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"

	"otus/internal/config"
	"otus/internal/cpu"
	"otus/internal/loadavg"
	"otus/internal/myerr"
	"otus/internal/netstat"
	"otus/internal/pb"
	"otus/internal/process"

	"google.golang.org/grpc"
)

var (
	ctxW context.Context
)

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

	// Gracefull shutdown
	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	return nil
}

func (s *server) LoadAvgGetMon(in *pb.Request, out pb.Monitoring_LoadAvgGetMonServer) error {
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

func (s *server) CpuGetMon(in *pb.Request, out pb.Monitoring_CpuGetMonServer) error {
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
			if err == myerr.ErrEmpty {
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
			if err == myerr.ErrEmpty {
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
