package mq

import (
	"context"
	"github.com/nats-io/nats.go"
	"log"
)

func FetchToChannel(ctx context.Context, msgCh chan<- []byte) nats.MsgHandler {
	//Writes nats messages to msgCh
	const op = `mq.handlers.FetchToChannel`
	return func(msg *nats.Msg) {
		select {
		case msgCh <- msg.Data:
			if err := msg.Ack(); err != nil {
				msg.Nak()
				log.Printf("%s Failed to confirm message: %v", op, err)
			}
		case <-ctx.Done():
			close(msgCh)
			log.Printf("%s: context canceled, stopped fetching msg", op)
			return

		}
	}
}
