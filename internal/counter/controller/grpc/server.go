package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ghkadim/highload_architect/generated/counter/go_proto"
	"github.com/ghkadim/highload_architect/internal/models"
)

type counter interface {
	CounterAdd(ctx context.Context, userID models.UserID, id string, value int64) (int64, error)
}

type Controller struct {
	pb.UnsafeCounterServiceServer

	counter counter
}

func NewController(counter counter) *Controller {
	return &Controller{
		counter: counter,
	}
}

func (c *Controller) Register(grpcServer grpc.ServiceRegistrar) {
	pb.RegisterCounterServiceServer(grpcServer, c)
}

func (c *Controller) Add(ctx context.Context, r *pb.AddRequest) (*pb.AddReply, error) {
	amount, err := c.counter.CounterAdd(ctx, models.UserID(r.UserID), r.CounterName, r.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.AddReply{Value: amount}, nil
}
