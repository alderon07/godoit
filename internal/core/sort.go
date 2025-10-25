package core

import (
	"sort"
	"strings"
)

// SortKey represents the field to sort by
type SortKey string

const (
	SortByDue      SortKey = "due"
	SortByPriority SortKey = "priority"
	SortByCreated  SortKey = "created"
	SortByStatus   SortKey = "status"
	SortByTitle    SortKey = "title"
)

// SortedWith returns a filtered and sorted copy of tasks
func SortedWith(tasks []Task, showAll bool, grep, sortKey string) []Task {
	// First filter by status
	result := FilterByStatus(tasks, showAll)

	// Apply search filter
	if grep != "" {
		result = SearchTasks(result, grep)
	}

	// Sort
	SortTasks(result, SortKey(sortKey))

	return result
}

// SortTasks sorts tasks in place according to the specified key
func SortTasks(tasks []Task, key SortKey) {
	switch key {
	case SortByDue:
		sortByDue(tasks)
	case SortByPriority:
		sortByPriority(tasks)
	case SortByCreated:
		sortByCreated(tasks)
	case SortByStatus:
		sortByStatus(tasks)
	case SortByTitle:
		sortByTitle(tasks)
	default:
		sortByDue(tasks) // default
	}
}

// sortByDue sorts by due date (nil dates at the end), then by priority
func sortByDue(tasks []Task) {
	sort.SliceStable(tasks, func(i, j int) bool {
		// Both have due dates
		if tasks[i].Due != nil && tasks[j].Due != nil {
			if !tasks[i].Due.Equal(*tasks[j].Due) {
				return tasks[i].Due.Before(*tasks[j].Due)
			}
			// Same due date, sort by priority
			return tasks[i].Priority > tasks[j].Priority
		}

		// Only i has due date
		if tasks[i].Due != nil {
			return true
		}

		// Only j has due date
		if tasks[j].Due != nil {
			return false
		}

		// Neither has due date, sort by priority
		return tasks[i].Priority > tasks[j].Priority
	})
}

// sortByPriority sorts by priority (high to low), then by due date
func sortByPriority(tasks []Task) {
	sort.SliceStable(tasks, func(i, j int) bool {
		if tasks[i].Priority != tasks[j].Priority {
			return tasks[i].Priority > tasks[j].Priority
		}

		// Same priority, sort by due date
		if tasks[i].Due != nil && tasks[j].Due != nil {
			return tasks[i].Due.Before(*tasks[j].Due)
		}

		if tasks[i].Due != nil {
			return true
		}

		return false
	})
}

// sortByCreated sorts by creation date (newest first)
func sortByCreated(tasks []Task) {
	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})
}

// sortByStatus sorts by completion status (pending first, then completed)
func sortByStatus(tasks []Task) {
	sort.SliceStable(tasks, func(i, j int) bool {
		iDone := tasks[i].IsDone()
		jDone := tasks[j].IsDone()

		if iDone != jDone {
			return !iDone // pending tasks first
		}

		// Same status, sort by due date
		if tasks[i].Due != nil && tasks[j].Due != nil {
			return tasks[i].Due.Before(*tasks[j].Due)
		}

		if tasks[i].Due != nil {
			return true
		}

		return false
	})
}

// sortByTitle sorts alphabetically by title
func sortByTitle(tasks []Task) {
	sort.SliceStable(tasks, func(i, j int) bool {
		return strings.ToLower(tasks[i].Title) < strings.ToLower(tasks[j].Title)
	})
}

// MultiSort allows sorting by multiple criteria
func MultiSort(tasks []Task, keys []SortKey) {
	if len(keys) == 0 {
		return
	}

	// Apply sorts in reverse order (last key has least priority)
	for i := len(keys) - 1; i >= 0; i-- {
		SortTasks(tasks, keys[i])
	}
}

