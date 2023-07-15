package subscriber

import (
	"log"
	"project/internal/models"
	"time"

	"github.com/nats-io/stan.go"
)

type db interface {
	SaveDataFromOrder(order models.Order) error
	GetOrdersInLast(time.Duration) ([]models.Order, error)
}

type cache interface {
	GetOrder(orederID string) (models.Order, error)
	SetOrder(order models.Order) error
}

func (s *subscriber) msgHandler() stan.MsgHandler {
	return s.standardMessage()
}

func (s *subscriber) standardMessage() stan.MsgHandler {
	return func(msg *stan.Msg) {
		s.msgCount++
		printMsg(msg, s.msgCount)
	}
}

func printMsg(m *stan.Msg, i int) {
	log.Printf("[#%d] Received: %s\n", i, m)
}
