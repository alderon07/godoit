package store

// Store defines the interface for task storage operations
// This abstraction allows for future implementations (SQLite, PostgreSQL, etc.)
type Store interface {
	// Load retrieves all tasks from storage
	Load() ([]byte, error)

	// Save persists tasks to storage
	Save(data []byte) error

	// Close closes any open resources
	Close() error
}

