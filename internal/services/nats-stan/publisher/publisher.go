package publisher

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type publisher struct {
	sc stan.Conn
}

func NewPublisher(nc *nats.Conn) (*publisher, error) {
	const op = "services.nats-stan.publisher.NewPublisher"

	sc, err := stan.Connect("test-cluster", "stan-pub", stan.NatsConn(nc))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &publisher{
		sc: sc,
	}, nil
}

func (p *publisher) Publish(data []byte, subject string) error {
	const op = "services.nats-stan.publisher.Publish"

	if err := p.sc.Publish(subject, data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println("Published message with subject:", subject)
	return nil
}
