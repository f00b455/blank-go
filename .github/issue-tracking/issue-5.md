# Issue #5: feat: Add PostgreSQL service to docker-compose for local development

**Issue URL**: https://github.com/f00b455/blank-go/issues/5
**Created**: 2026-01-14T13:52:07Z
**Assignee**: Unassigned

## Description
## User Story: PostgreSQL für lokale Entwicklung

### Als
Entwickler

### Möchte ich
PostgreSQL als Service im docker-compose.yml haben

### Damit ich
die DAX-Import-API lokal mit einer echten Datenbank testen kann

---

## Akzeptanzkriterien

### Docker Compose Service
- [ ] PostgreSQL 16 Alpine als Service hinzufügen
- [ ] Health Check für Datenbankbereitschaft
- [ ] Volume für Datenpersistenz
- [ ] Netzwerk-Konfiguration für Kommunikation mit API

### Konfiguration
- [ ] Environment Variables:
  - `POSTGRES_USER=dax_user`
  - `POSTGRES_PASSWORD=dax_password`
  - `POSTGRES_DB=dax_db`
- [ ] Port Mapping: `5432:5432`
- [ ] Named Volume: `postgres_data`

### API-Server Integration
- [ ] API-Server Service `depends_on` PostgreSQL mit health condition
- [ ] Database URL als Environment Variable an API übergeben
- [ ] Automatische Migration beim Start

### Beispiel docker-compose.yml Erweiterung

```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: dax_user
      POSTGRES_PASSWORD: dax_password
      POSTGRES_DB: dax_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U dax_user -d dax_db"]
      interval: 5s
      timeout: 5s
      retries: 5

  api:
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      DATABASE_URL: postgres://dax_user:dax_password@postgres:5432/dax_db?sslmode=disable

volumes:
  postgres_data:
```

---

## Hinweis

Dieses Issue ist unabhängig von PR #4 (DAX Import API). Die BDD-Tests in PR #4 werden mit Mock Repositories laufen. Diese PostgreSQL-Integration ist für lokale Entwicklung und manuelle Tests.

---

@claude Bitte implementiere dieses Feature.

## Work Log
- Branch created: issue-5-feat-add-postgresql-service-to-docker-compose-for-
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
