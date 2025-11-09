# Pull Request

## Description
<!-- Provide a detailed description of your changes -->

## Related Issue
<!-- Link to the related issue(s) -->
Fixes #
Relates to #

## Type of Change
<!-- Mark the relevant option with an "x" -->
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Refactoring (no functional changes)
- [ ] Performance improvement
- [ ] Test addition/update

## Changes Made
<!-- Provide a bullet-point list of specific changes -->
- Change 1
- Change 2
- Change 3

## How to Test
<!-- Provide step-by-step instructions for testing this PR -->
1. Check out this branch: `git checkout [branch-name]`
2. Run the application: `[command to run]`
3. Navigate to `[URL or path]`
4. Perform action: `[specific action]`
5. **Expected Result**: [What should happen]

## Test-Driven Development
<!-- Confirm you followed TDD workflow -->
- [ ] Tests were written BEFORE implementation (RED)
- [ ] Implementation makes tests pass (GREEN)
- [ ] Code was refactored while keeping tests green (REFACTOR)
- [ ] Test coverage meets or exceeds targets (≥80% backend, ≥70% frontend)

## Testing Checklist
<!-- Mark completed items with an "x" -->
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated (if applicable)
- [ ] E2E tests added/updated (if applicable)
- [ ] All tests pass locally (`go test ./...` or `npm test`)
- [ ] Test coverage report reviewed

## Code Quality
<!-- Mark completed items with an "x" -->
- [ ] Code follows project conventions (see DEVELOPER_ONBOARDING.md)
- [ ] Code has been linted (`golangci-lint run` or `npm run lint`)
- [ ] No new warnings or errors introduced
- [ ] Comments added for complex logic
- [ ] Function/method documentation updated

## Security
<!-- Mark completed items with an "x" -->
- [ ] No sensitive data (passwords, API keys, tokens) committed
- [ ] SQL queries use parameterized statements (if applicable)
- [ ] Input validation implemented (if applicable)
- [ ] Authentication/authorization checks in place (if applicable)
- [ ] OWASP Top 10 considerations reviewed (if applicable)

## Documentation
<!-- Mark completed items with an "x" -->
- [ ] README.md updated (if needed)
- [ ] API documentation updated (if API changes)
- [ ] OpenAPI spec updated (if API changes)
- [ ] Code comments added/updated
- [ ] CHANGELOG.md updated (if exists)

## Database Changes
<!-- Mark if applicable -->
- [ ] Database migration included
- [ ] Migration tested (up and down)
- [ ] Migration documented in DATABASE_SCHEMA.md
- [ ] N/A - No database changes

## Breaking Changes
<!-- If this PR introduces breaking changes, describe them here -->
- None
<!-- OR -->
<!-- - Breaking change 1: Description and migration path -->
<!-- - Breaking change 2: Description and migration path -->

## Screenshots/Videos
<!-- If applicable, add screenshots or videos demonstrating the changes -->

## Performance Impact
<!-- Describe any performance implications of this PR -->
- [ ] No performance impact
- [ ] Performance improved
- [ ] Performance impact acceptable (explain below)
<!-- Explanation: -->

## Deployment Notes
<!-- Any special deployment considerations? -->
- [ ] No special deployment steps required
- [ ] Requires environment variable changes (document below)
- [ ] Requires database migration
- [ ] Requires data backfill/migration
- [ ] Other (describe below)

## Rollback Plan
<!-- How can this change be rolled back if needed? -->
- [ ] Standard git revert
- [ ] Database migration rollback required
- [ ] Other (describe below)

## Checklist
<!-- Final checklist before requesting review -->
- [ ] My code follows the project's TDD workflow
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## Additional Notes
<!-- Any additional information that reviewers should know -->
