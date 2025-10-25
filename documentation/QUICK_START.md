# godoit Quick Start Guide

## Quick Command Reference

### Adding Tasks

```bash
# Simple task
godoit add -title "Buy groceries"

# Task with description
godoit add -title "Write report" -desc "Q4 performance analysis report"

# Full featured task
godoit add -title "Project review" \
  -desc "Review project deliverables and timeline" \
  -due 2025-10-31 \
  -p 3 \
  -tags "work,urgent" \
  -repeat weekly
```

### Viewing Tasks

```bash
# Quick views
godoit list                    # All pending tasks
godoit list -today             # Just today 📅
godoit list -week              # This week 📅
godoit list -detailed          # Show descriptions & timestamps

# Filtered views
godoit list -tags "work"       # Work tasks only
godoit list -grep "meeting"    # Search for "meeting"
godoit list -sort priority     # Sort by priority

# Combined filters
godoit list -week -detailed -tags "work" -sort priority
```

### Managing Tasks

```bash
# Mark complete
godoit done 1

# Edit task
godoit edit 2 -desc "Updated description" -p 3

# Remove task
godoit rm 3
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
godoit list -today -detailed

# Add a new urgent task
godoit add -title "Handle client issue" -desc "Client reported login problems" -p 3 -tags "urgent,support"

# Mark yesterday's task done
godoit done 5
```

#### Weekly Planning

```bash
# Review this week
godoit list -week -detailed -sort priority

# Add tasks for next week
godoit add -title "Team 1:1s" -desc "Schedule and conduct team one-on-ones" -p 2 -tags "management" -due 2025-10-28
```

#### Project Management

```bash
# Add dependent tasks
godoit add -title "Design mockups" -p 3 -tags "project,design"
godoit add -title "Implement features" -after "1" -p 3 -tags "project,dev"
godoit add -title "QA testing" -after "2" -p 2 -tags "project,qa"

# View project tasks
godoit list -tags "project" -detailed
```

## Key Features at a Glance

| Feature           | Flag        | Example                                      |
| ----------------- | ----------- | -------------------------------------------- |
| Add description   | `-desc`     | `godoit add -title "Task" -desc "Details"`     |
| Today's tasks     | `-today`    | `godoit list -today`                           |
| This week's tasks | `-week`     | `godoit list -week`                            |
| Detailed view     | `-detailed` | `godoit list -detailed`                        |
| Priority (1-3)    | `-p`        | `godoit add -title "Task" -p 3`                |
| Tags              | `-tags`     | `godoit add -title "Task" -tags "work,urgent"` |
| Due date          | `-due`      | `godoit add -title "Task" -due 2025-10-31`     |
| Recurring         | `-repeat`   | `godoit add -title "Task" -repeat daily`       |
| Dependencies      | `-after`    | `godoit add -title "Task" -after "1,2"`        |
| Search            | `-grep`     | `godoit list -grep "meeting"`                  |
| Sort              | `-sort`     | `godoit list -sort priority`                   |

## Tips & Tricks

1. **Combine flags for power**: `godoit list -week -detailed -sort priority -tags "work"`

2. **Use descriptions for context**: They show up in detailed view and help you remember what each task is about

3. **Color-coded priorities**: Red (high), yellow (medium), green (low) make it easy to prioritize at a glance

4. **Time indicators**: Tasks show "OVERDUDE!" or "(soon)" to help you stay on track

5. **View creation dates**: Use `-detailed` to see when tasks were created

6. **Quick daily check-in**:
   ```bash
   godoit list -today -detailed    # What's due today?
   godoit alerts                   # Any urgent items?
   godoit stats                    # How am I doing?
   ```

## Common Patterns

### GTD (Getting Things Done) Style

```bash
# Inbox processing
godoit add -title "Review emails" -tags "inbox" -p 2
godoit add -title "Process meeting notes" -tags "inbox" -p 2

# Context-based views
godoit list -tags "work"
godoit list -tags "personal"
godoit list -tags "calls"
```

### Pomodoro Technique

```bash
# Add tasks with time estimates in description
godoit add -title "Write documentation" -desc "Est. 2 pomodoros (50 min)" -p 2

# View today's tasks
godoit list -today -detailed
```

### Agile Sprint Planning

```bash
# Sprint tasks
godoit add -title "Implement user auth" -desc "Sprint 12, Story #45" -p 3 -tags "sprint,backend"
godoit add -title "Add unit tests" -after "1" -p 2 -tags "sprint,testing"

# View sprint tasks
godoit list -tags "sprint" -sort priority -detailed
```

## Need More Help?

- Full documentation: [README.md](../README.md)
- API reference: [API.md](API.md)
- Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)
- Emoji reference: [EMOJI_REFERENCE.md](EMOJI_REFERENCE.md)
- Server mode: `godoit server -h`
- Command help: `godoit <command> -h`

Happy task managing! 🚀
