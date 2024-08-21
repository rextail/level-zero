package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"level-zero/config"
	"level-zero/internal/dto"
	http_server "level-zero/internal/http-server"
	"level-zero/internal/http-server/handlers/getorder"
	"level-zero/internal/http-server/handlers/index"
	"level-zero/internal/mq"
	"level-zero/internal/storage/pgdb"
	"level-zero/internal/validator"
	"level-zero/pkg/nats"
	"level-zero/pkg/postgres"
	"log"
	"net/http"
)

const noOrdersPath = `internal/http-server/front/no_orders.html`
const ordersPath = `internal/http-server/front/order.html`

func main() {
	cfg := config.MustLoad("config/config.yml")

	pg, err := postgres.New(cfg.PG.ConnStr, postgres.MaxPoolSize(20))
	if err != nil {
		log.Fatalf("can't create postgre connection: %s", err.Error())
	}

	ctx := context.TODO()

	storage, err := pgdb.New(ctx, pg)
	if err != nil {
		log.Fatalf("can't create orders storage: %s", err.Error())
	}

	nats, err := nats.New(cfg.Nats.ConnStr)
	if err != nil {
		log.Fatalf("nats conn failed %v", err)
	}

	stream, err := mq.NewStream(nats.Conn, cfg.Nats.StreamName, cfg.Nats.SubjectName)
	if err != nil {
		log.Fatalf("nats stream init failed: %v", err)
	}

	msgs := make(chan []byte, 100)

	streamCfg := mq.Config{
		Token:      "orders.new",
		MsgHandler: mq.FetchToChannel,
	}

	stream.SubscribeChannel(ctx, streamCfg, msgs)

	orders := make(chan dto.Order, 100)

	go validator.Validate(ctx, msgs, orders)

	go storage.ConsumeOrders(ctx, orders)

	router := chi.NewRouter()

	responser := &http_server.OrderResponser{
		ordersPath,
		noOrdersPath,
	}

	router.Get("/", index.New("internal/http-server/front/index.html"))

	router.Post("/submit", getorder.New(ctx, storage, responser))

	srv := http.Server{
		Addr:    cfg.HTTP.Address,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.HTTP.Address)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("failed to start server")
		}
	}()

	<-ctx.Done()
}
