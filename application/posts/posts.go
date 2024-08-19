package posts

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/charlie-pecora/new-reddit/application/login"
	"github.com/charlie-pecora/new-reddit/application/middleware"
	"github.com/charlie-pecora/new-reddit/database"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostsEndpoints struct {
	db *database.Queries
}

func NewPostsEndpoints(db *database.Queries) PostsEndpoints {
	return PostsEndpoints{
		db,
	}
}

func (p PostsEndpoints) ListPosts(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value(middleware.ProfileContextKey).(login.Profile)
	posts, err := p.db.ListPosts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = postsTemplate.Execute(w, PostsData{
		Name:  profile.Nickname,
		Posts: posts,
	})
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p PostsEndpoints) GetPostForm(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value(middleware.ProfileContextKey).(login.Profile)
	err := createPostTemplate.Execute(w, PostsData{
		Name:  profile.Nickname,
		Posts: []database.ListPostsRow{},
	})
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p PostsEndpoints) CreatePost(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value(middleware.ProfileContextKey).(login.Profile)

	newPost, err := validateNewPostData(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := p.db.CreatePosts(
		r.Context(),
		database.CreatePostsParams{
			Sub:     profile.Sub,
			Title:   newPost.Title,
			Content: pgtype.Text{String: newPost.Content},
		},
	)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(created)
	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

type Post struct {
	Title   string
	Author  string
	Created time.Time
}

type PostsData struct {
	Name  string
	Posts []database.ListPostsRow
}

type NewPostData struct {
	Title   string
	Content string
}

func validateNewPostData(r *http.Request) (NewPostData, error) {
	errs := []error{}
	var newPost NewPostData
	err := r.ParseForm()
	if err != nil {
		return newPost, errors.New("Invalid form data")
	}
	title := r.Form.Get("title")
	if len(title) < 3 {
		errs = append(errs, errors.New("Post title must contain at least 3 characters"))
	}
	newPost.Title = title
	content := r.Form.Get("content")
	if len(content) < 10 {
		errs = append(errs, errors.New("Post content must be at least 10 characters in length"))
	}
	newPost.Content = content
	if len(errs) != 0 {
		return newPost, errors.Join(errs...)
	}
	return newPost, nil
}

var postsTemplate = template.Must(template.New("base").ParseFiles("./templates/posts.html", "./templates/base.html"))
var createPostTemplate = template.Must(template.New("base").ParseFiles("./templates/createPost.html", "./templates/base.html"))
