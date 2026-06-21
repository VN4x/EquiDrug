# EquiDrug product roadmap

## Phase 1 — Foundation (current)

- [x] Monorepo scaffold (Go + SvelteKit + Postgres)
- [x] Domain model: products, equivalences, locker, trips, wiki, avoid rules
- [x] Seed: Allegra / fexofenadine US → DE / EE
- [x] Mobile-first PWA shell with i18n hooks
- [x] Podman + Caddy deployment on ports 16125–16127
- [ ] User authentication
- [ ] Photo upload pipeline

## Phase 2 — Intelligence

- [ ] OCR + LLM locker/diet scan (Tesseract / cloud vision)
- [ ] openFDA, RxNorm, EMA-style ingredient crosswalk
- [ ] Geolocation → nearest pharmacy / retailer hints
- [ ] Price aggregation (where legally available)
- [ ] PDF + HTML trip report with origin ↔ foreign photos

## Phase 3 — SaaS & scale

- [ ] Multi-tenant billing
- [ ] Curated community contributions (moderated wiki)
- [ ] Prescription customs guides by country
- [ ] Capacitor iOS/Android store builds

## Phase 4 — Extensions (optional)

- [ ] Survival edition — natural sources (heavy legal/medical disclaimers)
- [ ] Drug interaction warnings (third-party API)
- [ ] Offline PWA cache for saved trips

## Naming inspiration

The Allegra case: same need (non-drowsy antihistamine), different brands and sometimes different actives (fexofenadine vs bilastine). EquiDrug maps **equivalence** by active ingredient first, brand second.
