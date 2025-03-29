package processor

import (
	"math"
	"time"

	"github.com/samuellee/StravaGraph/internal/strava"
)

// MetricsCalculator calculates activity metrics
type MetricsCalculator struct {
	DailyData []*strava.DailyActivity
	StartDate time.Time
	EndDate   time.Time
}

// NewMetricsCalculator creates a new metrics calculator
func NewMetricsCalculator(dailyData []*strava.DailyActivity, startDate, endDate time.Time) *MetricsCalculator {
	return &MetricsCalculator{
		DailyData: dailyData,
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// CalculateOverallStats calculates overall activity statistics
func (m *MetricsCalculator) CalculateOverallStats() *strava.ActivityStats {
	stats := &strava.ActivityStats{
		ActivityTypes: make(map[string]int),
	}

	var currentStreak int
	var longestStreak int

	for _, day := range m.DailyData {
		if day.Count > 0 {
			stats.TotalActivities += day.Count
			stats.TotalDistance += day.TotalDistance / 1000 // Convert to kilometers
			stats.TotalDuration += day.TotalDuration / 3600 // Convert to hours
			stats.TotalElevation += day.TotalElevation
			stats.ActiveDays++

			if day.HasPR {
				stats.PRCount++
			}

			// Add activity types
			for t, count := range day.Types {
				stats.ActivityTypes[t] += count
			}

			// Update streak
			currentStreak++
			if currentStreak > longestStreak {
				longestStreak = currentStreak
			}
		} else {
			// Reset streak if no activities
			currentStreak = 0
		}
	}

	stats.LongestStreak = longestStreak

	return stats
}

// CalculatePeriodStats calculates statistics for specific time periods
func (m *MetricsCalculator) CalculatePeriodStats(periodType string) []*strava.DatePeriodStats {
	var stats []*strava.DatePeriodStats

	// Group days by period
	periods := make(map[string]*strava.DatePeriodStats)

	for _, day := range m.DailyData {
		if day.Count == 0 {
			continue
		}

		var periodKey string
		switch periodType {
		case "weekly":
			// Calculate the week number and year
			year, week := day.Date.ISOWeek()
			periodKey = formatPeriodKey(year, week)
		case "monthly":
			// Year and month
			periodKey = day.Date.Format("2006-01")
		case "yearly":
			// Just the year
			periodKey = day.Date.Format("2006")
		default:
			continue
		}

		// Create or update period stats
		period, exists := periods[periodKey]
		if !exists {
			period = &strava.DatePeriodStats{
				Period: periodKey,
			}
			periods[periodKey] = period
		}

		// Add day's stats to period
		period.TotalDistance += day.TotalDistance / 1000 // km
		period.TotalDuration += day.TotalDuration / 3600 // hours
		period.TotalElevation += day.TotalElevation
		period.ActivityCount += day.Count
	}

	// Convert map to slice
	for _, period := range periods {
		stats = append(stats, period)
	}

	return stats
}

// CalculateAverages calculates average metrics per active day
func (m *MetricsCalculator) CalculateAverages() map[string]float64 {
	stats := m.CalculateOverallStats()
	averages := make(map[string]float64)

	if stats.ActiveDays > 0 {
		averages["distancePerDay"] = stats.TotalDistance / float64(stats.ActiveDays)
		averages["durationPerDay"] = float64(stats.TotalDuration) / float64(stats.ActiveDays)
		averages["elevationPerDay"] = stats.TotalElevation / float64(stats.ActiveDays)
		averages["activitiesPerDay"] = float64(stats.TotalActivities) / float64(stats.ActiveDays)
	}

	// Calculate activity frequency (percentage of days with activity)
	totalDays := m.EndDate.Sub(m.StartDate).Hours() / 24
	if totalDays > 0 {
		averages["activityFrequency"] = float64(stats.ActiveDays) / totalDays
	}

	return averages
}

// formatPeriodKey formats a period key based on year and period number
func formatPeriodKey(year, period int) string {
	return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC).
		AddDate(0, 0, (period-1)*7).Format("2006-W02")
}

// CalculateEffortScore calculates an overall effort score
func (m *MetricsCalculator) CalculateEffortScore() float64 {
	stats := m.CalculateOverallStats()

	// Simple formula based on total distance, elevation, and duration
	// Normalized to produce a 0-100 score for typical activity levels

	// Base score from distance (km)
	distanceScore := math.Min(stats.TotalDistance/10, 100)

	// Elevation bonus (m)
	elevationBonus := math.Min(stats.TotalElevation/100, 50)

	// Duration factor (hours)
	durationFactor := math.Min(float64(stats.TotalDuration)/5, 100)

	// Frequency bonus from active days percentage
	totalDays := m.EndDate.Sub(m.StartDate).Hours() / 24
	frequencyBonus := 0.0
	if totalDays > 0 {
		frequencyBonus = math.Min(float64(stats.ActiveDays)/totalDays*50, 50)
	}

	// Streak bonus
	streakBonus := math.Min(float64(stats.LongestStreak), 30)

	// PR bonus
	prBonus := math.Min(float64(stats.PRCount)*2, 20)

	// Calculate total score and normalize to 0-100
	rawScore := distanceScore + elevationBonus + (durationFactor * 0.5) +
		frequencyBonus + (streakBonus * 0.5) + prBonus

	normalizedScore := math.Min(rawScore/3, 100)

	return math.Round(normalizedScore*10) / 10 // Round to 1 decimal place
}
