package handler

import (
	"net/http"

	"github.com/equidrug/equidrug/internal/sources"
)

func (h *Handler) listDestinationSources(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	region := r.URL.Query().Get("region")
	category := r.URL.Query().Get("category")

	resp, err := sources.Resolve(country, region, category)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) listSourceCountries(w http.ResponseWriter, r *http.Request) {
	list, err := sources.ListCountries()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"countries": list,
	})
}
