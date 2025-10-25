# Documentation

This folder contains all documentation for the godoit project.

Key architectural highlights:

- Storage uses a JSON file with cross-process file locking to prevent concurrent write conflicts.
- A `TaskService` orchestrates all task operations; both CLI and HTTP server use it.
- A `Clock` abstraction enables deterministic time in tests.
- Domain helpers normalize priority and repeat rules.

## Available Documentation

- **[QUICK_START.md](QUICK_START.md)** - Quick command reference and common workflows  
  Start here for a rapid introduction to godoit commands and usage patterns.

- **[API.md](API.md)** - Complete HTTP REST API documentation  
  Full reference for the HTTP server endpoints with request/response examples.

- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Contribution guidelines  
  Learn how to contribute to the project, including code style and development workflow.

- **[EMOJI_REFERENCE.md](EMOJI_REFERENCE.md)** - Emoji usage guide  
  Reference for all emojis used in the CLI output and their meanings.

- **[CHANGELOG.md](CHANGELOG.md)** - Version history and changes  
  Track all notable changes to the project across versions.

## Main Documentation

For the main project overview and getting started guide, see the [main README](../README.md) in the root directory.

## Quick Links

### For Users

- [Quick Start Guide](QUICK_START.md) - Get started quickly
- [Visual Indicators](EMOJI_REFERENCE.md) - Understand the emoji meanings
- Main [README](../README.md) - Project overview

### For API Developers

- [API Reference](API.md) - Complete API documentation
- [Main README](../README.md#http-api-reference) - API overview

### For Contributors

- [Contributing Guide](CONTRIBUTING.md) - How to contribute
- [Changelog](CHANGELOG.md) - Track changes

## Documentation Structure

```
documentation/
├── README.md           # This file
├── QUICK_START.md      # Quick reference guide
├── API.md              # HTTP API documentation
├── CONTRIBUTING.md     # Contribution guidelines
├── EMOJI_REFERENCE.md  # Emoji meanings
└── CHANGELOG.md        # Version history
```

## Keeping Documentation Updated

When adding features or making changes:

1. Update relevant documentation files
2. Add entries to CHANGELOG.md
3. Update examples if command syntax changes
4. Keep emoji usage consistent with EMOJI_REFERENCE.md
5. Reflect architectural changes (services, repositories, locking, clock) when applicable
