package domain

import (
	"time"

	"github.com/google/uuid"
)

type Category string

const (
	CategoryConsume Category = "consume"
	CategoryAvoid   Category = "avoid"
)

type ProductType string

const (
	ProductTypeDrug       ProductType = "drug"
	ProductTypeSupplement ProductType = "supplement"
	ProductTypeVitamin    ProductType = "vitamin"
	ProductTypeSuperfood  ProductType = "superfood"
	ProductTypeFood       ProductType = "food"
)

type ActiveIngredient struct {
	ID        uuid.UUID `json:"id"`
	INN       string    `json:"inn"`
	Name      string    `json:"name"`
	Strength  float64   `json:"strength"`
	Unit      string    `json:"unit"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID           uuid.UUID          `json:"id"`
	Type         ProductType        `json:"type"`
	Category     Category           `json:"category"`
	BrandName    string             `json:"brand_name"`
	DisplayName  string             `json:"display_name"`
	CountryCode  string             `json:"country_code"`
	Ingredients  []ActiveIngredient `json:"ingredients,omitempty"`
	ImageURL     string             `json:"image_url,omitempty"`
	RetailerHint string             `json:"retailer_hint,omitempty"`
	PriceHint    string             `json:"price_hint,omitempty"`
	CreatedAt    time.Time          `json:"created_at"`
}

type Equivalence struct {
	ID               uuid.UUID `json:"id"`
	OriginProductID  uuid.UUID `json:"origin_product_id"`
	ForeignProductID uuid.UUID `json:"foreign_product_id"`
	Confidence       float64   `json:"confidence"`
	Notes            string    `json:"notes,omitempty"`
	Source           string    `json:"source"`
	CreatedAt        time.Time `json:"created_at"`
}

type LockerItem struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	ProductID   uuid.UUID `json:"product_id,omitempty"`
	CustomName  string    `json:"custom_name"`
	DosePerDay  float64   `json:"dose_per_day"`
	DoseUnit    string    `json:"dose_unit"`
	Frequency   string    `json:"frequency"`
	Category    Category  `json:"category"`
	SourcePhoto string    `json:"source_photo,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type TripPlan struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Title           string    `json:"title"`
	OriginCountry   string    `json:"origin_country"`
	DestCountry     string    `json:"dest_country"`
	DestCity        string    `json:"dest_city,omitempty"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	SparePercent    float64   `json:"spare_percent"`
	PreferredBrands []string  `json:"preferred_brands,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type TripLineItem struct {
	ID               uuid.UUID `json:"id"`
	TripID           uuid.UUID `json:"trip_id"`
	LockerItemID     uuid.UUID `json:"locker_item_id"`
	OriginProductID  uuid.UUID `json:"origin_product_id,omitempty"`
	ForeignProductID uuid.UUID `json:"foreign_product_id,omitempty"`
	QuantityNeeded   float64   `json:"quantity_needed"`
	QuantityUnit     string    `json:"quantity_unit"`
	EstimatedPrice   string    `json:"estimated_price,omitempty"`
	WhereToBuy       string    `json:"where_to_buy,omitempty"`
	Bought           bool      `json:"bought"`
	Notes            string    `json:"notes,omitempty"`
}

type TripReportRow struct {
	LineItem     TripLineItem `json:"line_item"`
	LockerItem   LockerItem   `json:"locker_item"`
	Origin       *Product     `json:"origin,omitempty"`
	Foreign      *Product     `json:"foreign,omitempty"`
	Confidence   float64      `json:"confidence"`
	MatchNotes   string       `json:"match_notes,omitempty"`
	MatchStatus  string       `json:"match_status"`
}

type TripReportSummary struct {
	TotalItems    int `json:"total_items"`
	MatchedItems  int `json:"matched_items"`
	UnmatchedItems int `json:"unmatched_items"`
	BoughtItems   int `json:"bought_items"`
	TripDays      int `json:"trip_days"`
}

type LookupRequest struct {
	Query       string   `json:"query"`
	ImageURL    string   `json:"image_url,omitempty"`
	ProductURL  string   `json:"product_url,omitempty"`
	CountryCode string   `json:"country_code"`
	City        string   `json:"city,omitempty"`
	Latitude    float64  `json:"latitude,omitempty"`
	Longitude   float64  `json:"longitude,omitempty"`
	Category    Category `json:"category"`
}

type LookupMatch struct {
	Product     Product     `json:"product"`
	Equivalence Equivalence `json:"equivalence"`
}

type LookupResponse struct {
	Query       string        `json:"query"`
	CountryCode string        `json:"country_code"`
	Matches     []LookupMatch `json:"matches"`
	Disclaimer  string        `json:"disclaimer"`
}

