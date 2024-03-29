package auth

import (
	"fmt"
	"net/http"

	"github.com/whatis277/harvest/bean/internal/entity/model"
	"github.com/whatis277/harvest/bean/internal/entity/viewmodel"

	"github.com/whatis277/harvest/bean/internal/adapter/controller/base"
)

func (c *Controller) LoginPage() base.HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		sessionToken := c.getSessionToken(r)
		if sessionToken != nil {
			AuthedUserRedirect(w, r)
			return nil
		}

		return c.renderLogin(w, &viewmodel.LoginViewData{})
	}
}

func (c *Controller) LoginForm() base.HTTPHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		sessionToken := c.getSessionToken(r)
		if sessionToken != nil {
			AuthedUserRedirect(w, r)
			return nil
		}

		email := r.FormValue("email")
		if email == "" {
			return c.LoginView.Render(w, &viewmodel.LoginViewData{})
		}

		password := r.FormValue("password")
		if password != "" {
			return c.LoginView.Render(w, &viewmodel.LoginViewData{})
		}

		ctx := r.Context()

		err := c.Passwordless.Login(ctx, email)
		switch err.(type) {
		case nil, model.UserInputError:
			return c.renderLogin(w, &viewmodel.LoginViewData{
				Email: email,
			})

		default:
			renderErr := c.renderLogin(w, &viewmodel.LoginViewData{
				Email: email,
			})
			return &base.HTTPError{
				Message: fmt.Sprintf(
					("auth: login: error logging in user: %v |" +
						" auth: login: error rendering view: %v"),
					err,
					renderErr,
				),
			}
		}
	}
}

func (c *Controller) renderLogin(
	w http.ResponseWriter,
	data *viewmodel.LoginViewData,
) error {
	err := c.LoginView.Render(w, data)
	if err != nil {
		return &base.HTTPError{
			Status: http.StatusInternalServerError,

			Message: fmt.Sprintf(
				"auth: login: error rendering login view: %v",
				err,
			),
		}
	}

	return nil
}
