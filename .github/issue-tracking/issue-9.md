# Issue #9: feat: Add Weather API integration with Open-Meteo

**Issue URL**: https://github.com/f00b455/blank-go/issues/9
**Created**: 2026-01-15T10:24:15Z
**Assignee**: Unassigned

## Description
## User Story: Wetter-API Integration

### Als
API-Nutzer

### Möchte ich
aktuelle Wetterdaten über die API abrufen können

### Damit ich
Wetterdaten in meine Anwendungen integrieren kann

---

## Technische Entscheidung: Open-Meteo API

**Gewählt: [Open-Meteo](https://open-meteo.com/)**

Vorteile:
- ✅ Komplett kostenlos für nicht-kommerzielle Nutzung
- ✅ **Kein API-Key erforderlich**
- ✅ Keine Registrierung nötig
- ✅ Schnelle Response-Zeiten (<10ms)
- ✅ CORS unterstützt
- ✅ Server in Europa und Nordamerika
- ✅ Hochauflösende Daten (1-11 km)
- ✅ Historische Daten verfügbar (80+ Jahre)

API-Dokumentation: https://open-meteo.com/en/docs

---

## Akzeptanzkriterien

### API Endpoints
- [ ] `GET /api/v1/weather?lat={latitude}&lon={longitude}` - Aktuelles Wetter
- [ ] `GET /api/v1/weather/forecast?lat={latitude}&lon={longitude}&days={1-7}` - Vorhersage
- [ ] `GET /api/v1/weather/cities/{city}` - Wetter nach Stadtname (mit Geocoding)

### Response Format
```json
{
  "location": {
    "latitude": 52.52,
    "longitude": 13.41,
    "city": "Berlin",
    "timezone": "Europe/Berlin"
  },
  "current": {
    "temperature": 15.2,
    "humidity": 65,
    "wind_speed": 12.5,
    "weather_code": 3,
    "weather_description": "Partly cloudy"
  },
  "units": {
    "temperature": "°C",
    "wind_speed": "km/h",
    "humidity": "%"
  }
}
```

### Forecast Response
```json
{
  "location": {...},
  "forecast": [
    {
      "date": "2025-01-15",
      "temperature_max": 18.5,
      "temperature_min": 8.2,
      "precipitation_probability": 20,
      "weather_code": 2,
      "weather_description": "Partly cloudy"
    }
  ]
}
```

### Open-Meteo API Integration
- [ ] HTTP Client für Open-Meteo API
- [ ] Geocoding für Städtenamen (Open-Meteo Geocoding API)
- [ ] Weather Code zu Beschreibung Mapping
- [ ] Error Handling für API-Fehler
- [ ] Caching (optional, für Rate-Limiting Schutz)

### Beispiel Open-Meteo API Calls
```bash
# Aktuelles Wetter Berlin
curl "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&current=temperature_2m,relative_humidity_2m,wind_speed_10m,weather_code"

# 7-Tage Vorhersage
curl "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&daily=temperature_2m_max,temperature_2m_min,precipitation_probability_max,weather_code&timezone=Europe/Berlin"

# Geocoding (Stadt zu Koordinaten)
curl "https://geocoding-api.open-meteo.com/v1/search?name=Berlin&count=1"
```

### Code-Struktur
```
pkg/weather/
├── model.go          # Weather, Forecast, Location structs
├── client.go         # Open-Meteo HTTP Client
├── service.go        # Business Logic
└── service_test.go   # Unit Tests mit Mock Client

internal/handlers/
└── weather.go        # HTTP Handler

features/
├── weather-current.feature   # BDD Tests aktuelles Wetter
└── weather-forecast.feature  # BDD Tests Vorhersage
```

### Testing
- [ ] Unit Tests mit Mock HTTP Client (kein echter API Call)
- [ ] BDD Feature Files mit Issue-Referenz
- [ ] Integration Test (optional, gegen echte API)

---

## Quellen
- [Open-Meteo API](https://open-meteo.com/)
- [Open-Meteo GitHub](https://github.com/open-meteo/open-meteo)
- [Open-Meteo Docs](https://open-meteo.com/en/docs)

---

@claude Bitte implementiere dieses Feature.

## Work Log
- Branch created: issue-9-feat-add-weather-api-integration-with-open-meteo
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
