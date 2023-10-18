package dialog

import (
	"context"

	"github.com/ghkadim/highload_architect/internal/models"
	"github.com/ghkadim/highload_architect/internal/saga"
)

const (
	unreadMessagesCounter = "unread_messages"
)

type storage interface {
	DialogSend(ctx context.Context, message models.DialogMessage) (models.DialogMessageID, error)
	DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error)
	DialogMessageRead(ctx context.Context, userID models.UserID, messageID models.DialogMessageID) error
}

type counter interface {
	CounterAdd(ctx context.Context, userID models.UserID, counter string, value int64) error
}

type Service interface {
	DialogSend(ctx context.Context, message models.DialogMessage) (models.DialogMessageID, error)
	DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error)
	DialogMessageRead(ctx context.Context, userID models.UserID, messageID models.DialogMessageID) error
}

type service struct {
	storage storage
	counter counter
}

func NewService(
	storage storage,
	counter counter,
) Service {
	return &service{
		storage: storage,
		counter: counter,
	}
}

func (s *service) DialogSend(ctx context.Context, message models.DialogMessage) (models.DialogMessageID, error) {
	var id models.DialogMessageID
	idPtr := &id
	err := saga.New([]saga.Step{
		{
			func() error {
				return s.counter.CounterAdd(ctx, message.From, unreadMessagesCounter, 1)
			},
			func() error {
				return s.counter.CounterAdd(ctx, message.From, unreadMessagesCounter, -1)
			},
		},
		{
			func() error {
				return s.counter.CounterAdd(ctx, message.To, unreadMessagesCounter, 1)
			},
			func() error {
				return s.counter.CounterAdd(ctx, message.To, unreadMessagesCounter, -1)
			},
		},
		{
			func() error {
				var err error
				*idPtr, err = s.storage.DialogSend(ctx, message)
				return err
			},
			nil,
		},
	}).Run()
	return id, err
}

func (s *service) DialogList(ctx context.Context, userID1, userID2 models.UserID) ([]models.DialogMessage, error) {
	messages, err := s.storage.DialogList(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *service) DialogMessageRead(ctx context.Context, userID models.UserID, messageID models.DialogMessageID) error {
	return saga.New([]saga.Step{
		{
			func() error {
				return s.counter.CounterAdd(ctx, userID, unreadMessagesCounter, -1)
			},
			func() error {
				return s.counter.CounterAdd(ctx, userID, unreadMessagesCounter, 1)
			},
		},
		{
			func() error {
				return s.storage.DialogMessageRead(ctx, userID, messageID)
			},
			nil,
		},
	}).Run()
}
