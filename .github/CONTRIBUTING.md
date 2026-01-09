# Contributing to SaaS Photo Listing Platform

<!-- DDD specific -->

## Architecture Principles
- Domain layer is independent
- No infrastructure dependencies in domain
- Aggregates enforce invariants
- Application layer orchestrates use cases only

<!-- General DDD / Normal Worflow -->

## Branch naming
- `main` → stable, production-ready
- `dev` → integration branch
- `feature/*` → normal changes
- `bugfix/*` → fix bugs
- `feature/ddd-*` → ddd changes changes


## Workflow
1. Create an issue first
2. Branch from `dev`
3. Commit small, focused changes
4. Open PR early
5. Review before merge
6. Squash & merge for clean history

## Commit Messages
Conventional commits:

- feat(domain):
- feat(app):
- fix:
- refactor(domain):
- chore:

Example:
    feat(domain): add tenant aggregate with invariants


## Code Review
- Assign reviewers via CODEOWNERS or manually
- Respond to comments
- Self-review if needed

## Testing
- Run `go test ./...` locally before PR
- Run `golangci-lint run ./...` for linting

