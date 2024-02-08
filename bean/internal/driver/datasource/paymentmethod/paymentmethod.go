package paymentmethod

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

func New(db *postgres.DB) interfaces.PaymentMethodDataSource {
	return &dataSource{db}
}

func (ds *dataSource) Create(
	userID string,
	label string,
	last4 string,
	brand entity.PaymentMethodBrand,
	expMonth int,
	expYear int,
) (*entity.PaymentMethod, error) {
	method := &entity.PaymentMethod{}

	err := ds.db.Pool.
		QueryRow(
			context.Background(),
			("INSERT INTO payment_methods (user_id, label, last4, brand, exp_month, exp_year)"+
				" VALUES ($1, $2, $3, $4, $5, $6)"+
				" RETURNING *"),
			userID, label, last4, brand, expMonth, expYear,
		).
		Scan(
			&method.ID, &method.UserID,
			&method.Label, &method.Last4, &method.Brand, &method.ExpMonth, &method.ExpYear,
			&method.CreatedAt, &method.UpdatedAt,
		)

	if err != nil {
		return nil, fmt.Errorf("failed to create payment method: %w", err)
	}

	return method, nil
}

func (ds *dataSource) FindByID(userID string, id string) (*entity.PaymentMethodWithSubscriptions, error) {
	rows, err := ds.db.Pool.
		Query(
			context.Background(),
			("SELECT" +
				" payment_methods.id, payment_methods.user_id," +
				" payment_methods.label, payment_methods.last4, payment_methods.brand, payment_methods.exp_month, payment_methods.exp_year," +
				" payment_methods.created_at, payment_methods.updated_at," +
				" subscriptions.id, subscriptions.label, subscriptions.amount, subscriptions.interval, subscriptions.period" +
				" FROM payment_methods" +
				" LEFT JOIN subscriptions" +
				" ON payment_methods.id = subscriptions.payment_method_id" +
				" WHERE" +
				" payment_methods.user_id = $1" +
				" AND payment_methods.id = $2"),
			userID, id,
		)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find payment method: %w", err)
	}

	defer rows.Close()

	method := &entity.PaymentMethod{}
	subs := []*entity.Subscription{}

	for rows.Next() {
		var (
			subID       *string
			subLabel    *string
			subAmount   *int
			subInterval *int
			subPeriod   *string
		)

		err = rows.Scan(
			&method.ID, &method.UserID,
			&method.Label, &method.Last4, &method.Brand, &method.ExpMonth, &method.ExpYear,
			&method.CreatedAt, &method.UpdatedAt,
			&subID, &subLabel, &subAmount, &subInterval, &subPeriod,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}

		sub := &entity.Subscription{}
		if subID != nil {
			sub.ID = *subID
			sub.Label = *subLabel
			sub.Amount = *subAmount
			sub.Interval = *subInterval
			sub.Period = entity.SubscriptionPeriod(*subPeriod)
		}

		if sub.ID != "" {
			subs = append(subs, sub)
		}
	}

	return &entity.PaymentMethodWithSubscriptions{
		PaymentMethod: method,
		Subscriptions: subs,
	}, nil
}

func (ds *dataSource) FindByUserID(userID string) ([]*entity.PaymentMethodWithSubscriptions, error) {
	rows, err := ds.db.Pool.
		Query(
			context.Background(),
			("SELECT" +
				" payment_methods.id, payment_methods.user_id," +
				" payment_methods.label, payment_methods.last4, payment_methods.brand, payment_methods.exp_month, payment_methods.exp_year," +
				" payment_methods.created_at, payment_methods.updated_at," +
				" subscriptions.id, subscriptions.label, subscriptions.amount, subscriptions.interval, subscriptions.period" +
				" FROM payment_methods" +
				" LEFT JOIN subscriptions" +
				" ON payment_methods.id = subscriptions.payment_method_id" +
				" WHERE" +
				" payment_methods.user_id = $1" +
				" ORDER BY payment_methods.created_at DESC"),
			userID,
		)

	if err != nil {
		return nil, fmt.Errorf("failed to find payment methods: %w", err)
	}

	defer rows.Close()

	methods := []*entity.PaymentMethodWithSubscriptions{}
	cache := &entity.PaymentMethodWithSubscriptions{
		PaymentMethod: &entity.PaymentMethod{},
		Subscriptions: []*entity.Subscription{},
	}

	for rows.Next() {
		method := &entity.PaymentMethod{}
		var (
			subID       *string
			subLabel    *string
			subAmount   *int
			subInterval *int
			subPeriod   *string
		)

		err = rows.Scan(
			&method.ID, &method.UserID,
			&method.Label, &method.Last4, &method.Brand, &method.ExpMonth, &method.ExpYear,
			&method.CreatedAt, &method.UpdatedAt,
			&subID, &subLabel, &subAmount, &subInterval, &subPeriod,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}

		sub := &entity.Subscription{}
		if subID != nil {
			sub.ID = *subID
			sub.Label = *subLabel
			sub.Amount = *subAmount
			sub.Interval = *subInterval
			sub.Period = entity.SubscriptionPeriod(*subPeriod)
		}

		if cache.PaymentMethod.ID == "" {
			cache.PaymentMethod = method
		}

		if cache.PaymentMethod.ID != method.ID {
			methods = append(methods, cache)

			cache = &entity.PaymentMethodWithSubscriptions{
				PaymentMethod: method,
				Subscriptions: []*entity.Subscription{},
			}
		}

		if sub.ID != "" {
			cache.Subscriptions = append(cache.Subscriptions, sub)
		}
	}

	if cache.PaymentMethod.ID != "" {
		methods = append(methods, cache)
	}

	return methods, nil
}

func (ds *dataSource) Delete(userID string, id string) error {
	_, err := ds.db.Pool.
		Exec(
			context.Background(),
			"DELETE FROM payment_methods WHERE user_id = $1 AND id = $2",
			userID, id,
		)

	if err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}

	return nil
}
