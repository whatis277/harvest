package paymentmethod

import (
	"fmt"
	"regexp"

	"github.com/whatis277/harvest/bean/internal/entity/model"
)

func validateLabel(label string) model.UserInputError {
	if len(label) > 255 {
		return fmt.Errorf("label must be less than 255 chars")
	}

	return nil
}

func validateLast4(last4 string) model.UserInputError {
	pattern, err := regexp.Compile(`^\d{4}$`)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %w", err)
	}

	if !pattern.MatchString(last4) {
		return fmt.Errorf("last4 must be 4 digits")
	}

	return nil
}

func validateBrand(brand model.PaymentMethodBrand) model.UserInputError {
	switch brand {
	case model.PaymentMethodBrandAmex,
		model.PaymentMethodBrandMastercard,
		model.PaymentMethodBrandVisa:
		return nil
	default:
		return fmt.Errorf(
			"brand must be: %s, %s, %s",
			model.PaymentMethodBrandAmex,
			model.PaymentMethodBrandMastercard,
			model.PaymentMethodBrandVisa,
		)
	}
}

func validateExpMonth(expMonth int) model.UserInputError {
	if expMonth < 1 {
		return fmt.Errorf("exp month must be greater than 0")
	}

	if expMonth > 12 {
		return fmt.Errorf("exp month must be less than 13")
	}

	return nil
}

func validateExpYear(expYear int) model.UserInputError {
	if expYear < 2000 {
		return fmt.Errorf("exp year must be greater than 2000")
	}

	if expYear > 2150 {
		return fmt.Errorf("exp year must be less than 2151")
	}

	return nil
}
