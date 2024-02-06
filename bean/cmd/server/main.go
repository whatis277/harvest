package main

import (
	"fmt"

	subscriptionUS "harvest/bean/internal/usecase/subscription"

	envAdapter "harvest/bean/internal/adapter/env"
	landingHandler "harvest/bean/internal/adapter/handler/landing"
	loginHandler "harvest/bean/internal/adapter/handler/login"
	subscriptionsHandler "harvest/bean/internal/adapter/handler/subscriptions"

	paymentMethodDS "harvest/bean/internal/driver/datasource/paymentmethod"
	subscriptionDS "harvest/bean/internal/driver/datasource/subscription"
	"harvest/bean/internal/driver/postgres"
	"harvest/bean/internal/driver/server"
	"harvest/bean/internal/driver/template"
	landingVD "harvest/bean/internal/driver/view/landing"
	loginVD "harvest/bean/internal/driver/view/login"
	subscriptionsVD "harvest/bean/internal/driver/view/subscriptions"
)

func main() {
	env, err := envAdapter.New()
	if err != nil {
		panic(
			fmt.Errorf("error reading env: %v", err),
		)
	}

	db, err := postgres.New(&postgres.DSNBuilder{
		Host:     env.DB.Host,
		Port:     env.DB.Port,
		Name:     env.DB.Name,
		Username: env.DB.Username,
		Password: env.DB.Password,
		SSLMode:  "disable",
	})
	if err != nil {
		panic(
			fmt.Errorf("error connecting db: %v", err),
		)
	}
	defer db.Close()

	paymentMethodRepo := paymentMethodDS.New(db)
	subscriptionRepo := subscriptionDS.New(db)

	landingView, err := landingVD.New(template.FS, template.LandingTemplate)
	if err != nil {
		panic(
			fmt.Errorf("error creating landing view: %v", err),
		)
	}

	loginView, err := loginVD.New(template.FS, template.LoginTemplate)
	if err != nil {
		panic(
			fmt.Errorf("error creating login view: %v", err),
		)
	}

	subscriptionsView, err := subscriptionsVD.New(template.FS, template.SubscriptionsTemplate)
	if err != nil {
		panic(
			fmt.Errorf("error creating subscriptions view: %v", err),
		)
	}

	s := server.New()

	s.Route("/", landingHandler.New(landingView))
	s.Route("/get-started", loginHandler.New(loginView))
	s.Route("/subscriptions", subscriptionsHandler.New(
		subscriptionUS.UseCase{
			PaymentMethods: paymentMethodRepo,
			Subscriptions:  subscriptionRepo,
		},
		subscriptionsView,
	))

	s.Listen(":8080")
}
