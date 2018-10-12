package providers

import "os"

type ProviderData struct {
	Provider     string
	ClientID     string
	ClientSecret string
}

func InitProviderData() *ProviderData {
	return &ProviderData{
		Provider:     os.Getenv("OAUTHPROVIDER"),
		ClientID:     os.Getenv("OAUTHCLIENTID"),
		ClientSecret: os.Getenv("OAUTHCLIENTSECRET"),
	}
}
