package validator

import (
	"context"
	"encoding/json"
	v10 "github.com/go-playground/validator/v10"
	"level-zero/internal/dto"
	"level-zero/internal/models"
	"log"
)

// go-playground/validator/v10
var valid = v10.New()

func Validate(ctx context.Context, msgCh <-chan []byte, ordersCh chan<- dto.Order) {
	const op = `validator.Validate`

	if err := valid.RegisterValidation("validateDate", ValidateDate); err != nil {
		log.Printf("%s: %v", op, err)
	}

	go func(msgCh <-chan []byte, ordersCh chan<- dto.Order) {
		defer close(ordersCh)
		for {
			select {
			case msg, ok := <-msgCh:
				if !ok {
					log.Printf("%s: msgCh is closed, validation was stopped", op)
					return
				}
				var order models.Order
				if err := json.Unmarshal(msg, &order); err != nil {
					log.Printf("%s failed to marshal %s: %v", op, string(msg), err)
					continue
				}
				if err := valid.Struct(order); err != nil {
					log.Printf("%s failed to validate %s: %v", op, string(msg), err)
					continue
				}
				ordersCh <- dto.Order{order.OrderUID, string(msg), order.DateCreated}
			case <-ctx.Done():
				log.Printf("%s: context was canceled, returning from validation", op)
				return
			}
		}
	}(msgCh, ordersCh)
}
