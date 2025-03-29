package github

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	startMarker = "<!-- STRAVA-HEATMAP-START -->"
	endMarker   = "<!-- STRAVA-HEATMAP-END -->"
)

// ReadmeUpdater handles updating the GitHub profile README
type ReadmeUpdater struct {
	FilePath string
	Debug    bool
}

// NewReadmeUpdater creates a new README updater
func NewReadmeUpdater(filePath string, debug bool) *ReadmeUpdater {
	return &ReadmeUpdater{
		FilePath: filePath,
		Debug:    debug,
	}
}

// UpdateReadme updates the README with the generated SVG
func (r *ReadmeUpdater) UpdateReadme(svgContent string) error {
	// Read the current README
	content, err := os.ReadFile(r.FilePath)
	if err != nil {
		return fmt.Errorf("error reading README: %w", err)
	}

	contentStr := string(content)

	// Check for markers
	if !strings.Contains(contentStr, startMarker) || !strings.Contains(contentStr, endMarker) {
		return fmt.Errorf("README does not contain required markers: %s and %s", startMarker, endMarker)
	}

	// Create the new content to insert
	newContent := fmt.Sprintf("%s\n%s\n%s", startMarker, svgContent, endMarker)

	// Replace the content between markers
	pattern := fmt.Sprintf("%s[\\s\\S]*?%s", regexp.QuoteMeta(startMarker), regexp.QuoteMeta(endMarker))
	re := regexp.MustCompile(pattern)
	updatedContent := re.ReplaceAllString(contentStr, newContent)

	// Write back to the file
	if err := os.WriteFile(r.FilePath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("error writing updated README: %w", err)
	}

	if r.Debug {
		fmt.Println("[DEBUG] Successfully updated README with Strava heatmap")
	}

	return nil
}

// ValidateReadme checks if the README has the required markers
func (r *ReadmeUpdater) ValidateReadme() (bool, error) {
	// Read the README
	content, err := os.ReadFile(r.FilePath)
	if err != nil {
		return false, fmt.Errorf("error reading README: %w", err)
	}

	contentStr := string(content)

	// Check for markers
	hasStartMarker := strings.Contains(contentStr, startMarker)
	hasEndMarker := strings.Contains(contentStr, endMarker)

	if !hasStartMarker && !hasEndMarker {
		return false, fmt.Errorf("README is missing both required markers: %s and %s", startMarker, endMarker)
	}

	if !hasStartMarker {
		return false, fmt.Errorf("README is missing the start marker: %s", startMarker)
	}

	if !hasEndMarker {
		return false, fmt.Errorf("README is missing the end marker: %s", endMarker)
	}

	return true, nil
}
