# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Task descriptions with `-desc` flag for adding context to tasks
- Today view with `-today` flag to show only today's tasks
- Week view with `-week` flag to show this week's tasks
- Detailed view with `-detailed` flag showing all task metadata
- Enhanced visual display with consistent emoji indicators:
  - Status: `[ ]` pending, `[✓]` completed
  - Priority: 🔴 high, 🟡 medium, 🟢 low
  - Task info: 📝 description, 📅 due date, ⏰ urgent, 🏷️ tags, 🔄 recurring, 🔗 dependencies, ⚠️ blocked, 🕐 created, ✅ completed
- Smart due date indicators ("OVERDUE!" and "soon")
- Visual header for today/week views with 📅 emoji
- [EMOJI_REFERENCE.md](EMOJI_REFERENCE.md) documenting all emoji usage
- [QUICK_START.md](QUICK_START.md) for quick command reference
- [CHANGELOG.md](CHANGELOG.md) for tracking changes

### Changed

- List view now shows richer information with color-coded priorities
- Task display format improved with better visual hierarchy
- Documentation updated with comprehensive examples
- All emojis standardized across codebase and documentation

### Fixed

- Linter error with non-constant format string
- Import consistency across modules

## [1.0.0] - Initial Release

### Added

- Core task management (add, list, done, remove, edit)
- Priority levels (1-3)
- Tags for organization
- Task dependencies
- Recurring tasks (daily, weekly, monthly)
- Cross-platform desktop notifications
- Alert system with watch mode
- Analytics and statistics
- HTTP REST API server
- Cross-platform support (Linux, macOS, Windows)
- JSON-based storage with atomic writes
- Advanced filtering and sorting
- Search functionality
- GitHub Actions CI/CD pipeline
- Comprehensive documentation
