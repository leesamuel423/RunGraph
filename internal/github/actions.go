package github

import (
	"fmt"
	"os"
	"time"
)

// ActionsHandler helps with GitHub Actions integration
type ActionsHandler struct {
	Debug bool
}

// NewActionsHandler creates a new GitHub Actions handler
func NewActionsHandler(debug bool) *ActionsHandler {
	return &ActionsHandler{
		Debug: debug,
	}
}

// SetOutput sets a GitHub Actions output variable
func (a *ActionsHandler) SetOutput(name, value string) error {
	// In GitHub Actions, outputs are set by writing to a specific file
	// or using a specific syntax in stdout
	// Here we'll just use the ::set-output syntax for simplicity
	fmt.Printf("::set-output name=%s::%s\n", name, value)

	if a.Debug {
		fmt.Printf("[DEBUG] Set GitHub Actions output: %s=%s\n", name, value)
	}

	return nil
}

// LogError logs an error in a GitHub Actions friendly format
func (a *ActionsHandler) LogError(msg string, err error) {
	// GitHub Actions specific error logging format
	fmt.Printf("::error::%s: %v\n", msg, err)

	if a.Debug {
		fmt.Printf("[DEBUG] Logged error: %s: %v\n", msg, err)
	}
}

// LogWarning logs a warning in a GitHub Actions friendly format
func (a *ActionsHandler) LogWarning(msg string) {
	// GitHub Actions specific warning logging format
	fmt.Printf("::warning::%s\n", msg)

	if a.Debug {
		fmt.Printf("[DEBUG] Logged warning: %s\n", msg)
	}
}

// LogInfo logs an info message in a GitHub Actions friendly format
func (a *ActionsHandler) LogInfo(msg string) {
	// GitHub Actions doesn't have a specific info format,
	// so we'll just print to stdout
	fmt.Println(msg)

	if a.Debug {
		fmt.Printf("[DEBUG] Logged info: %s\n", msg)
	}
}

// GetEnvWithFallback gets an environment variable with a fallback value
func (a *ActionsHandler) GetEnvWithFallback(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		if a.Debug {
			fmt.Printf("[DEBUG] Environment variable %s not set, using fallback: %s\n", key, fallback)
		}
		return fallback
	}
	return value
}

// IsRunningInActions checks if the code is running in GitHub Actions
func (a *ActionsHandler) IsRunningInActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

// RecordMetric records a metric for the GitHub Action
func (a *ActionsHandler) RecordMetric(name string, value interface{}) {
	// This would be used for GitHub Actions step summary or other metrics
	fmt.Printf("::notice title=%s::%v\n", name, value)

	if a.Debug {
		fmt.Printf("[DEBUG] Recorded metric: %s=%v\n", name, value)
	}
}

// CreateSummary adds content to the GitHub Actions step summary
func (a *ActionsHandler) CreateSummary(content string) error {
	// In actual GitHub Actions, this would write to $GITHUB_STEP_SUMMARY
	// For simplicity, we'll just print to stdout
	fmt.Println("\n--- Summary ---")
	fmt.Println(content)
	fmt.Println("---------------")

	return nil
}

// FormatTimestamp formats a timestamp for GitHub Actions logs
func (a *ActionsHandler) FormatTimestamp(t time.Time) string {
	return fmt.Sprintf("%s (UTC)", t.UTC().Format(time.RFC3339))
}
