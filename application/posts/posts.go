package posts

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/charlie-pecora/new-reddit/application/middleware"
	"github.com/charlie-pecora/new-reddit/application/login"
)

func ListPosts(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value(middleware.ProfileContextKey).(login.Profile)

	err := postsTemplate.Execute(w, PostsData{
		Name: profile.Nickname,
		Posts: []Post{
			{Title: "First post"},
			{Title: "Second post"},
			{Title: "Third post"},
		},
	})
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type Post struct {
	Title string;
	Author string;
	Created time.Time;
}

type PostsData struct {
	Name string;
	Posts []Post;
}

var postsTemplate = template.Must(template.New("base").ParseFiles("./templates/posts.html", "./templates/base.html"))


