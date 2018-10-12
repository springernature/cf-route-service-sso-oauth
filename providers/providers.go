package providers

import (
	"errors"
	"net/http"
	"os"
)

type ProviderData struct {
	Provider     string
	ClientID     string
	ClientSecret string
	AuthUri      string
	TokenUri     string
}

type Provider interface {
	SignIn(http.ResponseWriter, *http.Request)
	Callback(http.ResponseWriter, *http.Request)
}

func InitProviderData() *ProviderData {
	return &ProviderData{
		Provider:     os.Getenv("OAUTHPROVIDER"),
		ClientID:     os.Getenv("OAUTHCLIENTID"),
		ClientSecret: os.Getenv("OAUTHCLIENTSECRET"),
	}
}

func New(p *ProviderData) (Provider, error) {
	switch p.Provider {
	case "google":
		return NewGoogleProvider(p), nil
	case "github":
		return NewGitHubProvider(p), nil
	default:
		return nil, errors.New("No Oauth Provider specified")
	}
}
