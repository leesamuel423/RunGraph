package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/samuellee/StravaGraph/internal/auth"
	"github.com/samuellee/StravaGraph/internal/config"
	"github.com/samuellee/StravaGraph/internal/github"
	"github.com/samuellee/StravaGraph/internal/strava"
	"github.com/samuellee/StravaGraph/internal/svg"
)

const (
	configPath = "config.json"
	readmePath = "README.md"
	envFile    = ".env"
)

// loadEnvFile attempts to load variables from .env file
// It doesn't error if the file doesn't exist, as environment variables
// might be set through other means (especially in GitHub Actions)
func loadEnvFile() {
	// Try to find the .env file in the current directory or parent directories
	dir, err := os.Getwd()
	if err == nil {
		// Start with current directory
		envPath := filepath.Join(dir, envFile)
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
			return
		}

		// Check parent directory
		parentDir := filepath.Dir(dir)
		parentEnvPath := filepath.Join(parentDir, envFile)
		if _, err := os.Stat(parentEnvPath); err == nil {
			_ = godotenv.Load(parentEnvPath)
			return
		}
	}

	// As a last resort, try relative to executable
	_ = godotenv.Load(envFile)
}

func main() {
	// Define commands
	cmdAuth := flag.Bool("auth", false, "Generate authentication instructions")
	cmdUpdate := flag.Bool("update", false, "Update the heatmap in the README")
	cmdGenerate := flag.Bool("generate", false, "Generate SVG without updating README")
	cmdTest := flag.Bool("test", false, "Test configuration and authentication")

	// Parse command line arguments
	flag.Parse()

	// Load environment variables from .env file if it exists
	loadEnvFile()

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize GitHub Actions handler
	actionsHandler := github.NewActionsHandler(cfg.Debug)

	// Execute requested command
	switch {
	case *cmdAuth:
		// Generate authentication instructions
		handleAuthCommand(actionsHandler)

	case *cmdUpdate:
		// Update the heatmap in the README
		handleUpdateCommand(cfg, actionsHandler)

	case *cmdGenerate:
		// Generate SVG without updating README
		handleGenerateCommand(cfg, actionsHandler)

	case *cmdTest:
		// Test configuration and authentication
		handleTestCommand(cfg, actionsHandler)

	default:
		// No command specified
		fmt.Println("Please specify a command. Use -h for help.")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

// handleAuthCommand generates authentication instructions
func handleAuthCommand(actionsHandler *github.ActionsHandler) {
	// Get client ID and secret from environment variables
	clientID := actionsHandler.GetEnvWithFallback("STRAVA_CLIENT_ID", "")
	clientSecret := actionsHandler.GetEnvWithFallback("STRAVA_CLIENT_SECRET", "")

	if clientID == "" || clientSecret == "" {
		fmt.Println("Error: STRAVA_CLIENT_ID and STRAVA_CLIENT_SECRET environment variables must be set.")
		os.Exit(1)
	}

	// Generate and display instructions
	instructions := auth.GetInstructionsForUserAuth(clientID, clientSecret)
	fmt.Println(instructions)
}

// handleUpdateCommand updates the heatmap in the README
func handleUpdateCommand(cfg *config.Config, actionsHandler *github.ActionsHandler) {
	// Authenticate with Strava
	tokenManager, err := getTokenManager(actionsHandler)
	if err != nil {
		actionsHandler.LogError("Failed to authenticate with Strava", err)
		os.Exit(1)
	}

	// Create Strava client
	stravaClient := strava.NewClient(tokenManager, cfg.Debug)

	// Get activity date range
	startDate, endDate, err := cfg.GetDateRange()
	if err != nil {
		actionsHandler.LogError("Failed to get date range", err)
		os.Exit(1)
	}

	if cfg.Debug {
		fmt.Printf("Fetching activities from %s to %s\n",
			startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}

	// Fetch activities
	activities, err := stravaClient.GetAllActivities(startDate, endDate, cfg.ActivityTypes)
	if err != nil {
		actionsHandler.LogError("Failed to fetch activities", err)
		os.Exit(1)
	}

	if cfg.Debug {
		fmt.Printf("Found %d activities\n", len(activities))
	}

	// Generate SVG
	svgGenerator := svg.NewGenerator(cfg)
	svgContent, err := svgGenerator.GenerateHeatmap(activities)
	if err != nil {
		actionsHandler.LogError("Failed to generate heatmap SVG", err)
		os.Exit(1)
	}

	// Update README
	readmeUpdater := github.NewReadmeUpdater(readmePath, cfg.Debug)
	if err := readmeUpdater.UpdateReadme(svgContent); err != nil {
		actionsHandler.LogError("Failed to update README", err)
		os.Exit(1)
	}

	actionsHandler.LogInfo("Successfully updated README with Strava heatmap")

	// Record metrics if in GitHub Actions
	if actionsHandler.IsRunningInActions() {
		actionsHandler.RecordMetric("Activities", len(activities))
		actionsHandler.RecordMetric("UpdateTime", actionsHandler.FormatTimestamp(time.Now()))
	}
}

// handleGenerateCommand generates SVG without updating README
func handleGenerateCommand(cfg *config.Config, actionsHandler *github.ActionsHandler) {
	// Authenticate with Strava
	tokenManager, err := getTokenManager(actionsHandler)
	if err != nil {
		// Write errors to stderr, not stdout
		fmt.Fprintf(os.Stderr, "Error: Failed to authenticate with Strava: %v\n", err)
		os.Exit(1)
	}

	// Create Strava client
	stravaClient := strava.NewClient(tokenManager, cfg.Debug)

	// Get activity date range
	startDate, endDate, err := cfg.GetDateRange()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to get date range: %v\n", err)
		os.Exit(1)
	}

	// Fetch activities
	activities, err := stravaClient.GetAllActivities(startDate, endDate, cfg.ActivityTypes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to fetch activities: %v\n", err)
		os.Exit(1)
	}

	// Generate SVG
	svgGenerator := svg.NewGenerator(cfg)
	svgContent, err := svgGenerator.GenerateHeatmap(activities)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to generate heatmap SVG: %v\n", err)
		os.Exit(1)
	}

	// Verify the SVG starts with an opening tag
	if !strings.HasPrefix(svgContent, "<svg") {
		fmt.Fprintf(os.Stderr, "Warning: Generated SVG doesn't start with <svg> tag!\n")
		svgIndex := strings.Index(svgContent, "<svg")
		if svgIndex != -1 {
			// Extract just the SVG content
			fmt.Fprintf(os.Stderr, "Fixing SVG content...\n")
			svgContent = svgContent[svgIndex:]
		}
	}

	// Print just the SVG content to stdout with no additional output
	fmt.Print(svgContent)
}

// handleTestCommand tests configuration and authentication
func handleTestCommand(cfg *config.Config, actionsHandler *github.ActionsHandler) {
	fmt.Println("Testing configuration and authentication...")

	// Test configuration
	fmt.Println("\nConfiguration:")
	fmt.Printf("  Activity Types: %v\n", cfg.ActivityTypes)
	fmt.Printf("  Metric Type: %s\n", cfg.MetricType)
	fmt.Printf("  Date Range: %s\n", cfg.DateRange)

	// Test date range
	startDate, endDate, err := cfg.GetDateRange()
	if err != nil {
		fmt.Printf("  Date Range Error: %v\n", err)
	} else {
		fmt.Printf("  Start Date: %s\n", startDate.Format("2006-01-02"))
		fmt.Printf("  End Date: %s\n", endDate.Format("2006-01-02"))
	}

	// Test Strava authentication
	fmt.Println("\nStrava Authentication:")
	tokenManager, err := getTokenManager(actionsHandler)
	if err != nil {
		fmt.Printf("  Authentication Error: %v\n", err)
		return
	}

	// Test token refresh
	fmt.Println("  Refreshing token...")
	err = tokenManager.RefreshAccessToken()
	if err != nil {
		fmt.Printf("  Token Refresh Error: %v\n", err)
		return
	}

	// Create Strava client and test connection
	stravaClient := strava.NewClient(tokenManager, cfg.Debug)

	// Get athlete data
	fmt.Println("  Fetching athlete data...")
	athlete, err := stravaClient.GetAthlete()
	if err != nil {
		fmt.Printf("  API Error: %v\n", err)
		return
	}

	// Print athlete info
	if firstName, ok := athlete["firstname"].(string); ok {
		fmt.Printf("  Athlete: %s", firstName)
		if lastName, ok := athlete["lastname"].(string); ok {
			fmt.Printf(" %s", lastName)
		}
		fmt.Println()
	}

	// Test README markers if updating
	fmt.Println("\nREADME Validation:")
	readmeUpdater := github.NewReadmeUpdater(readmePath, cfg.Debug)
	valid, err := readmeUpdater.ValidateReadme()
	if err != nil {
		fmt.Printf("  README Error: %v\n", err)
	} else if valid {
		fmt.Println("  README markers are valid")
	}

	fmt.Println("\nTest completed successfully!")
}

// getTokenManager creates and initializes a token manager
func getTokenManager(actionsHandler *github.ActionsHandler) (*auth.TokenManager, error) {
	// Get credentials from environment variables
	clientID := actionsHandler.GetEnvWithFallback("STRAVA_CLIENT_ID", "")
	clientSecret := actionsHandler.GetEnvWithFallback("STRAVA_CLIENT_SECRET", "")
	refreshToken := actionsHandler.GetEnvWithFallback("STRAVA_REFRESH_TOKEN", "")

	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return nil, fmt.Errorf("STRAVA_CLIENT_ID, STRAVA_CLIENT_SECRET, and STRAVA_REFRESH_TOKEN environment variables must be set")
	}

	// Create token manager
	tokenManager := auth.NewTokenManager(clientID, clientSecret, refreshToken)

	return tokenManager, nil
}
