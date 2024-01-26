package paymentmethod

import (
	"fmt"

	"harvest/bean/internal/entity"

	"harvest/bean/internal/usecases/interfaces"
)

type UseCase struct {
	paymentMethods interfaces.PaymentMethodDataSource
}

func (u *UseCase) Create(
	userID string,
	label string,
	last4 string,
	brand entity.PaymentMethodBrand,
	expMonth int,
	expYear int,
) (*entity.PaymentMethod, error) {
	if err := validateLabel(label); err != nil {
		return nil, fmt.Errorf("invalid label: %w", err)
	}

	if err := validateLast4(last4); err != nil {
		return nil, fmt.Errorf("invalid last4: %w", err)
	}

	if err := validateBrand(brand); err != nil {
		return nil, fmt.Errorf("invalid brand: %w", err)
	}

	if err := validateExpMonth(expMonth); err != nil {
		return nil, fmt.Errorf("invalid exp month: %w", err)
	}

	if err := validateExpYear(expYear); err != nil {
		return nil, fmt.Errorf("invalid exp year: %w", err)
	}

	method, err := u.paymentMethods.Create(userID, label, last4, brand, expMonth, expYear)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment method: %w", err)
	}

	return method, nil
}

func (u *UseCase) Get(userID string, paymentMethodID string) (*entity.PaymentMethod, error) {
	method, err := u.paymentMethods.FindByID(userID, paymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method: %w", err)
	}

	return method, nil
}

func (u *UseCase) List(userID string) ([]*entity.PaymentMethod, error) {
	methods, err := u.paymentMethods.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment methods: %w", err)
	}

	return methods, nil
}

func (u *UseCase) Delete(userID string, paymentMethodID string) error {
	if err := u.paymentMethods.Delete(userID, paymentMethodID); err != nil {
		return fmt.Errorf("failed to delete payment method: %w", err)
	}

	return nil
}
