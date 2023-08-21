package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ghkadim/highload_architect/generated/dialog/go_proto"
	"github.com/ghkadim/highload_architect/internal/models"
)

type dialog interface {
	DialogSend(ctx context.Context, message models.DialogMessage) error
	DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error)
}

type Controller struct {
	pb.UnimplementedDialogServiceServer

	dialog dialog
}

func NewController(dialog dialog) *Controller {
	return &Controller{
		dialog: dialog,
	}
}

func (c *Controller) Register(grpcServer grpc.ServiceRegistrar) {
	pb.RegisterDialogServiceServer(grpcServer, c)
}

func (c *Controller) Send(ctx context.Context, req *pb.SendRequest) (*pb.SendReply, error) {
	err := c.dialog.DialogSend(ctx, models.DialogMessage{
		From: models.UserID(req.Message.FromUser),
		To:   models.UserID(req.Message.ToUser),
		Text: req.Message.Text,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.SendReply{}, nil
}

func (c *Controller) List(ctx context.Context, req *pb.ListRequest) (*pb.ListReply, error) {
	messages, err := c.dialog.DialogList(ctx, models.UserID(req.FromUser), models.UserID(req.ToUser))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &pb.ListReply{
		Messages: make([]*pb.DialogMessage, 0, len(messages)),
	}
	for _, m := range messages {
		resp.Messages = append(resp.Messages, &pb.DialogMessage{
			FromUser: string(m.From),
			ToUser:   string(m.To),
			Text:     m.Text,
		})
	}
	return resp, nil
}
