package svg

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/samuellee/StravaGraph/internal/strava"
)

// HeatmapCell represents a single cell in the heatmap
type HeatmapCell struct {
	Date      time.Time
	Intensity strava.HeatmapIntensity
	HasPR     bool
	Count     int
	Tooltip   string
}

// HeatmapData holds all data needed to generate the heatmap
type HeatmapData struct {
	StartDate  time.Time
	EndDate    time.Time
	Cells      [][]*HeatmapCell // [week][day]
	WeekLabels []string
	MonthLabels []struct {
		Month string
		X     int
	}
	ColorTheme      ColorTheme
	DarkModeTheme   ColorTheme
	CellSize        int
	CellSpacing     int
	WeekStart       string // "Sunday" or "Monday"
	DarkModeSupport bool
}

// NewHeatmapData creates a new heatmap data structure
func NewHeatmapData(
	activities []*strava.DailyActivity,
	startDate, endDate time.Time,
	colorScheme string,
	customColors []string,
	darkModeColors []string,
	cellSize int,
	weekStart string,
	darkModeSupport bool,
	metricType string,
) *HeatmapData {
	// Get color themes
	theme := GetTheme(colorScheme, customColors)
	darkTheme := GetDarkModeTheme(theme, darkModeColors)

	// Default values
	if cellSize < 5 {
		cellSize = 11 // GitHub default
	}
	if weekStart != "Sunday" && weekStart != "Monday" {
		weekStart = "Monday" // Default to Monday
	}
	cellSpacing := 2

	// Initialize heatmap data
	heatmap := &HeatmapData{
		StartDate:       startDate,
		EndDate:         endDate,
		ColorTheme:      theme,
		DarkModeTheme:   darkTheme,
		CellSize:        cellSize,
		CellSpacing:     cellSpacing,
		WeekStart:       weekStart,
		DarkModeSupport: darkModeSupport,
	}

	// Create week and day grid
	heatmap.createGrid(activities, metricType)
	heatmap.generateLabels()

	return heatmap
}

// dayOffset returns the day offset based on the configured week start
func (h *HeatmapData) dayOffset(day time.Weekday) int {
	if h.WeekStart == "Monday" {
		// Monday = 0, Tuesday = 1, ..., Sunday = 6
		if day == time.Sunday {
			return 6
		}
		return int(day) - 1
	}
	// Sunday = 0, Monday = 1, ..., Saturday = 6
	return int(day)
}

// createGrid creates the grid of cells for the heatmap
func (h *HeatmapData) createGrid(activities []*strava.DailyActivity, metricType string) {
	// Map of activities by date
	activityMap := make(map[string]*strava.DailyActivity)
	for _, activity := range activities {
		dateKey := activity.Date.Format("2006-01-02")
		activityMap[dateKey] = activity
	}

	// Calculate the number of weeks needed
	startOffset := h.dayOffset(h.StartDate.Weekday())
	endOffset := 6 - h.dayOffset(h.EndDate.Weekday())
	
	totalDays := int(h.EndDate.Sub(h.StartDate).Hours()/24) + 1
	totalWeeks := (totalDays + startOffset + endOffset) / 7
	if (totalDays+startOffset+endOffset)%7 > 0 {
		totalWeeks++
	}

	// Create the grid
	h.Cells = make([][]*HeatmapCell, totalWeeks)
	for i := range h.Cells {
		h.Cells[i] = make([]*HeatmapCell, 7)
	}

	// Fill the grid with days
	current := h.StartDate.AddDate(0, 0, -startOffset)
	for week := 0; week < totalWeeks; week++ {
		for day := 0; day < 7; day++ {
			dateKey := current.Format("2006-01-02")
			
			// Check if we have activity data for this day
			activity, exists := activityMap[dateKey]
			
			// Calculate intensity
			var intensity strava.HeatmapIntensity
			hasPR := false
			count := 0
			
			if exists && activity.Count > 0 {
				// Determine intensity based on metric type
				intensity = calculateIntensity(activity, metricType, activities)
				hasPR = activity.HasPR
				count = activity.Count
			}
			
			// Create tooltip
			tooltip := createTooltip(current, activity)
			
			// Create the cell
			h.Cells[week][day] = &HeatmapCell{
				Date:      current,
				Intensity: intensity,
				HasPR:     hasPR,
				Count:     count,
				Tooltip:   tooltip,
			}
			
			// Move to next day
			current = current.AddDate(0, 0, 1)
		}
	}
}

