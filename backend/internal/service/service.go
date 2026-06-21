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

func (s *Service) ConvertTrip(ctx context.Context, tripID uuid.UUID, userID uuid.UUID, req domain.ConvertTripRequest) (domain.ConvertTripResponse, error) {
	trip, err := s.repo.GetTrip(ctx, tripID)
	if err != nil {
		return domain.ConvertTripResponse{}, err
	}
	if trip.UserID != userID {
		return domain.ConvertTripResponse{}, ErrForbidden
	}

	days := int(math.Ceil(trip.EndDate.Sub(trip.StartDate).Hours() / 24))
	if days < 1 {
		days = 1
	}
	multiplier := 1 + trip.SparePercent/100

	locker, err := s.repo.ListLockerItems(ctx, userID)
	if err != nil {
		return domain.ConvertTripResponse{}, err
	}

	var lineItems []domain.TripLineItem
	var products []domain.Product

	for _, item := range locker {
		qty := item.DosePerDay * float64(days) * multiplier
		line := domain.TripLineItem{
			ID:             uuid.New(),
			TripID:         tripID,
			LockerItemID:   item.ID,
			QuantityNeeded: math.Ceil(qty*10) / 10,
			QuantityUnit:   item.DoseUnit,
		}

		if item.ProductID != uuid.Nil {
			line.OriginProductID = item.ProductID
			matches, err := s.repo.FindEquivalentsByOrigin(ctx, item.ProductID, trip.DestCountry)
			if err == nil && len(matches) > 0 {
				best := matches[0]
				line.ForeignProductID = best.Product.ID
				line.EstimatedPrice = best.Product.PriceHint
				line.WhereToBuy = best.Product.RetailerHint
				products = append(products, best.Product)
				_ = s.repo.SaveWikiEntry(ctx, userID, item.ProductID, best.Product.ID, best.Equivalence.Confidence, "")
			}
		}

		if err := s.repo.UpsertTripLineItem(ctx, line); err != nil {
			return domain.ConvertTripResponse{}, err
		}
		lineItems = append(lineItems, line)
	}

	trip.Status = "converted"
	return domain.ConvertTripResponse{
		Trip:      trip,
		LineItems: lineItems,
		Products:  products,
	}, nil
}

func (s *Service) ListWiki(ctx context.Context, userID uuid.UUID) ([]domain.LookupMatch, error) {
	return s.repo.ListWikiEntries(ctx, userID)
}

// ScanLocker is a stub for OCR+AI pipeline (phase 2).
func (s *Service) ScanLocker(_ context.Context, _ uuid.UUID, _ string) ([]domain.LockerItem, error) {
	return nil, ErrNotImplemented
}

var (
	ErrForbidden      = &ServiceError{Code: "forbidden", Message: "access denied"}
	ErrNotImplemented = &ServiceError{Code: "not_implemented", Message: "feature coming soon"}
)

type ServiceError struct {
	Code    string
	Message string
}

func (e *ServiceError) Error() string { return e.Message }

func ParseDemoUserID() uuid.UUID {
	// Demo user until auth is wired
	return uuid.MustParse("00000000-0000-0000-0000-000000000001")
}

func TripDays(start, end time.Time) int {
	d := int(math.Ceil(end.Sub(start).Hours() / 24))
	if d < 1 {
		return 1
	}
	return d
}
