package twitch

import (
	"github.com/codescot/go-common/httputil"
)

// Twitch API access for Twitch
type Twitch struct {
	ClientID     string
	ClientSecret string

	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// GetCredentials get credentials for authenticated routes
func (t *Twitch) GetCredentials() error {
	req := httputil.HTTP{
		TargetURL: "https://id.twitch.tv/oauth2/token",
		Method:    "POST",
		Headers: map[string]string{
			"Accept": "application/json",
		},
		Form: map[string]string{
			"client_id":     t.ClientID,
			"client_secret": t.ClientSecret,
			"grant_type":    "client_credentials",
		},
	}

	err := req.JSON(t)

	return err
}
