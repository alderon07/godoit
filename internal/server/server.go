package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"godo/internal/core"
	"godo/internal/store"
)

// Server represents the HTTP API server
type Server struct {
	store  *store.JSONStore
	mux    *http.ServeMux
	server *http.Server
}

// NewServer creates a new HTTP server
func NewServer(host string, port int, s *store.JSONStore) *Server {
	mux := http.NewServeMux()

	srv := &Server{
		store: s,
		mux:   mux,
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			Handler:      mux,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	srv.setupRoutes()

	return srv
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/tasks", s.corsMiddleware(s.handleTasks))
	s.mux.HandleFunc("/tasks/", s.corsMiddleware(s.handleTask))
	s.mux.HandleFunc("/stats", s.corsMiddleware(s.handleStats))
	s.mux.HandleFunc("/health", s.handleHealth)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Starting server on %s\n", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop

	log.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// handleTasks handles /tasks endpoint (list and create)
func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.listTasks(w, r)
	case "POST":
		s.createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTask handles /tasks/:id endpoint (get, update, delete)
func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Check for /tasks/:id/done endpoint
	if len(parts) > 1 && parts[1] == "done" && r.Method == "POST" {
		s.markDone(w, r, id)
		return
	}

	switch r.Method {
	case "GET":
		s.getTask(w, r, id)
	case "PUT":
		s.updateTask(w, r, id)
	case "DELETE":
		s.deleteTask(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listTasks returns all tasks with optional filtering
func (s *Server) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := store.LoadTasks[core.Task](s.store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply filters from query params
	query := r.URL.Query()

	// Filter by status
	showAll := query.Get("all") == "true"
	tasks = core.FilterByStatus(tasks, showAll)

	// Search
	if grep := query.Get("grep"); grep != "" {
		tasks = core.SearchTasks(tasks, grep)
	}

	// Filter by tags
	if tags := query.Get("tags"); tags != "" {
		tasks = core.FilterByTags(tasks, tags)
	}

	// Filter by date
	before := query.Get("before")
	after := query.Get("after")
	if before != "" || after != "" {
		tasks = core.FilterByDate(tasks, before, after)
	}

	// Sort
	sortKey := query.Get("sort")
	if sortKey == "" {
		sortKey = "due"
	}
	core.SortTasks(tasks, core.SortKey(sortKey))

	respondJSON(w, tasks)
}

// createTask creates a new task
func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Due         *string  `json:"due"`
		Priority    int      `json:"priority"`
		Tags        []string `json:"tags"`
		Repeat      string   `json:"repeat"`
		DependsOn   []int    `json:"depends_on"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	tasks, err := store.LoadTasks[core.Task](s.store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var due *time.Time
	if input.Due != nil && *input.Due != "" {
		t, err := time.Parse("2006-01-02", *input.Due)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
		due = &t
	}

	tasks = core.Add(tasks, input.Title, due)
	newTask := &tasks[len(tasks)-1]

	newTask.Description = input.Description
	newTask.Priority = input.Priority
	newTask.Tags = input.Tags
	newTask.Repeat = input.Repeat
	newTask.DependsOn = input.DependsOn

	if err := store.SaveTasks(s.store, tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, newTask)
}

// getTask returns a single task by ID
func (s *Server) getTask(w http.ResponseWriter, r *http.Request, id int) {
	tasks, err := store.LoadTasks[core.Task](s.store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task, err := core.GetByID(tasks, id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	respondJSON(w, task)
}

// updateTask updates an existing task
func (s *Server) updateTask(w http.ResponseWriter, r *http.Request, id int) {
	tasks, err := store.LoadTasks[core.Task](s.store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task, err := core.GetByID(tasks, id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	var input struct {
		Title       *string   `json:"title"`
		Description *string   `json:"description"`
		Due         *string   `json:"due"`
		Priority    *int      `json:"priority"`
		Tags        *[]string `json:"tags"`
		Repeat      *string   `json:"repeat"`
		DependsOn   *[]int    `json:"depends_on"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update fields if provided
	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Due != nil {
		if *input.Due == "" {
			task.Due = nil
		} else {
			t, err := time.Parse("2006-01-02", *input.Due)
			if err != nil {
				http.Error(w, "Invalid date format", http.StatusBadRequest)
				return
			}
			task.Due = &t
		}
	}
	if input.Priority != nil {
		task.Priority = *input.Priority
	}
	if input.Tags != nil {
		task.Tags = *input.Tags
	}
	if input.Repeat != nil {
		task.Repeat = *input.Repeat
	}
	if input.DependsOn != nil {
		task.DependsOn = *input.DependsOn
	}

	tasks, err = core.Update(tasks, *task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := store.SaveTasks(s.store, tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, task)
}

// deleteTask deletes a task by ID
func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	tasks, err := store.LoadTasks[core.Task](s.store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find the task in all tasks
	found := false
	newTasks := make([]core.Task, 0, len(tasks)-1)

	for _, t := range tasks {
		if t.ID == id {
			found = true
			continue
		}
		newTasks = append(newTasks, t)
	}

	if !found {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	if err := store.SaveTasks(s.store, newTasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// markDone marks a task as complete
func (s *Server) markDone(w http.ResponseWriter, r *http.Request, id int) {
	tasks, err := store.LoadTasks[core.Task](s.store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find task index
	taskIdx := -1
	for i, t := range tasks {
		if t.ID == id {
			taskIdx = i
			break
		}
	}

	if taskIdx == -1 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Create a visible list with just this task
	visible := []core.Task{tasks[taskIdx]}

	tasks, err = core.MarkDone(tasks, visible, 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := store.SaveTasks(s.store, tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find and return the updated task
	for _, t := range tasks {
		if t.ID == id {
			respondJSON(w, t)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// handleStats returns task statistics
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := store.LoadTasks[core.Task](s.store)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats := core.CalculateStats(tasks, time.Now())
	respondJSON(w, stats)
}

// handleHealth returns health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

