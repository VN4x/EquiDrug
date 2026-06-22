package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/equidrug/equidrug/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetMacroProfile(ctx context.Context, userID uuid.UUID) (domain.MacroProfile, error) {
	var p domain.MacroProfile
	err := r.pool.QueryRow(ctx, `
		SELECT user_id, protein_g, carbs_g, fat_g, calories_kcal, COALESCE(notes, ''), updated_at
		FROM macro_profiles WHERE user_id = $1
	`, userID).Scan(&p.UserID, &p.ProteinG, &p.CarbsG, &p.FatG, &p.CaloriesKcal, &p.Notes, &p.UpdatedAt)
	if err == pgx.ErrNoRows {
		p = domain.MacroProfile{
			UserID:   userID,
			ProteinG: 130,
			CarbsG:   80,
			FatG:     25,
		}
		p.CaloriesKcal = p.ProteinG*4 + p.CarbsG*4 + p.FatG*9
		return p, nil
	}
	return p, err
}

func (r *Repository) UpsertMacroProfile(ctx context.Context, p domain.MacroProfile) (domain.MacroProfile, error) {
	if p.CaloriesKcal == 0 {
		p.CaloriesKcal = p.ProteinG*4 + p.CarbsG*4 + p.FatG*9
	}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO macro_profiles (user_id, protein_g, carbs_g, fat_g, calories_kcal, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
			protein_g = EXCLUDED.protein_g,
			carbs_g = EXCLUDED.carbs_g,
			fat_g = EXCLUDED.fat_g,
			calories_kcal = EXCLUDED.calories_kcal,
			notes = EXCLUDED.notes,
			updated_at = now()
		RETURNING user_id, protein_g, carbs_g, fat_g, calories_kcal, COALESCE(notes, ''), updated_at
	`, p.UserID, p.ProteinG, p.CarbsG, p.FatG, p.CaloriesKcal, p.Notes,
	).Scan(&p.UserID, &p.ProteinG, &p.CarbsG, &p.FatG, &p.CaloriesKcal, &p.Notes, &p.UpdatedAt)
	return p, err
}

func (r *Repository) SearchFoodReference(ctx context.Context, query, countryCode, locale string) ([]domain.FoodReference, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, nil
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, name, COALESCE(name_local, ''), COALESCE(locale, ''), COALESCE(country_code, ''),
		       serving_label, protein_g, carbs_g, fat_g, COALESCE(calories_kcal, 0), tags, source
		FROM food_reference
		WHERE name ILIKE '%' || $1 || '%'
		   OR COALESCE(name_local, '') ILIKE '%' || $1 || '%'
		   OR EXISTS (SELECT 1 FROM unnest(COALESCE(tags, '{}')) t WHERE t ILIKE '%' || $1 || '%')
		ORDER BY
		  CASE WHEN $2 != '' AND country_code = $2 THEN 0 ELSE 1 END,
		  CASE WHEN name ILIKE $1 THEN 0
		       WHEN name ILIKE $1 || '%' THEN 1
		       ELSE 2 END,
		  name
		LIMIT 8
	`, q, countryCode)
	if err != nil {
		return nil, fmt.Errorf("search food: %w", err)
	}
	defer rows.Close()
	return scanFoodReferences(rows)
}

func scanFoodReferences(rows pgx.Rows) ([]domain.FoodReference, error) {
	var items []domain.FoodReference
	for rows.Next() {
		var f domain.FoodReference
		if err := rows.Scan(
			&f.ID, &f.Name, &f.NameLocal, &f.Locale, &f.CountryCode,
			&f.ServingLabel, &f.ProteinG, &f.CarbsG, &f.FatG, &f.CaloriesKcal, &f.Tags, &f.Source,
		); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}

func (r *Repository) ListFoodLog(ctx context.Context, userID uuid.UUID, day time.Time) ([]domain.FoodLogEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, logged_on::text, food_name, COALESCE(serving_label, ''),
		       protein_g, carbs_g, fat_g, COALESCE(calories_kcal, 0), confidence, source,
		       COALESCE(image_url, ''), COALESCE(notes, ''), created_at
		FROM food_log
		WHERE user_id = $1 AND logged_on = $2::date
		ORDER BY created_at DESC
	`, userID, day.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.FoodLogEntry
	for rows.Next() {
		var e domain.FoodLogEntry
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.LoggedOn, &e.FoodName, &e.ServingLabel,
			&e.ProteinG, &e.CarbsG, &e.FatG, &e.CaloriesKcal, &e.Confidence, &e.Source,
			&e.ImageURL, &e.Notes, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (r *Repository) CreateFoodLog(ctx context.Context, e domain.FoodLogEntry) (domain.FoodLogEntry, error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.Confidence == 0 {
		e.Confidence = 0.85
	}
	if e.Source == "" {
		e.Source = "manual"
	}
	if e.CaloriesKcal == 0 {
		e.CaloriesKcal = e.ProteinG*4 + e.CarbsG*4 + e.FatG*9
	}
	loggedOn := e.LoggedOn
	if loggedOn == "" {
		loggedOn = time.Now().Format("2006-01-02")
	}

	err := r.pool.QueryRow(ctx, `
		INSERT INTO food_log (id, user_id, logged_on, food_name, serving_label, protein_g, carbs_g, fat_g,
		                      calories_kcal, confidence, source, image_url, notes)
		VALUES ($1, $2, $3::date, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING logged_on::text, created_at
	`, e.ID, e.UserID, loggedOn, e.FoodName, e.ServingLabel, e.ProteinG, e.CarbsG, e.FatG,
		e.CaloriesKcal, e.Confidence, e.Source, e.ImageURL, e.Notes,
	).Scan(&e.LoggedOn, &e.CreatedAt)
	return e, err
}

func (r *Repository) DeleteFoodLog(ctx context.Context, userID, entryID uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM food_log WHERE id = $1 AND user_id = $2`, entryID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
