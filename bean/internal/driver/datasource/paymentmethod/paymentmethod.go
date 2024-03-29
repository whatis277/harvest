package paymentmethod

import (
	"context"
	"fmt"

	"github.com/whatis277/harvest/bean/internal/entity/model"

	"github.com/whatis277/harvest/bean/internal/usecase/interfaces"

	"github.com/whatis277/harvest/bean/internal/driver/postgres"
)

type dataSource struct {
	db *postgres.DB
}

func New(db *postgres.DB) interfaces.PaymentMethodDataSource {
	return &dataSource{db}
}

func (ds *dataSource) Create(
	ctx context.Context,
	userID string,
	label string,
	last4 string,
	brand model.PaymentMethodBrand,
	expMonth int,
	expYear int,
) (*model.PaymentMethod, error) {
	method := &model.PaymentMethod{}

	err := ds.db.Pool.
		QueryRow(
			ctx,
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

func (ds *dataSource) FindByID(
	ctx context.Context,
	userID string,
	id string,
) (*model.PaymentMethodWithSubscriptions, error) {
	rows, err := ds.db.Pool.
		Query(
			ctx,
			("SELECT" +
				" payment_methods.id, payment_methods.user_id," +
				" payment_methods.label, payment_methods.last4, payment_methods.brand, payment_methods.exp_month, payment_methods.exp_year," +
				" payment_methods.created_at, payment_methods.updated_at," +
				" subscriptions.id, subscriptions.payment_method_id," +
				"subscriptions.label, subscriptions.provider," +
				" subscriptions.amount, subscriptions.interval, subscriptions.period" +
				" FROM payment_methods" +
				" LEFT JOIN subscriptions" +
				" ON payment_methods.id = subscriptions.payment_method_id" +
				" WHERE" +
				" payment_methods.user_id = $1" +
				" AND payment_methods.id = $2"),
			userID, id,
		)

	if err != nil {
		return nil, fmt.Errorf("failed to find payment method: %w", err)
	}

	defer rows.Close()

	method := &model.PaymentMethod{}
	subs := []*model.Subscription{}

	for rows.Next() {
		var (
			subID              *string
			subPaymentMethodID *string
			subLabel           *string
			subProvider        *string
			subAmount          *int
			subInterval        *int
			subPeriod          *string
		)

		err = rows.Scan(
			&method.ID, &method.UserID,
			&method.Label, &method.Last4, &method.Brand, &method.ExpMonth, &method.ExpYear,
			&method.CreatedAt, &method.UpdatedAt,
			&subID, &subPaymentMethodID,
			&subLabel, &subProvider,
			&subAmount, &subInterval, &subPeriod,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}

		sub := &model.Subscription{}
		if subID != nil {
			sub.ID = *subID
			sub.PaymentMethodID = *subPaymentMethodID
			sub.Label = *subLabel
			sub.Provider = *subProvider
			sub.Amount = *subAmount
			sub.Interval = *subInterval
			sub.Period = model.SubscriptionPeriod(*subPeriod)
		}

		if sub.ID != "" {
			subs = append(subs, sub)
		}
	}

	if method.ID == "" {
		return nil, nil
	}

	return &model.PaymentMethodWithSubscriptions{
		PaymentMethod: method,
		Subscriptions: subs,
	}, nil
}

func (ds *dataSource) FindByUserID(
	ctx context.Context,
	userID string,
) ([]*model.PaymentMethodWithSubscriptions, error) {
	rows, err := ds.db.Pool.
		Query(
			ctx,
			("SELECT" +
				" payment_methods.id, payment_methods.user_id," +
				" payment_methods.label, payment_methods.last4, payment_methods.brand, payment_methods.exp_month, payment_methods.exp_year," +
				" payment_methods.created_at, payment_methods.updated_at," +
				" subscriptions.id, subscriptions.payment_method_id," +
				" subscriptions.label, subscriptions.provider," +
				" subscriptions.amount, subscriptions.interval, subscriptions.period" +
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

	methods := []*model.PaymentMethodWithSubscriptions{}
	cache := &model.PaymentMethodWithSubscriptions{
		PaymentMethod: &model.PaymentMethod{},
		Subscriptions: []*model.Subscription{},
	}

	for rows.Next() {
		method := &model.PaymentMethod{}
		var (
			subID              *string
			subPaymentMethodID *string
			subLabel           *string
			subProvider        *string
			subAmount          *int
			subInterval        *int
			subPeriod          *string
		)

		err = rows.Scan(
			&method.ID, &method.UserID,
			&method.Label, &method.Last4, &method.Brand, &method.ExpMonth, &method.ExpYear,
			&method.CreatedAt, &method.UpdatedAt,
			&subID, &subPaymentMethodID,
			&subLabel, &subProvider,
			&subAmount, &subInterval, &subPeriod,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %w", err)
		}

		sub := &model.Subscription{}
		if subID != nil {
			sub.ID = *subID
			sub.PaymentMethodID = *subPaymentMethodID
			sub.Label = *subLabel
			sub.Provider = *subProvider
			sub.Amount = *subAmount
			sub.Interval = *subInterval
			sub.Period = model.SubscriptionPeriod(*subPeriod)
		}

		if cache.PaymentMethod.ID == "" {
			cache.PaymentMethod = method
		}

		if cache.PaymentMethod.ID != method.ID {
			methods = append(methods, cache)

			cache = &model.PaymentMethodWithSubscriptions{
				PaymentMethod: method,
				Subscriptions: []*model.Subscription{},
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

func (ds *dataSource) Delete(
	ctx context.Context,
	userID string,
	id string,
) error {
	_, err := ds.db.Pool.
		Exec(
			ctx,
			"DELETE FROM payment_methods WHERE user_id = $1 AND id = $2",
			userID, id,
		)

	if err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}

	return nil
}
