# Issue #7: feat: Enforce 80% minimum test coverage in CI pipeline

**Issue URL**: https://github.com/f00b455/blank-go/issues/7
**Created**: 2026-01-15T10:11:29Z
**Assignee**: Unassigned

## Description
## User Story: Test Coverage Enforcement

### Als
Entwickler / Tech Lead

### Möchte ich
dass die CI-Pipeline eine Mindest-Testabdeckung von 80% erzwingt

### Damit
die Code-Qualität gewährleistet ist und keine ungetesteten Features in main gemergt werden

---

## Akzeptanzkriterien

### CI Pipeline Requirements
- [ ] `make validate` muss in der CI-Pipeline laufen (lint + test + test-bdd + build)
- [ ] Mindest-Testabdeckung: **80%**
- [ ] CI schlägt fehl wenn Coverage unter 80%
- [ ] Coverage-Report wird generiert und angezeigt

### Makefile Updates
- [ ] `make test` soll Coverage generieren (`-coverprofile=coverage.out`)
- [ ] `make test-cover` für Coverage-Report
- [ ] Coverage-Threshold Check integrieren

### GitHub Actions Workflow (go.yml)
- [ ] Coverage-Report generieren
- [ ] Coverage-Threshold prüfen (80% minimum)
- [ ] Coverage-Badge oder Summary in PR anzeigen
- [ ] Build schlägt fehl bei < 80% Coverage

### Beispiel Implementation

**Makefile:**
```makefile
COVERAGE_THRESHOLD=80

test:
	$(GOCMD) test -race -coverprofile=coverage.out -covermode=atomic ./...

test-cover: test
	$(GOCMD) tool cover -func=coverage.out
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$coverage < $(COVERAGE_THRESHOLD)" | bc) -eq 1 ]; then \
		echo "❌ Coverage $$coverage% is below threshold $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	else \
		echo "✅ Coverage $$coverage% meets threshold $(COVERAGE_THRESHOLD)%"; \
	fi
```

**go.yml Workflow:**
```yaml
- name: Run tests with coverage
  run: make test

- name: Check coverage threshold
  run: |
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    echo "Total coverage: ${COVERAGE}%"
    if (( $(echo "$COVERAGE < 80" | bc -l) )); then
      echo "::error::Coverage ${COVERAGE}% is below 80% threshold"
      exit 1
    fi
```

---

## Hinweis

Die pre-push Hooks mit Husky sind laut CLAUDE.md konfiguriert, aber die CI-Pipeline prüft aktuell keine Coverage-Schwelle. Dies muss synchronisiert werden.

---

@claude Bitte implementiere dieses Feature.

## Work Log
- Branch created: issue-7-feat-enforce-80-minimum-test-coverage-in-ci-pipeli
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
