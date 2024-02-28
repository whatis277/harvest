package auth

import "net/http"

func (c *Controller) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		session := SessionFromContext(ctx)
		if session != nil {
			c.Passwordless.Logout(ctx, session)
		}

		c.cleanupSessionToken(w)

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
