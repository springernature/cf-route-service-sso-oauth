package providers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type GoogleProvider struct {
	*ProviderData
}

func NewGoogleProvider(p *ProviderData) *GoogleProvider {
	// Declare default endpoints
	authUri := "https://accounts.google.com/o/oauth2/auth"
	tokenUri := "https://www.googleapis.com/oauth2/v3/token"
	// Check if there is an override for the default endpoints
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
	http.Redirect(w, r, p.AuthUri+"?"+
		"client_id="+p.ClientID+"&"+
		"response_type=code"+"&"+
		"scope=openid%20email"+"&"+
		"redirect_uri=https://"+r.Host+CallbackPath, 302)
}

func (p *GoogleProvider) Redeem(r *http.Request) ([]byte, error) {
	// Check if the callback contains an authorization code
	if code := r.FormValue("code"); code != "" {
		// Exchange code for access token and ID token
		v := url.Values{}
		v.Add("code", code)
		v.Add("client_id", p.ClientID)
		v.Add("client_secret", p.ClientSecret)
		v.Add("redirect_uri", "https://"+r.Host+CallbackPath)
		v.Add("grant_type", "authorization_code")
		resp, err := http.PostForm(p.TokenUri, v)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, errors.New(`The POST to the token endpoint did not return HTTP 200 OK. ` +
				`Returned body is: ` + string(body[:]))
		}
		// All is good. Return the body as byte array
		return body, nil
	}
	return nil, errors.New("Callback does not contain authorization code")
}

func (p *GoogleProvider) GetEmail(b []byte) (string, error) {
	type Token struct {
		IDToken string `json:"id_token"`
	}
	var t Token
	if err := json.Unmarshal(b, &t); err != nil {
		return "", err
	}
	// Get the second part (payload) of the JWT (IDToken)
	jwt := strings.TrimSuffix(strings.Split(t.IDToken, ".")[1], "=")
	// Decode this base64 string to byte array
	jwtBytes, err := base64.RawURLEncoding.DecodeString(jwt)
	if err != nil {
		return "", err
	}
	type Payload struct {
		Email string `json:"email"`
	}
	var payload Payload
	if err := json.Unmarshal(jwtBytes, &payload); err != nil {
		return "", err
	}
	// Return username
	return payload.Email, nil
}

func (p *GoogleProvider) Filter([]byte) (string, error) {
	// No additional authorization filter
	return "", nil
}
