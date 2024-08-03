package user

import (
	"html/template"
	"log"
	"net/http"

	"github.com/charlie-pecora/new-reddit/application/middleware"
)

func User(w http.ResponseWriter, r *http.Request) {
	profile := r.Context().Value(middleware.ProfileContextKey).(map[string]interface{})

	err := userTemplate.Execute(w, UserData{
		Nickname: profile["nickname"].(string),
		Picture:  profile["picture"].(string),
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
