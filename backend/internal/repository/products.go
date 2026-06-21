package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/equidrug/equidrug/internal/domain"
	"github.com/google/uuid"
)

func (r *Repository) SearchProducts(ctx context.Context, query, countryCode string, category domain.Category) ([]domain.Product, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, nil
	}

	rows, err := r.pool.Query(ctx, `
		SELECT DISTINCT p.id, p.type, p.category, p.brand_name, p.display_name,
		       p.country_code, p.image_url, p.retailer_hint, p.price_hint, p.created_at
		FROM products p
		LEFT JOIN product_ingredients pi ON pi.product_id = p.id
		LEFT JOIN active_ingredients ai ON ai.id = pi.ingredient_id
		WHERE p.country_code = $1
		  AND ($2::text = '' OR p.category = $2)
		  AND (
		    p.display_name ILIKE '%' || $3 || '%'
		    OR p.brand_name ILIKE '%' || $3 || '%'
		    OR ai.inn ILIKE '%' || $3 || '%'
		    OR ai.name ILIKE '%' || $3 || '%'
		  )
		ORDER BY p.display_name
		LIMIT 20
	`, countryCode, string(category), q)
	if err != nil {
		return nil, fmt.Errorf("search products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID, &p.Type, &p.Category, &p.BrandName, &p.DisplayName,
			&p.CountryCode, &p.ImageURL, &p.RetailerHint, &p.PriceHint, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *Repository) FindEquivalentsByOrigin(ctx context.Context, originProductID uuid.UUID, destCountry string) ([]domain.LookupMatch, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT e.id, e.origin_product_id, e.foreign_product_id, e.confidence, e.notes, e.source, e.created_at,
		       fp.id, fp.type, fp.category, fp.brand_name, fp.display_name,
		       fp.country_code, fp.image_url, fp.retailer_hint, fp.price_hint, fp.created_at
		FROM equivalences e
		JOIN products fp ON fp.id = e.foreign_product_id
		WHERE e.origin_product_id = $1 AND fp.country_code = $2
		ORDER BY e.confidence DESC
	`, originProductID, destCountry)
	if err != nil {
		return nil, fmt.Errorf("find equivalents: %w", err)
	}
	defer rows.Close()

	var matches []domain.LookupMatch
	for rows.Next() {
		var m domain.LookupMatch
		if err := rows.Scan(
			&m.Equivalence.ID, &m.Equivalence.OriginProductID, &m.Equivalence.ForeignProductID,
			&m.Equivalence.Confidence, &m.Equivalence.Notes, &m.Equivalence.Source, &m.Equivalence.CreatedAt,
			&m.Product.ID, &m.Product.Type, &m.Product.Category, &m.Product.BrandName, &m.Product.DisplayName,
			&m.Product.CountryCode, &m.Product.ImageURL, &m.Product.RetailerHint, &m.Product.PriceHint, &m.Product.CreatedAt,
		); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

func (r *Repository) LookupByQuery(ctx context.Context, req domain.LookupRequest) ([]domain.LookupMatch, error) {
	originRows, err := r.pool.Query(ctx, `
		SELECT DISTINCT p.id
		FROM products p
		LEFT JOIN product_ingredients pi ON pi.product_id = p.id
		LEFT JOIN active_ingredients ai ON ai.id = pi.ingredient_id
		WHERE p.display_name ILIKE '%' || $1 || '%'
		   OR p.brand_name ILIKE '%' || $1 || '%'
		   OR ai.inn ILIKE '%' || $1 || '%'
		LIMIT 5
	`, req.Query)
	if err != nil {
		return nil, err
	}
	defer originRows.Close()

	var all []domain.LookupMatch
	for originRows.Next() {
		var originID uuid.UUID
		if err := originRows.Scan(&originID); err != nil {
			return nil, err
		}
		matches, err := r.FindEquivalentsByOrigin(ctx, originID, req.CountryCode)
		if err != nil {
			return nil, err
		}
		all = append(all, matches...)
	}
	if len(all) == 0 {
		products, err := r.SearchProducts(ctx, req.Query, req.CountryCode, req.Category)
		if err != nil {
			return nil, err
		}
		for _, p := range products {
			all = append(all, domain.LookupMatch{
				Product: p,
				Equivalence: domain.Equivalence{
					Confidence: 0.5,
					Source:     "direct_search",
					Notes:      "No curated equivalence; ingredient match pending",
				},
			})
		}
	}
	return all, nil
}

func (r *Repository) ListLockerItems(ctx context.Context, userID uuid.UUID) ([]domain.LockerItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, product_id, custom_name, dose_per_day, dose_unit,
		       frequency, category, source_photo, created_at
		FROM locker_items WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.LockerItem
	for rows.Next() {
		var item domain.LockerItem
		var productID *uuid.UUID
		if err := rows.Scan(
			&item.ID, &item.UserID, &productID, &item.CustomName, &item.DosePerDay,
			&item.DoseUnit, &item.Frequency, &item.Category, &item.SourcePhoto, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		if productID != nil {
			item.ProductID = *productID
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) CreateLockerItem(ctx context.Context, item domain.LockerItem) (domain.LockerItem, error) {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	var productID *uuid.UUID
	if item.ProductID != uuid.Nil {
		productID = &item.ProductID
	}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO locker_items (id, user_id, product_id, custom_name, dose_per_day, dose_unit, frequency, category, source_photo)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at
	`, item.ID, item.UserID, productID, item.CustomName, item.DosePerDay, item.DoseUnit, item.Frequency, item.Category, item.SourcePhoto,
	).Scan(&item.CreatedAt)
	return item, err
}

