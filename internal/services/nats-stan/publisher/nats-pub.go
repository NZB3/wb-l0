package publisher

import (
	"github.com/nats-io/stan.go"
)

type publisher struct {
	sc   stan.Conn
	subj string
}
