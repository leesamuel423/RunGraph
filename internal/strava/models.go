package strava

import (
	"time"
)

// TokenResponse represents the response from Strava token endpoint
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Athlete      struct {
		ID int64 `json:"id"`
	} `json:"athlete"`
}

// SummaryActivity represents a summary of an activity from Strava API
type SummaryActivity struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Distance      float64   `json:"distance"`        // In meters
	MovingTime    int       `json:"moving_time"`     // In seconds
	ElapsedTime   int       `json:"elapsed_time"`    // In seconds
	TotalElevGain float64   `json:"total_elevation_gain"` // In meters
	Type          string    `json:"type"`
	StartDate     time.Time `json:"start_date"`
	StartDateLocal time.Time `json:"start_date_local"`
	Timezone      string    `json:"timezone"`
	AchievementCount int    `json:"achievement_count"`
	PR            bool      `json:"pr_count,omitempty"` // If PR count > 0, it's a PR
	AverageHeartrate float64 `json:"average_heartrate,omitempty"`
	MaxHeartrate float64     `json:"max_heartrate,omitempty"`
	StartLatlng []float64    `json:"start_latlng,omitempty"`
	EndLatlng   []float64    `json:"end_latlng,omitempty"`
	Map struct {
		SummaryPolyline string `json:"summary_polyline"`
	} `json:"map,omitempty"`
}

// DailyActivity represents aggregated activities for a single day
type DailyActivity struct {
	Date         time.Time
	Count        int
	TotalDistance float64     // In meters
	TotalDuration int         // In seconds
	TotalElevation float64    // In meters
	Activities   []int64      // IDs of activities on this day
	MaxHeartRate float64      // Max heart rate among all activities
	AvgHeartRate float64      // Average heart rate across all activities
	HasPR        bool         // True if any activity on this day has a PR
	Types        map[string]int // Count of each activity type
}

// HeatmapIntensity represents the intensity level for the heatmap cell
type HeatmapIntensity int

const (
	None HeatmapIntensity = iota
	Low
	Medium
	High
	VeryHigh
)

// ActivityStats represents summary statistics about activities
type ActivityStats struct {
	TotalActivities int
	TotalDistance   float64 // In kilometers
	TotalDuration   int     // In hours
	TotalElevation  float64 // In meters
	ActivityTypes   map[string]int
	PRCount         int
	ActiveDays      int
	LongestStreak   int
}

// DatePeriodStats represents statistics for a specific time period
type DatePeriodStats struct {
	Period        string // "weekly", "monthly", "yearly"
	TotalDistance float64
	TotalDuration int
	TotalElevation float64
	ActivityCount int
}