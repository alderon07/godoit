# godo HTTP API Documentation

## Overview

The godo HTTP API provides a RESTful interface for managing tasks. All endpoints return JSON responses.

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, no authentication is required. This is suitable for local development only.

## Content Type

All POST and PUT requests must include:

```
Content-Type: application/json
```

## Endpoints

### Health Check

Check if the server is running.

**Request:**

```
GET /health
```

**Response:**

```json
{
  "status": "ok",
  "time": "2025-10-24T12:00:00Z"
}
```

---

### List Tasks

Retrieve all tasks with optional filtering and sorting.

**Request:**

```
GET /tasks
```

**Query Parameters:**

- `all` (boolean): Include completed tasks (default: false)
- `grep` (string): Search keyword (case-insensitive)
- `tags` (string): Filter by tags (comma=OR, plus=AND)
- `sort` (string): Sort key - `due`, `priority`, `created`, `status`, `title` (default: due)
- `before` (string): Filter tasks before date (YYYY-MM-DD)
- `after` (string): Filter tasks after date (YYYY-MM-DD)

**Example:**

```
GET /tasks?all=true&sort=priority&tags=work
```

**Response:**

```json
[
  {
    "id": 1,
    "title": "Complete project proposal",
    "description": "Draft and submit Q4 proposal",
    "due": "2025-10-31T00:00:00Z",
    "done_at": null,
    "created_at": "2025-10-24T10:00:00Z",
    "priority": 3,
    "tags": ["work", "important"],
    "repeat": "",
    "depends_on": []
  },
  {
    "id": 2,
    "title": "Review pull requests",
    "description": "",
    "due": "2025-10-25T00:00:00Z",
    "done_at": null,
    "created_at": "2025-10-24T11:00:00Z",
    "priority": 2,
    "tags": ["work"],
    "repeat": "daily",
    "depends_on": [1]
  }
]
```

---

### Create Task

Create a new task.

**Request:**

```
POST /tasks
```

**Body:**

```json
{
  "title": "Task title",
  "description": "Optional description",
  "due": "2025-10-31",
  "priority": 2,
  "tags": ["work", "important"],
  "repeat": "weekly",
  "depends_on": [1, 2]
}
```

**Required Fields:**

- `title` (string): Task title

**Optional Fields:**

- `description` (string): Task description
- `due` (string): Due date in YYYY-MM-DD format
- `priority` (integer): Priority level (1=low, 2=medium, 3=high)
- `tags` (array of strings): Task tags
- `repeat` (string): Repeat rule - `daily`, `weekly`, or `monthly`
- `depends_on` (array of integers): IDs of tasks this task depends on

**Response:**

```json
{
  "id": 3,
  "title": "Task title",
  "description": "Optional description",
  "due": "2025-10-31T00:00:00Z",
  "done_at": null,
  "created_at": "2025-10-24T12:00:00Z",
  "priority": 2,
  "tags": ["work", "important"],
  "repeat": "weekly",
  "depends_on": [1, 2]
}
```

**Status Codes:**

- `200 OK`: Task created successfully
- `400 Bad Request`: Invalid input (missing title, invalid date format, etc.)
- `500 Internal Server Error`: Server error

---

### Get Task

Retrieve a single task by ID.

**Request:**

```
GET /tasks/:id
```

**Example:**

```
GET /tasks/5
```

**Response:**

```json
{
  "id": 5,
  "title": "Task title",
  "description": "Task description",
  "due": "2025-10-31T00:00:00Z",
  "done_at": null,
  "created_at": "2025-10-24T12:00:00Z",
  "priority": 2,
  "tags": ["work"],
  "repeat": "",
  "depends_on": []
}
```

**Status Codes:**

- `200 OK`: Task found
- `404 Not Found`: Task not found
- `400 Bad Request`: Invalid task ID

---

### Update Task

Update an existing task.

**Request:**

```
PUT /tasks/:id
```

**Body:**

```json
{
  "title": "Updated title",
  "description": "Updated description",
  "due": "2025-11-01",
  "priority": 3,
  "tags": ["work", "urgent"],
  "repeat": "monthly",
  "depends_on": [1]
}
```

**Notes:**

- All fields are optional
- Only provided fields will be updated
- To clear a field, set it to empty string (for `due`, `repeat`) or empty array (for `tags`, `depends_on`)

**Response:**

```json
{
  "id": 5,
  "title": "Updated title",
  "description": "Updated description",
  "due": "2025-11-01T00:00:00Z",
  "done_at": null,
  "created_at": "2025-10-24T12:00:00Z",
  "priority": 3,
  "tags": ["work", "urgent"],
  "repeat": "monthly",
  "depends_on": [1]
}
```

**Status Codes:**

- `200 OK`: Task updated successfully
- `404 Not Found`: Task not found
- `400 Bad Request`: Invalid input
- `500 Internal Server Error`: Server error

---

### Delete Task

Delete a task.

**Request:**

```
DELETE /tasks/:id
```

**Example:**

