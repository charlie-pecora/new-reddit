package sessions

import (
	"context"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
)

type SessionKey string

const SESSION_KEY SessionKey = "session"

// HTTP middleware providing the session to requessts
func NewSessionMiddleware() func(next http.Handler) http.Handler {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	s := sessions.NewSession(store, "auth-session")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//provide the "session" in the request context
			ctx := context.WithValue(r.Context(), SESSION_KEY, s)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSession(r *http.Request) *sessions.Session {
	session_value, ok := r.Context().Value(SESSION_KEY).(*sessions.Session)
	if !ok {
		panic("Session type is not known")
	}
	return session_value
}
