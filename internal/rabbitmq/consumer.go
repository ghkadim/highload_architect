package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/wagslane/go-rabbitmq"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
	"github.com/ghkadim/highload_architect/internal/result"
	"github.com/ghkadim/highload_architect/internal/utils/closer"
)

type Consumer interface {
	PostAdded(ctx context.Context, userID models.UserID, friends []models.UserID) (<-chan result.Result[models.Post], closer.Closer)
	FriendUpdated(ctx context.Context, userID models.UserID) (<-chan result.Result[models.FriendEvent], closer.Closer)
}

type consumer struct {
	queueLen             int
	postAddedConsumer    *queueConsumer[post]
	friendUpdateConsumer *queueConsumer[friendUpdate]
}

func NewConsumer(
	UserName string,
	Password string,
	Hostname string,
	QueueLen int,
) (
	Consumer,
	error,
) {
	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s", UserName, Password, Hostname),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, err
	}
	return &consumer{
		queueLen: QueueLen,
		postAddedConsumer: newRmqQueueConsumer[post](
			conn,
			postAddedQueueNamePrefix,
			postAddedExchangeName,
		),
		friendUpdateConsumer: newRmqQueueConsumer[friendUpdate](
			conn,
			friendUpdatedQueueNamePrefix,
			friendUpdatedExchangeName,
		),
	}, nil
}

func (c *consumer) PostAdded(ctx context.Context, userID models.UserID, friends []models.UserID) (<-chan result.Result[models.Post], closer.Closer) {
	logger.Debug("Consumer.PostAdded: new subscriberId=%s for authorId=%s", userID, friends)
	dataCh := make(chan result.Result[post], c.queueLen)
	closers := make([]closer.Closer, 0, len(friends))
	friendMap := make(map[models.UserID]struct{}, len(friends))

	for _, friend := range friends {
		friendMap[friend] = struct{}{}
		rk := postAddedRoutingKey(friend)
		closers = append(closers, c.postAddedConsumer.Consume(rk, dataCh))
	}

	asyncPipeWG := sync.WaitGroup{}
	asyncPipeWG.Add(1)

	resCh := make(chan result.Result[models.Post], c.queueLen)
	clsr := closer.FromFunction(func() error {
		logger.Debug("Consumer.PostAdded: closing subscription for subscriberId=%s", userID)
		var resErr error
		for _, cl := range closers {
			err := cl.Close()
			if err != nil {
				resErr = errors.Join(resErr, err)
			}
		}
		close(dataCh)
		asyncPipeWG.Wait()
		close(resCh)
		return resErr
	})

	go func() {
		defer asyncPipeWG.Done()
		for {
			select {
			case val, ok := <-dataCh:
				if !ok {
					return
				}
				v, err := val.Value()
				if err != nil {
					resCh <- result.Error[models.Post](err)
				}
				if _, ok = friendMap[v.AuthorID]; ok {
					resCh <- result.Value(models.Post{
						ID:       v.ID,
						AuthorID: v.AuthorID,
						Text:     v.Text,
					})
				}
			case <-ctx.Done():
				go func() {
					_ = clsr.Close()
				}()
				return
			}
		}
	}()

	return resCh, clsr
}

func (c *consumer) FriendUpdated(ctx context.Context, userID models.UserID) (<-chan result.Result[models.FriendEvent], closer.Closer) {
	logger.Debug("Consumer.FriendUpdated: subscribe userID=%s", userID)
	dataCh := make(chan result.Result[friendUpdate], c.queueLen)

	rk := friendUpdatedRoutingKey(userID)
	friendUpdateCloser := c.friendUpdateConsumer.Consume(rk, dataCh)

	asyncPipeWG := sync.WaitGroup{}
	asyncPipeWG.Add(1)

	resCh := make(chan result.Result[models.FriendEvent], c.queueLen)
	clsr := closer.FromFunction(func() error {
		logger.Debug("Consumer.FriendUpdated: closing subscription for subscriberId=%s", userID)
		err := friendUpdateCloser.Close()

		close(dataCh)
		asyncPipeWG.Wait()
		close(resCh)
		return err
	})
	go func() {
		defer asyncPipeWG.Done()
		for {
			select {
			case val, ok := <-dataCh:
				if !ok {
					return
				}
				v, err := val.Value()
				if err != nil {
					resCh <- result.Error[models.FriendEvent](err)
				}
				if v.UserID == userID {
					resCh <- result.Value(models.FriendEvent{
						Type:     v.Type,
						UserID:   v.UserID,
						FriendID: v.FriendID,
					})
				}
			case <-ctx.Done():
				go func() {
					_ = clsr.Close()
				}()
				return
			}
		}
	}()

	return resCh, clsr
}

func NewNopConsumer() Consumer {
	return nil
}