```
DELETE /tasks/5
```

**Response:**
No content (empty body)

**Status Codes:**

- `204 No Content`: Task deleted successfully
- `404 Not Found`: Task not found
- `500 Internal Server Error`: Server error

---

### Mark Task as Done

Mark a task as complete.

**Request:**

```
POST /tasks/:id/done
```

**Example:**

```
POST /tasks/5/done
```

**Response:**

```json
{
  "id": 5,
  "title": "Task title",
  "description": "Task description",
  "due": "2025-10-31T00:00:00Z",
  "done_at": "2025-10-24T15:30:00Z",
  "created_at": "2025-10-24T12:00:00Z",
  "priority": 2,
  "tags": ["work"],
  "repeat": "",
  "depends_on": []
}
```

**Notes:**

- If the task has a `repeat` rule, a new occurrence will be automatically created
- The `done_at` field will be set to the current timestamp
- Returns an error if dependencies are not met

**Status Codes:**

- `200 OK`: Task marked as done
- `404 Not Found`: Task not found
- `400 Bad Request`: Task already done or dependencies not met
- `500 Internal Server Error`: Server error

---

### Get Statistics

Retrieve task statistics and analytics.

**Request:**

```
GET /stats
```

**Response:**

```json
{
  "Total": 25,
  "Completed": 15,
  "Pending": 10,
  "Overdue": 2,
  "CompletionRate": 60.0,
  "ByPriority": {
    "1": 5,
    "2": 12,
    "3": 8
  },
  "ByTag": {
    "work": 15,
    "personal": 7,
    "urgent": 3
  },
  "AvgCompletionMS": 172800000,
  "CompletedToday": 3,
  "CompletedWeek": 8,
  "BlockedTasks": 1
}
```

**Field Descriptions:**

- `Total`: Total number of tasks
- `Completed`: Number of completed tasks
- `Pending`: Number of pending tasks
- `Overdue`: Number of overdue tasks
- `CompletionRate`: Percentage of completed tasks
- `ByPriority`: Task count by priority level
- `ByTag`: Task count by tag
- `AvgCompletionMS`: Average time to complete tasks (in milliseconds)
- `CompletedToday`: Tasks completed today
- `CompletedWeek`: Tasks completed this week
- `BlockedTasks`: Tasks blocked by dependencies

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

Common error codes:

- `400 Bad Request`: Invalid input or request
- `404 Not Found`: Resource not found
- `405 Method Not Allowed`: HTTP method not supported for endpoint
- `500 Internal Server Error`: Server-side error

---

## CORS

The API includes CORS headers allowing cross-origin requests from any domain. This is suitable for development but should be restricted in production.

---

## Examples with curl

### Create a task

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project",
    "priority": 3,
    "tags": ["work"],
    "due": "2025-10-31"
  }'
```

### List all tasks

```bash
curl http://localhost:8080/tasks
```

### List pending tasks sorted by priority

```bash
curl "http://localhost:8080/tasks?sort=priority"
```

### Search for tasks

```bash
curl "http://localhost:8080/tasks?grep=meeting"
```

### Get a specific task

```bash
curl http://localhost:8080/tasks/5
```

### Update a task

```bash
curl -X PUT http://localhost:8080/tasks/5 \
  -H "Content-Type: application/json" \
  -d '{
    "priority": 3,
    "tags": ["work", "urgent"]
  }'
```

### Mark task as done

```bash
curl -X POST http://localhost:8080/tasks/5/done
```

### Delete a task

```bash
curl -X DELETE http://localhost:8080/tasks/5
```

### Get statistics

```bash
curl http://localhost:8080/stats
```

---

## Examples with JavaScript/Fetch

### Create a task

```javascript
const response = await fetch("http://localhost:8080/tasks", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    title: "Complete project",
    priority: 3,
    tags: ["work"],
    due: "2025-10-31",
  }),
});
const task = await response.json();
console.log(task);
```

### List tasks

```javascript
const response = await fetch("http://localhost:8080/tasks?sort=priority");
const tasks = await response.json();
console.log(tasks);
```

### Mark task as done

```javascript
const response = await fetch("http://localhost:8080/tasks/5/done", {
  method: "POST",
});
const updatedTask = await response.json();
console.log(updatedTask);
```

---

## Rate Limiting

Currently, no rate limiting is implemented. This should be added for production use.

---

## Best Practices

1. **Always include Content-Type header** for POST/PUT requests
2. **Handle errors gracefully** - check response status codes
3. **Use appropriate HTTP methods** - GET for reading, POST for creating, PUT for updating, DELETE for deleting
4. **Validate dates** before sending (use YYYY-MM-DD format)
5. **Check for dependencies** before marking tasks as done
6. **Use filtering and sorting** to reduce data transfer

---

## Future Enhancements

Planned API improvements:

- Authentication and authorization
- Rate limiting
- Webhooks for task updates
- Batch operations
- Task search with advanced queries
- Export/import endpoints
- WebSocket support for real-time updates
