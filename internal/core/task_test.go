package core

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	tasks := []Task{}

	// Add first task
	tasks = Add(tasks, "Test task 1", nil)

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Title != "Test task 1" {
		t.Errorf("Expected title 'Test task 1', got '%s'", tasks[0].Title)
	}

	if tasks[0].ID != 1 {
		t.Errorf("Expected ID 1, got %d", tasks[0].ID)
	}

	// Add second task
	tasks = Add(tasks, "Test task 2", nil)

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	if tasks[1].ID != 2 {
		t.Errorf("Expected ID 2, got %d", tasks[1].ID)
	}
}

func TestTaskIsDone(t *testing.T) {
	now := time.Now()
	task := Task{
		ID:    1,
		Title: "Test",
	}

	if task.IsDone() {
		t.Error("New task should not be done")
	}

	task.DoneAt = &now

	if !task.IsDone() {
		t.Error("Task with DoneAt should be done")
	}
}

func TestTaskIsOverdue(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	tests := []struct {
		name     string
		task     Task
		expected bool
	}{
		{
			name:     "No due date",
			task:     Task{ID: 1, Title: "Test"},
			expected: false,
		},
		{
			name:     "Due tomorrow",
			task:     Task{ID: 1, Title: "Test", Due: &tomorrow},
			expected: false,
		},
		{
			name:     "Due yesterday",
			task:     Task{ID: 1, Title: "Test", Due: &yesterday},
			expected: true,
		},
		{
			name:     "Completed task that's overdue",
			task:     Task{ID: 1, Title: "Test", Due: &yesterday, DoneAt: &now},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.task.IsOverdue(now)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTaskHasTag(t *testing.T) {
	task := Task{
		ID:    1,
		Title: "Test",
		Tags:  []string{"work", "important"},
	}

	if !task.HasTag("work") {
		t.Error("Task should have 'work' tag")
	}

	if !task.HasTag("WORK") {
		t.Error("Tag matching should be case-insensitive")
	}

	if task.HasTag("personal") {
		t.Error("Task should not have 'personal' tag")
	}
}

func TestParseTags(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", nil},
		{"work", []string{"work"}},
		{"work,personal", []string{"work", "personal"}},
		{"work, personal, urgent", []string{"work", "personal", "urgent"}},
		{"  work  ,  personal  ", []string{"work", "personal"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseTags(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d tags, got %d", len(tt.expected), len(result))
				return
			}
			for i, tag := range result {
				if tag != tt.expected[i] {
					t.Errorf("Expected tag '%s', got '%s'", tt.expected[i], tag)
				}
			}
		})
	}
}

func TestParseIDs(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{"", nil},
		{"1", []int{1}},
		{"1,2,3", []int{1, 2, 3}},
		{"1, 2, 3", []int{1, 2, 3}},
		{"  1  ,  2  ", []int{1, 2}},
		{"1,invalid,3", []int{1, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseIDs(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d IDs, got %d", len(tt.expected), len(result))
				return
			}
			for i, id := range result {
				if id != tt.expected[i] {
					t.Errorf("Expected ID %d, got %d", tt.expected[i], id)
				}
			}
		})
	}
}

func TestMarkDone(t *testing.T) {
	tasks := []Task{
		{ID: 1, Title: "Task 1"},
		{ID: 2, Title: "Task 2"},
	}

	visible := tasks

	// Mark first task as done
	tasks, err := MarkDone(tasks, visible, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !tasks[0].IsDone() {
		t.Error("Task should be marked as done")
	}

	// Try to mark already done task
	_, err = MarkDone(tasks, visible, 1)
	if err == nil {
		t.Error("Expected error when marking already done task")
	}
}

func TestAllDependenciesMet(t *testing.T) {
	now := time.Now()
	tasks := []Task{
		{ID: 1, Title: "Task 1", DoneAt: &now},
		{ID: 2, Title: "Task 2"},
		{ID: 3, Title: "Task 3", DependsOn: []int{1}},
		{ID: 4, Title: "Task 4", DependsOn: []int{1, 2}},
	}

	// Task 3 depends on Task 1 which is done
	if !AllDependenciesMet(tasks, tasks[2]) {
		t.Error("Task 3's dependencies should be met")
	}

	// Task 4 depends on Task 1 (done) and Task 2 (not done)
	if AllDependenciesMet(tasks, tasks[3]) {
		t.Error("Task 4's dependencies should not be met")
	}
}

