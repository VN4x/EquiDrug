# Worldwide destination sources reference

Machine-readable catalog: `backend/internal/sources/data/destination_sources.json`  
API: `GET /api/v1/sources?country=VN&category=drug`  
Optional region override: `?region=CN-XZ` (Tibet)

EquiDrug resolves **where to look** in priority order. Google Lens/Maps are adapters like any other — often last, sometimes skipped.

## Google tier legend

| Tier | Meaning |
|------|---------|
| **primary** | Google Shopping/Lens/Maps OK as late step after local sources |
| **secondary** | Prefer local marketplaces & registries; Google for confirmation |
| **limited** | Google often degraded (e.g. Russia) — use local search engines |
| **avoid** | Do not depend on Google (China, Laos, Tibet) |

## Resolution order (every country)

1. EquiDrug curated equivalence DB  
2. EquiDrug My Wiki (community buys)  
3. Official drug/food registry (if exists)  
4. National pharmacy chains & marketplaces  
5. Local search engine (Baidu, Naver, Yandex…)  
6. OpenStreetMap pharmacy POIs (global fallback)  
7. Open Food Facts / USDA (food macros)  
8. Google Lens — **confirmation only**  
9. Google Maps — **store locator last**

---

## Americas

### United States (`US`) — google: primary
| P | Source | Type | URL |
|---|--------|------|-----|
| 1 | EquiDrug curated | curated | internal |
| 3 | openFDA | registry | https://api.fda.gov/ |
| 4 | DailyMed | registry | https://dailymed.nlm.nih.gov/ |
| 5 | RxNorm | registry | https://rxnav.nlm.nih.gov/ |
| 6 | Walmart | pharmacy | https://www.walmart.com/ |
| 7 | CVS | pharmacy | https://www.cvs.com/ |
| 8 | Google Shopping | marketplace | late step |
| 9 | USDA FoodData | nutrition | food macros |
| 10–11 | Google Lens / Maps | lens/maps | confirmation |

### Canada (`CA`) — google: primary
Health Canada DPD → Shoppers Drug Mart → Open Food Facts → Google (late)

### Mexico (`MX`) — google: primary
COFEPRIS → Farmacias Guadalajara → Open Food Facts → Google (late)

### Brazil (`BR`) — google: primary
ANVISA → Drogasil/Raia → Open Food Facts → Google (late)

---

## Europe

### European Union fallback (`EU` region)
EMA medicines DB → Open Food Facts → OSM pharmacies → Google (late)

### Germany (`DE`)
EMA → Shop-Apotheke → dm-drogerie → Open Food Facts → Google (late)

### United Kingdom (`GB`)
EMA → Boots UK → Open Food Facts → Google (late)

### France / Spain (`FR`, `ES`)
EMA → Open Food Facts → OSM → Google (late)

### Estonia (`EE`)
Ravimiamet → Benu → Apotheka → Open Food Facts → Google (late)

### Russia (`RU`) — google: limited
GRLS state register → Eapteka → **Yandex** → OSM → Google Lens (degraded)

---

## East Asia

### Japan (`JP`) — google: secondary
PMDA → Matsumoto Kiyoshi → Sugi → Rakuten → Amazon.co.jp → Open Food Facts → Google (late)

### South Korea (`KR`) — google: secondary
MFDS → Olive Young → Coupang → **Naver Shopping** → Google (late)

### China (`CN`) — google: **avoid**
NMPA → **JD.com** → Taobao/Tmall → **Baidu** → WeChat pharmacy mini-programs → OSM → **skip Google Lens**

### Tibet (`CN-XZ` region) — google: **avoid**
EquiDrug offline Tibet pack → Wiki → JD → Baidu → OSM → **do not use Google**

---

## Southeast Asia

### Vietnam (`VN`) — google: secondary
MoH portal → Pharmacity → Long Châu → **Shopee VN** → **Tiki** → Lazada → Open Food Facts → Google (late)

### Thailand (`TH`) — google: secondary
Thai FDA → Boots TH → Shopee TH → Lazada TH → Open Food Facts → Google (late)

### Laos (`LA`) — google: **avoid**
EquiDrug offline Laos pack → Wiki → MoH Laos → OSM → Lens/Maps (low confidence only)

---

## South Asia

### India (`IN`) — google: secondary
CDSCO → **1mg** → PharmEasy → Open Food Facts → Google (late)

---

## Middle East

### Saudi Arabia (`SA`) — google: secondary
SFDA → Al-Nahdi Pharmacy → Google (late)

---

## Oceania

### Australia (`AU`) — google: primary
TGA → Chemist Warehouse → Open Food Facts → Google (late)

---

## Global fallback (`DEFAULT`)

Any country not explicitly listed uses:

EquiDrug curated → Wiki → Open Food Facts → OSM pharmacies → Google Lens/Maps (last)

---

## Adapter type glossary

| Type | Use for |
|------|---------|
| `curated` | INN-based equivalence (EquiDrug DB) |
| `registry` | Government drug/food approval databases |
| `marketplace` | E-commerce (Shopee, JD, 1mg…) |
| `pharmacy_chain` | Boots, dm, Pharmacity… |
| `search_engine` | Baidu, Naver, Yandex (not Google) |
| `maps` | Store locator |
| `nutrition` | Macro databases (OFF, USDA) |
| `lens` | Photo identification (confirm only) |
| `wiki` | User-confirmed buys abroad |
| `offline_pack` | Pre-downloaded data for low connectivity |

---

## Adding a new country

1. Edit `backend/internal/sources/data/destination_sources.json` — add `countries.XX.sources[]`  
2. Reuse existing `adapters` or add new adapter definitions  
3. Set `google_tier` honestly  
4. Deploy — catalog is embedded in API binary via `go:embed`  
5. Optional: add human summary rows to this doc  

## API examples

```bash
# Vietnam drugs — note Shopee/Tiki before Google
curl 'http://localhost:16125/api/v1/sources?country=VN&category=drug'

# Tibet — offline pack, no Google
curl 'http://localhost:16125/api/v1/sources?region=CN-XZ'

# Laos — wiki + offline first
curl 'http://localhost:16125/api/v1/sources?country=LA'

# List all configured countries
curl 'http://localhost:16125/api/v1/sources/countries'
```

See also: [DATA_SOURCES.md](./DATA_SOURCES.md) for architecture rationale.
