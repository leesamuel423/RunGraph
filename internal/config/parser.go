package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config represents the application configuration
type Config struct {
	ActivityTypes        []string `json:"activityTypes"`
	MetricType           string   `json:"metricType"`
	ColorScheme          string   `json:"colorScheme"`
	CustomColors         []string `json:"customColors"`
	ShowStats            bool     `json:"showStats"`
	StatTypes            []string `json:"statTypes"`
	DateRange            string   `json:"dateRange"`
	CustomDateRange      struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"customDateRange"`
	CellSize              int      `json:"cellSize"`
	IncludePRs            bool     `json:"includePRs"`
	IncludeLocationHeatmap bool     `json:"includeLocationHeatmap"`
	LocationPrivacyRadius int      `json:"locationPrivacyRadius"`
	DarkModeSupport       bool     `json:"darkModeSupport"`
	DarkModeColors        []string `json:"darkModeColors"`
	WeekStart             string   `json:"weekStart"`
	Language              string   `json:"language"`
	TimeZone              string   `json:"timeZone"`
	Debug                 bool     `json:"debug"`
}

// LoadConfig loads the configuration from the specified file
func LoadConfig(filePath string) (*Config, error) {
	// Read the config file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse the configuration
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// Validate the configuration
	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the specified file
func SaveConfig(config *Config, filePath string) error {
	// Marshal the configuration to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	// Write to the file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// GetTimeZoneLocation returns the time.Location for the configured timezone
func (c *Config) GetTimeZoneLocation() (*time.Location, error) {
	loc, err := time.LoadLocation(c.TimeZone)
	if err != nil {
		// Default to UTC if the timezone is invalid
		return time.UTC, fmt.Errorf("invalid timezone: %s, defaulting to UTC", c.TimeZone)
	}
	return loc, nil
}

// GetDateRange returns the start and end time for the configured date range
func (c *Config) GetDateRange() (time.Time, time.Time, error) {
	loc, err := c.GetTimeZoneLocation()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	now := time.Now().In(loc)
	end := now

	switch c.DateRange {
	case "1year":
		// One year ago from today
		start := now.AddDate(-1, 0, 0)
		return start, end, nil
	case "ytd":
		// Start of current year
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
		return start, end, nil
	case "all":
		// Use a far past date for "all" - Strava was founded in 2009
		start := time.Date(2009, 1, 1, 0, 0, 0, 0, loc)
		return start, end, nil
	case "custom":
		// Parse custom date range
		start, err := time.ParseInLocation("2006-01-02", c.CustomDateRange.Start, loc)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid custom start date: %w", err)
		}

		end, err := time.ParseInLocation("2006-01-02", c.CustomDateRange.End, loc)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid custom end date: %w", err)
		}

		// Set end time to end of day
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, loc)

		return start, end, nil
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date range: %s", c.DateRange)
	}
}