// generateLabels creates week and month labels for the heatmap
func (h *HeatmapData) generateLabels() {
	// Week labels (for y-axis)
	h.WeekLabels = make([]string, len(h.Cells))
	for i := range h.Cells {
		if i%2 == 0 { // Only label every other week to avoid clutter
			week := h.Cells[i][0].Date
			h.WeekLabels[i] = week.Format("Jan 2")
		}
	}

	// Month labels (for x-axis)
	var monthLabels []struct {
		Month string
		X     int
	}

	// Group by month and find the first week of each month
	currentMonth := -1
	for week := 0; week < len(h.Cells); week++ {
		for day := 0; day < 7; day++ {
			cell := h.Cells[week][day]
			month := cell.Date.Month()
			
			// If we found a new month, add a label
			if int(month) != currentMonth {
				currentMonth = int(month)
				monthLabels = append(monthLabels, struct {
					Month string
					X     int
				}{
					Month: cell.Date.Format("Jan"),
					X:     week,
				})
				break
			}
		}
	}

	h.MonthLabels = monthLabels
}

// RenderSVG generates the SVG for the heatmap
func (h *HeatmapData) RenderSVG() string {
	// Make the heatmap much wider by displaying more days per row and spacing them out more
	cellsPerRow := 26 // Increase this number to fit more cells in a row
	widthPadding := 60 // Increase side padding for labels
	
	// Increase spacing between cells for better readability
	h.CellSpacing = 4
	
	totalWidth := (cellsPerRow * (h.CellSize + h.CellSpacing)) + widthPadding
	totalHeight := ((totalDaysInYear() / cellsPerRow) + 1) * (h.CellSize + h.CellSpacing) + 50 // +50 for labels

	var sb strings.Builder

	// SVG header
	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		totalWidth, totalHeight, totalWidth, totalHeight))

	// Add style
	h.writeStyle(&sb)

	// Write month labels
	h.writeMonthLabels(&sb)

	// Write week labels
	h.writeWeekLabels(&sb)

	// Write cells
	h.writeCells(&sb, totalWidth)

	// Add legend
	h.writeLegend(&sb, totalWidth)

	// Close SVG
	sb.WriteString(`</svg>`)

	return sb.String()
}

// writeStyle adds the CSS style to the SVG
func (h *HeatmapData) writeStyle(sb *strings.Builder) {
	sb.WriteString(`<style>
  .heatmap-cell { rx: 2; }
  .heatmap-label { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 10px; fill: #767676; }
  .heatmap-month-label { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 14px; font-weight: bold; fill: #24292e; }
  .heatmap-legend-text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 12px; fill: #24292e; font-weight: bold; }
  .heatmap-tooltip { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 12px; pointer-events: none; filter: drop-shadow(0px 0px 2px rgba(0,0,0,0.2)); opacity: 0; transition: opacity 0.2s; }
  .heatmap-cell:hover + .heatmap-tooltip { opacity: 1; }
  .heatmap-tooltip-rect { fill: white; stroke: #ddd; rx: 3; }
  .heatmap-tooltip-text { font-size: 11px; fill: #333; }
  .heatmap-tooltip-header { font-weight: bold; }
  .pr-marker { fill: #ff8c00; }`)

	// Add dark mode support if enabled
	if h.DarkModeSupport {
		sb.WriteString(`
  @media (prefers-color-scheme: dark) {
    .heatmap-label { fill: #8b949e; }
    .heatmap-legend-text { fill: #8b949e; }
    .heatmap-tooltip-rect { fill: #161b22; stroke: #30363d; }
    .heatmap-tooltip-text { fill: #c9d1d9; }
  }`)
	}

	// Add color classes based on theme
	for i := 0; i < 5; i++ {
		sb.WriteString(fmt.Sprintf(`
  .intensity-%d { fill: %s; }`, i, h.ColorTheme.Colors[i]))
	}

	// Add dark mode color classes if enabled
	if h.DarkModeSupport {
		sb.WriteString(`
  @media (prefers-color-scheme: dark) {`)
		for i := 0; i < 5; i++ {
			sb.WriteString(fmt.Sprintf(`
    .intensity-%d { fill: %s; }`, i, h.DarkModeTheme.Colors[i]))
		}
		sb.WriteString(`
  }`)
	}

	sb.WriteString(`
</style>`)
}

