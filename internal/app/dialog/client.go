package dialog

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/ghkadim/highload_architect/generated/dialog/go_proto"
	grpc2 "github.com/ghkadim/highload_architect/internal/grpc"
	"github.com/ghkadim/highload_architect/internal/models"
)

type Client struct {
	conn         *grpc.ClientConn
	dialogClient pb.DialogServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc2.NewClient(addr)
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:         conn,
		dialogClient: pb.NewDialogServiceClient(conn),
	}
	return c, nil
}

func (c *Client) DialogSend(ctx context.Context, message models.DialogMessage) error {
	_, err := c.dialogClient.Send(ctx, &pb.SendRequest{
		Message: &pb.DialogMessage{
			FromUser: string(message.From),
			ToUser:   string(message.To),
			Text:     message.Text,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error) {
	resp, err := c.dialogClient.List(ctx, &pb.ListRequest{
		FromUser: string(userID1),
		ToUser:   string(userID2),
	})
	if err != nil {
		return nil, err
	}

	messages := make([]models.DialogMessage, 0, len(resp.Messages))
	for _, msg := range resp.Messages {
		messages = append(messages, models.DialogMessage{
			From: models.UserID(msg.FromUser),
			To:   models.UserID(msg.ToUser),
			Text: msg.Text,
		})
	}
	return messages, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
