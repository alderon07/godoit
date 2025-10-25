package alerts

import (
	"fmt"
	"time"

	"godo/internal/core"
	"godo/internal/notifications"
)

// AlertType represents the type of alert
type AlertType string

const (
	AlertOverdue    AlertType = "Overdue"
	AlertDueSoon    AlertType = "Due Soon"
	AlertBlocked    AlertType = "Blocked"
	AlertDependency AlertType = "Dependency"
)

// Alert represents a task alert
type Alert struct {
	Task    core.Task
	Type    AlertType
	Message string
}

// Scanner scans tasks for alerts
type Scanner struct {
	notifier notifications.Notifier
}

// NewScanner creates a new alert scanner
func NewScanner(notifier notifications.Notifier) *Scanner {
	return &Scanner{
		notifier: notifier,
	}
}

// Scan scans tasks and returns alerts
func (s *Scanner) Scan(tasks []core.Task, now time.Time, lookahead time.Duration) []Alert {
	alerts := make([]Alert, 0)

	for _, task := range tasks {
		if task.IsDone() {
			continue
		}

		// Check for overdue tasks
		if task.IsOverdue(now) {
			alerts = append(alerts, Alert{
				Task:    task,
				Type:    AlertOverdue,
				Message: fmt.Sprintf("Task is overdue: %s", task.Title),
			})
			continue
		}

		// Check for tasks due soon
		if task.IsDueSoon(now, lookahead) {
			dueIn := task.Due.Sub(now)
			alerts = append(alerts, Alert{
				Task:    task,
				Type:    AlertDueSoon,
				Message: fmt.Sprintf("Task due in %s: %s", formatDuration(dueIn), task.Title),
			})
			continue
		}

		// Check for blocked tasks
		if !core.AllDependenciesMet(tasks, task) {
			alerts = append(alerts, Alert{
				Task:    task,
				Type:    AlertBlocked,
				Message: fmt.Sprintf("Task blocked by dependencies: %s", task.Title),
			})
		}
	}

	return alerts
}

// ScanAndNotify scans for alerts and sends notifications
func (s *Scanner) ScanAndNotify(tasks []core.Task, now time.Time, lookahead time.Duration) []Alert {
	alerts := s.Scan(tasks, now, lookahead)

	for _, alert := range alerts {
		_ = s.notifier.Send(string(alert.Type), alert.Message)
	}

	return alerts
}

// Watch continuously monitors tasks and sends alerts
func (s *Scanner) Watch(tasks []core.Task, interval, lookahead time.Duration, loadFunc func() ([]core.Task, error)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Printf("Watching for alerts (interval: %s, lookahead: %s)\n", interval, lookahead)
	fmt.Println("Press Ctrl+C to stop...")

	// Initial scan
	s.printAlerts(tasks, time.Now(), lookahead)

	for range ticker.C {
		// Reload tasks
		if loadFunc != nil {
			newTasks, err := loadFunc()
			if err != nil {
				fmt.Printf("Error loading tasks: %v\n", err)
				continue
			}
			tasks = newTasks
		}

		alerts := s.ScanAndNotify(tasks, time.Now(), lookahead)
		if len(alerts) > 0 {
			s.printAlerts(tasks, time.Now(), lookahead)
		}
	}
}

// printAlerts prints alerts to console
func (s *Scanner) printAlerts(tasks []core.Task, now time.Time, lookahead time.Duration) {
	alerts := s.Scan(tasks, now, lookahead)

	if len(alerts) == 0 {
		fmt.Println("No alerts")
		return
	}

	fmt.Printf("\n=== Alerts (%s) ===\n", now.Format("2006-01-02 15:04:05"))

	for _, alert := range alerts {
		fmt.Printf("[%s] %s\n", alert.Type, alert.Message)
		if alert.Task.Due != nil {
			fmt.Printf("  Due: %s\n", alert.Task.Due.Format("2006-01-02 15:04"))
		}
		if alert.Task.Priority > 1 {
			fmt.Printf("  Priority: %d\n", alert.Task.Priority)
		}
	}
	fmt.Println()
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "less than a minute"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours < 1 {
		return fmt.Sprintf("%d minute(s)", minutes)
	}

	if hours < 24 {
		return fmt.Sprintf("%d hour(s) %d minute(s)", hours, minutes)
	}

	days := hours / 24
	hours = hours % 24

	if days == 1 {
		return "1 day"
	}

	if hours > 0 {
		return fmt.Sprintf("%d days %d hours", days, hours)
	}

	return fmt.Sprintf("%d days", days)
}

// GetAlertSummary returns a summary string of alerts
func GetAlertSummary(tasks []core.Task, now time.Time, lookahead time.Duration) string {
	scanner := NewScanner(notifications.NewNoOpNotifier())
	alerts := scanner.Scan(tasks, now, lookahead)

	if len(alerts) == 0 {
		return "No alerts"
	}

	overdue := 0
	dueSoon := 0
	blocked := 0

	for _, alert := range alerts {
		switch alert.Type {
		case AlertOverdue:
			overdue++
		case AlertDueSoon:
			dueSoon++
		case AlertBlocked:
			blocked++
		}
	}

	summary := fmt.Sprintf("%d alerts: ", len(alerts))
	parts := make([]string, 0, 3)

	if overdue > 0 {
		parts = append(parts, fmt.Sprintf("%d overdue", overdue))
	}
	if dueSoon > 0 {
		parts = append(parts, fmt.Sprintf("%d due soon", dueSoon))
	}
	if blocked > 0 {
		parts = append(parts, fmt.Sprintf("%d blocked", blocked))
	}

	for i, part := range parts {
		if i > 0 {
			summary += ", "
		}
		summary += part
	}

	return summary
}

