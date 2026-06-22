package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/equidrug/equidrug/internal/config"
	"github.com/equidrug/equidrug/internal/domain"
	"github.com/equidrug/equidrug/internal/repository"
	"github.com/equidrug/equidrug/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

type Handler struct {
	cfg *config.Config
	svc *service.Service
}

func New(cfg config.Config, repo *repository.Repository) *Handler {
	return &Handler{
		cfg: &cfg,
		svc: service.New(repo),
	}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   h.cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", h.health)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/lookup", h.lookup)
		r.Get("/locker", h.listLocker)
		r.Post("/locker", h.createLockerItem)
		r.Post("/locker/scan", h.scanLocker)
		r.Post("/trips", h.createTrip)
		r.Get("/trips/{id}", h.getTrip)
		r.Post("/trips/{id}/convert", h.convertTrip)
		r.Patch("/trips/{id}/items/{itemId}", h.updateLineItem)
		r.Get("/trips/{id}/export", h.exportTrip)
		r.Get("/wiki", h.listWiki)
		r.Get("/diet/profile", h.getDietProfile)
		r.Put("/diet/profile", h.updateDietProfile)
		r.Get("/diet/today", h.getDietToday)
		r.Post("/diet/analyze", h.analyzeFood)
		r.Post("/diet/log", h.logFood)
		r.Delete("/diet/log/{id}", h.deleteFoodLog)
		r.Get("/sources", h.listDestinationSources)
		r.Get("/sources/countries", h.listSourceCountries)
	})

	return r
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "equidrug-api",
	})
}

func (h *Handler) lookup(w http.ResponseWriter, r *http.Request) {
	var req domain.LookupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Query == "" && req.ImageURL == "" && req.ProductURL == "" {
		writeError(w, http.StatusBadRequest, "query, image_url, or product_url required")
		return
	}
	if req.CountryCode == "" {
		writeError(w, http.StatusBadRequest, "country_code required")
		return
	}
	if req.Category == "" {
		req.Category = domain.CategoryConsume
	}

	resp, err := h.svc.Lookup(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) listLocker(w http.ResponseWriter, r *http.Request) {
	userID := service.ParseDemoUserID()
	items, err := h.svc.ListLocker(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if items == nil {
		items = []domain.LockerItem{}
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) createLockerItem(w http.ResponseWriter, r *http.Request) {
	var item domain.LockerItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	item.UserID = service.ParseDemoUserID()
	created, err := h.svc.AddLockerItem(r.Context(), item)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) scanLocker(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ImageURL string `json:"image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	_, err := h.svc.ScanLocker(r.Context(), service.ParseDemoUserID(), body.ImageURL)
	if err != nil {
		var se *service.ServiceError
		if errors.As(err, &se) && se.Code == "not_implemented" {
			writeJSON(w, http.StatusAccepted, map[string]string{
				"status":  "queued",
				"message": "OCR+AI locker scan will be available in a future release",
			})
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) createTrip(w http.ResponseWriter, r *http.Request) {
	var trip domain.TripPlan
	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	trip.UserID = service.ParseDemoUserID()
	created, err := h.svc.CreateTrip(r.Context(), trip)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) getTrip(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}
	report, err := h.svc.GetTripReport(r.Context(), id, service.ParseDemoUserID())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *Handler) convertTrip(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}
	var req domain.ConvertTripRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.svc.ConvertTrip(r.Context(), id, service.ParseDemoUserID(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) updateLineItem(w http.ResponseWriter, r *http.Request) {
	tripID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}
	itemID, err := uuid.Parse(chi.URLParam(r, "itemId"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}
	var req domain.UpdateLineItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	report, err := h.svc.UpdateLineItem(r.Context(), tripID, itemID, service.ParseDemoUserID(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (h *Handler) exportTrip(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid trip id")
		return
	}
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "html"
	}

	report, err := h.svc.GetTripReport(r.Context(), id, service.ParseDemoUserID())
	if err != nil {
		writeServiceError(w, err)
		return
	}

	switch format {
	case "html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(service.RenderTripHTML(report)))
	case "json":
		writeJSON(w, http.StatusOK, report)
	default:
		writeError(w, http.StatusBadRequest, "format must be html or json")
	}
}

func (h *Handler) listWiki(w http.ResponseWriter, r *http.Request) {
	entries, err := h.svc.ListWiki(r.Context(), service.ParseDemoUserID())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if entries == nil {
		entries = []domain.LookupMatch{}
	}
	writeJSON(w, http.StatusOK, entries)
}

func (h *Handler) getDietProfile(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.GetMacroProfile(r.Context(), service.ParseDemoUserID())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) updateDietProfile(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdateMacroProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	p, err := h.svc.UpdateMacroProfile(r.Context(), service.ParseDemoUserID(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) getDietToday(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	summary, err := h.svc.GetDietDay(r.Context(), service.ParseDemoUserID(), date)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func (h *Handler) analyzeFood(w http.ResponseWriter, r *http.Request) {
	var req domain.AnalyzeFoodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Query == "" && req.MenuText == "" && req.ImageURL == "" {
		writeError(w, http.StatusBadRequest, "query, menu_text, or image required")
		return
	}
	resp, err := h.svc.AnalyzeFood(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) logFood(w http.ResponseWriter, r *http.Request) {
	var req domain.LogFoodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.FoodName == "" {
		writeError(w, http.StatusBadRequest, "food_name required")
		return
	}
	summary, err := h.svc.LogFood(r.Context(), service.ParseDemoUserID(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, summary)
}

func (h *Handler) deleteFoodLog(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	date := r.URL.Query().Get("date")
	summary, err := h.svc.DeleteFoodLog(r.Context(), service.ParseDemoUserID(), id, date)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func writeServiceError(w http.ResponseWriter, err error) {
	var se *service.ServiceError
	if errors.As(err, &se) {
		switch se.Code {
		case "forbidden":
			writeError(w, http.StatusForbidden, se.Message)
		case "not_found":
			writeError(w, http.StatusNotFound, se.Message)
		default:
			writeError(w, http.StatusInternalServerError, se.Message)
		}
		return
	}
	writeError(w, http.StatusInternalServerError, err.Error())
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
