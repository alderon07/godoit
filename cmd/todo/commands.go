package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"godo/internal/alerts"
	"godo/internal/core"
	"godo/internal/notifications"
	"godo/internal/server"
	"godo/internal/store"
)

const (
	GetItDoneNow  = 3
	IGotTime      = 2
	WeAreChilling = 1
)

// Helper for consistent error handling
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// getStore returns the default store instance
func getStore() *store.JSONStore {
	s, err := store.DefaultStore()
	must(err)
	return s
}

// RunAdd adds a new task
func RunAdd(title, description string, dueStr, repeat string, priority int, tags, after string) {
	if title == "" {
		log.Fatal("Error: -title is required")
	}

	s := getStore()
	tasks, err := store.LoadTasks[core.Task](s)
	must(err)

	var due *time.Time
	if dueStr != "" {
		if t, err := time.Parse("2006-01-02", dueStr); err == nil {
			due = &t
		} else {
			log.Fatalf("Invalid -due: %v", err)
		}
	}

	tasks = core.Add(tasks, title, due)
	i := len(tasks) - 1

	// Set description
	tasks[i].Description = description

	// Set priority
	if priority < 1 || priority > 3 {
		priority = 1
	}
	tasks[i].Priority = priority

	// Set repeat
	if repeat != "" {
		tasks[i].Repeat = repeat
	}

	// Set tags
	tasks[i].Tags = core.ParseTags(tags)

	// Set dependencies
	tasks[i].DependsOn = core.ParseIDs(after)

	must(store.SaveTasks(s, tasks))
	fmt.Printf("Added: %s (ID: %d)\n", title, tasks[i].ID)
}

// RunList lists tasks with optional filters
func RunList(showAll, today, week, detailed bool, grep, tags, sortKey, before, after string) {
	s := getStore()
	tasks, err := store.LoadTasks[core.Task](s)
	must(err)

	now := time.Now()

	// Apply time-based filters
	if today {
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endOfDay := startOfDay.AddDate(0, 0, 1)
		before = endOfDay.Format("2006-01-02")
		after = startOfDay.Format("2006-01-02")
	} else if week {
		startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		// Move to start of week (Sunday)
		for startOfWeek.Weekday() != time.Sunday {
			startOfWeek = startOfWeek.AddDate(0, 0, -1)
		}
		endOfWeek := startOfWeek.AddDate(0, 0, 7)
		before = endOfWeek.Format("2006-01-02")
		after = startOfWeek.Format("2006-01-02")
	}

	visible := core.SortedWith(tasks, showAll, grep, sortKey)
	visible = core.FilterByTags(visible, tags)
	visible = core.FilterByDate(visible, before, after)

	if len(visible) == 0 {
		fmt.Println("(no tasks)")
		return
	}

	// Print header based on view
	if today {
		fmt.Println("📅 Today's Tasks")
		fmt.Println("================")
	} else if week {
		fmt.Println("📅 This Week's Tasks")
		fmt.Println("====================")
	}

	for i, t := range visible {
		status := " "
		if t.IsDone() {
			status = "✓"
		}

		// Determine priority indicator
		priorityStr := ""
		switch t.Priority {
		case 3:
			priorityStr = "🔴"
		case 2:
			priorityStr = "🟡"
		case 1:
			priorityStr = "🟢"
		}

		fmt.Printf("\n%2d. [%s] %s %s\n", i+1, status, priorityStr, t.Title)

		// Show description if present
		if detailed && t.Description != "" {
			fmt.Printf("    📝 %s\n", t.Description)
		}

		// Show due date with status indicator
		if t.Due != nil {
			dueStr := t.Due.Format("2006-01-02 15:04")
			if t.IsOverdue(now) {
				fmt.Printf("    ⏰ Due: %s (OVERDUE!)\n", dueStr)
			} else if t.IsDueSoon(now, 24*time.Hour) {
				fmt.Printf("    ⏰ Due: %s (soon)\n", dueStr)
			} else {
				fmt.Printf("    📅 Due: %s\n", dueStr)
			}
		}

		// Show tags
		if len(t.Tags) > 0 {
			fmt.Printf("    🏷️  ")
			for _, tag := range t.Tags {
				fmt.Printf("#%s ", tag)
			}
			fmt.Println()
		}

		// Show repeat info
		if t.Repeat != "" {
			fmt.Printf("    🔄 Repeats: %s\n", t.Repeat)
		}

		// Show dependencies
		if len(t.DependsOn) > 0 {
			fmt.Printf("    🔗 Depends on: %v\n", t.DependsOn)
			if !core.AllDependenciesMet(tasks, t) {
				fmt.Printf("    ⚠️  BLOCKED (dependencies not met)\n")
			}
		}

		// Show creation date in detailed view
		if detailed {
			fmt.Printf("    🕐 Created: %s\n", t.CreatedAt.Format("2006-01-02 15:04"))
		}

		// Show completion date if done
		if t.IsDone() && t.DoneAt != nil {
			fmt.Printf("    ✅ Completed: %s\n", t.DoneAt.Format("2006-01-02 15:04"))
		}
	}

	fmt.Print("\n" + strings.Repeat("-", 50) + "\n")
	fmt.Printf("Total: %d task(s)\n", len(visible))
}

