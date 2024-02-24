package interfaces

import (
	"net/http"

	"github.com/whatis277/harvest/bean/internal/entity/viewmodel"
)

// --- Views ---

type View[T any] interface {
	Render(http.ResponseWriter, *T) error
}

type LandingView View[viewmodel.LandingViewData]

type LoginView View[viewmodel.LoginViewData]

type HomeView View[viewmodel.HomeViewData]
type OnboardingView View[viewmodel.OnboardingViewData]

type CreatePaymentMethodView View[viewmodel.CreatePaymentMethodViewData]
type DeletePaymentMethodView View[viewmodel.DeletePaymentMethodViewData]

type CreateSubscriptionView View[viewmodel.CreateSubscriptionViewData]
type DeleteSubscriptionView View[viewmodel.DeleteSubscriptionViewData]
