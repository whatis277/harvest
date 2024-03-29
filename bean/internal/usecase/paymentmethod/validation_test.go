package paymentmethod

import (
	"strings"
	"testing"

	"github.com/whatis277/harvest/bean/internal/entity/model"
)

func TestValidateLabel(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		if err := validateLabel("label"); err != nil {
			t.Errorf("expected nil, got: %s", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		if err := validateLabel(
			strings.Repeat("a", 256),
		); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestValidateLast4(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		if err := validateLast4("1234"); err != nil {
			t.Errorf("expected nil, got: %s", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		tests := []string{"123", "12345", "abc", "123a", "1 23"}

		for _, test := range tests {
			if err := validateLast4(test); err == nil {
				t.Errorf("expected error, got nil: %s", test)
			}
		}
	})
}

func TestValidateBrand(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tests := []model.PaymentMethodBrand{
			model.PaymentMethodBrandAmex,
			model.PaymentMethodBrandMastercard,
			model.PaymentMethodBrandVisa,
		}

		for _, test := range tests {
			if err := validateBrand(test); err != nil {
				t.Errorf("expected nil, got: %s", err)
			}
		}
	})

	t.Run("invalid", func(t *testing.T) {
		if err := validateBrand("invalid"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestValidateExpMonth(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		for i := 1; i <= 12; i++ {
			if err := validateExpMonth(i); err != nil {
				t.Errorf("expected nil, got: %s", err)
			}
		}
	})

	t.Run("invalid", func(t *testing.T) {
		tests := []int{0, 13}

		for _, test := range tests {
			if err := validateExpMonth(test); err == nil {
				t.Errorf("expected error, got nil: %d", test)
			}
		}
	})
}

func TestValidateExpYear(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		for i := 2000; i <= 2150; i += 10 {
			if err := validateExpYear(i); err != nil {
				t.Errorf("expected nil, got: %s", err)
			}
		}
	})

	t.Run("invalid", func(t *testing.T) {
		tests := []int{1999, 2151}

		for _, test := range tests {
			if err := validateExpYear(test); err == nil {
				t.Errorf("expected error, got nil: %d", test)
			}
		}
	})
}
