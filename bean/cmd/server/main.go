package main

import (
	"fmt"

	estimatorUC "harvest/bean/internal/usecase/estimator"
	paymentMethodUC "harvest/bean/internal/usecase/paymentmethod"

	envAdapter "harvest/bean/internal/adapter/env"
	homeHandler "harvest/bean/internal/adapter/handler/home"
	landingHandler "harvest/bean/internal/adapter/handler/landing"
	loginHandler "harvest/bean/internal/adapter/handler/login"

	paymentMethodDS "harvest/bean/internal/driver/datasource/paymentmethod"
	"harvest/bean/internal/driver/postgres"
	"harvest/bean/internal/driver/server"
	"harvest/bean/internal/driver/template"
	homeVD "harvest/bean/internal/driver/view/home"
	landingVD "harvest/bean/internal/driver/view/landing"
	loginVD "harvest/bean/internal/driver/view/login"
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

	homeView, err := homeVD.New(template.FS, template.HomeTemplate)
	if err != nil {
		panic(
			fmt.Errorf("error creating payment methods view: %v", err),
		)
	}

	s := server.New()

	s.Route("/", landingHandler.New(landingView))
	s.Route("/get-started", loginHandler.New(loginView))

	s.Route("/home", homeHandler.New(
		estimatorUC.UseCase{},
		paymentMethodUC.UseCase{
			PaymentMethods: paymentMethodRepo,
		},
		homeView,
	))

	s.Listen(":8080")
}
