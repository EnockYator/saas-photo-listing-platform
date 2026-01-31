
# Contributing to SaaS Photo Listing 

Thank you for your interest in contributing to Photo Listing SaaS! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md) to foster an open and welcoming environment.

## Getting Started

### Prerequisites
- Go 1.21 or later
- Docker and Docker Compose
- Git
- Make (optional, but recommended)

### First-Time Setup

1. **Fork the Repository**
   
   ```bash

   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/photo-listing-saas.git
   cd photo-listing-saas
   
   # Add upstream remote
   git remote add upstream https://github.com/yourusername/photo-listing-saas.git

       Set Up Development Environment
    bash

    # Copy environment variables
    cp .env.example .env

    # Start dependencies
    docker-compose up -d postgres redis nats minio

    # Run migrations
    go run cmd/migrate/main.go up

    # Seed development data (optional)
    go run cmd/seed/main.go

    # Start the development server
    go run cmd/server/main.go

    Verify Setup
    bash

    # Run tests to verify everything works
    go test ./... -v

    # Check API is running
    curl http://localhost:8080/health

    ```

    **Project Structure**

    ```text

    saas-photo-listing-platform/
    ├── backend/
    │   ├── cmd/                                    # Application entry points
    │   │   ├── api/                                # Main API server
    │   │   ├── migrate/                            # Database migrations
    │   │   ├── tools/
    │   │   └── worker/                             # Background job processor
    │   │   └── seed/                               # Development data seeding
    │   └── internal/                               # Private application code
    │       ├── auth/
    │       ├── config/
    │       ├── docs/
    │       ├── domains/                            # Aggregates / entities
    │       │   ├── audit/
    │       │   │   ├── application/                # Use cases and services
    │       │   │   ├── domain/             
    │       │   │   └── infrastructure/
    │       │   │       └── repository/
    │       │   │           └── sqlc/
    │       │   │               ├── audit_logs/
    │       │   │               └── usage_stats/
    │       │   ├── auth/
    │       │   │   ├── application/                # Use cases and services
    │       │   │   ├── domain/
    │       │   │   └── infrastructure/
    │       │   │       └── repository/
    │       │   │           └── sqlc/
    │       │   ├── notification/
    │       │   │   └── email/
    │       │   │       ├── application/
    │       │   │       ├── domain/
    │       │   │       └── infrastructure/
    │       │   │           ├── provider/
    │       │   │           └── repository/
    │       │   │               └── sqlc/
    │       │   │                   └── notifications/
    │       │   ├── payment/
    │       │   │   ├── application/
    │       │   │   ├── domain/
    │       │   │   └── infrastructure/
    │       │   │       ├── repository/
    │       │   │       │   └── sqlc/
    │       │   │       │       ├── invoices/
    │       │   │       │       ├── payments/
    │       │   │       │       └── refunds/
    │       │   │       └── stripe/
    │       │   ├── sharing/
    │       │   │   ├── application/
    │       │   │   ├── domain/
    │       │   │   └── infrastructure/
    │       │   │       └── repository/
    │       │   │           └── sqlc/
    │       │   │               └── share_links/
    │       │   ├── subscription/
    │       │   │   ├── application/
    │       │   │   ├── domain/
    │       │   │   └── infrastructure/
    │       │   │       └── repository/
    │       │   │           └── sqlc/
    │       │   │               ├── plan_limits/
    │       │   │               ├── plans/
    │       │   │               └── subscriptions/
    │       │   └── tenant/
    │       │       ├── application/
    │       │       ├── domain/
    │       │       │   └── entity/
    │       │       └── infrastructure/
    │       │           └── repository/
    │       │               └── sqlc/
    │       │                   ├── files/
    │       │                   ├── listing_photos/
    │       │                   ├── listings/
    │       │                   ├── tenants/
    │       │                   ├── tenant_settings/
    │       │                   ├── tenant_storage_usage/
    │       │                   └── tenant_users/
    │       ├── infrastructure/                             # External implementations
    │       │   ├── database/
    │       │   │   └── postgres/
    │       │   │       ├── migrations/
    │       │   │       └── queries/
    │       │   │           ├── audit/
    │       │   │           ├── auth/
    │       │   │           ├── notification/
    │       │   │           ├── payment/
    │       │   │           ├── sharing/
    │       │   │           ├── subscription/
    │       │   │           └── tenant/
    │       │   ├── nats/
    │       │   ├── observability/
    │       │   ├── redis/
    │       │   └── storage/
    │       │       └── s3_storage/
    │       ├── interfaces/
    │       │   ├── http/
    │       │   │   ├── dto/
    │       │   │   └── handlers/                   # HTTP handlers
    │       │   └── messaging/
    │       │       └── dto/
    │       ├── logger/
    │       ├── middleware/
    │       └── util/
    ├── docker/                                     # Docker configurations
    ├── docs/
    │   ├── architecture/
    │   └── images/
    ├── pkg/                                        # Public library code
    ├── scripts/                                    # Build and deployment scripts
    ├── tests/                                      # Test suites
    ├── nginx/
    └── react-frontend/

    ```

