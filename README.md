# godo

A feature-rich, cross-platform command-line TODO list manager written in Go with support for alerts, recurring tasks, priorities, filtering, search, tagging, task dependencies, analytics, and an HTTP REST API.

## Features

- ✅ Simple and intuitive CLI interface
- 📝 Create tasks with titles, descriptions, due dates, and priorities
- 🏷️ Tag system for task organization
- 🔍 Advanced filtering and search capabilities
- 📊 Task analytics and statistics
- 🔗 Task dependencies (block tasks until dependencies are complete)
- 🔄 Recurring tasks (daily, weekly, monthly)
- 🔔 Desktop notifications for due/overdue tasks
- 👀 Watch mode for continuous monitoring
- 🌐 HTTP REST API server
- 💾 JSON-based storage with atomic writes
- 🔀 Cross-platform support (Linux, macOS, Windows)

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/yourusername/godo/releases).

Available platforms:

- Linux (amd64, arm64)
- macOS (amd64/Intel, arm64/Apple Silicon)
- Windows (amd64)

### Build from Source

Requires Go 1.25 or higher.

```bash
# Clone the repository
git clone https://github.com/yourusername/godo.git
cd godo

# Build
make build

# Or build for all platforms
make build-all
```

The binary will be available at `bin/godo`.

### Install to System

```bash
make install
```

This installs the binary to `$GOPATH/bin/godo`.

## Usage

### Basic Commands

```bash
godo <command> [options]
```

### Add a Task

```bash
godo add -title "Task description" [options]
```

**Options:**

- `-title` (required): Task title/description
- `-desc`: Short description with more details (optional)
- `-due YYYY-MM-DD`: Set due date
- `-p <1-3>`: Set priority level (1=low, 2=medium, 3=high)
- `-tags "tag1,tag2"`: Add comma-separated tags
- `-repeat "daily|weekly|monthly"`: Set repeat rule for recurring tasks
- `-after "1,2,3"`: Comma-separated dependency task IDs

**Examples:**

```bash
# Simple task
godo add -title "Buy groceries"

# Task with description and priority
godo add -title "Review pull requests" -desc "Check and merge all pending PRs" -due 2025-10-25 -p 3

# Task with tags
godo add -title "Team meeting" -desc "Weekly sync with the team" -tags "work,meeting" -p 2

# Recurring task
godo add -title "Weekly report" -desc "Submit weekly progress report" -repeat weekly -due 2025-10-27

# Task with dependencies (can't start until tasks 1 and 2 are done)
godo add -title "Deploy to production" -desc "Deploy latest build to prod environment" -after "1,2" -p 3
```

### List Tasks

```bash
godo list [options]
```

**Options:**

- `-all`: Show completed tasks as well
- `-today`: Show only today's tasks
- `-week`: Show only this week's tasks
- `-detailed`: Show detailed information including descriptions and timestamps
- `-grep "keyword"`: Filter by substring (case-insensitive)
- `-tags "tag1,tag2"`: Filter by tags
  - Comma-separated = OR logic (task has ANY of these tags)
  - Plus-separated = AND logic (task has ALL of these tags)
- `-sort <key>`: Sort by: `due`, `priority`, `created`, `status`, or `title` (default: `due`)
- `-before YYYY-MM-DD`: Show tasks before date
- `-after YYYY-MM-DD`: Show tasks after date

**Examples:**

```bash
# List all pending tasks
godo list

# Show detailed information (includes descriptions)
godo list -detailed

# Show only today's tasks
godo list -today

# Show this week's tasks with details
godo list -week -detailed

# Show all tasks including completed
godo list -all

# Filter by tag (OR logic)
godo list -tags "work,urgent"

# Filter by tag (AND logic - must have both tags)
godo list -tags "work+urgent"

# Search tasks
godo list -grep "meeting"

# Sort by priority
godo list -sort priority

# Show tasks due in the next week
godo list -before 2025-10-31

# Combine filters
godo list -tags "work" -sort priority -grep "review" -detailed
```

### Understanding Task Display

