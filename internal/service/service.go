package service

import (
	"context"
	"fmt"
	"time"

	"godoit/internal/clock"
	"godoit/internal/core"
	"godoit/internal/repository"
)

type AddTaskInput struct {
    Title       string
    Description string
    Due         *time.Time
    Priority    int
    Tags        []string
    Repeat      string
    DependsOn   []int
}

type UpdateTaskInput struct {
    Title       *string
    Description *string
    Due         *string // YYYY-MM-DD or empty to clear
    Priority    *int
    Tags        *[]string
    Repeat      *string
    DependsOn   *[]int
}

type Query struct {
    ShowAll bool
    Grep    string
    SortKey string
    Tags    string // raw form; reused from existing semantics
    Before  *time.Time
    After   *time.Time
}

type TaskService struct {
    repo  repository.TaskRepository
    clock clock.Clock
}

func NewTaskService(repo repository.TaskRepository, clk clock.Clock) *TaskService {
    return &TaskService{repo: repo, clock: clk}
}

func (s *TaskService) AddTask(ctx context.Context, in AddTaskInput) (core.Task, error) {
    if in.Title == "" {
        return core.Task{}, fmt.Errorf("title is required")
    }
    tasks, err := s.repo.LoadTasks(ctx)
    if err != nil { return core.Task{}, err }

    // use injected clock for deterministic CreatedAt
    now := s.clock.Now()
    tasks = core.AddAt(tasks, in.Title, in.Due, now)
    t := &tasks[len(tasks)-1]
    t.Description = in.Description
    t.Priority = core.NormalizePriority(in.Priority)
    t.Tags = in.Tags
    t.Repeat = core.NormalizeRepeat(in.Repeat)
    t.DependsOn = in.DependsOn

    if err := s.repo.SaveTasks(ctx, tasks); err != nil { return core.Task{}, err }
    return *t, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int, in UpdateTaskInput) (core.Task, error) {
    tasks, err := s.repo.LoadTasks(ctx)
    if err != nil { return core.Task{}, err }
    task, err := core.GetByID(tasks, id)
    if err != nil { return core.Task{}, err }

    if in.Title != nil { task.Title = *in.Title }
    if in.Description != nil { task.Description = *in.Description }
    if in.Due != nil {
        if *in.Due == "" { task.Due = nil } else if t, err := time.Parse("2006-01-02", *in.Due); err == nil { task.Due = &t } else { return core.Task{}, fmt.Errorf("invalid due date") }
    }
    if in.Priority != nil {
        task.Priority = core.NormalizePriority(*in.Priority)
    }
    if in.Tags != nil { task.Tags = *in.Tags }
    if in.Repeat != nil { task.Repeat = core.NormalizeRepeat(*in.Repeat) }
    if in.DependsOn != nil { task.DependsOn = *in.DependsOn }

    tasks, err = core.Update(tasks, *task)
    if err != nil { return core.Task{}, err }
    if err := s.repo.SaveTasks(ctx, tasks); err != nil { return core.Task{}, err }
    return *task, nil
}

func (s *TaskService) RemoveTask(ctx context.Context, visible []core.Task, idx int) error {
    tasks, err := s.repo.LoadTasks(ctx)
    if err != nil { return err }
    tasks, err = core.Remove(tasks, visible, idx)
    if err != nil { return err }
    return s.repo.SaveTasks(ctx, tasks)
}

func (s *TaskService) MarkDone(ctx context.Context, visible []core.Task, idx int) (core.Task, error) {
    tasks, err := s.repo.LoadTasks(ctx)
    if err != nil { return core.Task{}, err }
    before := tasks
    // use injected clock for deterministic DoneAt and recurrence
    tasks, err = core.MarkDoneAt(tasks, visible, 1, s.clock.Now())
    if err != nil { return core.Task{}, err }
    if err := s.repo.SaveTasks(ctx, tasks); err != nil { return core.Task{}, err }
    // Return the updated task (visible[idx-1] maps by ID)
    updated, _ := core.GetByID(tasks, visible[idx-1].ID)
    if updated != nil { return *updated, nil }
    // Fallback: in recurrence cases the original might be done; return done task by matching title/time
    if len(before) > 0 { return before[0], nil }
    return core.Task{}, nil
}

func (s *TaskService) QueryTasks(ctx context.Context, q Query) ([]core.Task, error) {
    tasks, err := s.repo.LoadTasks(ctx)
    if err != nil { return nil, err }
    // Apply layered filters similar to existing code
    result := core.SortedWith(tasks, q.ShowAll, q.Grep, q.SortKey)
    result = core.FilterByTags(result, q.Tags)
    if q.Before != nil || q.After != nil {
        before := ""; after := ""
        if q.Before != nil { before = q.Before.Format("2006-01-02") }
        if q.After != nil { after = q.After.Format("2006-01-02") }
        result = core.FilterByDate(result, before, after)
    }
    return result, nil
}

func (s *TaskService) GetTask(ctx context.Context, id int) (core.Task, error) {
    tasks, err := s.repo.LoadTasks(ctx)
    if err != nil { return core.Task{}, err }
    t, err := core.GetByID(tasks, id)
    if err != nil { return core.Task{}, err }
    return *t, nil
}

func (s *TaskService) DeleteTaskByID(ctx context.Context, id int) error {
    tasks, err := s.repo.LoadTasks(ctx)
    if err != nil { return err }
    found := false
    newTasks := make([]core.Task, 0, len(tasks))
    for _, t := range tasks {
        if t.ID == id { found = true; continue }
        newTasks = append(newTasks, t)
    }
    if !found { return fmt.Errorf("task not found") }
    return s.repo.SaveTasks(ctx, newTasks)
}

func (s *TaskService) MarkDoneByID(ctx context.Context, id int) (core.Task, error) {
    tasks, err := s.repo.LoadTasks(ctx)
    if err != nil { return core.Task{}, err }
    // create a visible slice containing the specific task
    idx := -1
    for i, t := range tasks { if t.ID == id { idx = i; break } }
    if idx == -1 { return core.Task{}, fmt.Errorf("task not found") }
    visible := []core.Task{tasks[idx]}
    tasks, err = core.MarkDone(tasks, visible, 1)
    if err != nil { return core.Task{}, err }
    if err := s.repo.SaveTasks(ctx, tasks); err != nil { return core.Task{}, err }
    updated, err := core.GetByID(tasks, id)
    if err != nil { return core.Task{}, err }
    return *updated, nil
}


