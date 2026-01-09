# DDD Feature PR Template
# For feature/ddd-* branches

## Summary
<!-- Describe what this PR is for, focusing on domain/business logic -->

## Motivation
<!-- Why is this change needed from a business/domain perspective? -->

## Domain Context
- Which bounded context is affected?
- Which aggregate(s) are involved?

## Design Decisions
- Key architectural or domain decisions
- Invariants enforced or modified
- Trade-offs considered

## Implementation Notes
- New entities / value objects
- Repository or service changes
- Persistence or infrastructure impact

## Related Issues
Closes #

## Checklist / Acceptance Criteria
- [ ] Domain invariants enforced
- [ ] Domain tests pass
- [ ] Business rules validated
- [ ] No domain logic leaked to application/infrastructure
- [ ] Naming aligns with ubiquitous language

## Notes
<!-- Any additional info for reviewers -->