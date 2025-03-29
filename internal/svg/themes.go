package svg

// ColorTheme represents a set of colors for the heatmap
type ColorTheme struct {
	Name   string
	Colors []string // From lowest to highest intensity, starting with "none"
}

// GetTheme returns a color theme by name or the default theme if not found
func GetTheme(name string, customColors []string) ColorTheme {
	switch name {
	case "github":
		return ColorTheme{
			Name:   "github",
			Colors: []string{"#ebedf0", "#9be9a8", "#40c463", "#30a14e", "#216e39"},
		}
	case "strava":
		return ColorTheme{
			Name:   "strava",
			Colors: []string{"#494950", "#ffd4d1", "#ffad9f", "#fc7566", "#e34a33"},
		}
	case "blue":
		return ColorTheme{
			Name:   "blue",
			Colors: []string{"#ebedf0", "#c0dbf1", "#7ab3e5", "#3282ce", "#0a60b6"},
		}
	case "purple":
		return ColorTheme{
			Name:   "purple",
			Colors: []string{"#ebedf0", "#d9c6ec", "#b888e0", "#9c4acf", "#7222bc"},
		}
	case "custom":
		// Validate custom colors
		if len(customColors) == 5 {
			return ColorTheme{
				Name:   "custom",
				Colors: customColors,
			}
		}
		// If custom colors are invalid, fall back to GitHub theme
		return GetTheme("github", nil)
	default:
		// Default to GitHub theme
		return GetTheme("github", nil)
	}
}

// GetDarkModeTheme returns the dark mode variant of a color theme
func GetDarkModeTheme(lightTheme ColorTheme, customDarkColors []string) ColorTheme {
	// If custom dark mode colors are provided, use them
	if len(customDarkColors) == 5 {
		return ColorTheme{
			Name:   lightTheme.Name + "-dark",
			Colors: customDarkColors,
		}
	}

	// Default dark mode variants for built-in themes
	switch lightTheme.Name {
	case "github":
		return ColorTheme{
			Name:   "github-dark",
			Colors: []string{"#161b22", "#0e4429", "#006d32", "#26a641", "#39d353"},
		}
	case "strava":
		return ColorTheme{
			Name:   "strava-dark",
			Colors: []string{"#36363c", "#7c2c2a", "#a63b33", "#d64c3b", "#fc7566"},
		}
	case "blue":
		return ColorTheme{
			Name:   "blue-dark",
			Colors: []string{"#161b22", "#0d2c4a", "#164879", "#2368a9", "#3282ce"},
		}
	case "purple":
		return ColorTheme{
			Name:   "purple-dark",
			Colors: []string{"#161b22", "#2a184a", "#422873", "#61359c", "#8047c9"},
		}
	case "custom":
		// For custom light theme without custom dark theme, create a darkened version
		// In a real implementation, we'd use color manipulation to create dark variants
		// For simplicity, default to GitHub dark theme
		return ColorTheme{
			Name:   "custom-dark",
			Colors: []string{"#161b22", "#0e4429", "#006d32", "#26a641", "#39d353"},
		}
	default:
		return ColorTheme{
			Name:   "github-dark",
			Colors: []string{"#161b22", "#0e4429", "#006d32", "#26a641", "#39d353"},
		}
	}
}

// ActivityTypeColors provides color mapping for different activity types
func ActivityTypeColors() map[string]string {
	return map[string]string{
		"Run":             "#fc5200", // Strava orange
		"Ride":            "#1eaedb", // Cycling blue
		"Swim":            "#007bff", // Swimming blue
		"Walk":            "#6c757d", // Grey
		"Hike":            "#28a745", // Hiking green
		"WeightTraining":  "#dc3545", // Red
		"Workout":         "#fd7e14", // Orange
		"VirtualRide":     "#17a2b8", // Teal
		"Yoga":            "#6f42c1", // Purple
		"AlpineSki":       "#0056b3", // Dark blue
		"BackcountrySki":  "#563d7c", // Purple blue
		"Canoeing":        "#20c997", // Teal green
		"Crossfit":        "#f8f9fa", // Light grey
		"EBikeRide":       "#b2bec3", // Silver
		"Elliptical":      "#ffcc00", // Yellow
		"Golf":            "#4caf50", // Green
		"Handcycle":       "#9c27b0", // Purple
		"IceSkate":        "#00bcd4", // Cyan
		"InlineSkate":     "#ff4081", // Pink
		"Kayaking":        "#3f51b5", // Indigo
		"NordicSki":       "#673ab7", // Deep purple
		"RockClimbing":    "#ff9800", // Orange
		"RollerSki":       "#9e9e9e", // Medium grey
		"Rowing":          "#2196f3", // Blue
		"Snowboard":       "#00acc1", // Cyan
		"Snowshoe":        "#78909c", // Blue grey
		"StairStepper":    "#ff5722", // Deep orange
		"StandUpPaddling": "#03a9f4", // Light blue
		"Surfing":         "#00bfa5", // Teal
		"VirtualRun":      "#ff7043", // Deep orange light
		// Fallback color for any other activity type
		"default":         "#6c757d", // Grey
	}
}