package user

import (
	"html/template"
	"log"
	"net/http"

	"github.com/charlie-pecora/new-reddit/application/login"
	"github.com/charlie-pecora/new-reddit/application/middleware"
)

func User(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value(middleware.ProfileContextKey).(login.Profile)
	log.Println(profile)

	err := userTemplate.Execute(w, UserData{
		Nickname: profile.Nickname,
		Picture:  profile.Picture,
	})
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type UserData struct {
	Nickname string
	Picture  string
}

var userTemplate = template.Must(template.New("base").ParseFiles("./templates/user.html", "./templates/base.html"))
