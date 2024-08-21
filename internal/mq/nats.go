package mq

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Stream struct {
	nats.JetStreamContext
}

type Config struct {
	Token        string
	RetryTimeout time.Duration
	MaxTimeout   time.Duration
	MsgHandler   func(ctx context.Context, msgCh chan<- []byte) nats.MsgHandler
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
	for {
		select {
		case <-ctx.Done():
			log.Printf("%s: Context was canceled, subscription was not initialized", op)
			return
		default:
			_, err := s.Subscribe(cfg.Token, cfg.MsgHandler(ctx, msgCh))
			if err != nil {
				log.Printf("%s: Can't initialize subscription, retry after: %s", op, cfg.RetryTimeout.String())
				time.Sleep(cfg.RetryTimeout)
				if cfg.RetryTimeout*2 <= cfg.MaxTimeout {
					cfg.RetryTimeout *= 2
				}
			} else {
				log.Printf("%s: successfully subcribed to: %s", op, cfg.Token)
				return
			}
		}
	}
}
