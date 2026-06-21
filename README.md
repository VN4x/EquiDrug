# EquiDrug

**Your medicines, supplements, and diet — found locally, anywhere.**

Travel-first SaaS that helps you buy equivalent products abroad when you cannot bring a month's supply (or forgot to pack). Born from a real problem: Allegra (fexofenadine) in the US, bilastine-based alternatives in Germany, and no Allegra on shelves in Estonia and much of the EU.

## What it does

| Flow | Description |
|------|-------------|
| **On-the-go lookup** | Enter name, URL, or photo → get local equivalents near you |
| **Locker** | OCR+AI scan of your medicine cabinet / supplement shelf (stub) |
| **Trip planner** | Duration + destination → quantities with 5% spare, shopping report |
| **Avoid** | Allergen and diet naming traps across countries |
| **My Wiki** | Personal dictionary of *your stuff in other countries* |

## Stack

- **Frontend**: SvelteKit PWA (port `16126`), Capacitor-ready
- **Backend**: Go API (port `16125`)
- **Database**: PostgreSQL 16 (port `16127`)
- **Deploy**: Podman rootful, Caddy, Tailscale-friendly
- **i18n**: en, ru, de, fr, es, zh, ja, ko, th

## Quick start (development)

### Database + API

```bash
cp .env.example deploy/.env
cd deploy && podman-compose up -d db
# wait for init, then:
cd ../backend && go run ./cmd/server
```

### Frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
# → http://localhost:16126
```

Try lookup: **Allegra** → destination **EE** (Estonia). Seed data returns Fexofast (fexofenadine).

### Full stack (Podman)

```bash
cd deploy
cp ../.env.example .env   # edit secrets
podman-compose up --build
```

## Project layout

```
backend/          Go API, migrations, seed (Allegra example)
frontend/         SvelteKit PWA + Capacitor config
deploy/           podman-compose, Caddyfile
docs/             Architecture & roadmap
```

## Roadmap

1. **Now** — Scaffold, seed data, lookup + locker + planner APIs
2. **Next** — OCR+AI locker scan, rich trip PDF/HTML export, shopping list UI
3. **Later** — openFDA / RxNorm ingestion, customs guides, auth & multi-tenant SaaS
4. **Future** — Survival edition (natural sources), prescription import guides

## Disclaimer

EquiDrug provides ingredient and naming information for travelers. **This is not medical advice.** Consult a healthcare professional before switching products.

## License

TBD
