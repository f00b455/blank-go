# Test Coverage Enforcement - Implementation Guide

## Overview

This document describes the test coverage enforcement implementation for the blank-go project, as specified in [Issue #7](https://github.com/f00b455/blank-go/issues/7).

## Current Implementation Status

### ✅ Completed Components

#### 1. Coverage Check Script (`scripts/check-coverage.sh`)
- **Status**: ✅ Implemented and working
- **Location**: `/scripts/check-coverage.sh`
- **Functionality**:
  - Takes threshold and coverage file as parameters
  - Parses total coverage from Go coverage output
  - Exits with error code 1 if coverage is below threshold
  - Provides clear success/failure messages

#### 2. Makefile Targets
- **Status**: ✅ Implemented with 80% threshold

**`make test-coverage-check`** (Line 144-148):
```makefile
test-coverage-check: deps
	@echo "Running tests with coverage validation..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic $(SRC_DIR)
	@chmod +x scripts/check-coverage.sh
	@./scripts/check-coverage.sh 80.0 coverage.out
```

**`make validate`** (Line 121):
```makefile
validate: lint test-coverage-check test-bdd build
	@echo "✅ All validation checks passed!"
	@echo "🎉 Safe to push to remote repository"
```

This means `make validate` automatically enforces 80% coverage.

#### 3. Pre-Push Hook Integration
- **Status**: ✅ Active via Husky
- **Hook**: `.husky/pre-push` runs `make validate`
- **Enforcement**: Local coverage check before every push

### ⚠️ Manual Action Required

#### GitHub Actions Workflow Update

**File**: `.github/workflows/go.yml`
**Line**: 67
**Current**: `./scripts/check-coverage.sh 0.0 coverage.out`
**Required**: `./scripts/check-coverage.sh 80.0 coverage.out`

**Why manual update needed**: GitHub App permissions prevent automated workflow modifications.

**How to update**:

1. Open `.github/workflows/go.yml` in your editor
2. Find line 67 in the `test` job:
```yaml
- name: Check test coverage threshold
  run: |
    chmod +x scripts/check-coverage.sh
    ./scripts/check-coverage.sh 0.0 coverage.out  # ← Change this line
```

3. Change `0.0` to `80.0`:
```yaml
- name: Check test coverage threshold
  run: |
    chmod +x scripts/check-coverage.sh
    ./scripts/check-coverage.sh 80.0 coverage.out  # ← Updated
```

4. Commit and push the change:
```bash
git add .github/workflows/go.yml
git commit -m "feat: Enforce 80% test coverage in CI pipeline

- Update coverage threshold from 0.0% to 80.0%
- Aligns CI with local validation (make validate)
- Implements requirement from issue #7

Closes #7"
git push
```

## Current Test Coverage

**Total Coverage**: 31.7%
**Target**: 80.0%
**Gap**: 48.3%

### Coverage Breakdown (by package)
- `cmd/api`: 0.0%
- `cmd/cli`: 0.0%
- `cmd/web`: 0.0%
- `internal/*`: Various (needs improvement)
- `pkg/*`: Various (needs improvement)
- `features/*`: BDD tests passing (good coverage for tested scenarios)

## Next Steps to Achieve 80% Coverage

1. **Immediate**: Update workflow file threshold to 80.0 (see above)
2. **Short-term**: Improve test coverage for critical packages:
   - Add unit tests for `internal/` packages
   - Add unit tests for `pkg/` packages
   - Focus on business logic and handlers
3. **Medium-term**: Add integration tests for CLI and API entry points
4. **Continuous**: Monitor coverage in CI and address gaps before merging

## Validation Commands

### Local Development
```bash
# Quick validation (lint + test)
make validate-quick

# Full validation (lint + test-coverage-check + test-bdd + build)
make validate

# View coverage report
make test-cover
open coverage.html  # Opens HTML coverage report
```

### CI/CD Pipeline
After updating the workflow file, the CI will:
1. Run linting (golangci-lint)
2. Run tests with coverage (go test -coverprofile=coverage.out)
3. Check coverage threshold (80% minimum)
4. Upload coverage to Codecov
5. Build binaries

**Pipeline will fail** if coverage drops below 80%.

## Benefits

1. **Quality Assurance**: Ensures code is adequately tested before merging
2. **Prevents Regressions**: Coverage threshold blocks untested code
3. **Visible Metrics**: Coverage reports available in CI and Codecov
4. **Developer Feedback**: Local validation catches issues before push
5. **Automated Enforcement**: No manual coverage checks needed

## Troubleshooting

### Coverage Check Fails Locally
```bash
# Run tests and see coverage breakdown
make test-cover

# View which files have low coverage
go tool cover -func=coverage.out | grep -v 100.0%
```

### Coverage Check Fails in CI
1. Check the CI logs for the coverage percentage
2. Run `make test-cover` locally to identify gaps
3. Add tests for uncovered code
4. Verify coverage locally with `make validate`
5. Push changes

### Script Permission Issues
```bash
# Ensure script is executable
chmod +x scripts/check-coverage.sh

# Test script manually
./scripts/check-coverage.sh 80.0 coverage.out
```

## References

- [Issue #7: Enforce 80% minimum test coverage](https://github.com/f00b455/blank-go/issues/7)
- [CLAUDE.md](../CLAUDE.md) - Repository development guidelines
- [Go Coverage Tool Documentation](https://go.dev/blog/cover)
