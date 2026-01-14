# Golang Template

A clean, modern Go project template featuring:

- **HTTP API Server** with Gin framework
- **Web Frontend** with RSS headlines and terminal UI
- **CLI Application** with Cobra framework
- **Docker Support** with multi-service orchestration
- **Container Security** with Trivy vulnerability scanning
- **Clean Architecture** with separated concerns
- **Comprehensive Testing** with unit tests and BDD
- **CI/CD Pipeline** with GitHub Actions
- **Code Quality** with golangci-lint

## Project Structure

```
golang-template/
├── cmd/
│   ├── api/          # HTTP API server
│   └── cli/          # CLI application
├── pkg/
│   ├── shared/       # Shared utilities and types
│   └── core/         # Core business logic
├── internal/
│   ├── config/       # Configuration management
│   ├── handlers/     # HTTP handlers
│   └── middleware/   # HTTP middleware
├── features/         # BDD test features
├── bin/             # Built binaries
└── docs/            # Documentation
```

## Quick Start

### Prerequisites

- Go 1.23 or higher
- Make (optional, for using Makefile commands)
- Docker & Docker Compose (for containerized deployment)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/f00b455/golang-template.git
cd golang-template
```

2. Install dependencies:
```bash
make deps
# or
go mod download
```

### Development

#### Running the API Server

```bash
# Development mode
make dev
# or
go run cmd/api/main.go

# Build and run
make build
./bin/api-server
```

The API server will start on `http://localhost:3002` with:
- API endpoints at `/api/`
- Swagger documentation at `/documentation/`

#### Running the CLI Tool

```bash
# Development mode
go run cmd/cli/main.go --name "Your Name"

# Build and run
make build
./bin/cli-tool --name "Your Name"
```

### Available Commands

```bash
# Development
make dev              # Run API server in development mode
make build            # Build all binaries
make clean            # Clean build artifacts

# Testing
make test             # Run unit tests
make test-cover       # Run tests with coverage
make test-bdd         # Run BDD tests

# Code Quality
make lint             # Run linter
make format           # Format code

# Validation
make validate         # Full validation pipeline
make validate-quick   # Quick validation (lint + test)

# Setup
make setup            # Install dev tools and dependencies
```

## API Endpoints

### Greet API

- **GET** `/api/greet?name=World` - Get greeting message

### RSS API

- **GET** `/api/rss/spiegel/latest` - Get latest SPIEGEL headline
- **GET** `/api/rss/spiegel/top5?limit=3` - Get top N headlines (max 5)

## CLI Usage

```bash
# Basic greeting
./bin/cli-tool

# Custom name
./bin/cli-tool --name "Alice"

# Help
./bin/cli-tool --help
```

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### BDD Tests

```bash
# Run BDD features
go test ./features/...

# Run specific feature
go test ./features/ -godog.tags="@issue-1"
```

## Architecture

### Clean Code Principles

- **Single Responsibility**: Each package/function has one purpose
- **Dependency Injection**: Handlers receive their dependencies
- **Pure Functions**: Business logic is stateless where possible
- **Explicit Error Handling**: All errors are handled explicitly
- **Immutable Data**: Prefer immutable data structures

### Package Organization

- `cmd/` - Application entry points
- `pkg/` - Public APIs and libraries
- `internal/` - Private application code
- `features/` - BDD test specifications

### Testing Strategy

- **Unit Tests**: Standard Go testing with testify
- **BDD Tests**: Godog for behavior-driven development
- **Mocking**: Always use mocks for external dependencies
- **Test Database**: Use test database for integration tests

## Configuration

Set environment variables:

```bash
PORT=3002                    # API server port
ENV=development             # Environment (development/production)
SPIEGEL_RSS_URL=https://...  # RSS feed URL
GO_ENV=test                 # For testing (shorter delays)
```

## CI/CD

The project includes GitHub Actions workflows for:

### Go CI/CD Pipeline (`.github/workflows/go.yml`)

1. **Lints** code with golangci-lint
2. **Tests** with unit tests and BDD tests
3. **Builds** all binaries (API, Web, CLI)
4. **Uploads** coverage to Codecov

### Docker Security Pipeline (`.github/workflows/docker-security.yml`)

Comprehensive container security and validation:

1. **Vulnerability Scanning** with Trivy
   - Scans all Docker images (api, web, cli)
   - Detects OS and dependency vulnerabilities
   - Fails build on CRITICAL CVEs
   - Uploads findings to GitHub Security tab