type ConvertTripRequest struct {
	PreferredBrands []string `json:"preferred_brands,omitempty"`
}

type TripReport struct {
	Trip        TripPlan          `json:"trip"`
	Summary     TripReportSummary `json:"summary"`
	Rows        []TripReportRow   `json:"rows"`
	Disclaimer  string            `json:"disclaimer"`
}

type ConvertTripResponse = TripReport

type UpdateLineItemRequest struct {
	Bought *bool  `json:"bought,omitempty"`
	Notes  *string `json:"notes,omitempty"`
}

type AvoidRule struct {
	ID            uuid.UUID `json:"id"`
	LocaleFrom    string    `json:"locale_from"`
	TermFrom      string    `json:"term_from"`
	LocaleTo      string    `json:"locale_to"`
	TermTo        string    `json:"term_to"`
	SafeAlternate string    `json:"safe_alternate,omitempty"`
	Severity      string    `json:"severity"`
	Notes         string    `json:"notes,omitempty"`
}

type MacroProfile struct {
	UserID       uuid.UUID `json:"user_id"`
	ProteinG     float64   `json:"protein_g"`
	CarbsG       float64   `json:"carbs_g"`
	FatG         float64   `json:"fat_g"`
	CaloriesKcal float64   `json:"calories_kcal,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FoodReference struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	NameLocal    string    `json:"name_local,omitempty"`
	Locale       string    `json:"locale,omitempty"`
	CountryCode  string    `json:"country_code,omitempty"`
	ServingLabel string    `json:"serving_label"`
	ProteinG     float64   `json:"protein_g"`
	CarbsG       float64   `json:"carbs_g"`
	FatG         float64   `json:"fat_g"`
	CaloriesKcal float64   `json:"calories_kcal,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	Source       string    `json:"source"`
}

type FoodLogEntry struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	LoggedOn     string    `json:"logged_on"`
	FoodName     string    `json:"food_name"`
	ServingLabel string    `json:"serving_label,omitempty"`
	ProteinG     float64   `json:"protein_g"`
	CarbsG       float64   `json:"carbs_g"`
	FatG         float64   `json:"fat_g"`
	CaloriesKcal float64   `json:"calories_kcal,omitempty"`
	Confidence   float64   `json:"confidence"`
	Source       string    `json:"source"`
	ImageURL     string    `json:"image_url,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type MacroTotals struct {
	ProteinG     float64 `json:"protein_g"`
	CarbsG       float64 `json:"carbs_g"`
	FatG         float64 `json:"fat_g"`
	CaloriesKcal float64 `json:"calories_kcal"`
}

type MacroRemaining struct {
	ProteinG float64 `json:"protein_g"`
	CarbsG   float64 `json:"carbs_g"`
	FatG     float64 `json:"fat_g"`
}

type DietDaySummary struct {
	Date      string         `json:"date"`
	Profile   MacroProfile   `json:"profile"`
	Consumed  MacroTotals    `json:"consumed"`
	Remaining MacroRemaining `json:"remaining"`
	Entries   []FoodLogEntry `json:"entries"`
}

type AnalyzeFoodRequest struct {
	Query       string `json:"query"`
	MenuText    string `json:"menu_text,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Locale      string `json:"locale,omitempty"`
}

type AnalyzeFoodMatch struct {
	Reference   FoodReference `json:"reference"`
	Confidence  float64       `json:"confidence"`
	MatchReason string        `json:"match_reason"`
}

type AnalyzeFoodResponse struct {
	Query       string             `json:"query"`
	Matches     []AnalyzeFoodMatch `json:"matches"`
	Suggested   *FoodReference     `json:"suggested,omitempty"`
	VisionNote  string             `json:"vision_note,omitempty"`
	Disclaimer  string             `json:"disclaimer"`
}

type LogFoodRequest struct {
	FoodName     string  `json:"food_name"`
	ServingLabel string  `json:"serving_label,omitempty"`
	ProteinG     float64 `json:"protein_g"`
	CarbsG       float64 `json:"carbs_g"`
	FatG         float64 `json:"fat_g"`
	CaloriesKcal float64 `json:"calories_kcal,omitempty"`
	Confidence   float64 `json:"confidence,omitempty"`
	Source       string  `json:"source,omitempty"`
	ImageURL     string  `json:"image_url,omitempty"`
	Notes        string  `json:"notes,omitempty"`
	LoggedOn     string  `json:"logged_on,omitempty"`
}

type UpdateMacroProfileRequest struct {
	ProteinG     float64 `json:"protein_g"`
	CarbsG       float64 `json:"carbs_g"`
	FatG         float64 `json:"fat_g"`
	CaloriesKcal float64 `json:"calories_kcal,omitempty"`
	Notes        string  `json:"notes,omitempty"`
}

