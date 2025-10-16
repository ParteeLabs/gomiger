# Contributing to Gomiger

Thank you for your interest in contributing to Gomiger! This document provides guidelines and information for contributors.

## 🎯 How to Contribute

### Types of Contributions

We welcome all types of contributions:

- 🐛 **Bug Reports**: Help us identify and fix issues
- 💡 **Feature Requests**: Suggest new features or improvements
- 📝 **Documentation**: Improve docs, add examples, fix typos
- 🔧 **Code Contributions**: Bug fixes, new features, optimizations
- 🧪 **Testing**: Add tests, improve test coverage
- 🔌 **Plugins**: Create new database plugins
- 💬 **Community**: Help others in discussions and issues

### Getting Started

1. **Fork the Repository**

   ```bash
   git clone https://github.com/your-username/gomiger.git
   cd gomiger
   ```

2. **Set Up Development Environment**

   ```bash
   go work use ./core ./mongomiger ./examples
   go mod download
   ```

3. **Install Development Tools**

   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

   # Install testing tools
   go install github.com/stretchr/testify@latest
   ```

4. **Run Tests**

   ```bash
   # Test all modules
   go test github.com/ParteeLabs/gomiger/core
   go test github.com/ParteeLabs/gomiger/mongomiger

   # Run with coverage
   go test -cover ./...
   ```

## 🛠️ Development Guidelines

### Code Style

We follow standard Go conventions:

- **Formatting**: Use `gofmt` and `goimports`
- **Linting**: Pass `golangci-lint` checks
- **Naming**: Follow Go naming conventions
- **Comments**: Document public APIs and complex logic

### Code Quality Standards

- **Test Coverage**: Maintain >80% test coverage
- **Error Handling**: Always handle errors appropriately
- **Documentation**: Document all public functions and types
- **Performance**: Consider performance implications
- **Security**: Follow security best practices

### Project Structure

```
gomiger/
├── core/                 # Core migration engine
│   ├── *.go             # Core functionality
│   ├── *_test.go        # Unit tests
│   └── cmd/             # CLI tools
├── mongomiger/          # MongoDB plugin
│   ├── *.go             # Plugin implementation
│   └── *_test.go        # Plugin tests
├── examples/            # Example projects
├── docs/               # Documentation
└── .github/            # GitHub workflows and templates
```

## 🔄 Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-number-description
```

### 2. Make Changes

- Write clean, readable code
- Add tests for new functionality
- Update documentation as needed
- Follow existing code patterns

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run linting
golangci-lint run

# Test specific functionality
go test -run TestSpecificFunction ./core

