package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/ghkadim/highload_architect/internal/logger"
)

func NewClient(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(clientLogMiddleware),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func clientLogMiddleware(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)

	st := status.Convert(err)

	if st.Code() == codes.OK {
		logger.Debug(
			"%s %d %s",
			method,
			st.Code(),
			time.Since(start),
		)
	} else {
		logger.Debug(
			"%s %d %s: %s",
			method,
			st.Code(),
			time.Since(start),
			st.Message(),
		)
	}
	return err
}
