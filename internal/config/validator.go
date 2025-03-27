package config

import (
	"fmt"
	"strings"
)

// ValidMetricTypes contains all valid metric types
var ValidMetricTypes = []string{"distance", "duration", "elevation", "effort", "heart_rate"}

// ValidColorSchemes contains all valid color schemes
var ValidColorSchemes = []string{"github", "strava", "blue", "purple", "custom"}

// ValidDateRanges contains all valid date ranges
var ValidDateRanges = []string{"1year", "all", "ytd", "custom"}

// ValidWeekStarts contains all valid week start days
var ValidWeekStarts = []string{"Sunday", "Monday"}

// ValidStatTypes contains all valid statistic types
var ValidStatTypes = []string{"weekly", "monthly", "yearly"}

// ValidateConfig validates the configuration
func ValidateConfig(config *Config) error {
	// Validate required fields
	if len(config.ActivityTypes) == 0 {
		return fmt.Errorf("activityTypes cannot be empty")
	}

	// Validate metric type
	if !contains(ValidMetricTypes, config.MetricType) {
		return fmt.Errorf("invalid metricType: %s, must be one of %v", config.MetricType, ValidMetricTypes)
	}

	// Validate color scheme
	if !contains(ValidColorSchemes, config.ColorScheme) {
		return fmt.Errorf("invalid colorScheme: %s, must be one of %v", config.ColorScheme, ValidColorSchemes)
	}

	// If custom color scheme, validate custom colors
	if config.ColorScheme == "custom" {
		if len(config.CustomColors) != 5 {
			return fmt.Errorf("customColors must contain exactly 5 colors")
		}
		
		for i, color := range config.CustomColors {
			if !isValidHexColor(color) {
				return fmt.Errorf("invalid hex color at position %d: %s", i, color)
			}
		}
	}

	// Validate date range
	if !contains(ValidDateRanges, config.DateRange) {
		return fmt.Errorf("invalid dateRange: %s, must be one of %v", config.DateRange, ValidDateRanges)
	}

	// Validate custom date range if needed
	if config.DateRange == "custom" {
		if config.CustomDateRange.Start == "" || config.CustomDateRange.End == "" {
			return fmt.Errorf("customDateRange must specify both start and end dates")
		}
	}

	// Validate cell size
	if config.CellSize < 5 || config.CellSize > 20 {
		return fmt.Errorf("cellSize must be between 5 and 20")
	}

	// Validate location privacy radius if location heatmap is enabled
	if config.IncludeLocationHeatmap && config.LocationPrivacyRadius < 0 {
		return fmt.Errorf("locationPrivacyRadius cannot be negative")
	}

	// Validate week start
	if !contains(ValidWeekStarts, config.WeekStart) {
		return fmt.Errorf("invalid weekStart: %s, must be one of %v", config.WeekStart, ValidWeekStarts)
	}

	// Validate dark mode colors if dark mode is enabled
	if config.DarkModeSupport {
		if len(config.DarkModeColors) != 5 {
			return fmt.Errorf("darkModeColors must contain exactly 5 colors")
		}
		
		for i, color := range config.DarkModeColors {
			if !isValidHexColor(color) {
				return fmt.Errorf("invalid dark mode hex color at position %d: %s", i, color)
			}
		}
	}

	// Validate stat types if stats are enabled
	if config.ShowStats {
		if len(config.StatTypes) == 0 {
			return fmt.Errorf("statTypes cannot be empty when showStats is true")
		}
		
		for _, statType := range config.StatTypes {
			if !contains(ValidStatTypes, statType) {
				return fmt.Errorf("invalid statType: %s, must be one of %v", statType, ValidStatTypes)
			}
		}
	}

	return nil
}

// Helper function to check if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper function to validate hex color codes
func isValidHexColor(color string) bool {
	// Check if it starts with # and has either 3 or 6 hex characters
	if !strings.HasPrefix(color, "#") {
		return false
	}
	
	hex := strings.TrimPrefix(color, "#")
	return (len(hex) == 6 || len(hex) == 3) && isHex(hex)
}

// Helper function to check if a string is valid hexadecimal
func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}