// writeMonthLabels adds month labels to the SVG
func (h *HeatmapData) writeMonthLabels(sb *strings.Builder) {
	sb.WriteString(`<g class="heatmap-month-labels">`)
	
	// Get the number of cells per row from totalWidth
	cellsPerRow := (25 * (h.CellSize + h.CellSpacing))
	
	// Define all months for clear labeling
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	
	// Calculate spacing to distribute months evenly
	spacing := cellsPerRow / 12
	leftPadding := 50
	
	// Add all months with even spacing
	for i, month := range months {
		x := leftPadding + (i * spacing)
		y := 15 // Top margin for month labels
		
		// Make the font larger and bolder for better visibility
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="heatmap-month-label">%s</text>`,
			x, y, month))
	}
	
	sb.WriteString(`</g>`)
}

// writeWeekLabels adds week labels to the SVG
func (h *HeatmapData) writeWeekLabels(sb *strings.Builder) {
	sb.WriteString(`<g class="heatmap-week-labels">`)
	
	// We're removing the confusing week labels since they weren't clear
	// Instead, we're using increased left padding and focusing on the month labels
	
	sb.WriteString(`</g>`)
}

// writeCells adds all cells to the SVG
func (h *HeatmapData) writeCells(sb *strings.Builder, totalWidth int) {
	sb.WriteString(`<g class="heatmap-cells">`)
	
	// Calculate total days to display
	totalDays := len(h.Cells) * 7
	
	// Get the number of cells per row from the width
	cellsPerRow := (totalWidth - 60) / (h.CellSize + h.CellSpacing)
	
	// To display the cells in a more spread out horizontal layout
	for i := 0; i < totalDays; i++ {
		// Determine which cell from the original grid to use
		week := i / 7
		day := i % 7
		
		// Skip if outside the array bounds
		if week >= len(h.Cells) || day >= 7 {
			continue
		}
		
		cell := h.Cells[week][day]
		
		// Skip days outside our date range
		if cell.Date.Before(h.StartDate) || cell.Date.After(h.EndDate) {
			continue
		}
		
		// Calculate x,y position for a more stretched layout
		newRow := i / cellsPerRow
		newCol := i % cellsPerRow
		
		// Add more left padding for a better layout
		leftPadding := 50
		
		x := (newCol * (h.CellSize + h.CellSpacing)) + leftPadding
		y := (newRow * (h.CellSize + h.CellSpacing)) + 30 // Increased top padding for month labels
			
		// Determine fill color based on intensity
		colorClass := fmt.Sprintf("intensity-%d", cell.Intensity)
		
		// Add cell
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" class="heatmap-cell %s" data-date="%s" data-count="%d">`,
			x, y, h.CellSize, h.CellSize, colorClass, cell.Date.Format("2006-01-02"), cell.Count))
		sb.WriteString(fmt.Sprintf(`<title>%s</title></rect>`, cell.Tooltip))
		
		// Add PR marker if applicable
		if cell.HasPR {
			prX := x + (h.CellSize * 3/4)
			prY := y + (h.CellSize * 1/4)
			prRadius := h.CellSize / 6
			
			sb.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" class="pr-marker" />`,
				prX, prY, prRadius))
		}
		
		// Add tooltip for hover
		tooltipWidth := 200
		tooltipHeight := 80
		tooltipX := x + h.CellSize + 5
		tooltipY := y
		
		// If tooltip would go off right edge, place it to the left of the cell
		if tooltipX + tooltipWidth > totalWidth {
			tooltipX = x - tooltipWidth - 5
		}
		
		sb.WriteString(fmt.Sprintf(`<g class="heatmap-tooltip" transform="translate(%d, %d)">`,
			tooltipX, tooltipY))
		
		sb.WriteString(fmt.Sprintf(`<rect x="0" y="0" width="%d" height="%d" class="heatmap-tooltip-rect" />`,
			tooltipWidth, tooltipHeight))
		
		// Only add detailed tooltip content if there are activities
		if cell.Count > 0 {
			// We'll use a simplified tooltip for now
			sb.WriteString(fmt.Sprintf(`<text x="10" y="15" class="heatmap-tooltip-text heatmap-tooltip-header">%s</text>`,
				cell.Date.Format("January 2, 2006")))
			
			sb.WriteString(fmt.Sprintf(`<text x="10" y="35" class="heatmap-tooltip-text">%d activities</text>`,
				cell.Count))
			
			if cell.HasPR {
				sb.WriteString(`<text x="10" y="55" class="heatmap-tooltip-text" fill="#ff8c00">Personal Record!</text>`)
			}
		} else {
			sb.WriteString(fmt.Sprintf(`<text x="10" y="25" class="heatmap-tooltip-text">No activities on %s</text>`,
				cell.Date.Format("January 2, 2006")))
		}
		
		sb.WriteString(`</g>`)
	}
	
	sb.WriteString(`</g>`)
}

