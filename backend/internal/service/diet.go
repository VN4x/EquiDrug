package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/equidrug/equidrug/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const dietDisclaimer = "Macro estimates are approximate. Restaurant portions vary. Not medical or dietary advice — verify labels and consult a professional for restricted diets."

func (s *Service) GetMacroProfile(ctx context.Context, userID uuid.UUID) (domain.MacroProfile, error) {
	return s.repo.GetMacroProfile(ctx, userID)
}

func (s *Service) UpdateMacroProfile(ctx context.Context, userID uuid.UUID, req domain.UpdateMacroProfileRequest) (domain.MacroProfile, error) {
	return s.repo.UpsertMacroProfile(ctx, domain.MacroProfile{
		UserID:       userID,
		ProteinG:     req.ProteinG,
		CarbsG:       req.CarbsG,
		FatG:         req.FatG,
		CaloriesKcal: req.CaloriesKcal,
		Notes:        req.Notes,
	})
}

func (s *Service) GetDietDay(ctx context.Context, userID uuid.UUID, date string) (domain.DietDaySummary, error) {
	profile, err := s.repo.GetMacroProfile(ctx, userID)
	if err != nil {
		return domain.DietDaySummary{}, err
	}

	day, err := parseDay(date)
	if err != nil {
		return domain.DietDaySummary{}, err
	}

	entries, err := s.repo.ListFoodLog(ctx, userID, day)
	if err != nil {
		return domain.DietDaySummary{}, err
	}
	if entries == nil {
		entries = []domain.FoodLogEntry{}
	}

	consumed := sumMacros(entries)
	return domain.DietDaySummary{
		Date:     day.Format("2006-01-02"),
		Profile:  profile,
		Consumed: consumed,
		Remaining: domain.MacroRemaining{
			ProteinG: round1(profile.ProteinG - consumed.ProteinG),
			CarbsG:   round1(profile.CarbsG - consumed.CarbsG),
			FatG:     round1(profile.FatG - consumed.FatG),
		},
		Entries: entries,
	}, nil
}

func (s *Service) AnalyzeFood(ctx context.Context, req domain.AnalyzeFoodRequest) (domain.AnalyzeFoodResponse, error) {
	query := strings.TrimSpace(req.Query)
	if query == "" {
		query = strings.TrimSpace(req.MenuText)
	}

	resp := domain.AnalyzeFoodResponse{
		Query:      query,
		Disclaimer: dietDisclaimer,
	}

	if req.ImageURL != "" && query == "" {
		resp.VisionNote = "Photo received — vision macro estimation ships in the next release. Add a dish name or menu text for accurate lookup now."
		resp.Matches = []domain.AnalyzeFoodMatch{}
		return resp, nil
	}

	if query == "" {
		return resp, nil
	}

	refs, err := s.repo.SearchFoodReference(ctx, query, req.CountryCode, req.Locale)
	if err != nil {
		return domain.AnalyzeFoodResponse{}, err
	}

	var matches []domain.AnalyzeFoodMatch
	for i, ref := range refs {
		conf := 0.92 - float64(i)*0.05
		if conf < 0.55 {
			conf = 0.55
		}
		reason := "Curated travel food database"
		if ref.NameLocal != "" {
			reason = "Matched local name: " + ref.NameLocal
		}
		matches = append(matches, domain.AnalyzeFoodMatch{
			Reference:   ref,
			Confidence:  conf,
			MatchReason: reason,
		})
	}

	if len(matches) > 0 {
		best := matches[0].Reference
		resp.Suggested = &best
	}

	if req.ImageURL != "" {
		resp.VisionNote = "Photo noted — macros from text/menu match; full plate vision coming soon."
	}

	if req.MenuText != "" && len(matches) == 0 {
		resp.VisionNote = "Menu text logged; no exact match — adjust macros manually before adding."
	}

	return resp, nil
}

func (s *Service) LogFood(ctx context.Context, userID uuid.UUID, req domain.LogFoodRequest) (domain.DietDaySummary, error) {
	entry := domain.FoodLogEntry{
		UserID:       userID,
		LoggedOn:     req.LoggedOn,
		FoodName:     req.FoodName,
		ServingLabel: req.ServingLabel,
		ProteinG:     req.ProteinG,
		CarbsG:       req.CarbsG,
		FatG:         req.FatG,
		CaloriesKcal: req.CaloriesKcal,
		Confidence:   req.Confidence,
		Source:       req.Source,
		ImageURL:     req.ImageURL,
		Notes:        req.Notes,
	}
	if entry.Source == "" {
		entry.Source = "user_log"
	}
	if _, err := s.repo.CreateFoodLog(ctx, entry); err != nil {
		return domain.DietDaySummary{}, err
	}
	day := req.LoggedOn
	if day == "" {
		day = time.Now().Format("2006-01-02")
	}
	return s.GetDietDay(ctx, userID, day)
}

func (s *Service) DeleteFoodLog(ctx context.Context, userID, entryID uuid.UUID, date string) (domain.DietDaySummary, error) {
	if err := s.repo.DeleteFoodLog(ctx, userID, entryID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.DietDaySummary{}, ErrNotFound
		}
		return domain.DietDaySummary{}, err
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return s.GetDietDay(ctx, userID, date)
}

func sumMacros(entries []domain.FoodLogEntry) domain.MacroTotals {
	var t domain.MacroTotals
	for _, e := range entries {
		t.ProteinG += e.ProteinG
		t.CarbsG += e.CarbsG
		t.FatG += e.FatG
		t.CaloriesKcal += e.CaloriesKcal
	}
	t.ProteinG = round1(t.ProteinG)
	t.CarbsG = round1(t.CarbsG)
	t.FatG = round1(t.FatG)
	t.CaloriesKcal = round1(t.CaloriesKcal)
	return t
}

func round1(v float64) float64 {
	return float64(int(v*10+0.5)) / 10
}

func parseDay(date string) (time.Time, error) {
	if date == "" {
		return time.Now().UTC().Truncate(24 * time.Hour), nil
	}
	return time.Parse("2006-01-02", date)
}
