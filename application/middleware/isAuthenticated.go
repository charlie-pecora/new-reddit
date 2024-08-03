package middleware

import (
	"net/http"
	"github.com/charlie-pecora/new-reddit/sessions"
)

// IsAuthenticated is a middleware that checks if
// the user has already been authenticated previously.
func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := sessions.GetSession(r).Values["profile"]; !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

