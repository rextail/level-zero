package getorder

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"level-zero/internal/models"
	"level-zero/internal/storage/strgerrs"
	resp "level-zero/pkg/api/response"
	"log"
	"net/http"
)

type Request struct {
	UID string `json:"id"`
}

type Response struct {
	resp.Response
	order string
}

func renderError(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, resp.Error(message))
}

type OrderGetter interface {
	OrderByID(ctx context.Context, ID string) (string, error)
}

type OrderResponser interface {
	OrderResponse(w http.ResponseWriter, order *models.Order) error
}

func New(ctx context.Context, orderGetter OrderGetter, orderResponser OrderResponser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = `http-server.getorder.New`

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Printf("%s failed to decode id %s: %v", op, r.Body, err)

			render.JSON(w, r, resp.Error("failed to decode uid"))

			return
		}

		order, err := orderGetter.OrderByID(ctx, req.UID)

		var Order models.Order

		json.Unmarshal([]byte(order), &Order)

		if err != nil {
			if errors.Is(err, strgerrs.ErrZeroRecordsFound) {
				log.Printf("%s: found no orders for requested id %s", op, req.UID)

				if err = orderResponser.OrderResponse(w, &Order); err != nil {

					log.Printf("%s: failed to execute no orders html", op)

					renderError(w, r, "server failed to form response")

					return
				}
				log.Printf("should be executed")
				return
			}
			log.Printf("%s: failed to decode id %v", op, err)

			renderError(w, r, "failed to decode uid")

			return
		}

		if err = orderResponser.OrderResponse(w, nil); err != nil {
			log.Printf("%s: failed to execute <orders.html>%v %v", op, order, err)

			renderError(w, r, "server failed to form response")

			return
		}

	}
}