The list view uses visual indicators (emojis) to make information easy to scan:

**Status**:

- `[ ]` - Pending task
- `[✓]` - Completed task

**Priority** (color-coded):

- 🔴 - High priority (p3)
- 🟡 - Medium priority (p2)
- 🟢 - Low priority (p1)

**Task Information**:

- 📝 - Task description (shown in `-detailed` view)
- 📅 - Due date (normal)
- ⏰ - Due soon or overdue (with indicators: "soon" or "OVERDUE!")
- 🏷️ - Tags
- 🔄 - Recurring task (daily/weekly/monthly)
- 🔗 - Task has dependencies
- ⚠️ - Task is blocked (dependencies not met)
- 🕐 - Created timestamp (shown in `-detailed` view)
- ✅ - Completion timestamp (shown in `-detailed` view)

**Example Output**:

```
 1. [ ] 🔴 Complete project proposal
    📝 Write and submit the Q4 project proposal
    ⏰ Due: 2025-10-26 00:00 (soon)
    🏷️  #work #important
    🕐 Created: 2025-10-23 23:49

 2. [✓] 🟡 Review pull requests
    📝 Check and merge pending PRs
    🏷️  #work #code-review
    ✅ Completed: 2025-10-23 15:30
```

### Mark Task as Done

```bash
godo done <index>
```

The index is shown in the list command. If the task is recurring, a new occurrence will be automatically created based on the repeat rule.

**Example:**

```bash
godo done 3
```

### Edit a Task

```bash
godo edit <index> [options]
```

**Options:**

- `-title "new title"`: Update title
- `-desc "new description"`: Update description (use "none" to clear)
- `-due YYYY-MM-DD`: Update due date (use "none" to clear)
- `-p <1-3>`: Update priority
- `-tags "tag1,tag2"`: Update tags (use "none" to clear)
- `-repeat "daily|weekly|monthly"`: Update repeat rule (use "none" to clear)
- `-after "1,2"`: Update dependencies (use "none" to clear)

**Examples:**

```bash
# Change title
godo edit 2 -title "Updated task name"

# Update description
godo edit 2 -desc "Updated description with more details"

# Update due date and priority
godo edit 1 -due 2025-11-01 -p 3

# Remove due date
godo edit 1 -due none

# Remove description
godo edit 2 -desc none

# Add tags
godo edit 3 -tags "important,urgent"
```

### Remove a Task

```bash
godo remove <index>
# or
godo rm <index>
```

**Example:**

```bash
godo rm 3
```

### View Alerts

Show due/overdue tasks and blocked tasks (tasks waiting on dependencies):

```bash
godo alerts [options]
```

**Options:**

- `-watch`: Continuously monitor for upcoming tasks
- `-interval <duration>`: Polling interval for watch mode (default: 60s)
- `-ahead <duration>`: Lookahead window for alerts (default: 24h)

**Examples:**

```bash
# One-time scan
godo alerts

# Watch mode with desktop notifications
godo alerts -watch

# Custom intervals
godo alerts -watch -interval 5m -ahead 48h
```

### View Statistics

```bash
godo stats
```

Shows:

- Total, completed, pending, overdue, and blocked tasks
- Completion rate
- Tasks completed today and this week
- Average completion time
- Breakdown by priority
- Breakdown by tag

### HTTP API Server

Start an HTTP REST API server:

```bash
godo server [options]
```

**Options:**

- `-host <hostname>`: Host to bind to (default: localhost)
- `-port <port>`: Port to listen on (default: 8080)

**Example:**

```bash
godo server -port 8080
```

## HTTP API Reference

The HTTP server provides a RESTful API for managing tasks.

### Endpoints

#### List Tasks

```
GET /tasks?all=true&grep=keyword&tags=work&sort=priority&before=2025-12-31&after=2025-01-01
```

Query parameters (all optional):

- `all`: Include completed tasks (true/false)
- `grep`: Search keyword
- `tags`: Filter by tags
- `sort`: Sort key (due, priority, created, status, title)
- `before`: Filter before date (YYYY-MM-DD)
- `after`: Filter after date (YYYY-MM-DD)

