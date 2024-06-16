package application

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/charlie-pecora/new-reddit/application/login"
	"github.com/charlie-pecora/new-reddit/application/user"
	"github.com/charlie-pecora/new-reddit/authenticator"
	"github.com/charlie-pecora/new-reddit/sessions"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *chi.Mux {
	router := chi.NewRouter()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})
	router.Use(sessions.NewSessionMiddleware())

	//register middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", Index)
	authEndpoints := login.NewAuthEndpoints(auth)
	router.Get("/login", authEndpoints.Login)
	router.Get("/callback", authEndpoints.Callback)
	router.Get("/logout", authEndpoints.Logout)

	router.Get("/user", user.User)

	return router
}

type IndexData struct {
	Name string
}

var indexTemplate = template.Must(template.New("base").ParseFiles("./templates/index.html", "./templates/base.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.Execute(w, IndexData{"Charlie"})
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
