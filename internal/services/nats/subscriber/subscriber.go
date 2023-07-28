package subscriber

import (
	"fmt"

	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
)

type subscriber struct {
	sc       stan.Conn
	db       db
	cache    cache
	startOpt stan.SubscriptionOption
	msgCount int
	subscription
}

type subscription struct {
	Subj string
	stan.Subscription
}

var (
	ErrSubscriptionNotExists error = fmt.Errorf("subscription not exists")
)

var (
	clusterID, clientID string = "test-cluster", "stan-sub"
	URL                 string = stan.DefaultNatsURL
	qgroup              string = ""
	startSeq            uint64 = 0
	startDelta          string = ""
	deliverAll          bool   = true
	newOnly             bool   = false
	deliverLast         bool   = false
	durable             string = ""
)

func NewSubscriber(nc *nats.Conn, db db, cache cache) (*subscriber, error) {
	const op = "nats.subscriber.NewSubscriber"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc), connectionLostHandler())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	startOpt := stan.StartAt(pb.StartPosition_NewOnly)
	if startSeq != 0 {
		startOpt = stan.StartAtSequence(startSeq)
	} else if deliverLast {
		startOpt = stan.StartWithLastReceived()
	} else if deliverAll && !newOnly {
		startOpt = stan.DeliverAllAvailable()
	} else if startDelta != "" {
		ago, err := time.ParseDuration(startDelta)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		startOpt = stan.StartAtTimeDelta(ago)
	}

	return &subscriber{
		sc:       sc,
		db:       db,
		cache:    cache,
		startOpt: startOpt,
		msgCount: 0,
	}, nil
}

func (s *subscriber) Subscribe(subj string) error {
	const op = "nats.subscriber.Subscribe"

	sub, err := s.sc.QueueSubscribe(subj, qgroup, s.MsgHandler(), s.startOpt, stan.DurableName(durable))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.subscription = subscription{
		Subj:         subj,
		Subscription: sub,
	}
	return nil
}
