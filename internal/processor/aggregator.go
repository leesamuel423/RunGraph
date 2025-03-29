package processor

import (
	"sort"
	"time"

	"github.com/samuellee/StravaGraph/internal/strava"
)

// ActivityAggregator processes and aggregates activity data
type ActivityAggregator struct {
	Activities []strava.SummaryActivity
	TimeZone   *time.Location
	DailyData  map[string]*strava.DailyActivity // key: YYYY-MM-DD
}

// NewActivityAggregator creates a new activity aggregator
func NewActivityAggregator(activities []strava.SummaryActivity, location *time.Location) *ActivityAggregator {
	return &ActivityAggregator{
		Activities: activities,
		TimeZone:   location,
		DailyData:  make(map[string]*strava.DailyActivity),
	}
}

// Aggregate processes activities and aggregates them by day
func (a *ActivityAggregator) Aggregate() map[string]*strava.DailyActivity {
	for _, activity := range a.Activities {
		// Convert to the configured timezone
		localDate := activity.StartDate.In(a.TimeZone)
		dateKey := localDate.Format("2006-01-02")

		// Create or update the daily activity entry
		dailyActivity, exists := a.DailyData[dateKey]
		if !exists {
			dailyActivity = &strava.DailyActivity{
				Date:       localDate,
				Types:      make(map[string]int),
				Activities: []int64{},
			}
			a.DailyData[dateKey] = dailyActivity
		}

		// Update counts and totals
		dailyActivity.Count++
		dailyActivity.TotalDistance += activity.Distance
		dailyActivity.TotalDuration += activity.MovingTime
		dailyActivity.TotalElevation += activity.TotalElevGain
		dailyActivity.Activities = append(dailyActivity.Activities, activity.ID)

		// Record activity type
		dailyActivity.Types[activity.Type]++

		// Update PR status
		if activity.PRCount > 0 {
			dailyActivity.HasPR = true
		}

		// Update heart rate if available
		if activity.AverageHeartrate > 0 {
			// If this is the first activity with heart rate data
			if dailyActivity.AvgHeartRate == 0 {
				dailyActivity.AvgHeartRate = activity.AverageHeartrate
			} else {
				// Calculate running average
				total := dailyActivity.AvgHeartRate * float64(dailyActivity.Count-1)
				dailyActivity.AvgHeartRate = (total + activity.AverageHeartrate) / float64(dailyActivity.Count)
			}
		}

		if activity.MaxHeartrate > dailyActivity.MaxHeartRate {
			dailyActivity.MaxHeartRate = activity.MaxHeartrate
		}
	}

	return a.DailyData
}

// GetOrderedDates returns daily activities ordered by date
func (a *ActivityAggregator) GetOrderedDates(startDate, endDate time.Time) []*strava.DailyActivity {
	var result []*strava.DailyActivity

	// Fill in all dates in the range for continuity
	current := startDate
	for !current.After(endDate) {
		dateKey := current.Format("2006-01-02")

		// Check if we have data for this date
		dailyActivity, exists := a.DailyData[dateKey]
		if !exists {
			// Create an empty record for this date
			dailyActivity = &strava.DailyActivity{
				Date:  current,
				Types: make(map[string]int),
			}
		}

		result = append(result, dailyActivity)
		current = current.AddDate(0, 0, 1) // Next day
	}

	// Sort by date
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Before(result[j].Date)
	})

	return result
}

// CalculateIntensity determines the heat intensity level for a given metric value
func (a *ActivityAggregator) CalculateIntensity(metricType string, day *strava.DailyActivity) strava.HeatmapIntensity {
	if day.Count == 0 {
		return strava.None
	}

	// Get all non-zero values for this metric to calculate percentiles
	var values []float64
	for _, data := range a.DailyData {
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

	// Sort values to calculate percentiles
	sort.Float64s(values)

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

	// Determine which percentile the day falls into
	percentile := getPercentileRank(values, dayValue)

	// Map percentile to intensity level
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

// Helper function to calculate percentile rank of a value in a sorted array
func getPercentileRank(sortedValues []float64, value float64) float64 {
	if len(sortedValues) == 0 {
		return 0
	}

	// Find position using binary search
	pos := sort.SearchFloat64s(sortedValues, value)

	// Calculate percentile rank
	return float64(pos) / float64(len(sortedValues))
}
