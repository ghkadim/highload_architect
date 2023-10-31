package tarantool

import (
	"context"

	"github.com/tarantool/go-tarantool/v2"

	"github.com/ghkadim/highload_architect/internal/models"
)

type Storage struct {
	address string
	opts    tarantool.Opts
}

func NewStorage(user, password, address string) *Storage {
	return &Storage{
		address: address,
		opts:    tarantool.Opts{User: user, Pass: password},
	}
}

func (s *Storage) do(ctx context.Context, fn func(conn *tarantool.Connection) error) error {
	conn, err := tarantool.Connect(ctx, s.address, s.opts)
	if err != nil {
		return err
	}
	defer conn.Close()
	return fn(conn)
}

func (s *Storage) CounterAdd(ctx context.Context, userID models.UserID, id string, value int64) (int64, error) {
	result := make([]int64, 0)
	err := s.do(ctx, func(conn *tarantool.Connection) error {
		f := conn.Do(tarantool.NewCallRequest("counter:add").
			Context(ctx).
			Args([]interface{}{userID, id, value}))
		return f.GetTyped(&result)
	})
	if err != nil {
		return 0, err
	}
	return result[0], nil
}

func (s *Storage) CounterRead(ctx context.Context, userID models.UserID, id string) (int64, error) {
	result := make([]struct {
		UserID string
		ID     string
		Value  int64
	}, 0)
	err := s.do(ctx, func(conn *tarantool.Connection) error {
		f := conn.Do(tarantool.NewCallRequest("counter:read").
			Context(ctx).
			Args([]interface{}{userID, id}))
		return f.GetTyped(&result)
	})
	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, models.ErrCounterNotFound
	}
	return result[0].Value, nil
}
