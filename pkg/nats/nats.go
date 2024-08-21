package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Nats struct {
	Conn *nats.Conn
}

func New(connStr string) (*Nats, error) {
	const op = `pkg.nats.New`

	opts := []nats.Option{
		nats.MaxReconnects(5),
		nats.ReconnectWait(10 * time.Second),
		nats.RetryOnFailedConnect(true),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("%s Reconnected to %v\n", op, nc.ConnectedUrl())
		}),
		nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
			log.Printf("%s: Connection lost, retrying in 2s, error: %v", op, err)
		}),
		nats.ConnectHandler(func(nc *nats.Conn) {
			log.Printf("%s: Connected to nats server", op)
		}),
	}

	conn, err := nats.Connect(connStr, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Nats{conn}, nil
}
