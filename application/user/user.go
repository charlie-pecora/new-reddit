package user

import (
	"html/template"
	"log"
	"net/http"

	"github.com/charlie-pecora/new-reddit/sessions"
)

func User(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	_profile, ok := session.Values["profile"]
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	profile := _profile.(map[string]interface{})
	log.Printf("%+v", profile)

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
