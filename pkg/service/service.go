package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Yaska1706/rates-improved/pkg/converter"
	"github.com/Yaska1706/rates-improved/pkg/database"
	"github.com/Yaska1706/rates-improved/pkg/models"
)

type RateService struct {
	logger   *slog.Logger
	ctx      context.Context
	rateRepo *database.RateRepo
}

func NewRateService(ctx context.Context, logger *slog.Logger, rateRepo *database.RateRepo) *RateService {
	return &RateService{
		logger:   logger,
		ctx:      ctx,
		rateRepo: rateRepo,
	}
}

func (rs *RateService) FetchRates(url string) ([]models.CurrencyRate, error) {
	var currencyRates []models.CurrencyRate
	rs.logger.InfoContext(rs.ctx, "Fetching data ...")

	envelope, err := models.FetchXML(url)
	if err != nil {
		rs.logger.ErrorContext(rs.ctx, "failed to fetch xml", "error", err.Error())
		return nil, err
	}

	if envelope != nil {
		currencyRates = converter.XMLtoCurrencyRate(*envelope)
	}

	rs.logger.InfoContext(rs.ctx, "Fetching successful")

	return currencyRates, nil
}

func (rs *RateService) GetLatestRates() (map[string]float64, error) {
	currencyRates, err := rs.rateRepo.GetLatestRates()
	if err != nil {
		return nil, fmt.Errorf("GetLatestRates(): %w", err)
	}

	rates := make(map[string]float64)

	for _, currencyRate := range currencyRates {
		rates[currencyRate.Currency] = currencyRate.Rate
	}

	return rates, nil
}

func (rs *RateService) GetCurrencyStats() (map[string]models.Stat, error) {
	currencyStats, err := rs.rateRepo.GetCurrencyStats()
	if err != nil {
		return nil, fmt.Errorf("GetCurrencyStats(): %w", err)
	}

	stats := make(map[string]models.Stat)

	for _, currencyStat := range currencyStats {
		stats[currencyStat.Currency] = models.Stat{
			MinRate: currencyStat.MinRate,
			MaxRate: currencyStat.MaxRate,
			AvgRate: currencyStat.AvgRate,
		}
	}

	return stats, nil
}

func (rs *RateService) GetCurrencyByDate(date string) (map[string]float64, error) {

	currencyRates, err := rs.rateRepo.GetRatesByDate(converter.FormatDate(date))
	if err != nil {
		return nil, fmt.Errorf("GetCurrencyByDate: %w", err)
	}

	rates := make(map[string]float64)

	for _, currencyRate := range currencyRates {
		rates[currencyRate.Currency] = currencyRate.Rate
	}

	return rates, nil

}