2. **Secret Detection**
   - Scans for hardcoded credentials
   - Detects API keys and passwords in images

3. **Docker Compose Validation**
   - Validates docker-compose.yml syntax
   - Tests multi-service startup
   - Verifies service health

4. **Build Optimization**
   - Matrix builds for parallel execution
   - GitHub Actions cache for faster builds
   - SBOM (Software Bill of Materials) generation

**Security Severity Handling:**
- 🔴 CRITICAL → Build fails
- 🟠 HIGH → Warning, build continues
- 🟡 MEDIUM/LOW → Info only

**View Security Findings:**
```bash
# GitHub Security tab shows all vulnerability reports
# Navigate to: Repository → Security → Code scanning alerts
```

## Docker

### Running with Docker Compose

```bash
# Build and start all services
docker-compose up -d --build

# View logs
docker-compose logs -f api
docker-compose logs -f web

# Check status
docker-compose ps

# Stop services
docker-compose down
```

**Services:**
- **API** (Port 8080): REST API + Swagger docs
- **Web** (Port 3000): Frontend with RSS headlines
- **CLI**: Command-line tool

**Environment Variables:**
```yaml
# docker-compose.yml includes defaults
PORT=8080              # API server port
API_URL=http://api:8080  # Web → API communication
```

### Building Individual Containers

```bash
# API server
docker build -f cmd/api/Dockerfile -t golang-template-api .
docker run -p 8080:8080 golang-template-api

# Web server
docker build -f cmd/web/Dockerfile -t golang-template-web .
docker run -p 3000:3000 -e API_URL=http://api:8080 golang-template-web

# CLI tool
docker build -f cmd/cli/Dockerfile -t golang-template-cli .
docker run golang-template-cli --name "Docker User"
```

### Security Scanning Locally

```bash
# Install Trivy
brew install aquasecurity/trivy/trivy  # macOS
# or: wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | sudo apt-key add -

# Scan Docker image
trivy image golang-template-api:latest

# Scan for secrets only
trivy image --scanners secret golang-template-api:latest

# Scan Dockerfile
trivy config cmd/api/Dockerfile
```

## Auto-Fix Workflow Setup

The repository includes an automated CI failure fix workflow that uses Claude Code Action to automatically detect and fix failing CI checks.

### Required Secrets

You need to configure two secrets for the auto-fix workflow to function:

#### 1. CLAUDE_CODE_OAUTH_TOKEN

Claude Code OAuth token for @mention triggers:

```bash
# Setup: Run once to generate token
claude setup-token

# Add to repository secrets
gh secret set CLAUDE_CODE_OAUTH_TOKEN
# (paste token when prompted)
```

#### 2. CLAUDE_PAT

GitHub Personal Access Token for automated comments:

```bash
# Create PAT with 'repo' scope at: https://github.com/settings/tokens
gh secret set CLAUDE_PAT
# (paste PAT when prompted)
```

**Note:** The PAT must have `repo` scope and cannot start with `GITHUB_` (GitHub restriction).

### How Auto-Fix Works

1. **CI Failure Detection**: When CI fails on a PR, `auto-comment-claude.yml` automatically posts an @claude mention
2. **Error Analysis**: Claude Code Action (`claude.yml`) responds to @mention and analyzes the error logs
3. **Automated Fixes**: Claude fixes issues (lint errors, test failures, coverage gaps) and pushes commits
4. **CI Re-run**: CI automatically re-runs with the fixes applied
5. **Safety Limit**: Maximum 3 auto-fix attempts per PR to prevent infinite loops

### Workflow Configuration

Current settings in `.github/workflows/claude.yml`:

- **Max turns**: 50 (number of tool calls Claude can make)
- **Allowed tools**: Read, Edit, Write, MultiEdit, Glob, Grep, Bash (go, make, git, gh)
- **Model**: claude-opus-4-1-20250805
- **Permissions**: Can read CI logs, edit code, commit, and push changes

### Manual Triggers

You can also manually trigger Claude on any PR or issue:

```bash
# In PR or issue comments:
@claude Please fix the CI failures
@claude Review this code
@claude Help me understand this error
```

Claude can:
- Read CI logs and analyze errors
- Fix lint errors, test failures, and coverage issues
- Review code and provide suggestions
- Answer questions about the codebase

### Troubleshooting

