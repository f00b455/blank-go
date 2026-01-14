# Issue #1: User Story: Task/Todo REST API implementieren

**Issue URL**: https://github.com/f00b455/blank-go/issues/1
**Created**: 2026-01-14T12:20:55Z
**Assignee**: Unassigned

## Description
# User Story

**Als** Entwickler/API-Nutzer  
**möchte ich** eine REST API für Tasks/Todos  
**damit ich** Aufgaben erstellen, verwalten und tracken kann

---

## Beschreibung

Implementierung einer vollständigen CRUD REST API für Task-Management mit Status-Tracking, Prioritäten und Filteroptionen.

---

## Akzeptanzkriterien

### Funktionale Anforderungen

- [ ] **Create Task** - `POST /api/v1/tasks`
  - Title (required), Description (optional)
  - Priority: low, medium, high (default: medium)
  - Status: pending, in_progress, completed (default: pending)
  - Due Date (optional)
  - Tags (optional, array)

- [ ] **Get All Tasks** - `GET /api/v1/tasks`
  - Pagination (limit, offset)
  - Filter by status, priority, tags
  - Sort by created_at, due_date, priority

- [ ] **Get Single Task** - `GET /api/v1/tasks/:id`
  - Returns 404 wenn nicht gefunden

- [ ] **Update Task** - `PUT /api/v1/tasks/:id`
  - Partial updates erlaubt
  - Returns 404 wenn nicht gefunden

- [ ] **Delete Task** - `DELETE /api/v1/tasks/:id`
  - Soft delete oder hard delete
  - Returns 404 wenn nicht gefunden

### Technische Anforderungen

- [ ] Clean Architecture (Handler → Service → Repository)
- [ ] In-Memory Storage (später erweiterbar auf SQLite/Postgres)
- [ ] Input Validation mit sinnvollen Fehlermeldungen
- [ ] Proper HTTP Status Codes (200, 201, 400, 404, 500)
- [ ] JSON Request/Response Format
- [ ] Unit Tests für Service Layer (>80% Coverage)
- [ ] BDD Feature Files gemäß CLAUDE.md

### API Response Format

```json
{
  "id": "uuid",
  "title": "Task Title",
  "description": "Optional description",
  "status": "pending|in_progress|completed",
  "priority": "low|medium|high",
  "due_date": "2024-01-15T10:00:00Z",
  "tags": ["work", "urgent"],
  "created_at": "2024-01-10T08:00:00Z",
  "updated_at": "2024-01-10T08:00:00Z"
}
```

### Error Response Format

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Title is required",
    "details": [...]
  }
}
```

---

## Definition of Done

- [ ] Alle Endpoints implementiert und getestet
- [ ] Unit Tests mit >80% Coverage
- [ ] BDD Feature Files vorhanden
- [ ] Code Review bestanden
- [ ] CI/CD Pipeline grün
- [ ] API Documentation (Swagger/OpenAPI)

---

## Technische Notizen

- Storage Interface für spätere DB-Erweiterung
- UUID für Task IDs (github.com/google/uuid)
- Gin Framework bereits vorhanden
- Folgt Clean Code Principles aus CLAUDE.md

---

**Labels:** enhancement, api, backend  
**Estimate:** Medium

## Work Log
- Branch created: issue-1-user-story-task-todo-rest-api-implementieren
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