# Integration tests (requires MongoDB)
export GOMIGER_URI="mongodb://localhost:27017/test_db"
go test ./mongomiger
```

### 4. Commit Your Changes

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```bash
git commit -m "feat: add PostgreSQL plugin support"
git commit -m "fix: resolve connection timeout issue"
git commit -m "docs: improve getting started guide"
git commit -m "test: add integration tests for MongoDB plugin"
```

### 5. Create Pull Request

- Fill out the PR template completely
- Reference any related issues
- Include tests for new features
- Update documentation if needed

## 🧪 Testing Guidelines

### Unit Tests

- Test all public functions
- Test error conditions
- Use table-driven tests when appropriate
- Mock external dependencies

Example:

```go
func TestMigrationValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   Migration
		wantErr bool
	}{
		{
			name: "valid migration",
			input: Migration{
				Version: "20241015_valid",
				Up:      func(ctx context.Context) error { return nil },
				Down:    func(ctx context.Context) error { return nil },
			},
			wantErr: false,
		},
		// More test cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMigration(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
```

### Integration Tests

- Test database interactions
- Use test containers when possible
- Clean up resources after tests
- Test migration up/down cycles

### Test Coverage

- Maintain >80% coverage for core packages
- Aim for >90% coverage for critical paths
- Use `go test -cover` to check coverage

## 📝 Documentation Guidelines

### Code Documentation

```go
// Package core provides the main migration engine and interfaces.
package core

// Gomiger defines the interface for migration operations.
// Implementations should handle database connections, schema tracking,
// and migration execution.
type Gomiger interface {
    // Up applies migrations up to the specified version.
    // If toVersion is empty, applies all pending migrations.
    Up(ctx context.Context, toVersion string) error

    // Down reverts migrations down to the specified version.
    Down(ctx context.Context, atVersion string) error
}
```

### README and Guides

- Use clear, concise language
- Include working code examples
- Add diagrams when helpful
- Keep documentation up to date

## 🔌 Plugin Development

### Creating a New Database Plugin

1. **Implement the Interface**

   ```go
   type MyDBPlugin struct {
       core.BaseMigrator
       // Your plugin fields
   }

   func (p *MyDBPlugin) Connect(ctx context.Context) error {
       // Implement connection logic
   }

   func (p *MyDBPlugin) GetSchema(ctx context.Context, version string) (*core.Schema, error) {
       // Implement schema retrieval
   }

   // Implement other required methods...
   ```

2. **Add Tests**

   - Unit tests for all methods
   - Integration tests with real database
   - Mock tests for error conditions

3. **Create Documentation**

   - Plugin-specific setup guide
   - Usage examples
   - Configuration options

4. **Update Main Repository**
   - Add plugin to workspace
   - Update CI/CD to test plugin
   - Add plugin to documentation

## 🚀 Release Process

### Version Management

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking API changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes, backward compatible

### Release Checklist

1. **Update Version Numbers**

   - Update version in relevant files
   - Update CHANGELOG.md

2. **Test Release**

   - Run full test suite
   - Test installation process
   - Test examples

3. **Create Release**
   - Tag release: `git tag v1.2.3`
   - Push tag: `git push origin v1.2.3`
   - GitHub Actions handles the rest

## 🤝 Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Help others learn and grow
- Focus on constructive feedback
- Follow our [Code of Conduct](CODE_OF_CONDUCT.md)

### Getting Help

- **Questions**: Use [GitHub Discussions](https://github.com/ParteeLabs/gomiger/discussions)
- **Bugs**: Create an [Issue](https://github.com/ParteeLabs/gomiger/issues)
- **Feature Requests**: Start a Discussion first
- **Security Issues**: Create an [Issue](https://github.com/ParteeLabs/gomiger/issues)

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Pull Requests**: Code contributions and reviews

## 📋 Pull Request Guidelines

### Before Submitting

- [ ] Tests pass locally
- [ ] Code is formatted (gofmt, goimports)
- [ ] Linting passes (golangci-lint)
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated (for significant changes)

### Pull Request Template

```markdown
## Description

Brief description of changes

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing

- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

## Checklist

- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added for new functionality
```

### Review Process

1. **Automated Checks**: CI/CD must pass
2. **Code Review**: At least one maintainer approval
3. **Testing**: Reviewers test changes when needed
4. **Documentation**: Ensure docs are accurate
5. **Merge**: Squash and merge for clean history

## 🏷️ Issue Guidelines

### Bug Reports

Use the bug report template:

```markdown
## Bug Description

Clear description of the bug

## Steps to Reproduce

1. Step one
2. Step two
3. Step three

## Expected Behavior

What should happen

## Actual Behavior

What actually happens

## Environment

- Go version:
- Gomiger version:
- Database:
- Operating System:
```

### Feature Requests

For new feature ideas:

1. **Check the [Roadmap](ROADMAP.md)** - Your idea might already be planned!
2. **Use the appropriate template**:
   - For roadmap items: Use the "🗺️ Roadmap Item" template
   - For immediate features: Use the "💡 Feature Request" template
3. **Start a Discussion** for complex features before opening an issue

```markdown
## Feature Description

Clear description of the feature

## Use Case

Why is this feature needed?

## Proposed Solution

How should this feature work?

## Alternatives Considered

What other approaches were considered?
```

## 🎯 Good First Issues

New contributors should look for issues labeled:

- `good first issue`: Perfect for beginners
- `help wanted`: Community help needed
- `documentation`: Docs improvements
- `testing`: Test coverage improvements

## 📚 Learning Resources

- [Go Documentation](https://golang.org/doc/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [MongoDB Go Driver](https://pkg.go.dev/go.mongodb.org/mongo-driver)

## 🙏 Recognition

Contributors are recognized in:

- [CONTRIBUTORS.md](CONTRIBUTORS.md)
- Release notes
- Project README
- Special recognition for major contributions

---

Thank you for contributing to Gomiger! 🚀

For questions about contributing, please open a discussion on GitHub or contact the maintainers.
