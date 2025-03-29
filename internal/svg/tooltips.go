package svg

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/samuellee/StravaGraph/internal/strava"
)

// TooltipData holds data needed for a tooltip
type TooltipData struct {
	Date           time.Time
	ActivityCount  int
	TotalDistance  float64
	TotalDuration  int
	TotalElevation float64
	ActivityTypes  map[string]int
	HasPR          bool
	CustomFields   map[string]string
}

// NewTooltipData creates tooltip data from a daily activity
func NewTooltipData(activity *strava.DailyActivity) *TooltipData {
	if activity == nil {
		return &TooltipData{
			ActivityCount: 0,
			ActivityTypes: make(map[string]int),
			CustomFields:  make(map[string]string),
		}
	}

	return &TooltipData{
		Date:           activity.Date,
		ActivityCount:  activity.Count,
		TotalDistance:  activity.TotalDistance,
		TotalDuration:  activity.TotalDuration,
		TotalElevation: activity.TotalElevation,
		ActivityTypes:  activity.Types,
		HasPR:          activity.HasPR,
		CustomFields:   make(map[string]string),
	}
}

// AddCustomField adds a custom field to the tooltip
func (t *TooltipData) AddCustomField(key, value string) {
	t.CustomFields[key] = value
}

// GenerateTooltipSVG creates an SVG tooltip
func GenerateTooltipSVG(data *TooltipData) string {
	// If no activities, generate empty day tooltip
	if data.ActivityCount == 0 {
		return generateEmptyTooltip(data.Date)
	}

	var sb strings.Builder

	// Tooltip size - will adjust based on content
	width := 200
	padding := 10
	lineHeight := 18

	// Calculate number of lines for sizing
	lines := 3 // Date and activity count + 1 empty line
	if data.TotalDistance > 0 {
		lines++
	}
	if data.TotalDuration > 0 {
		lines++
	}
	if data.TotalElevation > 0 {
		lines++
	}
	if data.HasPR {
		lines++
	}
	if len(data.ActivityTypes) > 0 {
		lines += min(len(data.ActivityTypes), 3)
	}
	for range data.CustomFields {
		lines++
	}

	height := (lines * lineHeight) + (padding * 2)

	// Start SVG tooltip
	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		width, height, width, height))

	// Add style
	sb.WriteString(`<style>
  .tooltip-bg { fill: white; stroke: #ddd; rx: 4; }
  .tooltip-title { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 13px; font-weight: bold; fill: #24292e; }
  .tooltip-text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 12px; fill: #586069; }
  .tooltip-highlight { fill: #ff8c00; }
  @media (prefers-color-scheme: dark) {
    .tooltip-bg { fill: #161b22; stroke: #30363d; }
    .tooltip-title { fill: #c9d1d9; }
    .tooltip-text { fill: #8b949e; }
  }
</style>`)

	// Background
	sb.WriteString(fmt.Sprintf(`<rect x="0" y="0" width="%d" height="%d" class="tooltip-bg" />`, width, height))

	// Title - date
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="tooltip-title">%s</text>`,
		padding, padding+lineHeight, data.Date.Format("Monday, January 2, 2006")))

	// Activity count
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="tooltip-text">%d %s</text>`,
		padding, padding+(lineHeight*2),
		data.ActivityCount, pluralize("activity", data.ActivityCount)))

	currentLine := 3

	// Distance
	if data.TotalDistance > 0 {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="tooltip-text">%.1f km total distance</text>`,
			padding, padding+(lineHeight*currentLine), data.TotalDistance/1000))
		currentLine++
	}

	// Duration
	if data.TotalDuration > 0 {
		hours := data.TotalDuration / 3600
		minutes := (data.TotalDuration % 3600) / 60

		durationText := ""
		if hours > 0 {
			durationText = fmt.Sprintf("%d %s %d %s",
				hours, pluralize("hour", hours),
				minutes, pluralize("minute", minutes))
		} else {
			durationText = fmt.Sprintf("%d %s", minutes, pluralize("minute", minutes))
		}

		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="tooltip-text">%s total time</text>`,
			padding, padding+(lineHeight*currentLine), durationText))
		currentLine++
	}

	// Elevation
	if data.TotalElevation > 0 {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="tooltip-text">%.0f m elevation gain</text>`,
			padding, padding+(lineHeight*currentLine), data.TotalElevation))
		currentLine++
	}

	// Personal Record
	if data.HasPR {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="tooltip-text tooltip-highlight">Personal Record!</text>`,
			padding, padding+(lineHeight*currentLine)))
		currentLine++
	}

	// Activity types
	if len(data.ActivityTypes) > 0 {
		// Convert to sorted slice for consistent display
		type activityTypeCount struct {
			Type  string
			Count int
		}

		var activityTypes []activityTypeCount
		for t, count := range data.ActivityTypes {
			activityTypes = append(activityTypes, activityTypeCount{t, count})
		}

		// Sort by count descending
		sort.Slice(activityTypes, func(i, j int) bool {
			return activityTypes[i].Count > activityTypes[j].Count
		})

		// Show up to 3 activity types
		maxTypes := min(len(activityTypes), 3)
		for i := 0; i < maxTypes; i++ {
			typeData := activityTypes[i]
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="tooltip-text">%d %s</text>`,
				padding, padding+(lineHeight*currentLine),
				typeData.Count, typeData.Type))
			currentLine++
		}
	}

	// Custom fields
	for key, value := range data.CustomFields {
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="tooltip-text">%s: %s</text>`,
			padding, padding+(lineHeight*currentLine), key, value))
		currentLine++
	}

	sb.WriteString(`</svg>`)

	return sb.String()
}

// generateEmptyTooltip creates a tooltip for days with no activities
func generateEmptyTooltip(date time.Time) string {
	var sb strings.Builder

	// Tooltip size
	width := 200
	height := 60
	padding := 10
	lineHeight := 18

	// Start SVG
	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		width, height, width, height))

	// Add style
	sb.WriteString(`<style>
  .tooltip-bg { fill: white; stroke: #ddd; rx: 4; }
  .tooltip-title { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 13px; font-weight: bold; fill: #24292e; }
  .tooltip-text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 12px; fill: #586069; }
  @media (prefers-color-scheme: dark) {
    .tooltip-bg { fill: #161b22; stroke: #30363d; }
    .tooltip-title { fill: #c9d1d9; }
    .tooltip-text { fill: #8b949e; }
  }
</style>`)

	// Background
	sb.WriteString(fmt.Sprintf(`<rect x="0" y="0" width="%d" height="%d" class="tooltip-bg" />`, width, height))

	// Date
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="tooltip-title">%s</text>`,
		padding, padding+lineHeight, date.Format("Monday, January 2, 2006")))

	// No activities message
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="tooltip-text">No activities on this day</text>`,
		padding, padding+(lineHeight*2)))

	sb.WriteString(`</svg>`)

	return sb.String()
}

// Helper function for minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
