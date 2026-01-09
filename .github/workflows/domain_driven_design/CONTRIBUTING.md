# Contributing Guidelines

## Architecture Principles
- Domain layer is independent
- No infrastructure dependencies in domain
- Aggregates enforce invariants
- Application layer orchestrates use cases only

## Branching Strategy
- `main` → stable, production-ready
- `dev` → integration branch
- `feature/*` → all changes

## Workflow
1. Create an issue first
2. Branch from `dev` to `feature/*`
3. Commit small, focused changes
4. Open PR early
5. Review before merge
6. Squash & merge

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
Every PR must be reviewed, even for solo work.

