package service

import (
	"context"
	"math"
	"time"

	"github.com/equidrug/equidrug/internal/domain"
	"github.com/equidrug/equidrug/internal/repository"
	"github.com/google/uuid"
)

const disclaimer = "EquiDrug provides ingredient and naming information for travelers. This is not medical advice. Consult a healthcare professional before switching products."

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Lookup(ctx context.Context, req domain.LookupRequest) (domain.LookupResponse, error) {
	matches, err := s.repo.LookupByQuery(ctx, req)
	if err != nil {
		return domain.LookupResponse{}, err
	}
	return domain.LookupResponse{
		Query:       req.Query,
		CountryCode: req.CountryCode,
		Matches:     matches,
		Disclaimer:  disclaimer,
	}, nil
}

func (s *Service) ListLocker(ctx context.Context, userID uuid.UUID) ([]domain.LockerItem, error) {
	return s.repo.ListLockerItems(ctx, userID)
}

func (s *Service) AddLockerItem(ctx context.Context, item domain.LockerItem) (domain.LockerItem, error) {
	return s.repo.CreateLockerItem(ctx, item)
}

func (s *Service) CreateTrip(ctx context.Context, trip domain.TripPlan) (domain.TripPlan, error) {
	return s.repo.CreateTrip(ctx, trip)
}

func (s *Service) GetTripReport(ctx context.Context, tripID, userID uuid.UUID) (domain.TripReport, error) {
	trip, err := s.repo.GetTrip(ctx, tripID)
	if err != nil {
		return domain.TripReport{}, err
	}
	if trip.UserID != userID {
		return domain.TripReport{}, ErrForbidden
	}
	return s.buildReportFromTrip(ctx, trip)
}

func (s *Service) ConvertTrip(ctx context.Context, tripID uuid.UUID, userID uuid.UUID, req domain.ConvertTripRequest) (domain.TripReport, error) {
	trip, err := s.repo.GetTrip(ctx, tripID)
	if err != nil {
		return domain.TripReport{}, err
	}
	if trip.UserID != userID {
		return domain.TripReport{}, ErrForbidden
	}

	if len(req.PreferredBrands) > 0 {
		trip.PreferredBrands = req.PreferredBrands
	}

	days := calcTripDays(trip.StartDate, trip.EndDate)
	multiplier := 1 + trip.SparePercent/100

	locker, err := s.repo.ListLockerItems(ctx, userID)
	if err != nil {
		return domain.TripReport{}, err
	}

	_ = s.repo.DeleteTripLineItems(ctx, tripID)

	for _, item := range locker {
		qty := item.DosePerDay * float64(days) * multiplier
		line := domain.TripLineItem{
			ID:             uuid.New(),
			TripID:         tripID,
			LockerItemID:   item.ID,
			QuantityNeeded: math.Ceil(qty*10) / 10,
			QuantityUnit:   item.DoseUnit,
		}

		originID := item.ProductID
		if originID == uuid.Nil {
			if p, err := s.repo.FindProductByName(ctx, item.CustomName, trip.OriginCountry); err == nil {
				originID = p.ID
			}
		}
		if originID != uuid.Nil {
			line.OriginProductID = originID
			matches, err := s.repo.FindEquivalentsByOrigin(ctx, originID, trip.DestCountry)
			if err == nil && len(matches) > 0 {
				best := matches[0]
				line.ForeignProductID = best.Product.ID
				line.EstimatedPrice = best.Product.PriceHint
				line.WhereToBuy = best.Product.RetailerHint
				_ = s.repo.SaveWikiEntry(ctx, userID, originID, best.Product.ID, best.Equivalence.Confidence, best.Equivalence.Notes)
			}
		}

		if err := s.repo.UpsertTripLineItem(ctx, line); err != nil {
			return domain.TripReport{}, err
		}
	}

	_ = s.repo.UpdateTripStatus(ctx, tripID, "converted")
	trip.Status = "converted"
	return s.buildReportFromTrip(ctx, trip)
}

func (s *Service) UpdateLineItem(ctx context.Context, tripID, itemID, userID uuid.UUID, req domain.UpdateLineItemRequest) (domain.TripReport, error) {
	trip, err := s.repo.GetTrip(ctx, tripID)
	if err != nil {
		return domain.TripReport{}, err
	}
	if trip.UserID != userID {
		return domain.TripReport{}, ErrForbidden
	}
	if _, err := s.repo.UpdateTripLineItem(ctx, tripID, itemID, req.Bought, req.Notes); err != nil {
		return domain.TripReport{}, err
	}
	return s.buildReportFromTrip(ctx, trip)
}

func (s *Service) buildReportFromTrip(ctx context.Context, trip domain.TripPlan) (domain.TripReport, error) {
	lineItems, err := s.repo.ListTripLineItems(ctx, trip.ID)
	if err != nil {
		return domain.TripReport{}, err
	}

	summary := domain.TripReportSummary{
		TotalItems: len(lineItems),
		TripDays:   calcTripDays(trip.StartDate, trip.EndDate),
	}
	var rows []domain.TripReportRow

	for _, line := range lineItems {
		lockerItem, _ := s.repo.GetLockerItem(ctx, line.LockerItemID)
		row := domain.TripReportRow{
			LineItem:   line,
			LockerItem: lockerItem,
		}

		if line.OriginProductID != uuid.Nil {
			if p, err := s.repo.GetProduct(ctx, line.OriginProductID); err == nil {
				row.Origin = &p
			}
		}
		if line.ForeignProductID != uuid.Nil {
			if p, err := s.repo.GetProduct(ctx, line.ForeignProductID); err == nil {
				row.Foreign = &p
			}
			if row.Origin != nil {
				matches, _ := s.repo.FindEquivalentsByOrigin(ctx, line.OriginProductID, trip.DestCountry)
				for _, m := range matches {
					if m.Product.ID == line.ForeignProductID {
						row.Confidence = m.Equivalence.Confidence
						row.MatchNotes = m.Equivalence.Notes
						break
					}
				}
			}
			row.MatchStatus = "matched"
			summary.MatchedItems++
		} else {
			row.MatchStatus = "unmatched"
			row.MatchNotes = "No curated equivalent yet — try on-the-go lookup at destination"
			summary.UnmatchedItems++
		}

		if line.Bought {
			summary.BoughtItems++
		}
		rows = append(rows, row)
	}

	return domain.TripReport{
		Trip:       trip,
		Summary:    summary,
		Rows:       rows,
		Disclaimer: disclaimer,
	}, nil
}

func (s *Service) ListWiki(ctx context.Context, userID uuid.UUID) ([]domain.LookupMatch, error) {
	return s.repo.ListWikiEntries(ctx, userID)
}

func (s *Service) ScanLocker(_ context.Context, _ uuid.UUID, _ string) ([]domain.LockerItem, error) {
	return nil, ErrNotImplemented
}

var (
	ErrForbidden      = &ServiceError{Code: "forbidden", Message: "access denied"}
	ErrNotImplemented = &ServiceError{Code: "not_implemented", Message: "feature coming soon"}
	ErrNotFound       = &ServiceError{Code: "not_found", Message: "not found"}
)

type ServiceError struct {
	Code    string
	Message string
}

func (e *ServiceError) Error() string { return e.Message }

func ParseDemoUserID() uuid.UUID {
	return uuid.MustParse("00000000-0000-0000-0000-000000000001")
}

func calcTripDays(start, end time.Time) int {
	d := int(math.Ceil(end.Sub(start).Hours() / 24))
	if d < 1 {
		return 1
	}
	return d
}
