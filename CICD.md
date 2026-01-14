# CI/CD Automated Fix Pipeline

## 🎯 Overview

This repository implements a **fully automated CI/CD pipeline** with intelligent auto-fix capabilities powered by Claude Code Action. When CI checks fail, the system automatically analyzes errors, fixes code issues, and re-triggers validation—all without manual intervention.

**Key Feature**: Zero-touch automation - developers push code, CI fails, Claude fixes it automatically.

---

## 🔄 How It Works

### Process Flow

```
Developer pushes to PR
         ↓
   CI runs (go.yml)
         ↓
    Checks fail? ──No──→ ✅ PR ready for review
         ↓ Yes
Auto-fix workflow triggers (auto-fix-on-failure.yml)
         ↓
  Fetch failed check logs
         ↓
Claude analyzes & fixes issues
         ↓
  Auto-commit & push fixes
         ↓
   CI re-runs automatically
         ↓
  Still failing? ──No──→ ✅ PR ready for review
         ↓ Yes
Retry (max 3 attempts)
         ↓
After 3 attempts → 🔔 Request human review
```

### Detailed Steps

1. **Trigger**: `auto-comment-claude.yml` triggers automatically when `go.yml` workflow completes with failures
2. **Detection**: Workflow identifies which checks failed and gets their log URLs
3. **Comment**: Workflow posts `@claude` comment with failure details and fix instructions
4. **Claude Activates**: Claude Code Action automatically responds to @mention
5. **Analysis**: Claude fetches detailed logs using `gh run view` and analyzes errors
6. **Fix**: Claude fixes all identified issues in the codebase
7. **Commit**: Claude commits with descriptive message
8. **Push**: Claude pushes using **OAuth token** (not GITHUB_TOKEN)
9. **Re-trigger**: Push automatically triggers new CI run ✅ (OAuth token enables this!)
10. **Retry**: If checks still fail, workflow posts new `@claude` comment (max 3 times)
11. **Escalation**: After 3 attempts, adds comment requesting human review

---

## 📋 Workflows

### 1. Standard CI/CD Pipeline (`go.yml`)

**Trigger**: Pull request creation/updates, pushes to main/develop

**Jobs**:
- **Lint**: Runs `golangci-lint` to check code quality
- **Test**: Runs unit tests with race detection and coverage analysis (95% threshold)
- **Build**: Builds all binaries (API server, CLI tool, Web server)

**File**: `.github/workflows/go.yml`

### 2. Automatic Fix on Failure (`auto-comment-claude.yml`)

**Trigger**: Automatically when `go.yml` completes with failure status

**Jobs**:
- Get PR number from failed workflow
- Check iteration count (max 3 attempts)
- Get failed check details and URLs
- Post `@claude` comment with:
  - List of failed checks with URLs
  - Specific instructions for each failure type
  - Request to fix and commit/push
- Claude Code Action automatically responds
- Claude fixes, commits, and pushes (using OAuth token)
- Push automatically triggers new CI run ✅

**File**: `.github/workflows/auto-comment-claude.yml`

**Key Advantage**: Claude Code Action's OAuth token (not `GITHUB_TOKEN`) triggers CI automatically!

### 3. Manual Fix Trigger (`simple-fix.yml`)

**Trigger**: Comment `/fix-ci` on PR (manual fallback)

**Jobs**: Same as auto-fix, but triggered manually

**File**: `.github/workflows/simple-fix.yml`

---

## 🔧 What Gets Fixed Automatically

### Lint Issues ✨
- **Unused imports**: `import "net/http"` removed if not used
- **Unused variables**: `cfg := &config.Config{}` removed if not referenced
- **Code formatting**: Automatic `gofmt` and `goimports` fixes
- **Naming conventions**: Fix non-idiomatic Go names
- **Type errors**: Resolve type mismatches and conversions
- **Inefficient patterns**: Replace with more efficient Go idioms

**Example Fix**:
```go
// Before (lint error: unused import)
import (
    "net/http"  // ❌ Not used
    "testing"
)

// After (auto-fixed)
import (
    "testing"
)
```

### Test Issues 🧪
- **Compilation errors**: Fix syntax errors in test files
- **Wrong signatures**: Update function calls to match current signatures
- **Missing imports**: Add required imports for test dependencies
- **Mock setup errors**: Fix mock configuration and initialization
- **Test data issues**: Correct test fixtures and expected values

**Example Fix**:
```go
// Before (error: undefined: http)
func TestHandler(t *testing.T) {
    req := http.NewRequest("GET", "/test", nil)  // ❌ http not imported
}

// After (auto-fixed)
import "net/http"

func TestHandler(t *testing.T) {
    req := http.NewRequest("GET", "/test", nil)  // ✅
}
```

