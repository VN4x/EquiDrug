package sources

import (
	_ "embed"
	"encoding/json"
	"sort"
	"strings"
	"sync"
)

//go:embed data/destination_sources.json
var catalogJSON []byte

type Catalog struct {
	Version      string                    `json:"version"`
	Updated      string                    `json:"updated"`
	AdapterTypes map[string]string         `json:"adapter_types"`
	Adapters     map[string]Adapter        `json:"adapters"`
	Countries    map[string]CountryConfig  `json:"countries"`
	Regions      map[string]RegionConfig   `json:"regions"`
}

type Adapter struct {
	ID           string   `json:"-"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	URL          string   `json:"url"`
	Categories   []string `json:"categories,omitempty"`
	APIAvailable bool     `json:"api_available,omitempty"`
	Google       bool     `json:"google,omitempty"`
	Notes        string   `json:"notes,omitempty"`
}

type SourceRef struct {
	AdapterID  string   `json:"adapter_id"`
	Priority   int      `json:"priority"`
	Categories []string `json:"categories,omitempty"`
	Notes      string   `json:"notes,omitempty"`
}

type CountryConfig struct {
	Name       string      `json:"name"`
	GoogleTier string      `json:"google_tier"`
	Sources    []SourceRef `json:"sources"`
}

type RegionConfig struct {
	Name          string      `json:"name"`
	ParentCountry string      `json:"parent_country,omitempty"`
	GoogleTier    string      `json:"google_tier"`
	Notes         string      `json:"notes,omitempty"`
	Sources       []SourceRef `json:"sources"`
}

type ResolvedSource struct {
	Priority     int      `json:"priority"`
	AdapterID    string   `json:"adapter_id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	TypeLabel    string   `json:"type_label"`
	URL          string   `json:"url"`
	Categories   []string `json:"categories"`
	APIAvailable bool     `json:"api_available"`
	Google       bool     `json:"google"`
	Notes        string   `json:"notes,omitempty"`
}

type ResolveResponse struct {
	CountryCode string           `json:"country_code"`
	CountryName string           `json:"country_name"`
	RegionCode  string           `json:"region_code,omitempty"`
	RegionName  string           `json:"region_name,omitempty"`
	GoogleTier  string           `json:"google_tier"`
	Category    string           `json:"category,omitempty"`
	Sources     []ResolvedSource `json:"sources"`
	Guidance    string           `json:"guidance"`
	CatalogVer  string           `json:"catalog_version"`
}

var (
	once    sync.Once
	catalog Catalog
	loadErr error
)

func CatalogData() (Catalog, error) {
	once.Do(func() {
		loadErr = json.Unmarshal(catalogJSON, &catalog)
		if loadErr != nil {
			return
		}
		for id, a := range catalog.Adapters {
			a.ID = id
			catalog.Adapters[id] = a
		}
	})
	return catalog, loadErr
}

