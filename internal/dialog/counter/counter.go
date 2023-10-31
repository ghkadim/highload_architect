package counter

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/ghkadim/highload_architect/generated/counter/go_proto"
	grpc2 "github.com/ghkadim/highload_architect/internal/grpc"
	"github.com/ghkadim/highload_architect/internal/models"
)

type Client struct {
	conn          *grpc.ClientConn
	counterClient pb.CounterServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc2.NewClient(addr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:          conn,
		counterClient: pb.NewCounterServiceClient(conn),
	}
	return c, nil
}

func (c *Client) CounterAdd(ctx context.Context, userID models.UserID, counter string, value int64) error {
	_, err := c.counterClient.Add(ctx, &pb.AddRequest{UserID: string(userID), CounterName: counter, Amount: value})
	if err != nil {
		return err
	}
	return nil
}
