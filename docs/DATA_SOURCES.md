# Local destination data strategy

How EquiDrug should find **where to buy** and **what the product looks like** abroad — including hard regions like Vietnam, Laos, and Tibet.

## Short answer

**Do not rely on Google Lens alone** as your primary data layer for Vietnam, Laos, or Tibet. Use a **tiered resolver**: curated equivalence DB first, then **region-appropriate local sources**, then visual/OCR assist — with Google as one adapter among several, not the default everywhere.

## Why Google Lens is insufficient alone

| Region | Google / Lens reality |
|--------|------------------------|
| **Vietnam** | Works in cities (Hanoi, HCMC) for many packaged goods; weaker on pharmacy-only brands, Vietnamese-only labels, and traditional medicine (thuốc nam). Google Shopping coverage is thin vs **Shopee**, **Lazada**, **Tiki**, local pharmacy chains. |
| **Laos** | Sparse retail indexing; many products sold in small **pharmacies (ຮ້ານຂາຍຢາ)** without web presence. Lens may identify generic packaging but rarely gives buy location or price. |
| **Tibet (TAR)** | Google services are **restricted or degraded**; Baidu/WeChat ecosystem dominates. Lens and Maps data are unreliable; plan on **offline curated data + human verification**. |

Lens is good at: “this photo looks like a box of X.”  
Lens is bad at: regulated drug equivalence, local brand names, pharmacy-only SKUs, and “where can I buy 30 tablets near me tonight.”

## Recommended architecture: Destination Resolver

```
User query (name / photo / URL) + GPS/country
        │
        ▼
┌───────────────────┐
│ 1. Curated graph  │  INN/active ingredient → country products (EquiDrug DB)
└─────────┬─────────┘
          │ miss
          ▼
┌───────────────────┐
│ 2. Locale router  │  Pick adapters by country/region
└─────────┬─────────┘
          │
    ┌─────┴─────┬─────────────┬──────────────┐
    ▼           ▼             ▼              ▼
 EU/regulatory  Asia marketplaces  Visual ID   Community wiki
 (openFDA,      (Shopee, Rakuten,   (Lens,      (user-confirmed
  EMA-style)     JD, local APIs)    on-device    buys abroad)
                                  OCR)
```

### Tier 1 — Curated equivalence (you own this)

- Map by **INN / active ingredient**, not brand.
- Store per country: brand, strength, Rx/OTC status, typical retailer *types* (pharmacy, supermarket, online).
- Your Allegra → Fexofast (EE) example belongs here.
- **This is the product moat** — travelers need trust, not just a Lens guess.

### Tier 2 — Locale router (local search first)

Before calling Google, route by destination:

| Region | Primary local sources | Secondary |
|--------|----------------------|-----------|
| **EU** | National drug registries, pharmacy chains (Benu, dm, Boots…) | Google Maps (where allowed) |
| **US** | openFDA, DailyMed, retailer APIs | Google Shopping |
| **Japan** | PMDA labels, **Matsumoto Kiyoshi / Sugi** store search, Rakuten | Lens for package match |
| **Vietnam** | **Shopee/Lazada/Tiki** search, pharmacy chain sites, Vietnamese MoH drug lists | Lens + OCR on Vietnamese label |
| **Laos** | Curated city pharmacy lists, NGO health guides, crowdsourced wiki | Lens (low confidence) |
| **China / Tibet** | **Baidu**, JD.com, local pharmacy apps; avoid assuming Google | On-device OCR (zh-Tibetan labels) |
| **Thailand** | FDA Thailand, Boots Thailand, Shopee TH | Lens |

**Implementation rule:** `resolve(country, city)` returns an ordered list of adapters. Google Lens/Google Shopping is **last resort** in CN/LA/Tibet; **first or second** in US/EU/JP cities.

### Tier 3 — Visual identification (Lens-class)

Use for:
- User photo of unknown package at shelf
- Confirming curated match (“does this box match Fexofast?”)

Do **not** use as sole source for:
- Dosage equivalence
- Allergy/diet safety
- Legal Rx status

Prefer **on-device OCR** (privacy, offline layovers) + server-side ingredient extraction.

### Tier 4 — Community + pro verification

- User marks “bought here” with photo + geotag → feeds wiki (moderated).
- Pharmacist/travel-clinic partners in hub cities.

## Google Lens specifically — when to use it

**Use Lens when:**
- User is in-store comparing a shelf product to their home item
- Country has good Google coverage (US, EU, JP urban, VN cities)
- Result is shown as **suggestion** with confidence + link to ingredient check

**Do not use Lens when:**
- Destination is Google-poor (Laos rural, Tibet, parts of CN)
- Product is Rx-only or label is not in Latin script without OCR pipeline
- Legal/medical liability requires curated source

## Data we should store per destination

```sql
destination_sources (
  country_code,
  adapter_type,      -- 'marketplace', 'registry', 'maps', 'lens', 'wiki'
  adapter_id,        -- 'shopee_vn', 'google_lens', 'openfda'
  priority,          -- lower = try first
  locale,
  notes
)
```

Example rows:
- `VN` → `shopee_vn` priority 1, `lens` priority 4
- `LA` → `wiki` priority 1, `lens` priority 3
- `CN` → `jd_cn` priority 1, **no** `google_maps`
- `US` → `openfda` priority 1, `google_shopping` priority 3

## Practical answer for your three examples

1. **Vietnam** — Start with Vietnamese INN/brand lists + Shopee/Tiki search by active ingredient; Lens as optional confirmation in app camera flow.
2. **Laos** — Mostly **curated + crowdsourced**; pre-download Vientiane/Luang Prabang pharmacy equivalents for common OTC items; Lens bonus only.
3. **Tibet** — Treat as **China-west ecosystem**; Baidu/JD + offline packs; do not depend on Google; warn user about connectivity and script differences.

## EquiDrug phased rollout

| Phase | Data |
|-------|------|
| **Now** | Curated seed (US/EU/EE), manual retailer hints |
| **Next** | Locale router table + Shopee/JP marketplace search adapters |
| **Later** | Lens/API visual confirm, geofenced “near me” via OSM + local chains |
| **Ongoing** | User wiki from trip shopping lists (your new report feature feeds this) |

See the full worldwide reference: [DESTINATION_SOURCES.md](./DESTINATION_SOURCES.md)  
Machine-readable catalog: `backend/internal/sources/data/destination_sources.json`  
API: `GET /api/v1/sources?country=VN&category=drug`

## Bottom line

**Google should not be step one globally.** EquiDrug should:

1. Match by **active ingredient** from your DB  
2. **Route to local search engines and marketplaces** for that country  
3. Use **Lens/OCR** to confirm what’s in the user’s hand  
4. Save successful buys back into **My Wiki** for the next traveler  

That sequence matches how travelers actually shop: they trust *what’s in the box*, not which search engine indexed it.