// writeLegend adds the color legend to the SVG
func (h *HeatmapData) writeLegend(sb *strings.Builder, totalWidth int) {
	// Calculate total days and rows for proper positioning
	totalDays := len(h.Cells) * 7
	cellsPerRow := (totalWidth - 60) / (h.CellSize + h.CellSpacing)
	numRows := (totalDays + cellsPerRow - 1) / cellsPerRow
	
	// Position legend just below the last row of cells with minimal gap
	legendY := (numRows * (h.CellSize + h.CellSpacing)) + 40
	
	// Center the legend
	leftPadding := 50
	legendWidth := 5 * (h.CellSize + 2) + 80 // space for boxes + labels
	
	sb.WriteString(fmt.Sprintf(`<g class="heatmap-legend" transform="translate(%d, %d)">`, 
		leftPadding + ((cellsPerRow * (h.CellSize + h.CellSpacing)) - legendWidth)/2, legendY))
	
	// Legend label
	sb.WriteString(`<text x="0" y="5" class="heatmap-legend-text" text-anchor="start">Less</text>`)
	
	// Legend boxes - increase size for better visibility
	boxSize := h.CellSize + 4 // Make boxes slightly larger
	for i := 0; i < 5; i++ {
		x := 40 + (i * (boxSize + 4))
		
		colorClass := fmt.Sprintf("intensity-%d", i)
		
		sb.WriteString(fmt.Sprintf(`<rect x="%d" y="0" width="%d" height="%d" class="heatmap-cell %s" />`,
			x, boxSize, boxSize, colorClass))
	}
	
	// More label
	sb.WriteString(fmt.Sprintf(`<text x="%d" y="5" class="heatmap-legend-text" text-anchor="start">More</text>`,
		40 + (5 * (boxSize + 4)) + 5))
	
	sb.WriteString(`</g>`)
}

