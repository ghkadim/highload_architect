package rabbitmq

import (
	"github.com/ghkadim/highload_architect/internal/models"
)

const (
	postAddedQueueNamePrefix     = "postAdded"
	postAddedExchangeName        = "postAdded"
	friendUpdatedQueueNamePrefix = "friendUpdated"
	friendUpdatedExchangeName    = "friendUpdated"
)

func postAddedRoutingKey(authorID models.UserID) string {
	return string(authorID)
}

func friendUpdatedRoutingKey(userID models.UserID) string {
	return string(userID)
}
