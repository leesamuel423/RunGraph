package strava

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseURL        = "https://www.strava.com/api/v3"
	activitiesPath = "/athlete/activities"
)

// TokenManager interface defines methods for token management
type TokenManager interface {
	GetAccessToken() (string, error)
	RefreshAccessToken() error
}

// Client handles API communication with Strava
type Client struct {
	httpClient   *http.Client
	tokenManager TokenManager
	debug        bool
}

// NewClient creates a new Strava API client
func NewClient(tokenManager TokenManager, debug bool) *Client {
	return &Client{
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		tokenManager: tokenManager,
		debug:        debug,
	}
}

// makeRequest makes an authenticated request to the Strava API
func (c *Client) makeRequest(method, path string, params url.Values) ([]byte, error) {
	// Get a valid access token
	accessToken, err := c.tokenManager.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Construct the request URL
	reqURL := baseURL + path
	if params != nil {
		reqURL += "?" + params.Encode()
	}

	// Create the request
	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add authorization header
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Check for rate limiting
	if resp.StatusCode == http.StatusTooManyRequests {
		// Extract rate limit reset time
		resetHeader := resp.Header.Get("X-RateLimit-Reset")
		if resetHeader != "" {
			resetTime, err := strconv.ParseInt(resetHeader, 10, 64)
			if err == nil {
				resetTimeFormatted := time.Unix(resetTime, 0).Format(time.RFC3339)
				return nil, fmt.Errorf("rate limit exceeded, reset at %s", resetTimeFormatted)
			}
		}
		return nil, fmt.Errorf("rate limit exceeded")
	}

	// Check for other error responses
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, bodyBytes)
	}

	// Read and return the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}

// GetAthlete gets the authenticated athlete's profile
func (c *Client) GetAthlete() (map[string]interface{}, error) {
	body, err := c.makeRequest("GET", "/athlete", nil)
	if err != nil {
		return nil, err
	}

	var athlete map[string]interface{}
	if err := json.Unmarshal(body, &athlete); err != nil {
		return nil, fmt.Errorf("error parsing athlete data: %w", err)
	}

	return athlete, nil
}

// logDebug logs debug information if debug mode is enabled
func (c *Client) logDebug(message string) {
	if c.debug {
		fmt.Println("[DEBUG]", message)
	}
}
