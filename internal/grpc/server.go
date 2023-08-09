package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ghkadim/highload_architect/internal/logger"
)

type registrable interface {
	Register(grpcServer grpc.ServiceRegistrar)
}

type Server struct {
	reg []registrable
}

func NewServer(reg ...registrable) *Server {
	return &Server{
		reg: reg,
	}
}

func (s *Server) ListenAndServe(addr string) {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(logMiddleware),
	)
	for _, r := range s.reg {
		r.Register(srv)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Fatal("listen: %s\n", err)
		}
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			logger.Fatal("listen: %s\n", err)
		}
	}()
	logger.Info("Server started")

	<-done
	logger.Info("Server stopping")

	srv.GracefulStop()
	logger.Info("Server exited properly")
}

func logMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error,
) {
	start := time.Now()
	resp, err = handler(ctx, req)

	st := status.Convert(err)

	if st.Code() == codes.OK {
		logger.Debug(
			"%s %d %s",
			info.FullMethod,
			st.Code(),
			time.Since(start),
		)
	} else {
		logger.Debug(
			"%s %d %s: %s",
			info.FullMethod,
			st.Code(),
			time.Since(start),
			st.Message(),
		)
	}
	return
}
