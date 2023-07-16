package subscriber

import (
	"fmt"
	"log"
	"project/internal/models"
	"time"

	"github.com/nats-io/stan.go"
)

type db interface {
	SaveDataFromOrder(order models.Order) error
	GetOrdersInLast(time.Duration) ([]models.Order, error)
	CheckDB() error
}

type cache interface {
	GetOrder(orederID string) (models.Order, error)
	SetOrder(order models.Order) error
	CheckCache() error
}

func (s *subscriber) msgHandler() (stan.MsgHandler, error) {
	op := "nats.subscriber.msgHandler"

	err := s.db.CheckDB()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = s.cache.CheckCache()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if s.subscription.Subj == "order" {
		return s.orderMessage(), nil
	}

	return s.standardMessage(), nil
}

func (s *subscriber) orderMessage() stan.MsgHandler {
	op := "nats.subscriber.orderMessage"
	return func(msg *stan.Msg) {
		s.msgCount++
		printMsg(msg, s.msgCount)
		_ = op
		// TODO: write logic to order message handler
	}
}

func (s *subscriber) standardMessage() stan.MsgHandler {
	return func(msg *stan.Msg) {
		s.msgCount++
		printMsg(msg, s.msgCount)
	}
}

func connectionLostHandler() stan.Option {
	return stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
		log.Printf("Connection lost, reason: %v", reason)
	})
}

func printMsg(m *stan.Msg, i int) {
	log.Printf("[#%d] Received: %s\n", i, m)
}
