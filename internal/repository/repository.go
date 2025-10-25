package repository

import (
	"context"
	"encoding/json"

	"godoit/internal/core"
	"godoit/internal/store"
)

// TaskRepository abstracts persistence for tasks.
type TaskRepository interface {
    LoadTasks(ctx context.Context) ([]core.Task, error)
    SaveTasks(ctx context.Context, tasks []core.Task) error
}

// JSONTaskRepository implements TaskRepository over store.Store (JSON file).
type JSONTaskRepository struct {
    store store.Store
}

func NewJSONTaskRepository(s store.Store) *JSONTaskRepository {
    return &JSONTaskRepository{store: s}
}

func (r *JSONTaskRepository) LoadTasks(_ context.Context) ([]core.Task, error) {
    data, err := r.store.Load()
    if err != nil {
        return nil, err
    }
    var tasks []core.Task
    if err := json.Unmarshal(data, &tasks); err != nil {
        return nil, err
    }
    return tasks, nil
}

func (r *JSONTaskRepository) SaveTasks(_ context.Context, tasks []core.Task) error {
    data, err := json.Marshal(tasks)
    if err != nil {
        return err
    }
    return r.store.Save(data)
}


