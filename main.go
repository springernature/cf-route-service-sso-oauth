package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	authUri  string
	tokenUri string
	clientID string
	clientS  string
)

func main() {
	log.SetOutput(os.Stdout)

	// Get mandatory ProviderData
	//p := providers.InitProviderData()

	authUri = os.Getenv("GOOGLEAUTHURI")
	tokenUri = os.Getenv("GOOGLETOKENURI")
	clientID = os.Getenv("GOOGLECLIENTID")
	clientS = os.Getenv("GOOGLECLIENTSECRET")

	// Sign in handler
	http.HandleFunc("/signin", signin)
	// Callback handler
	http.HandleFunc("/callback", callback)
	// Default handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		// Check if already logged in
		cookie, _ := r.Cookie("auth")
		if cookie != nil {
			fmt.Fprint(w, "You are authenticated!")
		} else {
			fmt.Fprint(w, "<html>You are NOT authenticated. <a href=\"/signin\">Sign in</a> Host: "+r.Host+"</html>")
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func signin(w http.ResponseWriter, r *http.Request) {
	scheme := "https"
	if r.Host == "localhost:8080" {
		scheme = "http"
	}
	http.Redirect(w, r, authUri+"?"+
		"client_id="+clientID+"&"+
		"response_type=code"+"&"+
		"scope=openid%20email"+"&"+
		"redirect_uri="+scheme+"://"+r.Host+"/callback", 302)
}

func callback(w http.ResponseWriter, r *http.Request) {
	// Check if the callback is a returning authorization code
	if code := r.FormValue("code"); code != "" {
		//fmt.Fprint(w, "Code is: "+code)
		// Exchange code for access token and ID token
		v := url.Values{}
		v.Add("code", code)
		v.Add("client_id", clientID)
		v.Add("client_secret", clientS)
		scheme := "https"
		if r.Host == "localhost:8080" {
			scheme = "http"
		}
		v.Add("redirect_uri", scheme+"://"+r.Host+"/callback")
		v.Add("grant_type", "authorization_code")
		resp, err := http.PostForm(tokenUri, v)
		if err != nil {
			// Output error
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