func (r *Repository) CreateTrip(ctx context.Context, trip domain.TripPlan) (domain.TripPlan, error) {
	if trip.ID == uuid.Nil {
		trip.ID = uuid.New()
	}
	if trip.SparePercent == 0 {
		trip.SparePercent = 5
	}
	if trip.Status == "" {
		trip.Status = "draft"
	}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO trip_plans (id, user_id, title, origin_country, dest_country, dest_city,
		                        start_date, end_date, spare_percent, preferred_brands, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING created_at
	`, trip.ID, trip.UserID, trip.Title, trip.OriginCountry, trip.DestCountry, trip.DestCity,
		trip.StartDate, trip.EndDate, trip.SparePercent, trip.PreferredBrands, trip.Status,
	).Scan(&trip.CreatedAt)
	return trip, err
}

func (r *Repository) GetTrip(ctx context.Context, id uuid.UUID) (domain.TripPlan, error) {
	var trip domain.TripPlan
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, title, origin_country, dest_country, dest_city,
		       start_date, end_date, spare_percent, preferred_brands, status, created_at
		FROM trip_plans WHERE id = $1
	`, id).Scan(
		&trip.ID, &trip.UserID, &trip.Title, &trip.OriginCountry, &trip.DestCountry, &trip.DestCity,
		&trip.StartDate, &trip.EndDate, &trip.SparePercent, &trip.PreferredBrands, &trip.Status, &trip.CreatedAt,
	)
	return trip, err
}

func (r *Repository) ListTripLineItems(ctx context.Context, tripID uuid.UUID) ([]domain.TripLineItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, trip_id, locker_item_id, origin_product_id, foreign_product_id,
		       quantity_needed, quantity_unit, estimated_price, where_to_buy, bought, notes
		FROM trip_line_items WHERE trip_id = $1
	`, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.TripLineItem
	for rows.Next() {
		var item domain.TripLineItem
		var originID, foreignID *uuid.UUID
		if err := rows.Scan(
			&item.ID, &item.TripID, &item.LockerItemID, &originID, &foreignID,
			&item.QuantityNeeded, &item.QuantityUnit, &item.EstimatedPrice, &item.WhereToBuy, &item.Bought, &item.Notes,
		); err != nil {
			return nil, err
		}
		if originID != nil {
			item.OriginProductID = *originID
		}
		if foreignID != nil {
			item.ForeignProductID = *foreignID
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) UpsertTripLineItem(ctx context.Context, item domain.TripLineItem) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	var originID, foreignID *uuid.UUID
	if item.OriginProductID != uuid.Nil {
		originID = &item.OriginProductID
	}
	if item.ForeignProductID != uuid.Nil {
		foreignID = &item.ForeignProductID
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO trip_line_items (id, trip_id, locker_item_id, origin_product_id, foreign_product_id,
		                             quantity_needed, quantity_unit, estimated_price, where_to_buy, bought, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			foreign_product_id = EXCLUDED.foreign_product_id,
			quantity_needed = EXCLUDED.quantity_needed,
			estimated_price = EXCLUDED.estimated_price,
			where_to_buy = EXCLUDED.where_to_buy,
			bought = EXCLUDED.bought,
			notes = EXCLUDED.notes
	`, item.ID, item.TripID, item.LockerItemID, originID, foreignID,
		item.QuantityNeeded, item.QuantityUnit, item.EstimatedPrice, item.WhereToBuy, item.Bought, item.Notes)
	return err
}

func (r *Repository) ListWikiEntries(ctx context.Context, userID uuid.UUID) ([]domain.LookupMatch, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT w.id, w.origin_product_id, w.foreign_product_id, w.confidence, w.notes, w.source, w.created_at,
		       op.id, op.type, op.category, op.brand_name, op.display_name, op.country_code, op.image_url,
		       op.retailer_hint, op.price_hint, op.created_at,
		       fp.id, fp.type, fp.category, fp.brand_name, fp.display_name, fp.country_code, fp.image_url,
		       fp.retailer_hint, fp.price_hint, fp.created_at
		FROM wiki_entries w
		JOIN products op ON op.id = w.origin_product_id
		JOIN products fp ON fp.id = w.foreign_product_id
		WHERE w.user_id = $1
		ORDER BY w.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.LookupMatch
	for rows.Next() {
		var m domain.LookupMatch
		var origin domain.Product
		if err := rows.Scan(
			&m.Equivalence.ID, &m.Equivalence.OriginProductID, &m.Equivalence.ForeignProductID,
			&m.Equivalence.Confidence, &m.Equivalence.Notes, &m.Equivalence.Source, &m.Equivalence.CreatedAt,
			&origin.ID, &origin.Type, &origin.Category, &origin.BrandName, &origin.DisplayName,
			&origin.CountryCode, &origin.ImageURL, &origin.RetailerHint, &origin.PriceHint, &origin.CreatedAt,
			&m.Product.ID, &m.Product.Type, &m.Product.Category, &m.Product.BrandName, &m.Product.DisplayName,
			&m.Product.CountryCode, &m.Product.ImageURL, &m.Product.RetailerHint, &m.Product.PriceHint, &m.Product.CreatedAt,
		); err != nil {
			return nil, err
		}
		_ = origin
		entries = append(entries, m)
	}
	return entries, rows.Err()
}

func (r *Repository) SaveWikiEntry(ctx context.Context, userID, originID, foreignID uuid.UUID, confidence float64, notes string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO wiki_entries (id, user_id, origin_product_id, foreign_product_id, confidence, notes, source)
		VALUES ($1, $2, $3, $4, $5, $6, 'user_lookup')
		ON CONFLICT (user_id, origin_product_id, foreign_product_id) DO UPDATE SET
			confidence = EXCLUDED.confidence, notes = EXCLUDED.notes
	`, uuid.New(), userID, originID, foreignID, confidence, notes)
	return err
}
