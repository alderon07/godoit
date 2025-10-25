# Emoji Reference Guide

This document defines the standard emojis used throughout the godo project for consistency.

## Task List Display Emojis

### Status Indicators

- `[ ]` - Pending task (unchecked)
- `[✓]` - Completed task (checked)

### Priority Indicators

- 🔴 - High priority (p3)
- 🟡 - Medium priority (p2)
- 🟢 - Low priority (p1)

### Task Properties

- 📝 - Description
- 📅 - Due date (normal)
- ⏰ - Due date (soon or overdue)
- 🏷️ - Tags
- 🔄 - Recurring task
- 🔗 - Dependencies
- ⚠️ - Blocked (dependencies not met)
- 🕐 - Created timestamp
- ✅ - Completed timestamp

### Time-based Views

- 📅 - Today's/Week's tasks header

## Feature Documentation Emojis

### Core Features

- ✅ - Completed/Available feature
- 📝 - Tasks and descriptions
- 🏷️ - Tag system
- 🔍 - Search and filtering
- 📊 - Analytics and statistics
- 🔗 - Task dependencies
- 🔄 - Recurring tasks
- 🔔 - Notifications and alerts
- 👀 - Watch mode
- 🌐 - HTTP API server
- 💾 - Data storage
- 🔀 - Cross-platform support
- 🎯 - Goals and priorities
- ⚙️ - Configuration

## Notification Emojis

- 🔔 - Notification bell (used in console fallback notifications)

## Usage Guidelines

1. **Consistency**: Always use the same emoji for the same concept
2. **Context**: Use the appropriate emoji based on context (e.g., 📅 for date displays, ⏰ for urgent dates)
3. **Clarity**: Emojis should enhance readability, not replace text
4. **Accessibility**: Always provide text alongside emojis for screen readers

## Examples

### Task List Output

```
 1. [ ] 🔴 Important meeting
    📝 Discuss Q4 roadmap with team
    ⏰ Due: 2025-10-24 09:00 (OVERDUE!)
    🏷️  #work #meeting
    🔗 Depends on: [3]
    ⚠️  BLOCKED (dependencies not met)
```

### Completed Task

```
 1. [✓] 🟡 Code review
    📝 Review pull requests from team
    🏷️  #work #code-review
    ✅ Completed: 2025-10-23 15:30
```

### Recurring Task

```
 1. [ ] 🟢 Daily standup
    📝 Morning team sync
    🔄 Repeats: daily
    🏷️  #work #meeting
```

## File Usage

- **cmd/todo/commands.go**: All task display emojis
- **[README.md](../README.md)**: Feature emojis in the features section
- **[QUICK_START.md](QUICK_START.md)**: Visual indicators section
- **internal/notifications/notify.go**: Notification bell emoji (🔔)

## Do Not Use

Avoid using these emojis to prevent confusion:

- ❌ (don't use for task status, use `[✓]` or `[ ]`)
- ⭐ (reserved for potential future "favorite" feature)
- 📌 (reserved for potential future "pin" feature)
- 🚀 (use only in marketing/documentation, not in CLI output)
