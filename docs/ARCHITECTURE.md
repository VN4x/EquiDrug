# EquiDrug Architecture

Travel-first SaaS for finding local equivalents of medicines, supplements, vitamins, superfoods, and diet-safe alternatives across countries.

## Problem

Travelers cannot rely on the same brand names abroad. Allegra (fexofenadine) in the US differs from bilastine-based products in Germany, and Allegra may be unavailable in Estonia and much of the EU. Users need on-the-go and ahead-of-trip answers: *what can I buy locally that matches what I use at home?*

## Core Categories

| Category | Purpose |
|----------|---------|
| **Consume** | Drugs, supplements, vitamins, superfoods — find local equivalents |
| **Avoid** | Allergens, unsuitable foods, diet traps — map naming differences and find safe alternatives |

## User Flows

### On-the-go lookup
1. User enters product name, URL, or photo
2. App uses GPS or manual location (country + city)
3. Returns equivalent products with active ingredients, photos, quantity, price hints, where to buy

### Trip planner
1. Photo medicine cabinet / supplement shelf (OCR + AI) or manual list
2. Set dosage regimen and trip duration (+5% spare)
3. Choose destination + optional preferred retailers (e.g. Vitamin Shoppe in US)
4. **Convert** → report: origin ↔ foreign item, totals, shopping list
5. Export PDF / HTML / shareable link; mark items bought with notes

### Personal wiki
Persistent dictionary: *my stuff in other countries* — grows with each lookup and trip.

## Tech Stack

| Layer | Choice |
|-------|--------|
| Frontend | SvelteKit → PWA, Capacitor for native shells |
| Backend | Go (stdlib + chi router) |
| Database | PostgreSQL 16 |
| i18n | EU: en, ru, de, fr, es — Asia: zh, ja, ko, th |
| Hosting | Bare metal, Podman (rootful), Tailscale mesh, Caddy TLS |
| Ports | Non-standard: API `16125`, web `16126`, Postgres `16127` |

## External Data (phased)

- **Phase 1**: Manual/curated equivalence DB, OCR pipeline stubs
- **Phase 2**: openFDA, RxNorm, EMA-style ingredient crosswalk
- **Phase 3**: Drug.com-style lookups, customs/prescription guides
- **Future**: Survival edition (natural sources — clearly separated, heavy disclaimers)

## Domain Model (summary)

```
User
  └── Locker (photo-derived inventory)
        └── LockerItem (product, dose, frequency)
  └── TripPlan (destination, dates, spare %)
        └── TripLineItem (needed qty, equivalent, bought, notes)
  └── EquivalenceLookup (origin product → foreign matches)

Product (canonical by active ingredient / INN)
  └── ProductVariant (brand, country, SKU, retailer)
  └── ActiveIngredient (name, strength, unit)

AllergenProfile / DietProfile
  └── AvoidRule (term in locale A → trap in locale B)
```

## API Surface (v1)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness |
| POST | `/api/v1/lookup` | On-the-go equivalence search |
| POST | `/api/v1/locker/scan` | OCR + AI locker ingest (async job) |
| GET/POST | `/api/v1/locker` | User inventory |
| POST | `/api/v1/trips` | Create trip plan |
| POST | `/api/v1/trips/{id}/convert` | Generate destination report |
| GET | `/api/v1/trips/{id}/export` | PDF/HTML export |
| GET | `/api/v1/wiki` | Personal equivalence dictionary |

## Security & Compliance

- Not medical advice — prominent disclaimers
- No diagnosis; ingredient matching only
- GDPR-ready data residency on bare metal
- Prescription/customs content: jurisdiction-specific, user-verified

## Deployment Topology

```
Internet → Tailscale funnel (optional) → Caddy :443
              ├── /api/* → Go API :16125
              ├── /*     → SvelteKit static/PWA :16126
              └── internal → Postgres :16127 (not exposed)
```

Podman quadlet or compose in `deploy/`. Secrets via env files, never in git.