### Coverage Issues 📊
- **Add missing tests**: Generate tests for uncovered functions
- **Improve scenarios**: Add edge cases to increase coverage
- **Test completeness**: Ensure 95% threshold is met

### Build Issues 🔨
- **Dependency sync**: Run `go mod tidy` for module issues
- **Import cycles**: Restructure imports to break cycles
- **Compilation errors**: Fix package-level syntax errors

---

## 🔁 Retry Mechanism

The workflow automatically tracks fix attempts to prevent infinite loops:

| Attempt | Commit Message | Action |
|---------|----------------|--------|
| 1 | `fix: auto-fix CI failures (attempt 1/3)` | Analyze & fix issues |
| 2 | `fix: auto-fix CI failures (attempt 2/3)` | Retry with fresh analysis |
| 3 | `fix: auto-fix CI failures (attempt 3/3)` | Final attempt |
| 4+ | ⚠️ Comment requesting human review | Stop automatic fixes |

**Implementation**: Counts commits matching pattern `"fix: auto-fix CI failures (attempt"` in git history

**Escalation Message** (after 3 attempts):
> ⚠️ **Auto-Fix Limit Reached**
>
> The automatic fix workflow has attempted to fix CI failures 3 times without success.
>
> **Next Steps:**
> 1. Review the auto-fix commits to understand what was attempted
> 2. Manually investigate the remaining issues
> 3. Use `/fix-ci` command if you want to retry with manual oversight

---

## 🎮 Manual Triggers

### Use Manual Trigger When:
- Auto-fix reaches retry limit
- Complex issues require human judgment
- You want to control timing of fixes
- Testing fix workflow changes

### How to Trigger Manually:

**Option 1: Comment Command**
```
/fix-ci
```
Post this comment on any PR to trigger `simple-fix.yml` workflow

**Option 2: GitHub UI**
1. Go to Actions tab
2. Select "Simple Auto Fix" workflow
3. Click "Run workflow"
4. Enter PR number
5. Click "Run workflow"

**Option 3: GitHub CLI**
```bash
gh workflow run simple-fix.yml -f pr_number=44
```

---

## ⚙️ Configuration

### Required Permissions

Both workflows require these permissions:

```yaml
permissions:
  contents: write        # Push fixes to PR branch
  pull-requests: write   # Comment on PR with status
  checks: read          # Read CI check results
  actions: read         # Access workflow run logs
  statuses: read        # Read commit statuses
  id-token: write       # OIDC authentication for Claude
```

### Required Secrets

Configure in repository settings → Secrets and variables → Actions:

| Secret Name | Description | How to Get |
|-------------|-------------|------------|
| `CLAUDE_CODE_OAUTH_TOKEN` | OAuth token for Claude Code Action | [Get from Anthropic Console](https://console.anthropic.com/) |
| `GITHUB_TOKEN` | GitHub Actions token | Automatically provided by GitHub |

### Repository Settings

1. **Allow GitHub Actions to create and approve pull requests**: Settings → Actions → General → Workflow permissions → ✅ Allow GitHub Actions to create and approve pull requests

2. **Branch protection rules**: Configure to allow bot commits:
   - Settings → Branches → Branch protection rules
   - Exclude `github-actions[bot]` from status check requirements for auto-fix commits

---

## 💡 Example Scenarios

### Scenario 1: Unused Import Fix

**Initial State**: Developer adds `import "net/http"` but doesn't use it

```
1. Push to PR → Lint fails: "net/http imported and not used"
2. Auto-fix triggers automatically (30 seconds later)
3. Claude analyzes lint logs
4. Claude removes unused import
5. Auto-commit: "fix: auto-fix CI failures (attempt 1/3)"
6. Push triggers new CI run
7. CI passes ✅ (2-3 minutes total)
```

**Developer Action Required**: None - completely automatic

### Scenario 2: Multiple Issues

**Initial State**: Unused import + unused variable + wrong test signature

```
1. Push to PR → Lint & Test both fail
2. Auto-fix triggers (30 seconds later)
3. Claude analyzes both lint and test logs
4. Claude fixes all issues in one commit:
   - Removes unused import
   - Removes unused variable
   - Updates test function signature
5. Auto-commit with all fixes
6. CI passes ✅
```

**Developer Action Required**: None

### Scenario 3: Complex Issue (Max Attempts)

**Initial State**: Deep logic bug causing test failures

```
1. Push to PR → Test fails
2. Attempt 1: Claude tries fix → Still fails
3. Attempt 2: Claude tries different approach → Still fails
4. Attempt 3: Claude tries another fix → Still fails
5. Bot comments: "⚠️ Auto-Fix Limit Reached"
6. Developer manually investigates and fixes
```

**Developer Action Required**: Manual investigation after 3 attempts

---

## 📊 Monitoring

### Check Workflow Status

**List all auto-fix runs**:
```bash
gh run list --workflow=auto-fix-on-failure.yml
```

**View specific run**:
```bash
gh run view <run-id> --log
```

**Check PR comments for fix attempts**:
```bash
gh pr view <pr-number> --comments
```

### Check Fix History

**See all auto-fix commits in PR**:
```bash
git log --oneline --grep="fix: auto-fix CI failures"
```

**Count attempts**:
```bash
git log --oneline --grep="fix: auto-fix CI failures (attempt" | wc -l
```

---

## 🐛 Troubleshooting

### Auto-fix not triggering

**Symptoms**: CI fails but auto-fix workflow doesn't run

**Possible Causes**:
1. Workflow permissions not configured
2. `CLAUDE_CODE_OAUTH_TOKEN` secret missing
3. `workflow_run` trigger disabled
4. Not a pull request event

**Solutions**:
```bash
# Check workflow status
gh run list --workflow=auto-fix-on-failure.yml

# Verify secret exists (should show name, not value)
gh secret list | grep CLAUDE_CODE_OAUTH_TOKEN

# Trigger manually as fallback
gh pr comment <pr-number> --body "/fix-ci"
```

### Fixes not being applied

**Symptoms**: Workflow runs but no commits appear

**Possible Causes**:
1. Git configuration incorrect
2. Branch protection blocking bot commits
3. Claude found no fixable issues

**Solutions**:
```bash
# Check workflow logs
gh run view <run-id> --log

# Look for git errors in commit step
gh run view <run-id> --log | grep "git commit"

# Check branch protection rules
# Settings → Branches → Review protection rules
```

### Fix loop not stopping

**Symptoms**: More than 3 auto-fix attempts

**Possible Causes**:
1. Iteration count logic broken
2. Commit message pattern not matching

**Solutions**:
```bash
# Check commit messages
git log --oneline | grep "auto-fix"

# Manually stop with comment
gh pr comment <pr-number> --body "Stop auto-fixing - investigating manually"

# Use manual trigger with oversight
gh pr comment <pr-number> --body "/fix-ci"
```

### Claude API errors

**Symptoms**: Workflow fails at Claude Code Action step

**Possible Causes**:
1. Invalid OAuth token
2. Rate limit exceeded
3. API endpoint issues

**Solutions**:
```bash
# Verify token in secrets
gh secret list

# Check Claude Code Action logs
gh run view <run-id> --log | grep "claude-code-action"

# Wait and retry if rate limited
# Manual trigger after cooldown period
```

---

## 🎁 Benefits

| Benefit | Impact |
|---------|--------|
| 🚀 **Zero friction** | Developers push without worrying about trivial errors |
| ⚡ **Fast feedback** | Issues fixed within 2-5 minutes of detection |
| 🎯 **Consistent quality** | Claude applies same standards across all code |
| 📚 **Learning tool** | Developers see fixes in commit history and learn patterns |
| 🚫 **Prevents blocking** | PRs don't get stuck on lint/formatting issues |
| 🧠 **Focus on logic** | Developers focus on business logic, not syntax |
| 🤖 **24/7 availability** | Works around the clock, even on weekends |
| 📉 **Reduced review time** | Reviewers focus on architecture, not style issues |

---

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Claude Code Action Documentation](https://github.com/anthropics/claude-code-action)
- [golangci-lint Configuration](https://golangci-lint.run/)
- [Go Testing Best Practices](https://golang.org/doc/code.html#Testing)

---

## 🔐 Security Considerations

1. **OAuth Token**: Keep `CLAUDE_CODE_OAUTH_TOKEN` secret, never commit to code
2. **Branch Protection**: Configure rules to prevent force pushes
3. **Code Review**: Auto-fixes still require human review before merge
4. **Audit Trail**: All fixes tracked in git history with detailed commit messages
5. **Retry Limits**: Prevents runaway automation with 3-attempt maximum

---

## 📝 Summary

This automated CI/CD pipeline with Claude Code Action provides:
- ✅ **Automatic** fixing of lint, test, and build issues
- ✅ **Zero manual intervention** required for simple issues
- ✅ **Retry logic** with 3-attempt limit to prevent loops
- ✅ **Manual fallback** via `/fix-ci` command
- ✅ **Full transparency** with detailed logs and PR comments

**Result**: Developers focus on building features, not fixing syntax errors.