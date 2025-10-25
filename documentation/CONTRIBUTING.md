# Contributing to godo

Thank you for your interest in contributing to godo! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the project and community

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Git
- Make

### Setting Up Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/godo.git
   cd godo
   ```
3. Install dependencies:
   ```bash
   make deps
   ```
4. Build the project:
   ```bash
   make build
   ```
5. Run tests:
   ```bash
   make test
   ```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

Branch naming conventions:

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions or updates

### 2. Make Your Changes

- Write clear, readable code
- Follow Go conventions and idioms
- Add comments for complex logic
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Run linter
make vet

# Run all checks
make all
```

### 4. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add natural date parsing support"
```

Commit message format:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Test additions or updates
- `chore:` - Build process or auxiliary tool changes

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub with:

- Clear description of changes
- Reference to related issues (if any)
- Screenshots (if UI changes)
- Test results

## Code Style

### Go Style Guidelines

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for formatting (run `make fmt`)
- Keep functions small and focused
- Use meaningful variable names
- Add comments for exported functions and types

### Example

```go
// CalculateStats computes statistics from a list of tasks
func CalculateStats(tasks []Task, now time.Time) Stats {
    stats := Stats{
        ByPriority: make(map[int]int),
        ByTag:      make(map[string]int),
    }

    // Process tasks...

    return stats
}
```

## Testing

### Writing Tests

- Place tests in `*_test.go` files
- Use table-driven tests when appropriate
- Test both success and error cases
- Aim for good coverage

Example:

```go
func TestTaskHasTag(t *testing.T) {
    tests := []struct {
        name     string
        task     Task
        tag      string
        expected bool
    }{
        {
            name:     "Has tag",
            task:     Task{Tags: []string{"work"}},
            tag:      "work",
            expected: true,
        },
        {
            name:     "Does not have tag",
            task:     Task{Tags: []string{"work"}},
            tag:      "personal",
            expected: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := tt.task.HasTag(tt.tag)
            if result != tt.expected {
                t.Errorf("Expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

## Documentation

### Update Documentation When:

- Adding new features
- Changing existing functionality
- Adding new CLI commands
- Adding new API endpoints

### Documentation Files:

- [README.md](../README.md) - Main project documentation
- [API.md](API.md) - HTTP API reference
- [QUICK_START.md](QUICK_START.md) - Quick reference guide
- [EMOJI_REFERENCE.md](EMOJI_REFERENCE.md) - Emoji usage guidelines
- [CHANGELOG.md](CHANGELOG.md) - Version history
- `CONTRIBUTING.md` - This file
- Code comments - Inline documentation

## Project Structure

```
godo/
├── cmd/todo/           # Main application
├── internal/           # Internal packages
│   ├── core/          # Core logic
│   ├── store/         # Storage layer
│   ├── alerts/        # Alert system
│   ├── notifications/ # Notifications
│   └── server/        # HTTP server
├── documentation/     # All documentation files
│   ├── API.md
│   ├── CHANGELOG.md
│   ├── CONTRIBUTING.md
│   ├── EMOJI_REFERENCE.md
│   └── QUICK_START.md
├── .github/           # GitHub Actions
└── tests/             # Integration tests (future)
```

## Feature Requests

We welcome feature requests! Please:

1. Check existing issues first
2. Open a new issue with:
   - Clear description
   - Use cases
   - Expected behavior
   - Possible implementation approach

## Bug Reports

When reporting bugs, please include:

1. **Description**: Clear description of the issue
2. **Steps to reproduce**: Detailed steps
3. **Expected behavior**: What should happen
4. **Actual behavior**: What actually happens
5. **Environment**:
   - OS and version
   - Go version
   - godo version
6. **Logs/Error messages**: Any relevant output

Example:

```markdown
## Description

Task completion fails with recurring tasks

## Steps to Reproduce

1. Create a recurring task: `godo add -title "Test" -repeat daily`
2. Mark it as done: `godo done 1`
3. Error occurs

## Expected Behavior

Task should be marked done and new occurrence created

## Actual Behavior

Error: "cannot create recurring task"

## Environment

- OS: Ubuntu 22.04
- Go: 1.25
- godo: v1.0.0
```

## Pull Request Guidelines

### Before Submitting

- [ ] Code follows style guidelines
- [ ] Tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linter passes (`make vet`)
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] Branch is up to date with main

### PR Description Should Include

- Summary of changes
- Motivation and context
- Related issues
- Type of change:
  - [ ] Bug fix
  - [ ] New feature
  - [ ] Breaking change
  - [ ] Documentation update
- Testing details

### Review Process

1. Automated tests will run
2. Maintainers will review
3. Address feedback
4. Once approved, PR will be merged

## Areas for Contribution

### Good First Issues

- Add more tests
- Improve documentation
- Fix typos
- Add examples
- Improve error messages

### Feature Ideas

- Natural date parsing
- Task templates
- Subtasks
- Time tracking
- Export formats
- Calendar integration
- Mobile app

### Infrastructure

- Improve CI/CD
- Add more tests
- Performance optimization
- Code refactoring

## Getting Help

- Open an issue for questions
- Check existing documentation
- Review closed issues and PRs

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

## Recognition

Contributors will be recognized in:

- CONTRIBUTORS file
- Release notes
- Project README

Thank you for contributing to godo! 🎉