2. **Testing**

    **Test Structure**
        - Unit Tests: In *_test.go files alongside the code
        - Integration Tests: In tests/integration/ with test containers
        - E2E Tests: In tests/e2e/ for full workflow testing

    **Running Tests**

    ```bash

    # Run all tests
    make test

    # Run unit tests only
    make test-unit

    # Run integration tests
    make test-integration

    # Run E2E tests
    make test-e2e

    # Run tests with coverage
    make test-coverage

    # Run tests with race detector
    make test-race

    ```

    **Writing Tests**

    ```go

    // Example test structure
    func TestAlbumService_Create(t *testing.T) {
        t.Run("creates album successfully", func(t *testing.T) {
            // Arrange
            repo := NewMockAlbumRepository()
            service := NewAlbumService(repo)
            req := CreateAlbumRequest{
                Title: "Test Album",
                Description: "Test Description",
            }
            
            // Act
            album, err := service.Create(context.Background(), "tenant-123", req)
            
            // Assert
            assert.NoError(t, err)
            assert.Equal(t, "Test Album", album.Title)
            assert.Equal(t, "draft", album.Status)
        })
        
        t.Run("fails with empty title", func(t *testing.T) {
            // Test validation error
        })
    }

    ```

## Development Workflow
1. **Choose or Create an Issue**
    - Check GitHub Issues
    - Look for good first issue or help wanted labels
    - Create a new issue if you have a bug report or feature request

2. **Sync with upstream**

    ```bash
    git fetch upstream
    git checkout main
    git merge upstream/main
    ```

3. **Create feature branch**
    ```bash
        git checkout -b feat/your-feature-name
    ```
    **or**
    ```bash
        git checkout -b fix/issue-number-description
    ```

    **Branch Naming Convention**:

        feat/ - New features

        fix/ - Bug fixes

        docs/ - Documentation updates

        refactor/ - Code refactoring

        test/ - Test improvements

        chore/ - Maintenance tasks

4. **Make Your Changes**
    **Code Style Guidelines**
    
    ```go

    // Use gofmt and goimports
    go fmt ./...
    goimports -w .

    // Follow Go naming conventions
    // - Use camelCase for variables and functions
    // - Use PascalCase for exported identifiers
    // - Use short, descriptive names

    // Error handling
    if err != nil {
        return fmt.Errorf("context: %w", err)
    }

    // Documentation
    // Add comments for exported functions, types, and packages
    // Use complete sentences ending with periods

    ```

    **Commit Message Convention**
    ```text

    type(scope): subject

    body

    footer

    Types:

        feat: New feature

        fix: Bug fix

        docs: Documentation

        style: Formatting changes

        refactor: Code refactoring

        test: Adding tests

        chore: Maintenance tasks
    ```

    Example:
    ```text

    feat(album): add support for album sharing

    - Add share link generation endpoint
    - Implement token-based access control
    - Add expiration dates for share links

    Closes #123
    ```

