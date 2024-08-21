package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"level-zero/internal/dto"
	"level-zero/internal/storage/cache/2Q"
	"level-zero/internal/storage/strgerrs"
	"level-zero/pkg/postgres"
	"log"
)

const (
	ErrorNotUniqueCode = `23505`

	defaultCacheSize = 100
)

type OrdersStorage struct {
	cache *twoq.Cache
	*postgres.Postgres
}

func New(ctx context.Context, pg *postgres.Postgres) (*OrdersStorage, error) {
	//TODO: Options for setting cache
	const op = `storage.pgdb.storage.NewOrdersStorage`

	query := `CREATE TABLE IF NOT EXISTS orders(order_uid TEXT PRIMARY KEY, content TEXT, created_at TIMESTAMPTZ)`

	if _, err := pg.Pool.Exec(ctx, query); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storage := &OrdersStorage{
		cache:    twoq.New(100),
		Postgres: pg,
	}

	if err := storage.warmupCacheByLatest(ctx, defaultCacheSize/4); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return storage, nil
}

func (o *OrdersStorage) warmupCacheByLatest(ctx context.Context, size int) error {
	//consider that later orders will be the most requested (oldest were taken)
	const op = `storage.pgdb.storage.warmupCacheByLatest`

	query := `SELECT order_uid, content FROM orders ORDER BY created_at DESC LIMIT $1`

	orders, err := o.Pool.Query(ctx, query, size)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer orders.Close()

	for orders.Next() {
		var order dto.Order
		if err := orders.Scan(&order.UID, &order.Content); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		o.cache.AddOrder(order)
	}
	return nil
}

func (o *OrdersStorage) OrderByID(ctx context.Context, id string) (string, error) {
	const op = `storage.pgdb.storage.OrderByID`

	if content := o.cache.OrderByID(id); content != "" {
		//if we found order in cache
		return content, nil
	}
	//if we don't, look for order in db

	var uid, json string

	query := `SELECT order_uid, content FROM orders WHERE order_uid = $1`

	if err := o.Pool.QueryRow(ctx, query, id).Scan(&uid, &json); err != nil {
		//if an error occurred
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, strgerrs.ErrZeroRecordsFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	order := dto.Order{
		UID:     uid,
		Content: json,
	}

	//add to cache

	o.cache.AddOrder(order)

	return json, nil
}

func (o *OrdersStorage) ConsumeOrders(ctx context.Context, orders <-chan dto.Order) {
	const op = `storage.pgdb.storage.InsertOrder`
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("%s: context canceled, stopping orders consumption", op)
				return
			case order, ok := <-orders:
				if !ok {
					log.Printf("%s: ordersCh is closed, returning")
					return
				}
				if err := o.insertOrder(ctx, order); err != nil {
					log.Printf("%s: %v", op, err)
				}

			}
		}
	}()
}

func (o *OrdersStorage) insertOrder(ctx context.Context, order dto.Order) error {
	const op = `storage.pgdb.storage.InsertOrder`

	query := `INSERT INTO orders(order_uid, content, created_at) VALUES($1, $2, $3)`

	_, err := o.Pool.Exec(ctx, query, order.UID, order.Content, order.CreatedAt)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == ErrorNotUniqueCode {
			return fmt.Errorf("%s: %w", op, strgerrs.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	o.cache.AddOrder(order)

	return nil
}
