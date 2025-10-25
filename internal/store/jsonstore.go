package store

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofrs/flock"
)

// JSONStore implements Store interface using JSON files
type JSONStore struct {
	filePath string
	mu       sync.RWMutex
    lock     *flock.Flock
}

// NewJSONStore creates a new JSON-based store
func NewJSONStore(filePath string) (*JSONStore, error) {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, err
	}

    return &JSONStore{
        filePath: filePath,
        lock:     flock.New(filePath + ".lock"),
    }, nil
}

// Load reads data from the JSON file
func (s *JSONStore) Load() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// If file doesn't exist, return empty array
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []byte("[]"), nil
	}

	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// If file is empty, return empty array
	if len(data) == 0 {
		return []byte("[]"), nil
	}

	return data, nil
}

// Save writes data to the JSON file atomically
func (s *JSONStore) Save(data []byte) error {
    s.mu.Lock()
    defer s.mu.Unlock()

	// Validate JSON before writing
	var test interface{}
	if err := json.Unmarshal(data, &test); err != nil {
		return err
	}

	// Write to temporary file first
	tempFile := s.filePath + ".tmp"
	file, err := os.OpenFile(tempFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		file.Close()
		os.Remove(tempFile)
		return err
	}

	// Sync to disk
	if err := file.Sync(); err != nil {
		file.Close()
		os.Remove(tempFile)
		return err
	}

	if err := file.Close(); err != nil {
		os.Remove(tempFile)
		return err
	}

	// Atomic rename
	return os.Rename(tempFile, s.filePath)
}

// Close implements Store interface (no-op for file-based storage)
func (s *JSONStore) Close() error {
	return nil
}

// WithExclusive acquires a cross-process exclusive lock for the duration of fn.
func (s *JSONStore) WithExclusive(ctx context.Context, fn func() error) error {
    // Try to acquire the file lock with retry/backoff honoring ctx
    // Use a short poll interval to avoid busy waiting
    ticker := time.NewTicker(50 * time.Millisecond)
    defer ticker.Stop()
    for {
        locked, err := s.lock.TryLock()
        if err != nil {
            return err
        }
        if locked {
            break
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
        }
    }
    defer s.lock.Unlock()
    return fn()
}

// LoadTasks is a helper to load and unmarshal tasks
func LoadTasks[T any](store Store) ([]T, error) {
	data, err := store.Load()
	if err != nil {
		return nil, err
	}

	var tasks []T
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// SaveTasks is a helper to marshal and save tasks
func SaveTasks[T any](store Store, tasks []T) error {
	data, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	return store.Save(data)
}

// DefaultStore returns a JSONStore using the default data file path
func DefaultStore() (*JSONStore, error) {
	filePath, err := GetDataFile()
	if err != nil {
		return nil, err
	}

	return NewJSONStore(filePath)
}