#### Create Task

```
POST /tasks
Content-Type: application/json

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

#### Get Single Task

```
GET /tasks/:id
```

#### Update Task

```
PUT /tasks/:id
Content-Type: application/json

{
  "title": "Updated title",
  "priority": 3
}
```

#### Delete Task

```
DELETE /tasks/:id
```

#### Mark Task as Done

```
POST /tasks/:id/done
```

#### Get Statistics

```
GET /stats
```

#### Health Check

```
GET /health
```

### Example API Usage

```bash
# Create a task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"New task","priority":2,"tags":["work"]}'

# List all tasks
curl http://localhost:8080/tasks

# Mark task 5 as done
curl -X POST http://localhost:8080/tasks/5/done

# Get statistics
curl http://localhost:8080/stats
```

## Project Structure

```
godo/
├── cmd/
│   └── todo/
│       ├── main.go         # CLI entry point
│       └── commands.go     # Command implementations
├── internal/
│   ├── core/               # Core task management
│   │   ├── task.go         # Task struct and operations
│   │   ├── filter.go       # Filtering utilities
│   │   ├── sort.go         # Sorting utilities
│   │   └── stats.go        # Statistics generation
│   ├── store/              # Data persistence
│   │   ├── store.go        # Storage interface
│   │   ├── jsonstore.go    # JSON storage implementation
│   │   └── paths.go        # Cross-platform path management
│   ├── alerts/             # Alert/notification logic
│   │   └── alerts.go       # Alert scanner and watch mode
│   ├── notifications/      # Desktop notifications
│   │   └── notify.go       # Cross-platform notifications
│   └── server/             # HTTP API server
│       └── server.go       # REST API implementation
├── documentation/          # All documentation files
│   ├── API.md              # HTTP API reference
│   ├── CHANGELOG.md        # Change history
│   ├── CONTRIBUTING.md     # Contribution guidelines
│   ├── EMOJI_REFERENCE.md  # Emoji usage guide
│   └── QUICK_START.md      # Quick reference guide
├── .github/
│   └── workflows/
│       ├── ci.yml          # Continuous integration
│       └── release.yml     # Release automation
├── go.mod
├── Makefile
└── README.md
```

## Data Storage

Tasks are stored in JSON format in platform-specific directories:

- **Linux**: `~/.local/share/godo/tasks.json`
- **macOS**: `~/Library/Application Support/godo/tasks.json`
- **Windows**: `%APPDATA%/godo/tasks.json`

Storage uses atomic writes to prevent data corruption.

## Desktop Notifications

Notifications use OS-specific commands:

- **Linux**: `notify-send` (install with `sudo apt install libnotify-bin`)
- **macOS**: Built-in `osascript`
- **Windows**: PowerShell toast notifications

If system notifications are unavailable, alerts fall back to console output.

## Development

### Prerequisites

- Go 1.25 or higher
- Make

### Available Make Targets

```bash
# Display all available commands
make help

# Development
make build          # Build the binary
make run            # Run the application
make test           # Run tests
make test-coverage  # Run tests with coverage
make dev            # Run with auto-reload (requires entr)
make fmt            # Format code
make vet            # Run go vet
make lint           # Run golangci-lint
make all            # Format, vet, test, and build

# Cross-platform Builds
make build-all      # Build for all platforms
make build-linux    # Build for Linux (amd64, arm64)
make build-darwin   # Build for macOS (amd64, arm64)
make build-windows  # Build for Windows (amd64)
make build-platform # Build for specific platform
make package        # Create compressed archives
make checksums      # Generate checksums

# Other
make clean          # Remove built binaries
make install        # Install binary to GOPATH/bin
make deps           # Download dependencies
```

### Running Tests

```bash
make test
```

### Building for Multiple Platforms

```bash
# Build for all platforms
make build-all

# Build for specific platform
make build-platform GOOS=linux GOARCH=arm64

