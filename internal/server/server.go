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

	"godoit/internal/clock"
	"godoit/internal/core"
	"godoit/internal/repository"
	"godoit/internal/service"
	"godoit/internal/store"
)

// Server represents the HTTP API server
type Server struct {
    store  *store.JSONStore
	mux    *http.ServeMux
	server *http.Server
    svc    *service.TaskService
}

// NewServer creates a new HTTP server
func NewServer(host string, port int, s *store.JSONStore) *Server {
	mux := http.NewServeMux()

    // wire repository and service
    repo := repository.NewJSONTaskRepository(s)
	svc := service.NewTaskService(repo, clock.SystemClock{})

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
        svc: svc,
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
	q := r.URL.Query()
	showAll := q.Get("all") == "true"
	grep := q.Get("grep")
	tags := q.Get("tags")
	sortKey := q.Get("sort")
	if sortKey == "" { sortKey = "due" }
	var beforePtr, afterPtr *time.Time
	if bs := q.Get("before"); bs != "" {
		if t, err := time.Parse("2006-01-02", bs); err == nil { beforePtr = &t }
	}
	if as := q.Get("after"); as != "" {
		if t, err := time.Parse("2006-01-02", as); err == nil { afterPtr = &t }
	}
	result, err := s.svc.QueryTasks(r.Context(), service.Query{
		ShowAll: showAll,
		Grep:    grep,
		SortKey: sortKey,
		Tags:    tags,
		Before:  beforePtr,
		After:   afterPtr,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, result)
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

	var due *time.Time
	if input.Due != nil && *input.Due != "" {
		if t, err := time.Parse("2006-01-02", *input.Due); err == nil { due = &t } else {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}
	}
	created, err := s.svc.AddTask(r.Context(), service.AddTaskInput{
		Title:       input.Title,
		Description: input.Description,
		Due:         due,
		Priority:    input.Priority,
		Tags:        input.Tags,
		Repeat:      input.Repeat,
		DependsOn:   input.DependsOn,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, created)
}

// getTask returns a single task by ID
func (s *Server) getTask(w http.ResponseWriter, r *http.Request, id int) {
	task, err := s.svc.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	respondJSON(w, task)
}

// updateTask updates an existing task
func (s *Server) updateTask(w http.ResponseWriter, r *http.Request, id int) {
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

	updated, err := s.svc.UpdateTask(r.Context(), id, service.UpdateTaskInput{
		Title:       input.Title,
		Description: input.Description,
		Due:         input.Due,
		Priority:    input.Priority,
		Tags:        input.Tags,
		Repeat:      input.Repeat,
		DependsOn:   input.DependsOn,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, updated)
}

// deleteTask deletes a task by ID
func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	if err := s.svc.DeleteTaskByID(r.Context(), id); err != nil {
		if err.Error() == "task not found" {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// markDone marks a task as complete
func (s *Server) markDone(w http.ResponseWriter, r *http.Request, id int) {
	updated, err := s.svc.MarkDoneByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, updated)
}

// handleStats returns task statistics
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := s.svc.QueryTasks(r.Context(), service.Query{ShowAll: true, SortKey: "due"})
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

