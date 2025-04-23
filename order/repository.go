package order

import (
	"context"
	"database/sql"
	"log"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, o Order) error
	GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	tx.ExecContext(
		ctx,
		"INSERT INTO orders(id, created_at, account_id, total_price) VALUES($1, $2, $3, $4)",
		o.Id,
		o.CreatedAt,
		o.AccountId,
		o.TotalPrice)

	if err != nil {
		return err
	}

	stat, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))

	for _, p := range o.Products {
		_, err = stat.ExecContext(ctx, o.Id, p.Id, p.Quantity)

		if err != nil {
			return err
		}

		_, err = stat.ExecContext(ctx)

		if err != nil {
			return err
		}
	}

	stat.Close()
	return nil

}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT o.Id,
		o.created_at,
		o.account_id,
		o.total_price::money::numeric::float8,
		op.product_id,
		op.quantity 
		FROM orders o JOIN order_products op ON(o.id = op.order_id)
		WHERE o.account_id = $1,
		ORDER BY o.created_at
		`,
		accountId,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := []Order{}
	order := &Order{}
	lastOrder := &Order{}
	orderedProduct := &OrderedProduct{}
	products := []OrderedProduct{}

	for rows.Next() {
		if err := rows.Scan(
			&order.Id,
			&order.CreatedAt,
			&order.TotalPrice,
			&order.AccountId,
			&order.Products,
			&orderedProduct.Id,
			&orderedProduct.Quantity,
		); err != nil {
			return nil, err
		}

		if lastOrder.Id != "" && lastOrder.Id != order.Id {
			newOrder := Order{
				Id:         lastOrder.Id,
				AccountId:  lastOrder.AccountId,
				CreatedAt:  lastOrder.CreatedAt,
				TotalPrice: lastOrder.TotalPrice,
				Products:   lastOrder.Products,
			}

			orders = append(orders, newOrder)
			products = []OrderedProduct{}
		}

		products = append(products, OrderedProduct{
			Id:       orderedProduct.Id,
			Quantity: orderedProduct.Quantity,
		})
		*lastOrder = *order
	}

	if lastOrder != nil {
		newOrder := Order{
			Id:         lastOrder.Id,
			AccountId:  lastOrder.AccountId,
			CreatedAt:  lastOrder.CreatedAt,
			TotalPrice: lastOrder.TotalPrice,
			Products:   lastOrder.Products,
		}

		orders = append(orders, newOrder)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
