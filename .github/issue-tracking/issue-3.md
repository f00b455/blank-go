# Issue #3: feat: DAX Financial Data Import API with PostgreSQL Storage

**Issue URL**: https://github.com/f00b455/blank-go/issues/3
**Created**: 2026-01-14T13:31:01Z
**Assignee**: Unassigned

## Description
## User Story: DAX Financial Data Import Service

### Als
Finanzanalyst / API-Nutzer

### Möchte ich
Finanzdaten von DAX-Unternehmen per POST-Request in eine PostgreSQL-Datenbank importieren können

### Damit ich
historische Finanzkennzahlen persistent speichern und später für Analysen abrufen kann

---

## Akzeptanzkriterien

### Datenstruktur
- [ ] PostgreSQL Tabelle `dax` mit Spalten:
  - `id` (UUID, Primary Key)
  - `company` (VARCHAR, NOT NULL)
  - `ticker` (VARCHAR(10), NOT NULL)
  - `report_type` (VARCHAR(50)) - z.B. "income", "balance", "cashflow"
  - `metric` (VARCHAR(100), NOT NULL)
  - `year` (INTEGER, NOT NULL)
  - `value` (DECIMAL(20,2))
  - `currency` (VARCHAR(3), DEFAULT 'EUR')
  - `created_at` (TIMESTAMP)
  - `updated_at` (TIMESTAMP)

### API Endpoints
- [ ] `POST /api/v1/dax/import` - CSV-Daten importieren (multipart/form-data oder JSON-Array)
- [ ] `GET /api/v1/dax` - Alle Einträge abrufen (mit Pagination)
- [ ] `GET /api/v1/dax?ticker=SIE&year=2025` - Filtern nach Ticker/Jahr
- [ ] `GET /api/v1/dax/metrics?ticker=SIE` - Verfügbare Metriken für ein Unternehmen

### Technische Anforderungen
- [ ] PostgreSQL in docker-compose.yml integrieren
- [ ] Database Migrations (golang-migrate oder GORM AutoMigrate)
- [ ] Connection Pooling
- [ ] Bulk Insert für Performance (~1000+ Zeilen)
- [ ] Duplicate-Handling (UPSERT auf company+ticker+metric+year)

### Tests
- [ ] BDD Feature Tests für Import-Workflow
- [ ] Unit Tests für Repository Layer
- [ ] Integration Tests mit Test-PostgreSQL Container

---

## Beispiel CSV-Struktur

```csv
company,ticker,report_type,metric,year,value,currency
Siemens AG,SIE,income,EBITDA,2025,15859000000.0,EUR
Siemens AG,SIE,income,Net Income,2025,9620000000.0,EUR
```

## Beispiel API-Nutzung

```bash
# CSV importieren
curl -X POST http://localhost:3002/api/v1/dax/import \
  -F "file=@test-data/dax_2025_data.csv"

# Daten abrufen
curl "http://localhost:3002/api/v1/dax?ticker=SIE&year=2025"
```

---

@claude Bitte implementiere dieses Feature.

## Work Log
- Branch created: issue-3-feat-dax-financial-data-import-api-with-postgresql
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
