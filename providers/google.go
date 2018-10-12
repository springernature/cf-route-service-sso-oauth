package providers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type GoogleProvider struct {
	*ProviderData
}

func NewGoogleProvider(p *ProviderData) *GoogleProvider {
	// Check if there is an override for the default endpoints
	authUri := "https://accounts.google.com/o/oauth2/auth"
	tokenUri := "https://www.googleapis.com/oauth2/v3/token"
	if uri := os.Getenv("GOOGLEAUTHURI"); uri != "" {
		authUri = uri
	}
	if uri := os.Getenv("GOOGLETOKENURI"); uri != "" {
		tokenUri = uri
	}
	// Add endpoints to ProviderData
	p.AuthUri = authUri
	p.TokenUri = tokenUri
	// Return
	return &GoogleProvider{
		ProviderData: p,
	}
}

func (p *GoogleProvider) SignIn(w http.ResponseWriter, r *http.Request) {
	scheme := "https"
	if r.Host == "localhost:8080" {
		scheme = "http"
	}
	http.Redirect(w, r, p.AuthUri+"?"+
		"client_id="+p.ClientID+"&"+
		"response_type=code"+"&"+
		"scope=openid%20email"+"&"+
		"redirect_uri="+scheme+"://"+r.Host+"/callback", 302)
}

func (p *GoogleProvider) Callback(w http.ResponseWriter, r *http.Request) {
	// Check if the callback comntains an authorization code
	if code := r.FormValue("code"); code != "" {
		//fmt.Fprint(w, "Code is: "+code)
		// Exchange code for access token and ID token
		v := url.Values{}
		v.Add("code", code)
		v.Add("client_id", p.ClientID)
		v.Add("client_secret", p.ClientSecret)
		scheme := "https"
		if r.Host == "localhost:8080" {
			scheme = "http"
		}
		v.Add("redirect_uri", scheme+"://"+r.Host+"/callback")
		v.Add("grant_type", "authorization_code")
		resp, err := http.PostForm(p.TokenUri, v)
		if err != nil {
			// Output error
			fmt.Println(err)
			return
		} else if resp.StatusCode != 200 {
			fmt.Fprint(w, "The POST to the token endpoint did not return HTTP 200 OK")
			return
		}
		type Token struct {
			IDToken string `json:"id_token"`
		}
		var t Token
		if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
			// Output err
			fmt.Println(err)
			return
		}
		// Get the second part (payload) of the JWT (IDToken)
		// Decode this base64 string to byte array
		TokBytes, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(strings.Split(t.IDToken, ".")[1])
		if err != nil {
			// Output error
			fmt.Println(err)
			return
		}
		type Payload struct {
			Email string `json:"email"`
		}
		var p Payload
		if err = json.Unmarshal(TokBytes, &p); err != nil {
			// Output error
			fmt.Println(err)
			return
		}
		// Return username
		fmt.Fprintf(w, "Username: %v", p.Email)
		return
	}
	fmt.Fprint(w, "Not a valid callback!")
}
