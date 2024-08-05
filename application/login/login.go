package login

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"os"
	"log"

	"github.com/charlie-pecora/new-reddit/authenticator"
	"github.com/charlie-pecora/new-reddit/database"
	"github.com/charlie-pecora/new-reddit/sessions"
)

type AuthEndpoints struct {
	a *authenticator.Authenticator
	db *database.Queries
}

type Profile struct {
	Name string
	Nickname string
	Sub string
	Picture string
}

func ParseProfile(source map[string]any) (Profile, error) {
	var profile Profile
	name, ok := source["name"].(string)
	if !ok {
		return profile, errors.New("name was not present in user profile.")
	}
	profile.Name = name
	nickname, ok := source["nickname"].(string)
	if !ok {
		return profile, errors.New("nickname was not present in user profile.")
	}
	profile.Nickname = nickname
	sub, ok := source["sub"].(string)
	if !ok {
		return profile, errors.New("sub was not present in user profile.")
	}
	profile.Sub = sub
	picture, ok := source["picture"].(string)
	if !ok {
		return profile, errors.New("picture was not present in user profile.")
	}
	profile.Picture = picture

	return profile, nil
}

func NewAuthEndpoints(a *authenticator.Authenticator, db *database.Queries) AuthEndpoints {
	return AuthEndpoints{a, db}
}

// Handler for our login.
func (auth AuthEndpoints) Login(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the state inside the session.
	session := sessions.GetSession(r)
	session.Values["state"] = state
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, auth.a.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func (auth AuthEndpoints) Callback(w http.ResponseWriter, r *http.Request) {
	session := sessions.GetSession(r)
	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "Invalid state parameter.", http.StatusBadRequest)
		return
	}

	// Exchange an authorization code for a token.
	token, err := auth.a.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange an authorization code for a token", http.StatusUnauthorized)
		return
	}

	idToken, err := auth.a.VerifyIDToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Failed to verify ID Token.", http.StatusInternalServerError)
		return
	}

	var profile Profile
	if err := idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// save user to database or update info
	user, err := auth.db.CreateOrUpdateUser(r.Context(), database.CreateOrUpdateUserParams{
		Name: profile.Nickname,
		Sub: profile.Sub,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Received login fromm user %v", user.ID)

	// Redirect to logged in page.
	http.Redirect(w, r, "/user", http.StatusTemporaryRedirect)
}

func (auth AuthEndpoints) Logout(w http.ResponseWriter, r *http.Request) {
	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