// RunDone marks a task as complete
func RunDone(indexStr string) {
	idx, err := core.Atoi1(indexStr)
	must(err)

	s := getStore()
	tasks, err := store.LoadTasks[core.Task](s)
	must(err)

	visible := core.SortedWith(tasks, false, "", "due")

	tasks, err = core.MarkDone(tasks, visible, idx)
	must(err)

	must(store.SaveTasks(s, tasks))
	fmt.Println("Marked done:", visible[idx-1].Title)

	// Check if it was a recurring task
	for _, t := range tasks {
		if t.Title == visible[idx-1].Title && !t.IsDone() {
			fmt.Printf("Created next occurrence (ID: %d)\n", t.ID)
			break
		}
	}
}

// RunRemove removes a task
func RunRemove(indexStr string) {
	idx, err := core.Atoi1(indexStr)
	must(err)

	s := getStore()
	tasks, err := store.LoadTasks[core.Task](s)
	must(err)

	visible := core.SortedWith(tasks, true, "", "due")

	taskTitle := visible[idx-1].Title
	tasks, err = core.Remove(tasks, visible, idx)
	must(err)

	must(store.SaveTasks(s, tasks))
	fmt.Println("Removed:", taskTitle)
}

// RunEdit edits an existing task
func RunEdit(indexStr string, title, description, dueStr, repeat string, priority int, tags, after string) {
	idx, err := core.Atoi1(indexStr)
	must(err)

	s := getStore()
	tasks, err := store.LoadTasks[core.Task](s)
	must(err)

	visible := core.SortedWith(tasks, true, "", "due")

	if idx < 1 || idx > len(visible) {
		log.Fatalf("Invalid index: %d", idx)
	}

	targetID := visible[idx-1].ID
	var targetTask *core.Task

	for i := range tasks {
		if tasks[i].ID == targetID {
			targetTask = &tasks[i]
			break
		}
	}

	if targetTask == nil {
		log.Fatalf("Task %d not found", targetID)
	}

	// Update fields if provided
	if title != "" {
		targetTask.Title = title
	}

	if description != "" {
		if description == "none" {
			targetTask.Description = ""
		} else {
			targetTask.Description = description
		}
	}

	if dueStr != "" {
		if dueStr == "none" {
			targetTask.Due = nil
		} else if t, err := time.Parse("2006-01-02", dueStr); err == nil {
			targetTask.Due = &t
		} else {
			log.Fatalf("Invalid -due: %v", err)
		}
	}

	if repeat != "" {
		if repeat == "none" {
			targetTask.Repeat = ""
		} else {
			targetTask.Repeat = repeat
		}
	}

	if priority > 0 {
		if priority < 1 || priority > 3 {
			priority = 1
		}
		targetTask.Priority = priority
	}

	if tags != "" {
		if tags == "none" {
			targetTask.Tags = nil
		} else {
			targetTask.Tags = core.ParseTags(tags)
		}
	}

	if after != "" {
		if after == "none" {
			targetTask.DependsOn = nil
		} else {
			targetTask.DependsOn = core.ParseIDs(after)
		}
	}

	must(store.SaveTasks(s, tasks))
	fmt.Println("Updated:", targetTask.Title)
}

// RunAlerts shows alerts for due/overdue tasks
func RunAlerts(watch bool, interval, ahead time.Duration) {
	s := getStore()
	notifier := notifications.NewSystemNotifier(true)
	scanner := alerts.NewScanner(notifier)

	if watch {
		// Watch mode with continuous monitoring
		loadFunc := func() ([]core.Task, error) {
			return store.LoadTasks[core.Task](s)
		}

		tasks, err := store.LoadTasks[core.Task](s)
		must(err)

		scanner.Watch(tasks, interval, ahead, loadFunc)
	} else {
		// One-time scan
		tasks, err := store.LoadTasks[core.Task](s)
		must(err)

		alertList := scanner.Scan(tasks, time.Now(), ahead)

		if len(alertList) == 0 {
			fmt.Println("No alerts")
			return
		}

		fmt.Printf("=== Alerts ===\n\n")
		for _, alert := range alertList {
			fmt.Printf("[%s] %s\n", alert.Type, alert.Message)
			if alert.Task.Due != nil {
				fmt.Printf("  Due: %s\n", alert.Task.Due.Format("2006-01-02 15:04"))
			}
			if alert.Task.Priority > 1 {
				fmt.Printf("  Priority: %d\n", alert.Task.Priority)
			}
			fmt.Println()
		}
	}
}

// RunStats shows task statistics
func RunStats() {
	s := getStore()
	tasks, err := store.LoadTasks[core.Task](s)
	must(err)

	fmt.Print(core.StatsReport(tasks, time.Now()))
}

// RunServer starts the HTTP API server
func RunServer(host string, port int) {
	s := getStore()
	srv := server.NewServer(host, port, s)

	fmt.Printf("Starting HTTP server on %s:%d\n", host, port)
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println("\nEndpoints:")
	fmt.Println("  GET    /tasks          - List all tasks")
	fmt.Println("  POST   /tasks          - Create a task")
	fmt.Println("  GET    /tasks/:id      - Get a task")
	fmt.Println("  PUT    /tasks/:id      - Update a task")
	fmt.Println("  DELETE /tasks/:id      - Delete a task")
	fmt.Println("  POST   /tasks/:id/done - Mark task as done")
	fmt.Println("  GET    /stats          - Get statistics")
	fmt.Println("  GET    /health         - Health check")
	fmt.Println()

	must(srv.Start())
}
