package subscriber

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
)

type subscriber struct {
	Subscriptions map[string]subscription
	sc            stan.Conn
	db            db
	cache         cache
	startOpt      stan.SubscriptionOption
	msgCount      int
}

type subscription struct {
	Subj string
	stan.Subscription
}

var (
	clusterID, clientID string = "test-cluster", "stan-sub"
	URL                 string = stan.DefaultNatsURL
	showTime            bool   = false
	qgroup              string = ""
	unsubscribe         bool   = false
	startSeq            uint64 = 0
	startDelta          string = ""
	deliverAll          bool   = true
	newOnly             bool   = false
	deliverLast         bool   = false
	durable             string = ""
)

func NewSubscriber(nc *nats.Conn, db db, cache cache) (*subscriber, error) {
	op := "nats.subscriber.NewSubscriber"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc), connectionLostHendler())
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
		Subscriptions: make(map[string]subscription),
		sc:            sc,
		db:            db,
		cache:         cache,
		startOpt:      startOpt,
		msgCount:      0,
	}, nil
}

func (s *subscriber) Subscribe(subj string) error {
	op := "nats.subscriber.Subscribe"

	sub, err := s.subscribe(subj)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.Subscriptions[subj] = subscription{
		Subj:         subj,
		Subscription: sub,
	}
	return nil
}

func (s *subscriber) ServSubscription(sub *subscription) {
	log.Printf("Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n", sub.Subj, clientID, qgroup, durable)

	if showTime {
		log.SetFlags(log.LstdFlags)
	}

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			if durable == "" || unsubscribe {
				sub.Unsubscribe()
			}
			s.sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func (s *subscriber) subscribe(subj string) (stan.Subscription, error) {
	op := "nats.subscriber.subscribe"

	sub, err := s.sc.QueueSubscribe(subj, qgroup, s.msgHandler(), s.startOpt, stan.DurableName(durable))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}

func connectionLostHendler() stan.Option {
	return stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
		log.Printf("Connection lost, reason: %v", reason)
	})
}
