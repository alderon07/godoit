package core

import (
	"fmt"
	"strings"
	"time"
)

// Stats holds analytics data about tasks
type Stats struct {
	Total           int
	Completed       int
	Pending         int
	Overdue         int
	CompletionRate  float64
	ByPriority      map[int]int
	ByTag           map[string]int
	AvgCompletionMS int64
	CompletedToday  int
	CompletedWeek   int
	BlockedTasks    int
}

// CalculateStats computes statistics from a list of tasks
func CalculateStats(tasks []Task, now time.Time) Stats {
	stats := Stats{
		ByPriority: make(map[int]int),
		ByTag:      make(map[string]int),
	}

	var totalCompletionTime time.Duration
	var completedCount int

	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -int(todayStart.Weekday()))

	for _, task := range tasks {
		stats.Total++

		// Count by status
		if task.IsDone() {
			stats.Completed++

			// Calculate completion time
			if task.DoneAt != nil {
				duration := task.DoneAt.Sub(task.CreatedAt)
				totalCompletionTime += duration
				completedCount++

				// Count by time period
				if task.DoneAt.After(todayStart) {
					stats.CompletedToday++
				}
				if task.DoneAt.After(weekStart) {
					stats.CompletedWeek++
				}
			}
		} else {
			stats.Pending++

			// Check if overdue
			if task.IsOverdue(now) {
				stats.Overdue++
			}

			// Check if blocked
			if !AllDependenciesMet(tasks, task) {
				stats.BlockedTasks++
			}
		}

		// Count by priority
		stats.ByPriority[task.Priority]++

		// Count by tags
		for _, tag := range task.Tags {
			stats.ByTag[strings.ToLower(tag)]++
		}
	}

	// Calculate completion rate
	if stats.Total > 0 {
		stats.CompletionRate = float64(stats.Completed) / float64(stats.Total) * 100
	}

	// Calculate average completion time
	if completedCount > 0 {
		stats.AvgCompletionMS = totalCompletionTime.Milliseconds() / int64(completedCount)
	}

	return stats
}

// StatsReport generates a human-readable statistics report
func StatsReport(tasks []Task, now time.Time) string {
	stats := CalculateStats(tasks, now)

	var sb strings.Builder

	sb.WriteString("Task Statistics\n")
	sb.WriteString("===============\n\n")

	// Overall stats
	sb.WriteString(fmt.Sprintf("Total Tasks:     %d\n", stats.Total))
	sb.WriteString(fmt.Sprintf("Completed:       %d\n", stats.Completed))
	sb.WriteString(fmt.Sprintf("Pending:         %d\n", stats.Pending))
	sb.WriteString(fmt.Sprintf("Overdue:         %d\n", stats.Overdue))
	sb.WriteString(fmt.Sprintf("Blocked:         %d\n", stats.BlockedTasks))
	sb.WriteString(fmt.Sprintf("Completion Rate: %.1f%%\n\n", stats.CompletionRate))

	// Productivity
	sb.WriteString("Productivity\n")
	sb.WriteString("------------\n")
	sb.WriteString(fmt.Sprintf("Completed Today: %d\n", stats.CompletedToday))
	sb.WriteString(fmt.Sprintf("Completed This Week: %d\n", stats.CompletedWeek))
	if stats.AvgCompletionMS > 0 {
		avgDays := float64(stats.AvgCompletionMS) / (1000 * 60 * 60 * 24)
		sb.WriteString(fmt.Sprintf("Avg Completion Time: %.1f days\n", avgDays))
	}
	sb.WriteString("\n")

	// By priority
	if len(stats.ByPriority) > 0 {
		sb.WriteString("By Priority\n")
		sb.WriteString("-----------\n")
		for p := 3; p >= 1; p-- {
			if count, ok := stats.ByPriority[p]; ok {
				priorityName := "Low"
				if p == 2 {
					priorityName = "Medium"
				} else if p == 3 {
					priorityName = "High"
				}
				sb.WriteString(fmt.Sprintf("%s (p%d):  %d\n", priorityName, p, count))
			}
		}
		sb.WriteString("\n")
	}

	// By tag
	if len(stats.ByTag) > 0 {
		sb.WriteString("By Tag\n")
		sb.WriteString("------\n")

		// Sort tags by count
		type tagCount struct {
			tag   string
			count int
		}
		var tagCounts []tagCount
		for tag, count := range stats.ByTag {
			tagCounts = append(tagCounts, tagCount{tag, count})
		}

		// Simple bubble sort for small datasets
		for i := 0; i < len(tagCounts); i++ {
			for j := i + 1; j < len(tagCounts); j++ {
				if tagCounts[j].count > tagCounts[i].count {
					tagCounts[i], tagCounts[j] = tagCounts[j], tagCounts[i]
				}
			}
		}

		for _, tc := range tagCounts {
			sb.WriteString(fmt.Sprintf("#%-15s %d\n", tc.tag, tc.count))
		}
	}

	return sb.String()
}

// GetUpcomingTasks returns tasks due within the specified duration
func GetUpcomingTasks(tasks []Task, now time.Time, window time.Duration) []Task {
	result := make([]Task, 0)
	for _, task := range tasks {
		if task.IsDueSoon(now, window) {
			result = append(result, task)
		}
	}
	return result
}

// GetOverdueTasks returns all overdue tasks
func GetOverdueTasks(tasks []Task, now time.Time) []Task {
	result := make([]Task, 0)
	for _, task := range tasks {
		if task.IsOverdue(now) {
			result = append(result, task)
		}
	}
	return result
}

