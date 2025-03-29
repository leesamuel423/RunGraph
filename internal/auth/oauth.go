package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/samuellee/StravaGraph/internal/strava"
)

const (
	stravaAuthorizeURL = "https://www.strava.com/oauth/authorize"
)

// OAuthConfig holds OAuth configuration for Strava
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

// NewOAuthConfig creates a new OAuth configuration
func NewOAuthConfig(clientID, clientSecret, redirectURI string, scopes []string) *OAuthConfig {
	return &OAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Scopes:       scopes,
	}
}

// GetAuthorizationURL returns the URL to redirect the user for authorization
func (c *OAuthConfig) GetAuthorizationURL() string {
	params := url.Values{}
	params.Add("client_id", c.ClientID)
	params.Add("redirect_uri", c.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", strings.Join(c.Scopes, ","))

	return fmt.Sprintf("%s?%s", stravaAuthorizeURL, params.Encode())
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (c *OAuthConfig) ExchangeCodeForToken(code string) (*strava.TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", stravaTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating token request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response from token endpoint: %d", resp.StatusCode)
	}

	var tokenResp strava.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("error parsing token response: %w", err)
	}

	return &tokenResp, nil
}

// GetInstructionsForUserAuth returns instructions for manual token acquisition
func GetInstructionsForUserAuth(clientID, clientSecret string) string {
	authURL := fmt.Sprintf("https://www.strava.com/oauth/authorize?client_id=%s&redirect_uri=http://localhost&response_type=code&scope=activity:read_all", clientID)

	instructions := fmt.Sprintf(`
Follow these steps to get your Strava API tokens:

1. Visit this URL in your browser:
%s

2. Click "Authorize" to grant access to your Strava account

3. You'll be redirected to a URL like:
http://localhost?state=&code=AUTHORIZATION_CODE&scope=read,activity:read_all

4. Copy the authorization code from the URL (the AUTHORIZATION_CODE part)

5. Run this curl command to get your refresh token:
curl -X POST https://www.strava.com/oauth/token \
  -F client_id=%s \
  -F client_secret=%s \
  -F grant_type=authorization_code \
  -F code=YOUR_AUTHORIZATION_CODE

6. From the JSON response, copy the "refresh_token" value to use in your configuration
`, authURL, clientID, clientSecret)

	return instructions
}

