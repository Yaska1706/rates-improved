package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Yaska1706/rates-improved/pkg/models"
	"github.com/Yaska1706/rates-improved/pkg/service"
	"github.com/gorilla/mux"
)

type RateHandler struct {
	RateService *service.RateService
}

func NewRateHandler(rateService *service.RateService) *RateHandler {
	return &RateHandler{
		RateService: rateService,
	}
}

func (h *RateHandler) GetLatestRates(w http.ResponseWriter, r *http.Request) {
	rates, err := h.RateService.GetLatestRates()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch rates: %v", err), http.StatusInternalServerError)
		return
	}

	currencyRateResponse := models.CurrencyRateResponse{
		Base:  "EUR",
		Rates: rates,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(currencyRateResponse); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *RateHandler) GetCurrencyStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.RateService.GetCurrencyStats()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch rates: %v", err), http.StatusInternalServerError)
		return
	}

	currencyStatsResponse := models.CurrencyStatsResponse{
		Base:         "EUR",
		RatesAnalyze: stats,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(currencyStatsResponse); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *RateHandler) GetCurrencyByDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	if !isValidDate(date) {
		http.Error(w, "Invalid date format. Expected format: yyyy-mm-dd", http.StatusBadRequest)
		return
	}

	rates, err := h.RateService.GetCurrencyByDate(date)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch rates: %v", err), http.StatusInternalServerError)
		return
	}

	currencyRateResponse := models.CurrencyRateResponse{
		Base:  "EUR",
		Rates: rates,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(currencyRateResponse); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func isValidDate(date string) bool {
	// Define the expected layout
	layout := "2006-01-02"

	// Attempt to parse the date string
	_, err := time.Parse(layout, date)
	return err == nil
}
