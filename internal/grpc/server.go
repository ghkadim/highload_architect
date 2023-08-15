package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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
		grpc.ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.ChainUnaryInterceptor(logMiddleware),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	for _, r := range s.reg {
		r.Register(srv)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Fatalf("listen: %s\n", err)
		}
		if err := srv.Serve(ln); err != nil && err != grpc.ErrServerStopped {
			logger.Fatalf("listen: %s\n", err)
		}
	}()
	logger.Infof("Server started")

	<-done
	logger.Infof("Server stopping")

	srv.GracefulStop()
	logger.Infof("Server exited properly")
}

func logMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	resp interface{}, err error,
) {
	start := time.Now()
	resp, err = handler(ctx, req)

	st := status.Convert(err)

	if st.Code() == codes.OK {
		logger.FromContext(ctx).Debugf(
			"%s %d %s",
			info.FullMethod,
			st.Code(),
			time.Since(start),
		)
	} else {
		logger.FromContext(ctx).Debugf(
			"%s %d %s: %s",
			info.FullMethod,
			st.Code(),
			time.Since(start),
			st.Message(),
		)
	}
	return
}
