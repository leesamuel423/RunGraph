package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	stravaTokenURL = "https://www.strava.com/oauth/token"
)

// TokenResponse represents the response from Strava token endpoint
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Athlete      struct {
		ID int64 `json:"id"`
	} `json:"athlete"`
}

// TokenManager handles Strava token management
type TokenManager struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	AccessToken  string
	ExpiresAt    time.Time
}

// NewTokenManager creates a new token manager
func NewTokenManager(clientID, clientSecret, refreshToken string) *TokenManager {
	return &TokenManager{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
	}
}

// GetAccessToken returns a valid access token, refreshing if necessary
func (tm *TokenManager) GetAccessToken() (string, error) {
	// If we don't have an access token or it's expired, refresh it
	if tm.AccessToken == "" || time.Now().After(tm.ExpiresAt) {
		if err := tm.RefreshAccessToken(); err != nil {
			return "", fmt.Errorf("failed to refresh access token: %w", err)
		}
	}
	return tm.AccessToken, nil
}

// RefreshAccessToken refreshes the Strava access token using the refresh token
func (tm *TokenManager) RefreshAccessToken() error {
	data := url.Values{}
	data.Set("client_id", tm.ClientID)
	data.Set("client_secret", tm.ClientSecret)
	data.Set("refresh_token", tm.RefreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", stravaTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating token request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response from token endpoint: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("error parsing token response: %w", err)
	}

	// Update the token manager with the new tokens
	tm.AccessToken = tokenResp.AccessToken
	tm.RefreshToken = tokenResp.RefreshToken
	tm.ExpiresAt = time.Unix(tokenResp.ExpiresAt, 0)

	return nil
}