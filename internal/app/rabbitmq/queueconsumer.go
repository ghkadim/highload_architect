package rabbitmq

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"

	"github.com/wagslane/go-rabbitmq"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/result"
	"github.com/ghkadim/highload_architect/internal/utils/closer"
)

var nodeID = func() string {
	hostname := os.Getenv("HOSTNAME")
	if hostname != "" {
		return hostname
	}
	return strconv.FormatInt(rand.Int63(), 16)
}()

type queueConsumer[T any] struct {
	mtx               sync.Mutex
	conn              *rabbitmq.Conn
	queueNamePrefix   string
	exchangeName      string
	consumingChannels map[string]struct {
		channels  map[chan<- result.Result[T]]int
		rmqCloser closer.Closer
	}
}

func newRmqQueueConsumer[T any](
	conn *rabbitmq.Conn,
	queueNamePrefix string,
	exchangeName string,
) *queueConsumer[T] {
	return &queueConsumer[T]{
		conn:            conn,
		queueNamePrefix: queueNamePrefix,
		exchangeName:    exchangeName,
		consumingChannels: make(map[string]struct {
			channels  map[chan<- result.Result[T]]int
			rmqCloser closer.Closer
		}),
	}
}

func (c *queueConsumer[T]) Consume(routingKey string, outCh chan<- result.Result[T]) closer.Closer {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if _, ok := c.consumingChannels[routingKey]; !ok {
		rmqCons, err := c.startConsumer(routingKey)
		if err != nil {
			outCh <- result.Error[T](err)
			return closer.Nop()
		}
		c.consumingChannels[routingKey] = struct {
			channels  map[chan<- result.Result[T]]int
			rmqCloser closer.Closer
		}{
			channels:  make(map[chan<- result.Result[T]]int),
			rmqCloser: rmqCons,
		}
	}
	channels := c.consumingChannels[routingKey].channels
	channels[outCh] = channels[outCh] + 1

	return closer.FromFunction(func() error {
		c.mtx.Lock()
		defer c.mtx.Unlock()
		channels[outCh] = channels[outCh] - 1
		if channels[outCh] == 0 {
			delete(channels, outCh)
			if len(channels) == 0 {
				_ = c.consumingChannels[routingKey].rmqCloser.Close()
				delete(c.consumingChannels, routingKey)
				return nil
			}
		}
		return nil
	})
}

func (c *queueConsumer[T]) startConsumer(
	routingKey string,
) (closer.Closer, error) {
	queueName := fmt.Sprintf("%s.%s.%s", c.queueNamePrefix, routingKey, nodeID)
	logger.Debugf("Starting RMQ consumer queue=%s exchange=%s", queueName, c.exchangeName)
	cons, err := rabbitmq.NewConsumer(
		c.conn,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			logger.Debugf("RMQ_CONSUME queue=%s body=%s", queueName, string(d.Body))
			c.mtx.Lock()
			defer c.mtx.Unlock()
			channelsForRK, ok := c.consumingChannels[routingKey]
			if !ok {
				logger.Errorf("All channels for RK=%s are closed", routingKey)
				return rabbitmq.NackDiscard
			}

			var t T
			err := json.Unmarshal(d.Body, &t)
			if err != nil {
				for ch := range channelsForRK.channels {
					ch <- result.Error[T](err)
				}
				return rabbitmq.NackDiscard
			}
			for ch := range channelsForRK.channels {
				ch <- result.Value(t)
			}
			return rabbitmq.Ack
		},
		queueName,
		rabbitmq.WithConsumerOptionsQueueAutoDelete,
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsExchangeName(c.exchangeName),
		rabbitmq.WithConsumerOptionsExchangeKind("direct"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}
	return closer.FromFunction(func() error {
		logger.Debugf("Stopping RMQ consumer queue=%s exchange=%s", queueName, c.exchangeName)
		cons.Close()
		return nil
	}), nil
}
