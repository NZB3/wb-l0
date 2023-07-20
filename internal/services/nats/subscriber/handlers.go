// Not shure about rightness of this code. Maybe it should not be here.
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
	SetOrder(order models.Order) error
}

func (s *subscriber) MsgHandler() stan.MsgHandler {
	return func(msg *stan.Msg) {
		s.msgCount++
		printMsg(msg, s.msgCount)
		if s.subscription.Subj == "order" {
			s.orderMsgHandler(msg)
		}
	}
}

func (s *subscriber) orderMsgHandler(msg *stan.Msg) {
	const op = "nats-stan.subscriber.handlers.orderMsgHandler"

	order := models.Order{}
	err := order.Unmarshal(msg.Data)
	if err != nil {
		log.Printf("Not order message: %v", err)
		return
	}

	err = s.db.SaveDataFromOrder(order)
	if err != nil {
		log.Printf("%s: %v", op, err)
	}

	err = s.cache.SetOrder(order)
	if err != nil {
		log.Printf("%s: %v", op, err)
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
