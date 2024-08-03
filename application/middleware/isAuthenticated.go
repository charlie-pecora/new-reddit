package middleware

import (
	"log"
	"net/http"
	"github.com/charlie-pecora/new-reddit/sessions"
	"context"
)

type ProfileKey string
const ProfileContextKey ProfileKey = "profile"

// IsAuthenticated is a middleware that checks if
// the user has already been authenticated previously.
// It then adds the "profile" key to the request context for use in handlers
func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile, ok := sessions.GetSession(r).Values["profile"]
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			ctx := context.WithValue(r.Context(), ProfileContextKey, profile)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

