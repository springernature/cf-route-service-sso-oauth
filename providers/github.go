package providers

import "net/http"

type GitHubProvider struct {
	*ProviderData
}

func NewGitHubProvider(p *ProviderData) *GitHubProvider {
	return &GitHubProvider{
		ProviderData: p,
	}
}

func (p *GitHubProvider) SignIn(w http.ResponseWriter, r *http.Request) {

}

func (p *GitHubProvider) Callback(w http.ResponseWriter, r *http.Request) {

}

func (p *GitHubProvider) Redeem(*http.Request) ([]byte, error) {
	return nil, nil
}

func (p *GitHubProvider) GetEmail([]byte) (string, error) {
	return "", nil
}

func (p *GitHubProvider) Filter([]byte) (bool, error) {
	return false, nil
}