# Create release packages
make package VERSION=1.0.0
```

## Priority Levels

1. **Priority 1** (Low): Nice to have
2. **Priority 2** (Medium): Should do
3. **Priority 3** (High): Must do now

## Recurring Tasks

When you mark a recurring task as done:

1. The current task is marked complete
2. A new task is automatically created with:
   - Same title, tags, priority, and repeat rule
   - Due date calculated based on repeat rule
   - Fresh creation timestamp

Supported repeat patterns:

- `daily`: Every day
- `weekly`: Every 7 days
- `monthly`: Every month (same day)

## Task Dependencies

Tasks can depend on other tasks using the `-after` flag:

```bash
# Task 3 depends on tasks 1 and 2
godo add -title "Task 3" -after "1,2"
```

Benefits:

- Blocked tasks show up in alerts
- Can't mark a task done until its dependencies are complete
- Statistics show count of blocked tasks
- Helps organize workflow and project phases

## CI/CD

The project includes GitHub Actions workflows:

- **CI Workflow**: Runs on push/PR to test across Linux, macOS, and Windows
- **Release Workflow**: Automatically creates releases with binaries for all platforms when you push a version tag

To create a release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Examples

### Personal Task Management

```bash
# Add daily tasks with descriptions
godo add -title "Morning exercise" -desc "30 min cardio workout" -repeat daily -p 2
godo add -title "Review emails" -desc "Check and respond to important emails" -tags "work" -p 1

# Add project tasks with dependencies
godo add -title "Design API" -desc "Create API specification and design documents" -p 3 -tags "project,backend"
godo add -title "Implement API" -desc "Develop API endpoints and business logic" -after "1" -p 3 -tags "project,backend"
godo add -title "Write tests" -desc "Create unit and integration tests" -after "2" -p 2 -tags "project,backend"
godo add -title "Deploy" -desc "Deploy to production environment" -after "2,3" -p 3 -tags "project,ops"

# View today's tasks
godo list -today -detailed

# View this week's work tasks
godo list -week -tags "work" -sort priority

# View work tasks sorted by priority with details
godo list -tags "work" -sort priority -detailed

# Monitor for upcoming deadlines
godo alerts -watch -ahead 48h
```

### Team Workflow

```bash
# Start API server for team access
godo server -host 0.0.0.0 -port 8080

# Team members can use HTTP API
curl http://team-server:8080/tasks
```

## Contributing

This is a practice project, but suggestions and improvements are welcome!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

This is a practice project. Feel free to use and modify as needed.

## Troubleshooting

### Notifications not working

**Linux**: Install libnotify-bin:

```bash
sudo apt install libnotify-bin
```

**macOS**: Should work out of the box

**Windows**: Should work out of the box with PowerShell

### Data file location

To find where your tasks are stored:

- Linux: `~/.local/share/godo/tasks.json`
- macOS: `~/Library/Application Support/godo/tasks.json`
- Windows: `%APPDATA%/godo/tasks.json`

You can manually edit or backup this file if needed.

### Build errors

Ensure you have Go 1.25 or higher:

```bash
go version
```

Update dependencies:

```bash
make deps
```

## Roadmap

Future enhancements:

- Natural date parsing ("tomorrow", "next week", "in 3 days")
- Task templates/presets
- Subtasks
- Time tracking
- Export to various formats (CSV, PDF, Markdown)
- Web UI
- Mobile sync
- Multiple task lists/projects
- Collaboration features
- Calendar integration
- Email notifications

## Documentation

This README provides an overview and quick start guide. For more detailed information, see:

- **[Quick Start Guide](documentation/QUICK_START.md)** - Quick command reference and common workflows
- **[API Reference](documentation/API.md)** - Complete HTTP REST API documentation
- **[Contributing Guide](documentation/CONTRIBUTING.md)** - How to contribute to the project
- **[Emoji Reference](documentation/EMOJI_REFERENCE.md)** - Emoji meanings and usage guidelines
- **[Changelog](documentation/CHANGELOG.md)** - Version history and changes

## Support

For issues, questions, or feature requests, please open an issue on GitHub.
