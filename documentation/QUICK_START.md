# godo Quick Start Guide

## Quick Command Reference

### Adding Tasks

```bash
# Simple task
godo add -title "Buy groceries"

# Task with description
godo add -title "Write report" -desc "Q4 performance analysis report"

# Full featured task
godo add -title "Project review" \
  -desc "Review project deliverables and timeline" \
  -due 2025-10-31 \
  -p 3 \
  -tags "work,urgent" \
  -repeat weekly
```

### Viewing Tasks

```bash
# Quick views
godo list                    # All pending tasks
godo list -today             # Just today 📅
godo list -week              # This week 📅
godo list -detailed          # Show descriptions & timestamps

# Filtered views
godo list -tags "work"       # Work tasks only
godo list -grep "meeting"    # Search for "meeting"
godo list -sort priority     # Sort by priority

# Combined filters
godo list -week -detailed -tags "work" -sort priority
```

### Managing Tasks

```bash
# Mark complete
godo done 1

# Edit task
godo edit 2 -desc "Updated description" -p 3

# Remove task
godo rm 3
```

### Visual Indicators

The list view uses emojis to make information easy to scan:

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
- ⏰ - Due soon or overdue
- 🏷️ - Tags
- 🔄 - Recurring task (daily/weekly/monthly)
- 🔗 - Task has dependencies
- ⚠️ - Task is blocked (dependencies not met)
- 🕐 - Created timestamp (shown in `-detailed` view)
- ✅ - Completion timestamp (shown in `-detailed` view)

### Example Workflows

#### Daily Morning Routine

```bash
# Check today's tasks
godo list -today -detailed

# Add a new urgent task
godo add -title "Handle client issue" -desc "Client reported login problems" -p 3 -tags "urgent,support"

# Mark yesterday's task done
godo done 5
```

#### Weekly Planning

```bash
# Review this week
godo list -week -detailed -sort priority

# Add tasks for next week
godo add -title "Team 1:1s" -desc "Schedule and conduct team one-on-ones" -p 2 -tags "management" -due 2025-10-28
```

#### Project Management

```bash
# Add dependent tasks
godo add -title "Design mockups" -p 3 -tags "project,design"
godo add -title "Implement features" -after "1" -p 3 -tags "project,dev"
godo add -title "QA testing" -after "2" -p 2 -tags "project,qa"

# View project tasks
godo list -tags "project" -detailed
```

## Key Features at a Glance

| Feature           | Flag        | Example                                      |
| ----------------- | ----------- | -------------------------------------------- |
| Add description   | `-desc`     | `godo add -title "Task" -desc "Details"`     |
| Today's tasks     | `-today`    | `godo list -today`                           |
| This week's tasks | `-week`     | `godo list -week`                            |
| Detailed view     | `-detailed` | `godo list -detailed`                        |
| Priority (1-3)    | `-p`        | `godo add -title "Task" -p 3`                |
| Tags              | `-tags`     | `godo add -title "Task" -tags "work,urgent"` |
| Due date          | `-due`      | `godo add -title "Task" -due 2025-10-31`     |
| Recurring         | `-repeat`   | `godo add -title "Task" -repeat daily`       |
| Dependencies      | `-after`    | `godo add -title "Task" -after "1,2"`        |
| Search            | `-grep`     | `godo list -grep "meeting"`                  |
| Sort              | `-sort`     | `godo list -sort priority`                   |

## Tips & Tricks

1. **Combine flags for power**: `godo list -week -detailed -sort priority -tags "work"`

2. **Use descriptions for context**: They show up in detailed view and help you remember what each task is about

3. **Color-coded priorities**: Red (high), yellow (medium), green (low) make it easy to prioritize at a glance

4. **Time indicators**: Tasks show "OVERDUE!" or "(soon)" to help you stay on track

5. **View creation dates**: Use `-detailed` to see when tasks were created

6. **Quick daily check-in**:
   ```bash
   godo list -today -detailed    # What's due today?
   godo alerts                   # Any urgent items?
   godo stats                    # How am I doing?
   ```

## Common Patterns

### GTD (Getting Things Done) Style

```bash
# Inbox processing
godo add -title "Review emails" -tags "inbox" -p 2
godo add -title "Process meeting notes" -tags "inbox" -p 2

# Context-based views
godo list -tags "work"
godo list -tags "personal"
godo list -tags "calls"
```

### Pomodoro Technique

```bash
# Add tasks with time estimates in description
godo add -title "Write documentation" -desc "Est. 2 pomodoros (50 min)" -p 2

# View today's tasks
godo list -today -detailed
```

### Agile Sprint Planning

```bash
# Sprint tasks
godo add -title "Implement user auth" -desc "Sprint 12, Story #45" -p 3 -tags "sprint,backend"
godo add -title "Add unit tests" -after "1" -p 2 -tags "sprint,testing"

# View sprint tasks
godo list -tags "sprint" -sort priority -detailed
```

## Need More Help?

- Full documentation: [README.md](../README.md)
- API reference: [API.md](API.md)
- Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)
- Emoji reference: [EMOJI_REFERENCE.md](EMOJI_REFERENCE.md)
- Server mode: `godo server -h`
- Command help: `godo <command> -h`

Happy task managing! 🚀
