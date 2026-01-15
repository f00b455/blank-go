# 🚨 Manual Workflow Update Required

## Quick Fix for Issue #7

To complete the test coverage enforcement implementation, make this single change:

### File: `.github/workflows/go.yml`

**Line 67** - Change from:
```yaml
./scripts/check-coverage.sh 0.0 coverage.out
```

**To**:
```yaml
./scripts/check-coverage.sh 80.0 coverage.out
```

### Complete Command

```bash
# Update the workflow file
sed -i 's/check-coverage.sh 0.0/check-coverage.sh 80.0/g' .github/workflows/go.yml

# Commit and push
git add .github/workflows/go.yml
git commit -m "feat: Enforce 80% test coverage in CI pipeline

Closes #7"
git push
```

### Why This Change?

- **What**: Changes CI coverage threshold from 0% to 80%
- **Why**: Enforces minimum test coverage in CI pipeline
- **Impact**: CI will fail if coverage drops below 80%
- **Local**: Already enforced locally via `make validate` and pre-push hooks

### Full Documentation

See [COVERAGE_ENFORCEMENT.md](./COVERAGE_ENFORCEMENT.md) for complete details.

### Important Note

**Current coverage is 31.7%**. After making this change:
- CI will fail until coverage reaches 80%
- You may want to improve test coverage first
- Or temporarily set a lower threshold (e.g., 35%) and gradually increase it

### Alternative: Gradual Rollout

If you prefer a gradual approach:

```bash
# Start with current coverage level
sed -i 's/check-coverage.sh 0.0/check-coverage.sh 32.0/g' .github/workflows/go.yml

# Then gradually increase as you add tests:
# 32% → 40% → 50% → 60% → 70% → 80%
```
