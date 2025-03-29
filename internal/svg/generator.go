package svg

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/samuellee/StravaGraph/internal/config"
	"github.com/samuellee/StravaGraph/internal/processor"
	"github.com/samuellee/StravaGraph/internal/strava"
)

// Generator handles SVG generation
type Generator struct {
	Config *config.Config
	Debug  bool
}

// NewGenerator creates a new SVG generator
func NewGenerator(cfg *config.Config) *Generator {
	return &Generator{
		Config: cfg,
		Debug:  cfg.Debug,
	}
}

// GenerateHeatmap creates a heatmap SVG from activity data
func (g *Generator) GenerateHeatmap(activities []strava.SummaryActivity) (string, error) {
	// Get timezone location
	location, err := g.Config.GetTimeZoneLocation()
	if err != nil && g.Debug {
		// Use stderr to avoid polluting the SVG output
		fmt.Fprintf(os.Stderr, "[DEBUG] %v\n", err)
	}

	// Get date range
	startDate, endDate, err := g.Config.GetDateRange()
	if err != nil {
		return "", fmt.Errorf("error getting date range: %w", err)
	}

	// Create activity aggregator
	aggregator := processor.NewActivityAggregator(activities, location)
	aggregator.Aggregate()

	// Convert map to ordered slice
	orderedDailyData := aggregator.GetOrderedDates(startDate, endDate)

	// Create heatmap data
	heatmapData := NewHeatmapData(
		orderedDailyData,
		startDate,
		endDate,
		g.Config.ColorScheme,
		g.Config.CustomColors,
		g.Config.DarkModeColors,
		g.Config.CellSize,
		g.Config.WeekStart,
		g.Config.DarkModeSupport,
		g.Config.MetricType,
	)

	// Generate SVG
	svgContent := heatmapData.RenderSVG()

	// Add stats if enabled
	if g.Config.ShowStats {
		statsGenerator := processor.NewStatsGenerator(orderedDailyData, startDate, endDate, g.Config.MetricType)
		stats := statsGenerator.GenerateStats()
		
		statsSVG := g.generateStatsSVG(stats)
		
		// Combine heatmap and stats
		svgContent = g.combineHeatmapAndStats(svgContent, statsSVG)
	}

	// Sanity check to ensure we're returning valid SVG
	if !strings.HasPrefix(svgContent, "<svg") {
		if g.Debug {
			fmt.Fprintf(os.Stderr, "[DEBUG] Generated SVG does not start with <svg> tag!\n")
		}
		
		// Try to fix by extracting just the SVG content
		svgIndex := strings.Index(svgContent, "<svg")
		if svgIndex != -1 {
			if g.Debug {
				fmt.Fprintf(os.Stderr, "[DEBUG] Found <svg> tag at position %d, trimming content.\n", svgIndex)
			}
			svgContent = svgContent[svgIndex:]
		}
	}

	// Validate that we have a valid SVG
	if !strings.HasPrefix(svgContent, "<svg") {
		return "", fmt.Errorf("generated content is not a valid SVG (does not start with <svg> tag)")
	}

	return svgContent, nil
}

