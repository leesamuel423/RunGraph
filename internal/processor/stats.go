package processor

import (
	"sort"
	"time"

	"github.com/samuellee/StravaGraph/internal/strava"
)

// StatsGenerator generates comprehensive statistics
type StatsGenerator struct {
	DailyData  []*strava.DailyActivity
	StartDate  time.Time
	EndDate    time.Time
	MetricType string
}

// NewStatsGenerator creates a new stats generator
func NewStatsGenerator(dailyData []*strava.DailyActivity, startDate, endDate time.Time, metricType string) *StatsGenerator {
	return &StatsGenerator{
		DailyData:  dailyData,
		StartDate:  startDate,
		EndDate:    endDate,
		MetricType: metricType,
	}
}

// GenerateStats generates all statistics for the heatmap
func (sg *StatsGenerator) GenerateStats() map[string]interface{} {
	calculator := NewMetricsCalculator(sg.DailyData, sg.StartDate, sg.EndDate)

	stats := make(map[string]interface{})

	// Overall stats
	stats["overall"] = calculator.CalculateOverallStats()

	// Period stats
	stats["weekly"] = calculator.CalculatePeriodStats("weekly")
	stats["monthly"] = calculator.CalculatePeriodStats("monthly")
	stats["yearly"] = calculator.CalculatePeriodStats("yearly")

	// Averages
	stats["averages"] = calculator.CalculateAverages()

	// Effort score
	stats["effortScore"] = calculator.CalculateEffortScore()

	// Top days
	stats["topDays"] = sg.getTopDays(5)

	// Activity type breakdown
	stats["activityBreakdown"] = sg.getActivityTypeBreakdown()

	// Time period metadata
	stats["timePeriod"] = map[string]interface{}{
		"start":     sg.StartDate.Format("2006-01-02"),
		"end":       sg.EndDate.Format("2006-01-02"),
		"totalDays": int(sg.EndDate.Sub(sg.StartDate).Hours()/24) + 1,
	}

	return stats
}

// getTopDays returns the top N days based on the configured metric
func (sg *StatsGenerator) getTopDays(n int) []map[string]interface{} {
	// Create a slice to hold day data
	type dayData struct {
		day   *strava.DailyActivity
		value float64
	}

	var days []dayData

	// Calculate metric value for each day with activity
	for _, day := range sg.DailyData {
		if day.Count == 0 {
			continue
		}

		var value float64
		switch sg.MetricType {
		case "distance":
			value = day.TotalDistance / 1000 // km
		case "duration":
			value = float64(day.TotalDuration) / 3600 // hours
		case "elevation":
			value = day.TotalElevation
		case "heart_rate":
			value = day.AvgHeartRate
		case "effort":
			if day.TotalDuration > 0 {
				value = (day.TotalDistance * (1 + day.TotalElevation/100)) / float64(day.TotalDuration)
			}
		default:
			value = float64(day.Count)
		}

		days = append(days, dayData{day, value})
	}

	// Sort days by metric value in descending order
	sort.Slice(days, func(i, j int) bool {
		return days[i].value > days[j].value
	})

	// Take top N days
	result := make([]map[string]interface{}, 0, n)
	for i := 0; i < n && i < len(days); i++ {
		day := days[i]

		// Format the value based on metric type
		formattedValue := day.value
		unit := ""
		switch sg.MetricType {
		case "distance":
			unit = "km"
		case "duration":
			unit = "hours"
		case "elevation":
			unit = "m"
		case "heart_rate":
			unit = "bpm"
		}

		topDay := map[string]interface{}{
			"date":          day.day.Date.Format("2006-01-02"),
			"dayOfWeek":     day.day.Date.Format("Monday"),
			"value":         formattedValue,
			"unit":          unit,
			"activityCount": day.day.Count,
			"activities":    day.day.Activities,
			"types":         day.day.Types,
			"hasPR":         day.day.HasPR,
		}

		result = append(result, topDay)
	}

	return result
}

// getActivityTypeBreakdown returns the breakdown of activity types
func (sg *StatsGenerator) getActivityTypeBreakdown() map[string]interface{} {
	typeCounts := make(map[string]int)
	typeDistance := make(map[string]float64)
	typeDuration := make(map[string]int)

	for _, day := range sg.DailyData {
		if day.Count == 0 {
			continue
		}

		// We don't have per-activity breakdown in daily data,
		// so we'll distribute metrics proportionally by activity type
		for actType, count := range day.Types {
			typeCounts[actType] += count

			// Distribute metrics proportionally
			proportion := float64(count) / float64(day.Count)
			typeDistance[actType] += day.TotalDistance * proportion / 1000             // km
			typeDuration[actType] += int(float64(day.TotalDuration) * proportion / 60) // minutes
		}
	}

	// Convert to sorted slice for easier consumption
	type typeInfo struct {
		Type     string  `json:"type"`
		Count    int     `json:"count"`
		Distance float64 `json:"distance"`
		Duration int     `json:"duration"`
		Percent  float64 `json:"percent"`
	}

	var types []typeInfo
	totalActivities := 0
	for _, count := range typeCounts {
		totalActivities += count
	}

	for t, count := range typeCounts {
		percent := 0.0
		if totalActivities > 0 {
			percent = float64(count) / float64(totalActivities) * 100
		}

		types = append(types, typeInfo{
			Type:     t,
			Count:    count,
			Distance: typeDistance[t],
			Duration: typeDuration[t],
			Percent:  percent,
		})
	}

	// Sort by count descending
	sort.Slice(types, func(i, j int) bool {
		return types[i].Count > types[j].Count
	})

	return map[string]interface{}{
		"totalActivities": totalActivities,
		"types":           types,
	}
}
