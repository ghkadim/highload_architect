package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/wagslane/go-rabbitmq"

	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

type Publisher interface {
	PostAdd(newPost models.Post) error
	FriendAdd(userID, friendID models.UserID) error
	FriendDelete(userID, friendID models.UserID) error
	Close() error
}

type publisher struct {
	conn                  *rabbitmq.Conn
	postPublisher         *rabbitmq.Publisher
	friendUpdatePublisher *rabbitmq.Publisher
}

func NewPublisher(
	UserName string,
	Password string,
	Hostname string,
) (Publisher, error) {
	conn, err := rabbitmq.NewConn(
		fmt.Sprintf("amqp://%s:%s@%s", UserName, Password, Hostname),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		return nil, err
	}
	postPublisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(postAddedExchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}

	friendUpdatePublisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(friendUpdatedExchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, err
	}

	return &publisher{
		conn:                  conn,
		postPublisher:         postPublisher,
		friendUpdatePublisher: friendUpdatePublisher,
	}, nil
}

func (p *publisher) PostAdd(newPost models.Post) error {
	data, err := json.Marshal(post{ID: newPost.ID, Text: newPost.Text, AuthorID: newPost.AuthorID})
	if err != nil {
		return err
	}
	routingKey := postAddedRoutingKey(newPost.AuthorID)
	return p.publish(p.postPublisher, postAddedExchangeName, routingKey, data)
}

func (p *publisher) FriendAdd(userID, friendID models.UserID) error {
	return p.friendUpdate(models.FriendAddedEvent, userID, friendID)
}

func (p *publisher) FriendDelete(userID, friendID models.UserID) error {
	return p.friendUpdate(models.FriendDeletedEvent, userID, friendID)
}

func (p *publisher) friendUpdate(updateType models.FriendEventType, userID, friendID models.UserID) error {
	data, err := json.Marshal(friendUpdate{Type: updateType, UserID: userID, FriendID: friendID})
	if err != nil {
		return err
	}
	routingKey := friendUpdatedRoutingKey(userID)
	return p.publish(p.friendUpdatePublisher, friendUpdatedExchangeName, routingKey, data)
}

func (p *publisher) publish(rmqPub *rabbitmq.Publisher, exchangeName, routingKey string, data []byte) error {
	logger.Debugf("RMQ_PUBLISH exchange=%s RK=%v %s", exchangeName, routingKey, string(data))
	return rmqPub.Publish(
		data,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(exchangeName),
	)
}

func (p *publisher) Close() error {
	p.postPublisher.Close()
	p.friendUpdatePublisher.Close()
	return p.conn.Close()
}

type nopPublisher struct{}

func NewNopPublisher() Publisher {
	return &nopPublisher{}
}

func (p *nopPublisher) PostAdd(newPost models.Post) error {
	return nil
}

func (p *nopPublisher) FriendAdd(userID, friendID models.UserID) error {
	return nil
}

func (p *nopPublisher) FriendDelete(userID, friendID models.UserID) error {
	return nil
}

func (p *nopPublisher) Close() error {
	return nil
}
