package providers

import (
	"net/http"
	"os"
)

const (
	SigninPath   string = "/sso/signin"
	CallbackPath string = "/sso/callback"
)

type ProviderData struct {
	ClientID     string
	ClientSecret string
	AuthUri      string
	TokenUri     string
}

type Provider interface {
	Name() string
	SignIn(http.ResponseWriter, *http.Request)
	Redeem(*http.Request) ([]byte, error)
	GetEmail([]byte) (string, error)
	Filter([]byte) (string, error)
}

func InitProviderData() *ProviderData {
	return &ProviderData{
		ClientID:     os.Getenv("OAUTHCLIENTID"),
		ClientSecret: os.Getenv("OAUTHCLIENTSECRET"),
	}
}

// Returns an array with Provider interfaces
func EnabledProviders() []Provider {
	providerInterfaces := make([]Provider, 0)
	providers := []func() (Provider, error){
		// Array with a function for every provider which returns a Provider interface
		GoogleProviderInterface,
	}
	for _, getInterface := range providers {
		i, err := getInterface()
		if err == nil {
			providerInterfaces = append(providerInterfaces, i)
		}
	}
	return providerInterfaces
}

// func New(p *ProviderData) (Provider, error) {
// 	switch p.Provider {
// 	case "google":
// 		return NewGoogleProvider(p), nil
// 	case "github":
// 		return NewGitHubProvider(p), nil
// 	default:
// 		return nil, errors.New("No Oauth Provider specified")
// 	}
// }
