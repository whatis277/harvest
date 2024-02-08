package subscription

import (
	"context"
	"fmt"

	"harvest/bean/internal/entity"

	"harvest/bean/internal/usecase/interfaces"

	"harvest/bean/internal/driver/postgres"

	"github.com/jackc/pgx/v5"
)

type dataSource struct {
	db *postgres.DB
}

func New(db *postgres.DB) interfaces.SubscriptionDataSource {
	return &dataSource{
		db: db,
	}
}

func (ds *dataSource) Create(
	userID string,
	paymentMethodID string,
	label string,
	provider string,
	amount int,
	interval int,
	period entity.SubscriptionPeriod,
) (*entity.Subscription, error) {
	sub := &entity.Subscription{}

	err := ds.db.Pool.
		QueryRow(
			context.Background(),
			("INSERT INTO subscriptions"+
				" (user_id, payment_method_id, label, provider, amount, interval, period)"+
				" VALUES ($1, $2, $3, $4, $5, $6, $7)"+
				" RETURNING *"),
			userID, paymentMethodID, label, provider, amount, interval, period,
		).
		Scan(
			&sub.ID, &sub.UserID, &sub.PaymentMethodID,
			&sub.Label, &sub.Provider,
			&sub.Amount, &sub.Interval, &sub.Period,
			&sub.CreatedAt, &sub.UpdatedAt,
		)

	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return sub, nil
}

func (ds *dataSource) FindByID(userID string, id string) (*entity.SubscriptionWithPaymentMethod, error) {
	sub := &entity.Subscription{}
	method := &entity.PaymentMethod{}

	err := ds.db.Pool.
		QueryRow(
			context.Background(),
			("SELECT"+
				" subscriptions.*, payment_methods.*"+
				" FROM subscriptions"+
				" LEFT JOIN payment_methods"+
				" ON subscriptions.payment_method_id = payment_methods.id"+
				" WHERE"+
				" subscriptions.user_id = $1"+
				" AND subscriptions.id = $2"),
			userID, id,
		).
		Scan(
			&sub.ID, &sub.UserID, &sub.PaymentMethodID,
			&sub.Label, &sub.Provider,
			&sub.Amount, &sub.Interval, &sub.Period,
			&sub.CreatedAt, &sub.UpdatedAt,
			&method.ID, &method.UserID,
			&method.Label, &method.Last4, &method.Brand, &method.ExpMonth, &method.ExpYear,
			&method.CreatedAt, &method.UpdatedAt,
		)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find subscription: %w", err)
	}

	return &entity.SubscriptionWithPaymentMethod{
		Subscription:  sub,
		PaymentMethod: method,
	}, nil
}

func (ds *dataSource) FindByUserID(userID string) ([]*entity.SubscriptionWithPaymentMethod, error) {
	rows, err := ds.db.Pool.
		Query(
			context.Background(),
			("SELECT" +
				" subscriptions.*, payment_methods.*" +
				" FROM subscriptions" +
				" LEFT JOIN payment_methods" +
				" ON subscriptions.payment_method_id = payment_methods.id" +
				" WHERE" +
				" subscriptions.user_id = $1"),
			userID,
		)

	if err != nil {
		return nil, fmt.Errorf("failed to find subscriptions: %w", err)
	}

	defer rows.Close()

	subs := []*entity.SubscriptionWithPaymentMethod{}
	for rows.Next() {
		sub := &entity.Subscription{}
		method := &entity.PaymentMethod{}

		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.PaymentMethodID,
			&sub.Label, &sub.Provider,
			&sub.Amount, &sub.Interval, &sub.Period,
			&sub.CreatedAt, &sub.UpdatedAt,
			&method.ID, &method.UserID,
			&method.Label, &method.Last4, &method.Brand, &method.ExpMonth, &method.ExpYear,
			&method.CreatedAt, &method.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		subs = append(subs, &entity.SubscriptionWithPaymentMethod{
			Subscription:  sub,
			PaymentMethod: method,
		})
	}

	return subs, nil
}

func (ds *dataSource) Delete(userID string, id string) error {
	_, err := ds.db.Pool.
		Exec(
			context.Background(),
			"DELETE FROM subscriptions WHERE user_id = $1 AND id = $2",
			userID, id,
		)

	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}