func Resolve(countryCode, regionCode, category string) (ResolveResponse, error) {
	c, err := CatalogData()
	if err != nil {
		return ResolveResponse{}, err
	}

	cc := strings.ToUpper(strings.TrimSpace(countryCode))
	rc := strings.ToUpper(strings.TrimSpace(regionCode))
	cat := strings.ToLower(strings.TrimSpace(category))

	var cfg CountryConfig
	var regionName string
	var found bool

	if rc != "" {
		if r, ok := c.Regions[rc]; ok {
			found = true
			regionName = r.Name
			cfg = CountryConfig{Name: r.Name, GoogleTier: r.GoogleTier, Sources: r.Sources}
		}
	}
	if !found && cc != "" {
		if country, ok := c.Countries[cc]; ok {
			found = true
			cfg = country
		}
	}
	if !found {
		if r, ok := c.Regions["DEFAULT"]; ok {
			cfg = CountryConfig{Name: "Global fallback", GoogleTier: r.GoogleTier, Sources: r.Sources}
			if cc == "" {
				cc = "DEFAULT"
			}
		}
	}
	if !found && cc != "" {
		if r, ok := c.Regions["EU"]; ok && isEU(cc) {
			cfg = CountryConfig{Name: c.countriesName(cc), GoogleTier: r.GoogleTier, Sources: r.Sources}
			found = true
		}
	}

	refs := cfg.Sources
	sort.SliceStable(refs, func(i, j int) bool { return refs[i].Priority < refs[j].Priority })

	var out []ResolvedSource
	for _, ref := range refs {
		adapter, ok := c.Adapters[ref.AdapterID]
		if !ok {
			continue
		}
		cats := adapter.Categories
		if len(ref.Categories) > 0 {
			cats = ref.Categories
		}
		if cat != "" && adapter.Type != "curated" && adapter.Type != "wiki" && len(cats) > 0 && !categoryMatch(cats, cat) {
			continue
		}
		notes := ref.Notes
		if notes == "" {
			notes = adapter.Notes
		}
		out = append(out, ResolvedSource{
			Priority:     ref.Priority,
			AdapterID:    ref.AdapterID,
			Name:         adapter.Name,
			Type:         adapter.Type,
			TypeLabel:    c.AdapterTypes[adapter.Type],
			URL:          adapter.URL,
			Categories:   cats,
			APIAvailable: adapter.APIAvailable,
			Google:       adapter.Google,
			Notes:        notes,
		})
	}

	return ResolveResponse{
		CountryCode: cc,
		CountryName: cfg.Name,
		RegionCode:  rc,
		RegionName:  regionName,
		GoogleTier:  cfg.GoogleTier,
		Category:    cat,
		Sources:     out,
		Guidance:    guidanceForTier(cfg.GoogleTier),
		CatalogVer:  c.Version,
	}, nil
}

func ListCountries() ([]CountryListItem, error) {
	c, err := CatalogData()
	if err != nil {
		return nil, err
	}
	var list []CountryListItem
	for code, country := range c.Countries {
		list = append(list, CountryListItem{Code: code, Name: country.Name, GoogleTier: country.GoogleTier})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
	return list, nil
}

type CountryListItem struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	GoogleTier string `json:"google_tier"`
}

func (c Catalog) countriesName(code string) string {
	if country, ok := c.Countries[code]; ok {
		return country.Name
	}
	return code
}

func isEU(code string) bool {
	eu := map[string]bool{
		"AT": true, "BE": true, "BG": true, "HR": true, "CY": true, "CZ": true,
		"DK": true, "EE": true, "FI": true, "FR": true, "DE": true, "GR": true,
		"HU": true, "IE": true, "IT": true, "LV": true, "LT": true, "LU": true,
		"MT": true, "NL": true, "PL": true, "PT": true, "RO": true, "SK": true,
		"SI": true, "ES": true, "SE": true,
	}
	return eu[code]
}

func contains(slice []string, v string) bool {
	for _, s := range slice {
		if s == v {
			return true
		}
	}
	return false
}

func categoryMatch(adapterCats []string, query string) bool {
	if query == "" {
		return true
	}
	if contains(adapterCats, query) {
		return true
	}
	// drug registry useful for supplement/vitamin ingredient checks
	if query == "supplement" || query == "vitamin" {
		return contains(adapterCats, "drug") || contains(adapterCats, "supplement") || contains(adapterCats, "vitamin")
	}
	return false
}

func guidanceForTier(tier string) string {
	switch tier {
	case "avoid":
		return "Local search engines and EquiDrug offline packs first. Do not rely on Google Lens or Maps."
	case "limited":
		return "Google may be degraded. Prefer local search engines and pharmacy chains listed below."
	case "secondary":
		return "Use local marketplaces and registries before Google. Lens is for confirmation only."
	default:
		return "Curated EquiDrug data first, then official registries and local retailers. Google is a late step."
	}
}
