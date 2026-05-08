package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"url-shortener/internal/service"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler{
	return &URLHandler{
		service: service,
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct{
		OriginalURL string `json:"original_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.OriginalURL == "" {
		writeError(w, http.StatusBadRequest, "original url is required")
		return
	}

	resp, err := h.service.Shorten(req.OriginalURL)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract short code from path: /:code
	shortcode := strings.TrimPrefix(r.URL.Path, "/")
	if shortcode == "" {
		writeError(w, http.StatusBadRequest, "short code is required")
		return
	}

	originalURL, err := h.service.Redirect(shortcode)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func (h *URLHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet{
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract shortcode from path: /stats/:code
	shortcode := strings.TrimPrefix(r.URL.Path, "/stats/")
	if shortcode == "" {
		writeError(w, http.StatusBadRequest, "short code is required")
		return
	}

	stats, err := h.service.GetStats(shortcode)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}