5. **Run Tests and Linters**
    
    ```bash

    # Run linters
    make lint

    # Run all tests
    make test

    # Check code coverage
    make test-coverage

    # Build to verify compilation
    make build

    ```

6. **Update Documentation**
    - Update relevant documentation in /docs/
    - Update inline code comments if needed
    - Update API documentation if endpoints changed

7. **Create Pull Request**

    **Push your branch**:
    
    ```bash
    git push origin feat/your-feature-name
    ```

    **Create a PR on GitHub with**:
    - Clear title following commit convention
    - Description of changes and motivation
    - Reference to related issue (#123)
    - Screenshots for UI changes
    - Test coverage information

    **Fill out the PR template**:
        
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
    - [ ] E2E tests added/updated

    ## Checklist
    - [ ] Code follows style guidelines
    - [ ] Self-review completed
    - [ ] Documentation updated
    - [ ] Tests pass locally
    - [ ] No new warnings

    ```

8. **Address Review Feedback**
    - Respond to all review comments
    - Make requested changes
    - Push updates to your branch
    - Request re-review when ready

## Release Process

### Versioning

We follow Semantic Versioning:

    MAJOR: Breaking changes

    MINOR: New features (backwards compatible)

    PATCH: Bug fixes (backwards compatible)

**Release Checklist**

    All tests pass

    Documentation updated

    Migration scripts tested

    Changelog updated

    Version tags created

    Release notes written

## Documentation

### Writing Documentation
- Use clear, concise language
- Include code examples
- Update when code changes
- Use relative links within the repo

### Documentation Structure

```text

docs/
├── overview.md          # Product overview
├── setup.md            # Development setup
├── deployment.md       # Production deployment
├── security.md         # Security practices
├── api.md             # API reference
└── architecture/       # Architecture docs
    ├── context.md
    ├── events.md
    └── data-model.md

```

## Bug Reports

### Reporting Bugs

    Check if the bug already exists in issues

    Use the bug report template

    Include:

        Steps to reproduce

        Expected vs actual behavior

        Environment details

        Logs or error messages

        Screenshots if applicable

### Bug Report Template

```markdown

## Description
[Clear description of the bug]

## Steps to Reproduce
1. [First Step]
2. [Second Step]
3. [Third Step]

## Expected Behavior
[What should happen]

## Actual Behavior
[What actually happens]

## Environment
- OS: [e.g., macOS 12.0]
- Go Version: [e.g., 1.21.0]
- Docker Version: [e.g., 20.10.0]
- Browser: [if applicable]

## Additional Context
[Logs, screenshots, etc.]

```

## Feature Requests

### Suggesting Features

    Check existing feature requests

    Use the feature request template

    Explain the problem and proposed solution

    Include use cases and benefits

### Feature Request Template

```markdown

## Problem Statement
[Describe the problem you're trying to solve]

## Proposed Solution
[Describe your proposed solution]

## Use Cases
[Who would use this and how?]

## Alternatives Considered
[What other approaches did you consider?]

## Additional Context
[Screenshots, mockups, etc.]

```

## Recognition

Contributors are recognized in:

    GitHub contributors list

    Release notes

    Project documentation

Significant contributions may receive:

    Contributor badge

    Mention in blog posts

    Invitation to join the core team

## ❓ Getting Help

    GitHub Discussions: For questions and discussions

    GitHub Issues: For bug reports and feature requests

    Documentation: Check /docs/ first

    Code Comments: Read inline documentation

##  License

By contributing, you agree that your contributions will be licensed under the MIT License. See LICENSE for details.

Thank you for contributing to SaaS Photo Listing Platform! ❤️