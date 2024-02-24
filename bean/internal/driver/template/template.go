package template

import (
	"embed"
)

//go:embed *.html layout/*.html
//go:embed paymentmethod/*.html subscription/*.html
var FS embed.FS

const (
	baseTemplate = "layout/base.html"
)

var (
	LandingTemplate = []string{
		baseTemplate,
		"landing.html",
	}

	LoginTemplate = []string{
		baseTemplate,
		"login.html",
	}

	HomeTemplate = []string{
		baseTemplate,
		"subscription/item.html",
		"paymentmethod/item.html",
		"home.html",
	}

	OnboardingTemplate = []string{
		baseTemplate,
		"onboarding.html",
	}

	CreatePaymentMethodTemplate = []string{
		baseTemplate,
		"paymentmethod/create.html",
	}
	DeletePaymentMethodTemplate = []string{
		baseTemplate,
		"subscription/item.html",
		"paymentmethod/item.html",
		"paymentmethod/delete.html",
	}

	CreateSubscriptionTemplate = []string{
		baseTemplate,
		"subscription/create.html",
	}
	DeleteSubscriptionTemplate = []string{
		baseTemplate,
		"subscription/item.html",
		"subscription/delete.html",
	}
)
