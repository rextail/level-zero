package mq

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

type Stream struct {
	nats.JetStreamContext
}

type Config struct {
	Token      string
	MsgHandler func(ctx context.Context, msgCh chan<- []byte) nats.MsgHandler
}

func NewStream(conn *nats.Conn, streamName string, subjects string) (*Stream, error) {
	const op = `internal.mq.nats.NewStream`

	js, _ := conn.JetStream(nats.PublishAsyncMaxPending(128))

	stream, _ := js.StreamInfo(streamName)

	if stream == nil {
		log.Printf("Stream %s doesn't exist, creating", streamName)
		_, err := js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			Retention: nats.WorkQueuePolicy,
			Subjects:  []string{subjects},
		})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return &Stream{js}, nil
}

func (s *Stream) SubscribeChannel(ctx context.Context, cfg Config, msgCh chan<- []byte) {
	//func initialize subscription to nats jetstream server
	//callback function MsgHandler sends messages to msgCh
	const op = `mq.nats.ConsumeOrders`
	_, err := s.Subscribe(cfg.Token, cfg.MsgHandler(ctx, msgCh))
	if err != nil {
		log.Printf("%s: Can't initialize subscription, check server configuration, error %v", op, err)
		return
	}
}
