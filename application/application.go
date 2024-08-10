package application

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"time"
	"os"
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/charlie-pecora/new-reddit/application/login"
	"github.com/charlie-pecora/new-reddit/application/user"
	"github.com/charlie-pecora/new-reddit/application/posts"
	"github.com/charlie-pecora/new-reddit/authenticator"
	"github.com/charlie-pecora/new-reddit/database"
	"github.com/charlie-pecora/new-reddit/sessions"
	myMiddleware "github.com/charlie-pecora/new-reddit/application/middleware"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *chi.Mux {
	router := chi.NewRouter()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]any{})
	gob.Register(login.Profile{})
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

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db := database.New(pool)

	router.Get("/", Index)
	authEndpoints := login.NewAuthEndpoints(auth, db)
	router.Get("/login", authEndpoints.Login)
	router.Get("/callback", authEndpoints.Callback)
	router.Get("/logout", authEndpoints.Logout)

	postsEndpoints := posts.NewPostsEndpoints(db)
	router.Group(func(r chi.Router) {
		r.Use(myMiddleware.IsAuthenticated)
		r.Get("/user", user.User)
		r.Get("/posts", postsEndpoints.ListPosts)
		r.Get("/posts/create", postsEndpoints.GetPostForm)
		r.Post("/posts/create", postsEndpoints.CreatePost)
	})

	return router
}

type IndexData struct {
	Name string
	Picture string
}

var indexTemplate = template.Must(template.New("base").ParseFiles("./templates/index.html", "./templates/base.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	var user IndexData
	log.Printf("%+v\n", session)
	switch profile := session.Values["profile"].(type) {
	case (login.Profile):
		user = IndexData{
			Name: profile.Nickname,
			Picture:  profile.Picture,
		}
	default:
	}

	err := indexTemplate.Execute(w, user)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
