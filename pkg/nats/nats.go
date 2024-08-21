package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	Conn *nats.Conn
}

func New(connStr string) (*Nats, error) {
	const op = `pkg.nats.New`

	conn, err := nats.Connect(connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Nats{conn}, nil
}
