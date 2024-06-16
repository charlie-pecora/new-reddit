package application 

import (
	"log"
	"time"
	"net/http"
	"html/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/charlie-pecora/new-reddit/authenticator"

)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *chi.Mux {
	router := chi.NewRouter()

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

	return router
}

type IndexData struct {
	Name string
}

var indexTemplate = template.Must(template.New("index.html").ParseFiles("./templates/index.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.Execute(w, IndexData{"Charlie"})
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
