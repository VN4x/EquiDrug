package repository

import (
	"context"
	"fmt"

	"github.com/equidrug/equidrug/internal/domain"
	"github.com/google/uuid"
)

func scanProduct(row interface {
	Scan(dest ...any) error
}) (domain.Product, error) {
	var p domain.Product
	err := row.Scan(
		&p.ID, &p.Type, &p.Category, &p.BrandName, &p.DisplayName,
		&p.CountryCode, &p.ImageURL, &p.RetailerHint, &p.PriceHint, &p.CreatedAt,
	)
	return p, err
}

func (r *Repository) GetProduct(ctx context.Context, id uuid.UUID) (domain.Product, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, type, category, brand_name, display_name, country_code,
		       image_url, retailer_hint, price_hint, created_at
		FROM products WHERE id = $1
	`, id)
	p, err := scanProduct(row)
	if err != nil {
		return domain.Product{}, fmt.Errorf("get product: %w", err)
	}
	return p, nil
}

func (r *Repository) FindProductByName(ctx context.Context, name, countryCode string) (domain.Product, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, type, category, brand_name, display_name, country_code,
		       image_url, retailer_hint, price_hint, created_at
		FROM products
		WHERE country_code = $2
		  AND (display_name ILIKE '%' || $1 || '%' OR brand_name ILIKE '%' || $1 || '%')
		ORDER BY CASE WHEN brand_name ILIKE $1 THEN 0 ELSE 1 END
		LIMIT 1
	`, name, countryCode)
	return scanProduct(row)
}

func (r *Repository) GetLockerItem(ctx context.Context, id uuid.UUID) (domain.LockerItem, error) {
	var item domain.LockerItem
	var productID *uuid.UUID
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, product_id, custom_name, dose_per_day, dose_unit,
		       frequency, category, source_photo, created_at
		FROM locker_items WHERE id = $1
	`, id).Scan(
		&item.ID, &item.UserID, &productID, &item.CustomName, &item.DosePerDay,
		&item.DoseUnit, &item.Frequency, &item.Category, &item.SourcePhoto, &item.CreatedAt,
	)
	if err != nil {
		return domain.LockerItem{}, err
	}
	if productID != nil {
		item.ProductID = *productID
	}
	return item, nil
}

func (r *Repository) UpdateTripStatus(ctx context.Context, tripID uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE trip_plans SET status = $2 WHERE id = $1`, tripID, status)
	return err
}

func (r *Repository) UpdateTripLineItem(ctx context.Context, tripID, itemID uuid.UUID, bought *bool, notes *string) (domain.TripLineItem, error) {
	var item domain.TripLineItem
	var originID, foreignID *uuid.UUID

	// Fetch current
	err := r.pool.QueryRow(ctx, `
		SELECT id, trip_id, locker_item_id, origin_product_id, foreign_product_id,
		       quantity_needed, quantity_unit, estimated_price, where_to_buy, bought, notes
		FROM trip_line_items WHERE id = $1 AND trip_id = $2
	`, itemID, tripID).Scan(
		&item.ID, &item.TripID, &item.LockerItemID, &originID, &foreignID,
		&item.QuantityNeeded, &item.QuantityUnit, &item.EstimatedPrice, &item.WhereToBuy, &item.Bought, &item.Notes,
	)
	if err != nil {
		return domain.TripLineItem{}, err
	}
	if originID != nil {
		item.OriginProductID = *originID
	}
	if foreignID != nil {
		item.ForeignProductID = *foreignID
	}

	if bought != nil {
		item.Bought = *bought
	}
	if notes != nil {
		item.Notes = *notes
	}

	_, err = r.pool.Exec(ctx, `
		UPDATE trip_line_items SET bought = $3, notes = $4 WHERE id = $1 AND trip_id = $2
	`, itemID, tripID, item.Bought, item.Notes)
	return item, err
}

func (r *Repository) DeleteTripLineItems(ctx context.Context, tripID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM trip_line_items WHERE trip_id = $1`, tripID)
	return err
}
