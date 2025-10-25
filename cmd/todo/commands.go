package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"godoit/internal/alerts"
	"godoit/internal/clock"
	"godoit/internal/core"
	"godoit/internal/notifications"
	"godoit/internal/repository"
	"godoit/internal/server"
	"godoit/internal/service"
	"godoit/internal/store"
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

func getService() *service.TaskService {
  s := getStore()
  repo := repository.NewJSONTaskRepository(s)
  return service.NewTaskService(repo, clock.SystemClock{})
}

// RunAdd adds a new task
func RunAdd(title, description string, dueStr, repeat string, priority int, tags, after string) {
  if title == "" {
    log.Fatal("Error: -title is required")
  }

  var due *time.Time
  if dueStr != "" {
    if t, err := time.Parse("2006-01-02", dueStr); err == nil {
      due = &t
    } else {
      log.Fatalf("Invalid -due: %v", err)
    }
  }

  svc := getService()
  created, err := svc.AddTask(context.Background(), service.AddTaskInput{
    Title:       title,
    Description: description,
    Due:         due,
    Priority:    priority,
    Tags:        core.ParseTags(tags),
    Repeat:      repeat,
    DependsOn:   core.ParseIDs(after),
  })
  must(err)
  fmt.Printf("Added: %s (ID: %d)\n", created.Title, created.ID)
}

// RunList lists tasks with optional filters
func RunList(showAll, today, week, detailed bool, grep, tags, sortKey, before, after string) {
  svc := getService()
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

  var beforePtr, afterPtr *time.Time
  if before != "" { if t, err := time.Parse("2006-01-02", before); err == nil { beforePtr = &t } }
  if after != "" { if t, err := time.Parse("2006-01-02", after); err == nil { afterPtr = &t } }

  visible, err := svc.QueryTasks(context.Background(), service.Query{
    ShowAll: showAll,
    Grep:    grep,
    SortKey: sortKey,
    Tags:    tags,
    Before:  beforePtr,
    After:   afterPtr,
  })
  must(err)

  // also fetch all tasks to compute dependency info
  allTasks, err := svc.QueryTasks(context.Background(), service.Query{ShowAll: true, SortKey: sortKey})
  must(err)

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

    fmt.Print(strings.Repeat("=", 50), "\n")
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
      if !core.AllDependenciesMet(allTasks, t) {
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

  fmt.Print("\n" + strings.Repeat("=", 50) + "\n")
  fmt.Printf("Total: %d task(s)\n", len(visible))
}

// RunDone marks a task as complete
func RunDone(indexStr string) {
  idx, err := core.Atoi1(indexStr)
  must(err)

  svc := getService()
  visible, err := svc.QueryTasks(context.Background(), service.Query{ShowAll: false, SortKey: "due"})
  must(err)
  if idx < 1 || idx > len(visible) { log.Fatal("Invalid index") }
  updated, err := svc.MarkDone(context.Background(), visible, idx)
  must(err)
  fmt.Println("Marked done:", visible[idx-1].Title)
  if updated.Repeat != "" {
    fmt.Println("Created next occurrence")
  }
}

// RunRemove removes a task
func RunRemove(indexStr string) {
  idx, err := core.Atoi1(indexStr)
  must(err)

  svc := getService()
  visible, err := svc.QueryTasks(context.Background(), service.Query{ShowAll: true, SortKey: "due"})
  must(err)
  if idx < 1 || idx > len(visible) { log.Fatal("Invalid index") }
  taskTitle := visible[idx-1].Title
  must(svc.RemoveTask(context.Background(), visible, idx))
  fmt.Println("Removed:", taskTitle)
}

// RunEdit edits an existing task
func RunEdit(indexStr string, title, description, dueStr, repeat string, priority int, tags, after string) {
  idx, err := core.Atoi1(indexStr)
  must(err)

  svc := getService()
  visible, err := svc.QueryTasks(context.Background(), service.Query{ShowAll: true, SortKey: "due"})
  must(err)
  if idx < 1 || idx > len(visible) {
    log.Fatalf("Invalid index: %d", idx)
  }
  targetID := visible[idx-1].ID

  var titlePtr *string
  if title != "" { titlePtr = &title }

  var descPtr *string
  if description != "" {
    if description == "none" { empty := ""; descPtr = &empty } else { descPtr = &description }
  }

  var duePtr *string
  if dueStr != "" { duePtr = &dueStr }

  var prioPtr *int
  if priority > 0 { prioPtr = &priority }

  var tagsPtr *[]string
  if tags != "" {
    if tags == "none" { empty := []string{}; tagsPtr = &empty } else { v := core.ParseTags(tags); tagsPtr = &v }
  }

  var depsPtr *[]int
  if after != "" {
    if after == "none" { empty := []int{}; depsPtr = &empty } else { v := core.ParseIDs(after); depsPtr = &v }
  }

  updated, err := svc.UpdateTask(context.Background(), targetID, service.UpdateTaskInput{
    Title:       titlePtr,
    Description: descPtr,
    Due:         duePtr,
    Priority:    prioPtr,
    Tags:        tagsPtr,
    Repeat:      func() *string { if repeat == "" { return nil }; if repeat == "none" { empty := ""; return &empty }; return &repeat }(),
    DependsOn:   depsPtr,
  })
  must(err)
  fmt.Println("Updated:", updated.Title)
}

// RunAlerts shows alerts for due/overdue tasks
func RunAlerts(watch bool, interval, ahead time.Duration) {
  svc := getService()
  notifier := notifications.NewSystemNotifier(true)
  scanner := alerts.NewScanner(notifier)

  if watch {
    // Watch mode with continuous monitoring
    loadFunc := func() ([]core.Task, error) {
      return svc.QueryTasks(context.Background(), service.Query{ShowAll: true, SortKey: "due"})
    }

    tasks, err := svc.QueryTasks(context.Background(), service.Query{ShowAll: true, SortKey: "due"})
    must(err)

    scanner.Watch(tasks, interval, ahead, loadFunc)
  } else {
    // One-time scan
    tasks, err := svc.QueryTasks(context.Background(), service.Query{ShowAll: true, SortKey: "due"})
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
  svc := getService()
  tasks, err := svc.QueryTasks(context.Background(), service.Query{ShowAll: true, SortKey: "due"})
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
