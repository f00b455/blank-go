# GitHub Workflows Documentation

## Claude-Powered Workflows

This repository uses Claude Code Action to provide intelligent automation for various development tasks.

### 🤖 Interactive Claude (`claude.yml`)
- **Trigger:** Mention `@claude` in any issue, PR comment, or review
- **Purpose:** General assistance, code reviews, explanations, fixes
- **Example:** `@claude can you explain this function?`

### 🔧 PR Auto Fix (`pr-auto-fix.yml`)
- **Trigger:** Comment `@claude fix` or `/fix` on a PR
- **Purpose:** Automatically fix code issues mentioned in comments
- **Example:** `@claude fix the linting errors`

### 🚨 CI Failure Detector (`ci-failure-detector.yml`)
- **Trigger:** Automatically when CI fails on a PR
- **Purpose:** Comments on PR with failure details and fix options
- **Action:** Creates helpful comment with quick fix commands

### 🏷️ Auto Triage (`auto-triage.yml`)
- **Trigger:** Automatically when new issues are opened
- **Purpose:** Labels issues by type, priority, and component
- **Action:** Adds labels and welcoming comment

### 🛠️ Simple Auto Fix (`simple-fix.yml`)
- **Trigger:** Comment `/fix-ci` on a PR
- **Purpose:** Analyzes and fixes all CI failures
- **Example:** `/fix-ci` when tests or linting fails

### 📝 Claude Code Review (`claude-code-review.yml`)
- **Trigger:** Manual workflow dispatch
- **Purpose:** Comprehensive PR code review
- **Usage:** Run from Actions tab with PR number

## Quick Commands

| Command | Action | Where to Use |
|---------|--------|--------------|
| `@claude` | Get help or explanations | Any issue/PR comment |
| `@claude fix` | Fix code issues | PR comments |
| `/fix` | Alternative fix command | PR comments |
| `/fix-ci` | Fix all CI failures | PR comments |
| `@claude review` | Request code review | PR comments |

## CI/CD Workflows

### Go CI/CD (`go.yml`)
- Runs on every push and PR
- Tests, lints, and builds the project
- Checks code coverage (85% threshold)

### Auto Branch PR (`auto-branch-pr.yml`)
- Creates PRs automatically for issue branches
- Links PRs to corresponding issues

### Auto Rebase (`auto-rebase.yml`)
- Keeps PRs up to date with main branch
- Runs when main branch is updated

## Tips for Using Claude

1. **Be specific:** The more detail you provide, the better Claude can help
2. **Use examples:** Show Claude what you want with examples
3. **Ask for explanations:** Claude can explain complex code
4. **Request tests:** Ask Claude to write tests for new features
5. **Fix iteratively:** You can have conversations with Claude in comments

## Common Use Cases

### Fixing Test Failures
```
/fix-ci
```
or
```
@claude the tests are failing, can you fix them?
```

### Code Review
```
@claude can you review this PR for security issues?
```

### Adding Features
```
@claude can you add error handling to this function?
```

### Understanding Code
```
@claude what does this function do and how can I improve it?
```

## Security Notes

- Claude has restricted tool access for safety
- Sensitive operations require explicit approval
- All changes are tracked in git history
- Claude cannot access secrets or credentials