// generateStatsSVG creates an SVG for statistics
func (g *Generator) generateStatsSVG(stats map[string]interface{}) string {
	// This is a simplified version of the stats SVG generator
	var sb strings.Builder

	// Extract some key stats
	overall, _ := stats["overall"].(*strava.ActivityStats)
	
	// Create a simple stats panel
	width := 300
	height := 200
	
	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		width, height, width, height))
	
	// Add style
	sb.WriteString(`<style>
  .stats-panel { fill: #f6f8fa; stroke: #e1e4e8; rx: 6; }
  .stats-title { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 16px; font-weight: bold; fill: #24292e; }
  .stats-label { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 12px; fill: #586069; }
  .stats-value { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 14px; font-weight: bold; fill: #24292e; }
  .stats-unit { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; font-size: 12px; fill: #586069; }`)
	
	// Add dark mode support if enabled
	if g.Config.DarkModeSupport {
		sb.WriteString(`
  @media (prefers-color-scheme: dark) {
    .stats-panel { fill: #0d1117; stroke: #30363d; }
    .stats-title { fill: #c9d1d9; }
    .stats-label { fill: #8b949e; }
    .stats-value { fill: #c9d1d9; }
    .stats-unit { fill: #8b949e; }
  }`)
	}
	
	sb.WriteString(`
</style>`)
	
	// Stats panel background
	sb.WriteString(fmt.Sprintf(`<rect x="0" y="0" width="%d" height="%d" class="stats-panel" />`, width, height))
	
	// Title
	sb.WriteString(`<text x="15" y="30" class="stats-title">Activity Summary</text>`)
	
	// Stats grid
	if overall != nil {
		// Total activities
		sb.WriteString(`<text x="15" y="60" class="stats-label">Total Activities</text>`)
		sb.WriteString(fmt.Sprintf(`<text x="150" y="60" class="stats-value">%d</text>`, overall.TotalActivities))
		
		// Total distance
		sb.WriteString(`<text x="15" y="85" class="stats-label">Total Distance</text>`)
		sb.WriteString(fmt.Sprintf(`<text x="150" y="85" class="stats-value">%.1f</text>`, overall.TotalDistance))
		sb.WriteString(`<text x="185" y="85" class="stats-unit">km</text>`)
		
		// Total duration
		sb.WriteString(`<text x="15" y="110" class="stats-label">Total Duration</text>`)
		sb.WriteString(fmt.Sprintf(`<text x="150" y="110" class="stats-value">%d</text>`, overall.TotalDuration))
		sb.WriteString(`<text x="170" y="110" class="stats-unit">hours</text>`)
		
		// Active days
		sb.WriteString(`<text x="15" y="135" class="stats-label">Active Days</text>`)
		sb.WriteString(fmt.Sprintf(`<text x="150" y="135" class="stats-value">%d</text>`, overall.ActiveDays))
		
		// Longest streak
		sb.WriteString(`<text x="15" y="160" class="stats-label">Longest Streak</text>`)
		sb.WriteString(fmt.Sprintf(`<text x="150" y="160" class="stats-value">%d</text>`, overall.LongestStreak))
		sb.WriteString(`<text x="170" y="160" class="stats-unit">days</text>`)
		
		// Personal records
		sb.WriteString(`<text x="15" y="185" class="stats-label">Personal Records</text>`)
		sb.WriteString(fmt.Sprintf(`<text x="150" y="185" class="stats-value">%d</text>`, overall.PRCount))
	}
	
	sb.WriteString(`</svg>`)
	
	return sb.String()
}

// combineHeatmapAndStats combines the heatmap and stats SVGs into a single SVG
func (g *Generator) combineHeatmapAndStats(heatmapSVG, statsSVG string) string {
	// Extract width and height from heatmap
	heatmapWidth, heatmapHeight := extractSVGDimensions(heatmapSVG)
	
	// Extract width and height from stats
	statsWidth, statsHeight := extractSVGDimensions(statsSVG)
	
	// Calculate combined dimensions
	totalWidth := heatmapWidth + statsWidth + 10 // Add 10px margin between
	totalHeight := max(heatmapHeight, statsHeight)
	
	// Create combined SVG
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		totalWidth, totalHeight, totalWidth, totalHeight))
	
	// Extract and include heatmap content
	heatmapContent := extractSVGContent(heatmapSVG)
	sb.WriteString(fmt.Sprintf(`<g transform="translate(0, 0)">%s</g>`, heatmapContent))
	
	// Extract and include stats content
	statsContent := extractSVGContent(statsSVG)
	sb.WriteString(fmt.Sprintf(`<g transform="translate(%d, 0)">%s</g>`, heatmapWidth+10, statsContent))
	
	sb.WriteString(`</svg>`)
	
	return sb.String()
}

// Helper function to extract width and height from SVG
func extractSVGDimensions(svg string) (int, int) {
	width := 0
	height := 0
	
	// Use regex to extract dimensions
	widthRegex := regexp.MustCompile(`width="(\d+)"`)
	heightRegex := regexp.MustCompile(`height="(\d+)"`)
	
	widthMatches := widthRegex.FindStringSubmatch(svg)
	if len(widthMatches) > 1 {
		fmt.Sscanf(widthMatches[1], "%d", &width)
	}
	
	heightMatches := heightRegex.FindStringSubmatch(svg)
	if len(heightMatches) > 1 {
		fmt.Sscanf(heightMatches[1], "%d", &height)
	}
	
	return width, height
}

// Helper function to extract content from SVG
func extractSVGContent(svg string) string {
	startIdx := strings.Index(svg, ">")
	endIdx := strings.LastIndex(svg, "</svg>")
	
	if startIdx != -1 && endIdx != -1 {
		return svg[startIdx+1 : endIdx]
	}
	
	return ""
}

// Helper function for max of two ints
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GenerateLocationHeatmap creates a heatmap of activity locations
func (g *Generator) GenerateLocationHeatmap(activities []strava.SummaryActivity, privacyRadius int) (string, error) {
	// Placeholder for future implementation
	// This would generate a map visualization of activity locations
	
	// For now, return a placeholder SVG
	return `<svg width="400" height="300" viewBox="0 0 400 300" xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="400" height="300" fill="#f0f0f0" rx="6" />
  <text x="200" y="150" text-anchor="middle" font-family="sans-serif" font-size="16">
    Location heatmap to be implemented
  </text>
</svg>`, nil
}