// Helper function to calculate intensity for a day
func calculateIntensity(day *strava.DailyActivity, metricType string, allActivities []*strava.DailyActivity) strava.HeatmapIntensity {
	if day.Count == 0 {
		return strava.None
	}

	// Get all non-zero values for this metric to calculate percentiles
	var values []float64
	for _, data := range allActivities {
		if data.Count == 0 {
			continue
		}

		var value float64
		switch metricType {
		case "distance":
			value = data.TotalDistance
		case "duration":
			value = float64(data.TotalDuration)
		case "elevation":
			value = data.TotalElevation
		case "heart_rate":
			value = data.AvgHeartRate
		case "effort":
			// Simple effort formula: distance * elevation gain / duration
			// This rewards activities with higher distance, more elevation, but shorter time
			if data.TotalDuration > 0 {
				value = (data.TotalDistance * (1 + data.TotalElevation/100)) / float64(data.TotalDuration)
			}
		default:
			value = float64(data.Count) // Default to count-based intensity
		}

		if value > 0 {
			values = append(values, value)
		}
	}

	// If no values, return low intensity for any day with activity
	if len(values) == 0 {
		return strava.Low
	}

	// Get the value for this day
	var dayValue float64
	switch metricType {
	case "distance":
		dayValue = day.TotalDistance
	case "duration":
		dayValue = float64(day.TotalDuration)
	case "elevation":
		dayValue = day.TotalElevation
	case "heart_rate":
		dayValue = day.AvgHeartRate
	case "effort":
		if day.TotalDuration > 0 {
			dayValue = (day.TotalDistance * (1 + day.TotalElevation/100)) / float64(day.TotalDuration)
		}
	default:
		dayValue = float64(day.Count)
	}

	// Simple percentile-based binning
	// Here we're using a simple algorithm for demonstration
	// In a real implementation, we might use more sophisticated statistical methods
	
	// Sort values in ascending order
	sort.Float64s(values)
	
	// Find position of the day's value
	pos := sort.SearchFloat64s(values, dayValue)
	percentile := float64(pos) / float64(len(values))
	
	if percentile <= 0.25 {
		return strava.Low
	} else if percentile <= 0.5 {
		return strava.Medium
	} else if percentile <= 0.75 {
		return strava.High
	} else {
		return strava.VeryHigh
	}
}

// Helper function to create a tooltip for a day
func createTooltip(date time.Time, activity *strava.DailyActivity) string {
	if activity == nil || activity.Count == 0 {
		return fmt.Sprintf("No activities on %s", date.Format("Jan 2, 2006"))
	}

	// Format distance in km
	distance := activity.TotalDistance / 1000
	
	// Format duration in hours and minutes
	hours := activity.TotalDuration / 3600
	minutes := (activity.TotalDuration % 3600) / 60
	
	tooltip := fmt.Sprintf("%s: %d %s", 
		date.Format("Jan 2, 2006"),
		activity.Count,
		pluralize("activity", activity.Count))
	
	if distance > 0 {
		tooltip += fmt.Sprintf("\nTotal distance: %.1f km", distance)
	}
	
	if activity.TotalDuration > 0 {
		if hours > 0 {
			tooltip += fmt.Sprintf("\nTotal time: %d %s %d %s", 
				hours, pluralize("hour", hours), 
				minutes, pluralize("minute", minutes))
		} else {
			tooltip += fmt.Sprintf("\nTotal time: %d %s", 
				minutes, pluralize("minute", minutes))
		}
	}
	
	if activity.TotalElevation > 0 {
		tooltip += fmt.Sprintf("\nTotal elevation: %.0f m", activity.TotalElevation)
	}
	
	if activity.HasPR {
		tooltip += "\nPersonal Record!"
	}
	
	return tooltip
}

// Helper function to pluralize words
func pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

// Helper function to calculate total days in a year (approximate)
func totalDaysInYear() int {
	return 365
}