If auto-fix doesn't work:

1. **Check secrets are configured:**
   ```bash
   gh secret list | grep -E "CLAUDE_CODE_OAUTH_TOKEN|CLAUDE_PAT"
   ```

2. **Verify workflows exist:**
   ```bash
   ls .github/workflows/auto-comment-claude.yml .github/workflows/claude.yml
   ```

3. **Check workflow permissions:** Go to Settings → Actions → General → Workflow permissions
   - Must be set to "Read and write permissions"
   - "Allow GitHub Actions to create and approve pull requests" must be enabled

4. **Review recent workflow runs:**
   ```bash
   gh run list --workflow=claude.yml --limit 5
   gh run list --workflow=auto-comment-claude.yml --limit 5
   ```

## Product Management with GitHub Projects

This project uses GitHub Projects for agile product management with epics, sprints, user stories, and release planning.

### Quick Setup

```bash
# 1. Set up GitHub Project with custom fields
./scripts/setup-github-project.sh

# 2. Read the workflow documentation
cat docs/PRODUCT_WORKFLOW.md

# 3. Create your first epic
gh issue create --template epic.md --label epic
```

### Workflow Overview

We use a hierarchical structure for organizing work:

```
Milestones (v1.0.0, v1.1.0)
├── Epics (large features, 1-3 months)
│   ├── User Stories (1-2 sprints)
│   │   └── Tasks (1-3 days)
```

### Issue Templates

- **Epic** (`.github/ISSUE_TEMPLATE/epic.md`) - Large features with sub-issues
- **User Story** (`.github/ISSUE_TEMPLATE/user_story.md`) - BDD-driven feature development
- **Feature Request** (`.github/ISSUE_TEMPLATE/feature_request.md`) - External requests
- **Bug Report** (`.github/ISSUE_TEMPLATE/bug_report.md`) - Bug tracking

### Project Board Views

1. **Product Backlog**: All work items prioritized by product owner
2. **Sprint Board**: Current sprint work with status columns
3. **Roadmap**: Timeline visualization of epics and releases
4. **Epics Overview**: Progress tracking for all epics

### Custom Fields

Our project uses custom fields for tracking:

- **Priority**: 🔴 Critical, 🟠 High, 🟡 Medium, 🟢 Low
- **Story Points**: Fibonacci sequence (1, 2, 3, 5, 8, 13, 21)
- **Epic**: Link user stories to parent epic
- **Component**: Frontend, Backend, Infrastructure, Documentation, Testing
- **Size**: XS, S, M, L, XL
- **Target Date**: Deadline tracking
- **Iteration**: 2-week sprint cycles

### Sprint Cycle (2 weeks)

```
Week 1
├── Monday: Sprint Planning
├── Daily: Standups
└── Friday: Mid-sprint check

Week 2
├── Daily: Standups
├── Thursday: Sprint Review
└── Friday: Sprint Retrospective
```

### Common Operations

```bash
# Create an epic
gh issue create --template epic.md --label epic --title "[EPIC] Feature Name"

# Create a user story
gh issue create --template user_story.md --label user-story

# Add issue to project
gh project item-add <PROJECT_NUMBER> --owner "@me" --url <ISSUE_URL>

# View project
gh project view <PROJECT_NUMBER> --owner "@me" --web

# Create milestone
gh api repos/OWNER/REPO/milestones -X POST -f title="v1.0.0" -f due_on="2025-12-31T23:59:59Z"
```

### Best Practices

**For Product Owners:**
- Refine backlog continuously
- Prioritize by business value
- Keep "Ready" column populated
- Use epics for features >1 sprint

**For Developers:**
- Update issue status as work progresses
- Link PRs to issues
- Break down stories >13 points
- Comment with progress updates

**For Teams:**
- Estimate as a team during planning
- Track velocity over time
- Hold regular refinement sessions
- Celebrate sprint completions

### Documentation

- **Complete Workflow Guide**: [docs/PRODUCT_WORKFLOW.md](docs/PRODUCT_WORKFLOW.md)
- **Setup Script**: [scripts/setup-github-project.sh](scripts/setup-github-project.sh)
- **Validation Script**: [scripts/check-github-settings.sh](scripts/check-github-settings.sh)

## Contributing

1. Follow the existing code style
2. Write tests for new features
3. Update BDD features for new user stories
4. Run `make validate` before submitting PRs

## License

MIT License - see LICENSE file for details.