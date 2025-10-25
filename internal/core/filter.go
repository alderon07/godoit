package core

import (
	"strings"
	"time"
)

// FilterByStatus filters tasks by completion status
func FilterByStatus(tasks []Task, showAll bool) []Task {
	if showAll {
		return tasks
	}

	result := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		if !t.IsDone() {
			result = append(result, t)
		}
	}
	return result
}

// FilterByTags filters tasks by tags
// tagStr format: "tag1,tag2" (comma = OR), "tag1+tag2" (plus = AND)
func FilterByTags(tasks []Task, tagStr string) []Task {
	if tagStr == "" {
		return tasks
	}

	// Check if it's AND logic (contains +)
	if strings.Contains(tagStr, "+") {
		return filterByTagsAND(tasks, tagStr)
	}

	// Otherwise, use OR logic
	return filterByTagsOR(tasks, tagStr)
}

// filterByTagsOR returns tasks that have ANY of the specified tags
func filterByTagsOR(tasks []Task, tagStr string) []Task {
	tags := strings.Split(tagStr, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(strings.ToLower(tags[i]))
	}

	result := make([]Task, 0)
	for _, task := range tasks {
		for _, filterTag := range tags {
			if filterTag == "" {
				continue
			}
			if task.HasTag(filterTag) {
				result = append(result, task)
				break
			}
		}
	}
	return result
}

// filterByTagsAND returns tasks that have ALL of the specified tags
func filterByTagsAND(tasks []Task, tagStr string) []Task {
	tags := strings.Split(tagStr, "+")
	for i := range tags {
		tags[i] = strings.TrimSpace(strings.ToLower(tags[i]))
	}

	result := make([]Task, 0)
	for _, task := range tasks {
		hasAll := true
		for _, filterTag := range tags {
			if filterTag == "" {
				continue
			}
			if !task.HasTag(filterTag) {
				hasAll = false
				break
			}
		}
		if hasAll && len(tags) > 0 {
			result = append(result, task)
		}
	}
	return result
}

// FilterByDate filters tasks by due date range
func FilterByDate(tasks []Task, beforeStr, afterStr string) []Task {
	var before, after *time.Time

	if beforeStr != "" {
		if t, err := time.Parse("2006-01-02", beforeStr); err == nil {
			before = &t
		}
	}

	if afterStr != "" {
		if t, err := time.Parse("2006-01-02", afterStr); err == nil {
			after = &t
		}
	}

	if before == nil && after == nil {
		return tasks
	}

	result := make([]Task, 0)
	for _, task := range tasks {
		if task.Due == nil {
			continue
		}

		if before != nil && task.Due.After(*before) {
			continue
		}

		if after != nil && task.Due.Before(*after) {
			continue
		}

		result = append(result, task)
	}

	return result
}

// FilterByPriority filters tasks by priority level
func FilterByPriority(tasks []Task, priority int) []Task {
	if priority <= 0 {
		return tasks
	}

	result := make([]Task, 0)
	for _, task := range tasks {
		if task.Priority == priority {
			result = append(result, task)
		}
	}
	return result
}

// SearchTasks performs case-insensitive search in title and description
func SearchTasks(tasks []Task, query string) []Task {
	if query == "" {
		return tasks
	}

	query = strings.ToLower(query)
	result := make([]Task, 0)

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Title), query) ||
			strings.Contains(strings.ToLower(task.Description), query) {
			result = append(result, task)
		}
	}

	return result
}

// FilterByDependencies returns tasks that have all dependencies met
func FilterByDependencies(tasks []Task, onlyReady bool) []Task {
	if !onlyReady {
		return tasks
	}

	result := make([]Task, 0)
	for _, task := range tasks {
		if task.IsDone() {
			continue
		}
		if AllDependenciesMet(tasks, task) {
			result = append(result, task)
		}
	}
	return result
}

// GetBlockedTasks returns tasks that are blocked by incomplete dependencies
func GetBlockedTasks(tasks []Task) []Task {
	result := make([]Task, 0)
	for _, task := range tasks {
		if task.IsDone() {
			continue
		}
		if !AllDependenciesMet(tasks, task) {
			result = append(result, task)
		}
	}
	return result
}

