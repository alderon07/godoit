package core

import (
	"fmt"
	"strings"
	"time"
)

// Task represents a todo item with all its properties
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Due         *time.Time `json:"due,omitempty"`
	DoneAt      *time.Time `json:"done_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	Priority    int        `json:"priority"` // 1=low, 2=medium, 3=high
	Tags        []string   `json:"tags,omitempty"`
	Repeat      string     `json:"repeat,omitempty"` // daily, weekly, monthly
	DependsOn   []int      `json:"depends_on,omitempty"`
}

// IsDone returns true if the task is completed
func (t *Task) IsDone() bool {
	return t.DoneAt != nil
}

// IsOverdue returns true if the task has a due date in the past and is not done
func (t *Task) IsOverdue(now time.Time) bool {
	if t.IsDone() || t.Due == nil {
		return false
	}
	return t.Due.Before(now)
}

// IsDueSoon returns true if the task is due within the given duration
func (t *Task) IsDueSoon(now time.Time, window time.Duration) bool {
	if t.IsDone() || t.Due == nil {
		return false
	}
	deadline := now.Add(window)
	return t.Due.Before(deadline) && !t.Due.Before(now)
}

// HasTag returns true if the task has the given tag
func (t *Task) HasTag(tag string) bool {
	tag = strings.ToLower(strings.TrimSpace(tag))
	for _, t := range t.Tags {
		if strings.ToLower(t) == tag {
			return true
		}
	}
	return false
}

// Add creates a new task with the given title and returns the updated task list
func Add(tasks []Task, title string, due *time.Time) []Task {
	nextID := 1
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}

	task := Task{
		ID:        nextID,
		Title:     title,
		Due:       due,
		CreatedAt: time.Now(),
		Priority:  1, // default to low priority
	}

	return append(tasks, task)
}

// GetByID finds a task by its ID
func GetByID(tasks []Task, id int) (*Task, error) {
	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task %d not found", id)
}

// Update replaces a task with the same ID
func Update(tasks []Task, updated Task) ([]Task, error) {
	for i := range tasks {
		if tasks[i].ID == updated.ID {
			tasks[i] = updated
			return tasks, nil
		}
	}
	return tasks, fmt.Errorf("task %d not found", updated.ID)
}

// Remove deletes a task from the list by its position in visible list
func Remove(tasks []Task, visible []Task, idx int) ([]Task, error) {
	if idx < 1 || idx > len(visible) {
		return tasks, fmt.Errorf("invalid index: %d", idx)
	}

	targetID := visible[idx-1].ID
	result := make([]Task, 0, len(tasks)-1)
	found := false

	for _, t := range tasks {
		if t.ID == targetID {
			found = true
			continue
		}
		result = append(result, t)
	}

	if !found {
		return tasks, fmt.Errorf("task %d not found", targetID)
	}

	return result, nil
}

// MarkDone marks a task as complete, handling recurring tasks
func MarkDone(tasks []Task, visible []Task, idx int) ([]Task, error) {
	if idx < 1 || idx > len(visible) {
		return tasks, fmt.Errorf("invalid index: %d", idx)
	}

	targetID := visible[idx-1].ID
	now := time.Now()

	for i := range tasks {
		if tasks[i].ID == targetID {
			if tasks[i].IsDone() {
				return tasks, fmt.Errorf("task already completed")
			}

			// Check dependencies
			if !AllDependenciesMet(tasks, tasks[i]) {
				return tasks, fmt.Errorf("cannot complete task: dependencies not met")
			}

			tasks[i].DoneAt = &now

			// Handle recurring tasks
			if tasks[i].Repeat != "" {
				nextTask := createNextRecurrence(tasks[i])
				if nextTask != nil {
					tasks = append(tasks, *nextTask)
				}
			}

			return tasks, nil
		}
	}

	return tasks, fmt.Errorf("task %d not found", targetID)
}

// AllDependenciesMet returns true if all dependencies are completed
func AllDependenciesMet(tasks []Task, task Task) bool {
	for _, depID := range task.DependsOn {
		depTask, err := GetByID(tasks, depID)
		if err != nil || !depTask.IsDone() {
			return false
		}
	}
	return true
}

// createNextRecurrence creates the next occurrence of a recurring task
func createNextRecurrence(task Task) *Task {
	if task.Due == nil {
		return nil
	}

	nextDue := calculateNextDue(*task.Due, task.Repeat)
	if nextDue == nil {
		return nil
	}

	// Find the highest ID
	nextID := task.ID + 1

	nextTask := Task{
		ID:          nextID,
		Title:       task.Title,
		Description: task.Description,
		Due:         nextDue,
		CreatedAt:   time.Now(),
		Priority:    task.Priority,
		Tags:        append([]string{}, task.Tags...),
		Repeat:      task.Repeat,
		DependsOn:   append([]int{}, task.DependsOn...),
	}

	return &nextTask
}

// calculateNextDue calculates the next due date based on repeat rule
func calculateNextDue(current time.Time, repeat string) *time.Time {
	var next time.Time

	switch strings.ToLower(repeat) {
	case "daily":
		next = current.AddDate(0, 0, 1)
	case "weekly":
		next = current.AddDate(0, 0, 7)
	case "monthly":
		next = current.AddDate(0, 1, 0)
	default:
		return nil
	}

	return &next
}

// ParseTags parses a comma-separated string into a slice of tags
func ParseTags(tagStr string) []string {
	if tagStr == "" {
		return nil
	}

	parts := strings.Split(tagStr, ",")
	tags := make([]string, 0, len(parts))

	for _, p := range parts {
		tag := strings.TrimSpace(p)
		if tag != "" {
			tags = append(tags, tag)
		}
	}

	return tags
}

// ParseIDs parses a comma-separated string into a slice of integers
func ParseIDs(idStr string) []int {
	if idStr == "" {
		return nil
	}

	parts := strings.Split(idStr, ",")
	ids := make([]int, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		var id int
		if _, err := fmt.Sscanf(p, "%d", &id); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}

// Atoi1 converts a 1-indexed string to 0-indexed int, but returns the 1-indexed value
func Atoi1(s string) (int, error) {
	var idx int
	if _, err := fmt.Sscanf(s, "%d", &idx); err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}
	if idx < 1 {
		return 0, fmt.Errorf("index must be >= 1")
	}
	return idx, nil
}

