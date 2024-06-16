package login

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"

	"github.com/charlie-pecora/new-reddit/authenticator"
	"github.com/charlie-pecora/new-reddit/sessions"
)

type AuthEndpoints struct {
	a *authenticator.Authenticator
}

func NewAuthEndpoints(a *authenticator.Authenticator) AuthEndpoints {
	return AuthEndpoints{a}
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

	var profile map[string]interface{}
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
