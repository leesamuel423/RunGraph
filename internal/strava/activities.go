package strava

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// GetActivities retrieves activities for the authenticated athlete
func (c *Client) GetActivities(after, before time.Time, page, perPage int) ([]SummaryActivity, error) {
	if perPage <= 0 {
		perPage = 30 // Default per page
	}
	if page <= 0 {
		page = 1 // Default page
	}

	params := url.Values{}
	params.Add("after", strconv.FormatInt(after.Unix(), 10))
	params.Add("before", strconv.FormatInt(before.Unix(), 10))
	params.Add("page", strconv.Itoa(page))
	params.Add("per_page", strconv.Itoa(perPage))

	if c.debug {
		c.logDebug(fmt.Sprintf("Fetching activities page %d with %d items per page", page, perPage))
		c.logDebug(fmt.Sprintf("Date range: %s to %s", after.Format(time.RFC3339), before.Format(time.RFC3339)))
	}

	body, err := c.makeRequest("GET", activitiesPath, params)
	if err != nil {
		return nil, err
	}

	var activities []SummaryActivity
	if err := json.Unmarshal(body, &activities); err != nil {
		return nil, fmt.Errorf("error parsing activities data: %w", err)
	}

	if c.debug {
		c.logDebug(fmt.Sprintf("Retrieved %d activities", len(activities)))
	}

	return activities, nil
}

// GetAllActivities retrieves all activities within the given time range
func (c *Client) GetAllActivities(after, before time.Time, types []string) ([]SummaryActivity, error) {
	var allActivities []SummaryActivity
	var page int = 1
	const perPage int = 100 // Maximum allowed by Strava API

	if c.debug {
		c.logDebug(fmt.Sprintf("Fetching all activities between %s and %s",
			after.Format("2006-01-02"), before.Format("2006-01-02")))
	}

	// Use a map to quickly check if an activity type is included
	activityTypeMap := make(map[string]bool)
	for _, t := range types {
		activityTypeMap[t] = true
	}

	hasMorePages := true
	for hasMorePages {
		// Get a page of activities
		activities, err := c.GetActivities(after, before, page, perPage)
		if err != nil {
			return nil, fmt.Errorf("error fetching activities (page %d): %w", page, err)
		}

		// If we get fewer than perPage, we've reached the last page
		if len(activities) < perPage {
			hasMorePages = false
		}

		// Filter activities by type if needed
		if len(activityTypeMap) > 0 {
			for _, activity := range activities {
				if activityTypeMap[activity.Type] {
					allActivities = append(allActivities, activity)
				}
			}
		} else {
			// No filtering, add all activities
			allActivities = append(allActivities, activities...)
		}

		// Move to the next page
		page++

		// Implement rate limiting - Strava has a limit of 100 requests per 15 minutes
		// Sleep for 200ms between requests to stay comfortably within limits
		time.Sleep(200 * time.Millisecond)
	}

	if c.debug {
		c.logDebug(fmt.Sprintf("Retrieved a total of %d activities after filtering", len(allActivities)))
	}

	return allActivities, nil
}
